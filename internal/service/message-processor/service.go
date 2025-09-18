package service_client_api

import (
	"context"

	adapter_mysql "github.com/goregion/hexago/internal/adapter/mysql"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/pkg/database"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/goregion/must"
)

type config struct {
	RedisURL string `env:"REDIS_URL" required:"true"`
	MySQLDSN string `env:"MYSQL_DSN" required:"true"`
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

	databaseClient, databaseClose := must.Return2(
		database.NewClient(ctx, "mysql", serviceConfig.MySQLDSN),
	)
	defer databaseClose()
	logger.Info("database client connected")

	var messagePublisher = adapter_mysql.NewMessagePublisher(ctx, databaseClient)

	var messageProcessor = app.NewMessageProcessor(messagePublisher)

	var messageConsumer = adapter_redis.NewMessageConsumer(ctx, redisClient, messageProcessor)

	logger.Info("start read trade records from nats")
	for range tools.Uint64IteratorWithContext(ctx) {
		if err := messageConsumer.ReadMessage(ctx); err != nil {
			logger.Error("failed to read message", "error", err)
		}
	}

	return nil
}
