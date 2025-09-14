package tools

import (
	"github.com/caarlos0/env"
)

func ParseEnvConfig[ConfigType any]() (*ConfigType, error) {
	var cfg = new(ConfigType)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
