package api

import (
	"fmt"
	"log"

	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"
	natsgo "github.com/nats-io/nats.go"
)

// VmAsyncClient is the host-side (Star) client for VM management: it reads
// the desired VM set from the node record and reports observed VM status
// back. Both are request/reply over the subjects in subjects.go.
type VmAsyncClient struct {
	publisher         messaging.Publisher
	subscriberFactory func(subject string) messaging.Subscriber
}

func NewVmAsyncClient(natsAddress string) (*VmAsyncClient, error) {
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
	return &VmAsyncClient{
		publisher:         publisher,
		subscriberFactory: subscriberFactory,
	}, nil
}

func (c *VmAsyncClient) GetDesiredVms(hostId string, callback VmDesiredCallback) error {
	req := &VmDesiredGetReq{HostId: hostId}
	reqMarshalled, err := req.Marshal()
	if err != nil {
		return err
	}

	replySubject := c.publisher.GenerateReplySubject()
	subscriber := c.subscriberFactory(replySubject)
	err = subscriber.Subscribe(func(msg []byte, _ string) {
		defer func() { _ = subscriber.Unsubscribe() }()
		resp := &VmDesiredGetResp{}
		if err := resp.Unmarshal(msg); err != nil {
			log.Println(err)
			return
		}
		callback(resp)
	})
	if err != nil {
		return err
	}

	err = c.publisher.Request(reqMarshalled, VmDesiredGetSubject, replySubject)
	if err != nil {
		_ = subscriber.Unsubscribe()
		return err
	}
	return nil
}

func (c *VmAsyncClient) SetVmStatuses(hostId string, statuses []*VmStatus, callback VmStatusCallback) error {
	req := &VmStatusSetReq{HostId: hostId, Statuses: statuses}
	reqMarshalled, err := req.Marshal()
	if err != nil {
		return err
	}

	replySubject := c.publisher.GenerateReplySubject()
	subscriber := c.subscriberFactory(replySubject)
	err = subscriber.Subscribe(func(msg []byte, _ string) {
		defer func() { _ = subscriber.Unsubscribe() }()
		resp := &VmStatusSetResp{}
		if err := resp.Unmarshal(msg); err != nil {
			log.Println(err)
			return
		}
		callback(resp)
	})
	if err != nil {
		return err
	}

	err = c.publisher.Request(reqMarshalled, VmStatusSetSubject, replySubject)
	if err != nil {
		_ = subscriber.Unsubscribe()
		return err
	}
	return nil
}

type VmDesiredCallback func(resp *VmDesiredGetResp)

type VmStatusCallback func(resp *VmStatusSetResp)
