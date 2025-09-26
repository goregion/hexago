package app_ohlc_generator

import (
	"context"
	"time"

	adapter_mysql "github.com/goregion/hexago/internal/adapter/mysql"
	adapter_redis "github.com/goregion/hexago/internal/adapter/redis"
	service_ohlc "github.com/goregion/hexago/internal/service/ohlc"
	"github.com/goregion/hexago/pkg/config"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/redis"
	sqlgen_db "github.com/goregion/hexago/pkg/sqlgen-db"
	"github.com/goregion/must"
	"github.com/pkg/errors"
)

// serviceConfig holds the configuration for the service
type serviceConfig struct {
	RedisURL string   `env:"REDIS_URL" required:"true"`
	MysqlDSN string   `env:"MYSQL_DSN" required:"true"`
	Symbols  []string `env:"SYMBOLS" required:"true"`
}

func Launch(ctx context.Context) error {
	logger, logStopService := log.MustGetLoggerFromContext(ctx).
		StartService("ohlc-generator")
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

	databaseClient, databaseClose := must.Return2(
		sqlgen_db.NewClient(ctx, "mysql", cfg.MysqlDSN),
	)
	defer databaseClose()
	logger.Info("database mysql client connected")
	// - Initialize clients

	var ohlcTimeFrame = 1 * time.Minute

	// + Initialize publishers
	var ohlcRedisPublisher = adapter_redis.NewOHLCPublisher(redisClient, "m1")
	var ohlcDBRepository = adapter_mysql.NewOHLCRepository(databaseClient, "m1")
	// - Initialize publishers

	// + Initialize applications
	var ohlcProcessor = service_ohlc.NewOHLCCreator(service_ohlc.USE_BID_PRICE,
		ohlcRedisPublisher,
		databaseClient,
		ohlcDBRepository,
	)
	// - Initialize applications

	// + Initialize consumers
	var tickRangeConsumer = adapter_redis.NewTickRangeConsumer(
		redisClient,
		time.Now(),
		ohlcTimeFrame,
		cfg.Symbols,
		ohlcProcessor,
	)
	// - Initialize consumers

	// + Consume data
	if err := tickRangeConsumer.Launch(ctx); err != nil {
		return errors.Wrap(err, "ohlc generator stopped unexpectedly")
	}
	// - Consume data

	return nil
}
