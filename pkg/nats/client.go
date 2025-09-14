package natsmq

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

const DefaultTimeout = nats.DefaultTimeout

type Client struct {
	*nats.Conn
}

func NewClient(name, url string) (*Client, func(), error) {
	var client = &Client{}
	conn, err := nats.Connect(url,
		nats.Name(name),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, func() {}, err
	}

	client.Conn = conn
	return client,
		func() {
			conn.Drain()
			conn.Close()
		}, nil
}

type SubscriptionInterface interface {
	ReadNext(ctx context.Context, timeout time.Duration) (string, []byte, error)
}

func (c *Client) NewSubscription(subject string) (SubscriptionInterface, func(), error) {
	sub, err := c.Conn.SubscribeSync(subject)
	if err != nil {
		return nil, func() {}, err
	}
	return &subscription{Subscription: sub},
		func() {
			sub.Drain()
			sub.Unsubscribe()
		}, nil
}
