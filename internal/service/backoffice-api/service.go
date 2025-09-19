package service_backoffice_api

import (
	"context"
	"time"

	adapter_grpc_api "github.com/goregion/hexago/internal/adapter/grpc-api/impl"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/goregion/must"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Config holds the configuration for the service
type config struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func RunBlocked(ctx context.Context, logger *log.Logger) {
	logger, logStopService := logger.StartService("backoffice-api")
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

	var ohlcTimeFrame = 1 * time.Minute

	// + Initialize publishers
	var grpcServer = adapter_grpc_api.NewServer(":50051")
	// - Initialize publishers

	// + Initialize applications
	var ohlcPublisherApp = app.NewOHLCPublisher(grpcServer)
	// - Initialize applications

	// + Initialize consumers
	var ohlcConsumer = adapter_redis.NewOHLCConsumer(redisClient,
		ohlcPublisherApp,
		ohlcTimeFrame.String(),
		serviceConfig.Symbols,
	)
	// - Initialize consumers

	// + Consume data
	var errGroup = errgroup.Group{}
	errGroup.Go(func() error {
		return errors.Wrap(
			ohlcConsumer.RunBlocked(ctx),
			"ohlc consumer stopped unexpectedly",
		)
	})
	errGroup.Go(func() error {
		return errors.Wrap(
			grpcServer.RunBlocked(ctx),
			"grpc server stopped unexpectedly",
		)
	})
	logger.LogIfError(
		errGroup.Wait(),
	)
	// - Consume data
}
