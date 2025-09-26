package redis

import (
	"context"
	"time"

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, func() {}, err
	}

	return client,
		func() {
			client.Close()
		},
		nil
}
