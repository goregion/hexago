package service_token_generator

import (
	"context"
	"feeder/pkg/jwt"
	"feeder/pkg/log"
	"feeder/pkg/tools"
	"os"

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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := cmd.String("client")
			token, err := tokenManager.GenerateToken(client)
			if err != nil {
				logger.Error("error generating token", err)
				return err
			}
			logger.Info("generated token", "client", client, "token", token)
			return nil
		},
	}

	return cmd.Run(context.Background(), os.Args)
}
