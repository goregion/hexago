package service_client_api

import (
	"context"

	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/goregion/must"
)

type config struct {
	RedisURL string `env:"REDIS_URL" required:"true"`
}

func Run(ctx context.Context, logger *log.Logger) error {
	logger, logStopService := logger.StartService("client-api")
	defer logStopService()

	var serviceConfig = must.Return(
		tools.ParseEnvConfig[config](),
	)
	logger.Info("service config loaded", "config", serviceConfig)

	redisClient, redisClose := must.Return2(
		redis.NewClient(ctx, serviceConfig.RedisURL),
	)
	defer redisClose()
	logger.Info("redis client connected")

	var messagePublisher = adapter_redis.NewMessagePublisher(redisClient)

	var messageConsumerClient = app.NewMessageConsumerClient(messagePublisher)

	logger.Info("start read trade records from nats")
	for range tools.Uint64IteratorWithContext(ctx) {
	}

	return nil
}
