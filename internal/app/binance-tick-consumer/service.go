package app_binance_tick_consumer

import (
	"context"

	adapter_binance "github.com/goregion/hexago/internal/adapter/binance"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	service_tick "github.com/goregion/hexago/internal/service/tick"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/goregion/must"
)

// Config holds the configuration for the service
type config struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func RunBlocked(ctx context.Context, logger *log.Logger) {
	logger, logStopService := logger.StartService("binance-tick-consumer")
	defer logStopService()

	// + Load config
	var serviceConfig = must.Return(
		tools.ParseEnvConfig[config](),
	)
	logger.Info("service config loaded", "config", serviceConfig)
	// - Load config

	// + Initialize clients
	redisClient, redisClose := must.Return2(
		redis.NewClient(ctx, serviceConfig.RedisURL),
	)
	defer redisClose()
	logger.Info("redis client connected")
	// - Initialize clients

	// + Initialize publishers
	var tickPublisher = adapter_redis.NewTickPublisher(redisClient)
	// - Initialize publishers

	// + Initialize applications
	var tickProcessor = service_tick.NewTickProcessor(tickPublisher)
	// - Initialize applications

	// + Initialize consumers
	var binanceListener = adapter_binance.NewLPTickConsumer(serviceConfig.Symbols, tickProcessor,
		func(err error) {
			logger.Error("failed to handle tick event", "error", err)
		},
	)
	// - Initialize consumers

	// + Consume data
	logger.LogIfError(
		binanceListener.RunBlocked(ctx),
	)
	// - Consume data
}
