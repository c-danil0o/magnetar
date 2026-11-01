# magnetar — c12s node registry and desired-state store

The control plane's source of truth about nodes. It owns the node records
(labels, capabilities, resource capacities, OS status), the org ownership
model (pool → org claim), and the per-node desired state that nodes pull:
the OS image and the VM fleet. All storage is etcd; all node-facing traffic
is NATS.

## Role in c12s

| Pipeline | What Magnetar does |
|---|---|
| Enrollment | node registers over NATS → gets a UUID + a record in the pool |
| Configuration | placement queries by label (`QueryNodePool`, `QueryOrgOwnedNodes`) |
| OS updates | records the desired image per node (written by gravity, pulled by star) and stores every node's `os_status` (reported by star) |
| VM fleet | stores `desired_vms` / `vm_status` on a host's node record (set by the operator via the gateway, pulled/reported by the host's star) |

## What it stores (etcd)

| Key | Content |
|---|---|
| `nodes/pool/<node-id>` | unowned node record |
| `nodes/orgs/<org>/<node-id>` | owned node record |
| `labels/pool/...`, `labels/orgs/...` | label index entries for queries |
| `desired_os/<node-id>` | operator's declared OS image (separate key on purpose: status ticks rewrite the node record every ~30 s and would race with operator writes) |

The node record carries identity, labels, `capabilities`
(`config|gossip|os_update|run_vms`), `os_status`
(`healthy|switching|degraded|rolled_back` + booted/attempted digests),
`desired_vms` and `vm_status`.

## gRPC API (`pkg/api/proto/magnetar.proto`, reflection enabled)

| Group | RPCs | Authorization |
|---|---|---|
| Node reads | `GetFromNodePool`, `GetFromOrg`, `ListNodePool`, `ListOrgOwnedNodes`, `ListAllNodes` | pool reads unauthenticated; org reads require `node.get\|org\|<org>` |
| Queries | `QueryNodePool`, `QueryOrgOwnedNodes` | `node.get` on the org |
| Ownership | `ClaimOwnership` | `node.put` on the org |
| Labels | `PutBoolLabel`, `PutFloat64Label`, `PutStringLabel`, `DeleteLabel` | `node.label.put/delete` on the node |
| VMs | `SetDesiredVms`, `GetDesiredVms`, `SetVmStatuses`, `GetVmStatuses` | `node.put` / `node.get` on the **host** node |
| OS intent | `SetDesiredOsImage` (write) | `os.update` on the node |

Caller identity travels as `authz-token` metadata; the interceptor copies it
into the context and `AuthZService` checks the JWT's `permissions` claim
(`permission|kind|id` entries).

## NATS surface

| Subject | Kind | Purpose |
|---|---|---|
| `magnetar.registration` | queue `magnetar`, req/reply | enroll a node (returns UUID) |
| `magnetar.node.update` | req/reply | capabilities + OS status sync; reply carries `notFound` (field 2) so star re-registers only on an explicit answer |
| `magnetar.os.desired.get` | req/reply | a node pulls its declared OS image (`notFound` when unknown) |
| `magnetar.vms.desired.get` | req/reply | a host pulls its desired VM set |
| `magnetar.vms.status.set` | req/reply | a host reports observed VM status |

The `pkg/messaging/nats` helpers implement two reliability rules used across
the system: status-only messages (no-responders 503) are dropped, never
decoded as an empty reply; and `pkg/api.NatsOptions()` makes every client
retry the first connect and reconnect forever.

## Claim ownership (`internal/services/nodes.go`)

`ClaimOwnership` moves selected pool nodes into an org and wires up the rest
of the control plane: Oort (authorization inheritance), Meridian (default
namespace + quota), Gravity (`JoinCluster`, which makes star join the Serf
cluster). It is **not transactional**: per-node failures are logged and the
loop continues, so a claim can end partially applied.

## Configuration

| Variable | Purpose |
|---|---|
| `NATS_ADDRESS` | NATS server (host:port) |
| `ETCD_ADDRESS` | etcd (plaintext in the lab) |
| `MAGNETAR_ADDRESS` | gRPC listen address |
| `OORT_ADDRESS` | policy/authorization service |
| `MERIDIAN_ADDRESS` | namespace/quota service |
| `GRAVITY_ADDRESS` | agent queue (cluster join) |
| `SECRET_KEY` | JWT key used to validate `authz-token` permissions |

## Build, test, run

```bash
go build ./...
go test ./...        # os service, messaging (no-responders), NATS request, publisher
go run ./cmd         # requires NATS + etcd; claim also needs oort/meridian/gravity
```

In compose (`tools/docker-compose.yml`) it runs as the `magnetar` container
with those seven variables from `tools/.env`. Ports are defined per
environment (`MAGNETAR_HOSTNAME`/`MAGNETAR_PORT`).

## Status

Stable — the node registry and the declarative stores (desired OS image,
desired VMs) that every other component leans on are etcd-backed; the
robustness gaps below are known and accepted for the lab.

## Bugs, disabled and unfinished features

### Confirmed bugs

| Where | Effect |
|---|---|
| `internal/repos/node_etcd.go` | Node record and label index are written in separate etcd operations with no transaction/rollback: a failure between them desynchronizes the index (queries miss the node) (review: high). |
| `internal/services/nodes.go` — `ClaimOwnership` | Best-effort, not atomic: per-node delete/put and Oort/Meridian/Gravity side effects are logged and skipped on failure → partially claimed orgs (review: high). |
| `internal/servers/magnetar_grpc.go` | gRPC errors are mapped to status codes with `strings.Contains(err.Error(), "not found")` instead of the existing `domain.ErrNodeNotFound` type (review R12, low). |
| No liveness/lease for registrations | Nodes that died or were wiped stay in the registry forever; after a CP reset, nodes re-register (new ids) before their host removes the old records, leaving dead pool entries until `setup-fleet` claims them (documented, low). |

### Disabled / unfinished features

- **Pool reads are unauthenticated** (`GetFromNodePool`, `ListNodePool`,
  `ListAllNodes`, `QueryNodePool`): any caller that can reach the gRPC
  endpoint sees every pool node. By design for the enrollment flow, but it
  also means the pool listing is not org-scoped (review notes it as a gap).
- **Node liveness/cleanup job** does not exist (no TTL/lease, see above).
- The Oort evaluator client is wired but not used for effective
  authorization decisions beyond the JWT permission strings.

### PoC boundaries (by design, not bugs)

- Lab transport is plaintext: NATS unauthenticated (review S1), etcd
  plaintext, gRPC insecure.
- Registration has no attestation: any NATS client can enroll a node
  (review S7; flightctl-style CSR approval is the production answer).
- `NodeUpdateReq` semantics: nil field = unchanged (a status tick must not
  clobber capabilities — that bug is documented as DEV-10).

Design details and review findings: `../plan/architecture.md`,
`../plan/review.md`, `../plan/deviations.md` (DEV-10, DEV-13, DEV-14).
