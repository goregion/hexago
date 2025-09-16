package service_feeder

import (
	"context"
	adapter_nats "feeder/internal/adapter/nats"
	adapter_ws "feeder/internal/adapter/ws"
	"feeder/internal/app"
	"feeder/pkg/jwt"
	"feeder/pkg/log"
	natsmq "feeder/pkg/nats"
	"feeder/pkg/tools"
	"net/http"
	"time"

	"github.com/goregion/must"
)

type config struct {
	HttpAddr             string `env:"HTTP_ADDR" envDefault:":8080"`
	NatsURL              string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	NatsTradeRecordTopic string `env:"NATS_TRADE_RECORD_TOPIC" envDefault:"MT4.TRADE"`
	BlocklistPath        string `env:"BLOCKLIST_PATH" envDefault:"./blocklist.txt"`
}

const serviceName = "single-trade-feeder"

func Run(ctx context.Context, logger *log.Logger) error {
	logger, ctx, logStopServiceLog := logger.StartService(ctx, serviceName)
	defer logStopServiceLog()

	var serviceConfig = must.Return(
		tools.ParseEnvConfig[config](),
	)
	logger.Info("service config", "config", serviceConfig)

	natsClient, natsClientClose := must.Return2(natsmq.NewClient(serviceName, serviceConfig.NatsURL))
	defer natsClientClose()
	logger.Info("nats client connected", "url", serviceConfig.NatsURL)

	subscription, closeSubscription := must.Return2(natsClient.NewSubscription(serviceConfig.NatsTradeRecordTopic))
	defer closeSubscription()
	logger.Info("nats subscription created", "topic", serviceConfig.NatsTradeRecordTopic)

	var tokenManager = jwt.NewTokenManager("your-secret-key", serviceConfig.BlocklistPath)

	var wsHandler = adapter_ws.NewHandler(
		ctx,
		"/ws/v1",
		func(ctx context.Context, ip string, token string) (string, error) {
			client, err := tokenManager.ParseToken(token)
			if err != nil {
				return "", err
			}
			logger.Info("ws connection opened", "ip", ip, "client", client)
			return client, nil
		},
		func(ctx context.Context, ip, id string) {
			logger.Info("ws connection closed", "ip", ip, "client", id)
		},
	)

	var tradeRecordHandler = app.NewTradeRecordHandler(wsHandler)

	go func() {
		server := &http.Server{
			Addr:         serviceConfig.HttpAddr,
			Handler:      wsHandler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		defer server.Shutdown(ctx)
		defer logger.Info("http server stopped")
		logger.Info("http server start", "addr", serviceConfig.HttpAddr)
		if err := server.ListenAndServe(); err != nil {
			logger.Error("http server error", err)
		}
	}()

	var tradeRecordSubscriber = adapter_nats.NewTradeRecordSubscriber(
		subscription,
		tradeRecordHandler,
	)

	logger.Info("start read trade records from nats")
	for range tools.IteratorInt64WithContext(ctx) {
		if err := tradeRecordSubscriber.ReadNext(ctx); err != nil {
			logger.Error("failed to read next trade record", err)
		}
	}

	return nil
}
