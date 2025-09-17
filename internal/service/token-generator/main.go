package service_token_generator

import (
	"context"
	"feeder/pkg/jwt"
	"feeder/pkg/log"
	"feeder/pkg/tools"
	"os"
	"time"

	"github.com/goregion/must"
	"github.com/urfave/cli/v3"
)

type config struct {
	BlocklistPath string `env:"BLOCKLIST_PATH" envDefault:"./blocklist.txt"`
}

const serviceName = "token-generator"

func Run(ctx context.Context, logger *log.Logger) error {
	logger, ctx, logStopServiceLog := logger.StartService(ctx, serviceName)
	defer logStopServiceLog()

	var serviceConfig = must.Return(
		tools.ParseEnvConfig[config](),
	)
	logger.Info("service config", "config", serviceConfig)

	var tokenManager = jwt.NewTokenManager("your-secret-key", serviceConfig.BlocklistPath)

	const dateFormat = "02/01/2006"

	cmd := &cli.Command{
		Name:      "generate-token",
		Usage:     "Generate a JWT token for a client",
		UsageText: "token-generator --client <client-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "client",
				Required: true,
				Usage:    "client name to generate token for",
			},
			&cli.StringFlag{
				Name:  "expiration",
				Usage: "token expiration date (02/01/2006 format, date only, time will be set to 00:00:00). Default is 1 year from today.",
				Value: time.Now().AddDate(1, 0, 0).Truncate(24 * time.Hour).Format(dateFormat),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := cmd.String("client")
			if client == "" {
				logger.Error("client arg is required")
				return nil
			}
			expiration := cmd.String("expiration")
			expDate, err := time.Parse(dateFormat, expiration)
			if err != nil {
				logger.Error("invalid expiration date", err)
				return err
			}
			expDate = expDate.Truncate(24 * time.Hour)
			token, err := tokenManager.GenerateToken(client, expDate)
			if err != nil {
				logger.Error("error generating token", err)
				return err
			}
			logger.Info("generated token", "client", client, "token", token, "expiration", expDate)
			return nil
		},
	}

	return cmd.Run(context.Background(), os.Args)
}
