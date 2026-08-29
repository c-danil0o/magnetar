package api

const (
	CapabilityConfig   = "config"
	CapabilityOsUpdate = "os_update"
	CapabilityGossip   = "gossip"
	CapabilityRunVms   = "run_vms"
)

const (
	OsPhaseHealthy    = "healthy"
	OsPhaseSwitching  = "switching"
	OsPhaseDegraded   = "degraded"
	OsPhaseRolledBack = "rolled_back"
)

const (
	OsAttemptApplied    = "applied"
	OsAttemptRolledBack = "rolled_back"
	OsAttemptFailed     = "failed"
)
