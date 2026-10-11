package api

import (
	"context"
	"fmt"

	natsgo "github.com/nats-io/nats.go"
)

type protoMessage interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

func connect(natsAddress string) (*natsgo.Conn, error) {
	return natsgo.Connect(fmt.Sprintf("nats://%s", natsAddress))
}

// request does one NATS request/reply. The library manages the reply inbox
// (nothing to leak) and reports a missing service as natsgo.ErrNoResponders
// instead of an empty reply that would decode as a valid empty message.
// ctx must carry a deadline.
func request(ctx context.Context, conn *natsgo.Conn, subject string, req, resp protoMessage) error {
	data, err := req.Marshal()
	if err != nil {
		return err
	}
	msg, err := conn.RequestWithContext(ctx, subject, data)
	if err != nil {
		return fmt.Errorf("%s: %w", subject, err)
	}
	return resp.Unmarshal(msg.Data)
}
