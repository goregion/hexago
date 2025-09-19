package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type XAddArgs = redis.XAddArgs
type XReadArgs = redis.XReadArgs

type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, redisURL string) (*Client, func(), error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, func() {}, err
	}

	var client = &Client{
		Client: redis.NewClient(opt),
	}

	return client,
		func() {
			client.Close()
		},
		client.Ping(ctx).Err()
}
