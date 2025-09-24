package app_backoffice_api

import (
	"context"
	"time"

	adapter_grpc_api "github.com/goregion/hexago/internal/adapter/grpc-api/impl"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	service_ohlc "github.com/goregion/hexago/internal/service/ohlc"
	"github.com/goregion/hexago/pkg/config"
	"github.com/goregion/hexago/pkg/goture"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/must"
)

// serviceConfig holds the configuration for the service
type serviceConfig struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func Launch(ctx context.Context) error {
	logger, logStopService := log.MustGetLoggerFromContext(ctx).
		StartService("backoffice-api")
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

	var ohlcTimeFrame = 1 * time.Minute

	// + Initialize publishers
	var grpcServer = adapter_grpc_api.NewServer(":50051")
	// - Initialize publishers

	// + Initialize applications
	var ohlcPublisherApp = service_ohlc.NewOHLCPublisher(grpcServer)
	// - Initialize applications

	// + Initialize consumers
	var ohlcConsumer = adapter_redis.NewOHLCConsumer(redisClient,
		ohlcPublisherApp,
		ohlcTimeFrame.String(),
		cfg.Symbols,
	)
	// - Initialize consumers

	// + Consume data
	logger.LogIfError(
		goture.NewParallelGoture(ctx,
			ohlcConsumer.Launch,
			grpcServer.Launch,
		).Wait(),
	)
	// - Consume data

	return nil
}
