package app_binance_tick_consumer

import (
	"context"

	adapter_binance "github.com/goregion/hexago/internal/adapter/binance"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	service_tick "github.com/goregion/hexago/internal/service/tick"
	"github.com/goregion/hexago/pkg/config"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/must"
	"github.com/pkg/errors"
)

// serviceConfig holds the configuration for the service
type serviceConfig struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func Launch(ctx context.Context) error {
	logger, logStopService := log.MustGetLoggerFromContext(ctx).
		StartService("binance-tick-consumer")
	defer logStopService()

	// + Load config
	var cfg = must.Return(
		config.ParseEnv[serviceConfig](),
	)
	logger.Info("service config loaded", "config", cfg)
	// - Load config

	// + Initialize clients
	redisClient, redisClose := must.Return2(
		redis.NewClient(ctx, cfg.RedisURL),
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
	var binanceListener = adapter_binance.NewLPTickConsumer(cfg.Symbols, tickProcessor,
		func(err error) {
			logger.Error("failed to handle tick event", "error", err)
		},
	)
	// - Initialize consumers

	// + Consume data
	if err := binanceListener.RunBlocked(ctx); err != nil {
		return errors.Wrap(err, "failed to run binance tick consumer")
	}
	// - Consume data

	return nil
}
