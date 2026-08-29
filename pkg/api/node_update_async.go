package api

import (
	"fmt"
	"log"

	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"
	natsgo "github.com/nats-io/nats.go"
)

type NodeUpdateAsyncClient struct {
	publisher         messaging.Publisher
	subscriberFactory func(subject string) messaging.Subscriber
}

func NewNodeUpdateAsyncClient(natsAddress string) (*NodeUpdateAsyncClient, error) {
	conn, err := natsgo.Connect(fmt.Sprintf("nats://%s", natsAddress))
	if err != nil {
		return nil, err
	}
	publisher, err := nats.NewPublisher(conn)
	if err != nil {
		return nil, err
	}
	subscriberFactory := func(subject string) messaging.Subscriber {
		subscriber, _ := nats.NewSubscriber(conn, subject, "")
		return subscriber
	}
	return &NodeUpdateAsyncClient{
		publisher:         publisher,
		subscriberFactory: subscriberFactory,
	}, nil
}

func (n *NodeUpdateAsyncClient) Update(req *NodeUpdateReq, callback NodeUpdateCallback) error {
	reqMarshalled, err := req.Marshal()
	if err != nil {
		return err
	}

	replySubject := n.publisher.GenerateReplySubject()
	subscriber := n.subscriberFactory(replySubject)
	err = subscriber.Subscribe(func(msg []byte, _ string) {
		resp := &NodeUpdateResp{}
		err := resp.Unmarshal(msg)
		if err != nil {
			log.Println(err)
			return
		}
		callback(resp)
	})
	if err != nil {
		return err
	}

	err = n.publisher.Request(reqMarshalled, NodeUpdateSubject, replySubject)
	if err != nil {
		_ = subscriber.Unsubscribe()
		return err
	}
	return nil
}

type NodeUpdateCallback func(resp *NodeUpdateResp)
