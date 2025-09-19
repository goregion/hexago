package service_ohlc_generator

import (
	"context"
	"time"

	adapter_mysql "github.com/goregion/hexago/internal/adapter/mysql"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/pkg/database"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/goregion/must"
)

// Config holds the configuration for the service
type config struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	MysqlDSN string   `env:"MYSQL_DSN" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func RunBlocked(ctx context.Context, logger *log.Logger) {
	logger, logStopService := logger.StartService("ohlc-generator")
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

	databaseClient, databaseClose := must.Return2(
		database.NewClient(ctx, "mysql", serviceConfig.MysqlDSN),
	)
	defer databaseClose()
	logger.Info("database mysql client connected")
	// - Initialize clients

	var ohlcTimeFrame = 1 * time.Minute

	// + Initialize publishers
	var ohlcRedisPublisher = adapter_redis.NewOHLCPublisher(redisClient, ohlcTimeFrame.String())
	var ohlcDatabasePublisher = adapter_mysql.NewOHLCPublisher(ctx, databaseClient, ohlcTimeFrame.String())
	// - Initialize publishers

	// + Initialize applications
	var ohlcProcessor = app.NewOHLCCreator(app.USE_BID_PRICE,
		ohlcRedisPublisher,
		ohlcDatabasePublisher,
	)
	// - Initialize applications

	// + Initialize consumers
	var tickRangeConsumer = adapter_redis.NewTickRangeConsumer(redisClient, ohlcProcessor)
	// - Initialize consumers

	// + Consume data
	logger.LogIfError(
		tickRangeConsumer.RunBlocked(ctx,
			time.Now(),
			ohlcTimeFrame,
			serviceConfig.Symbols,
		),
	)
	// - Consume data
}
