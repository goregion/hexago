package log

import (
	"log/slog"
	"os"
)

const envEnableDebugLogLevel = "ENABLE_DEBUG_LOG_LEVEL"

func makeOptions() *slog.HandlerOptions {
	var logLevel = slog.LevelInfo
	if os.Getenv(envEnableDebugLogLevel) == "true" {
		logLevel = slog.LevelDebug
	}
	return &slog.HandlerOptions{
		Level: logLevel,
	}
}

func NewTextStdOutHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout,
		makeOptions(),
	)
}

func NewJsonStdOutHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout,
		makeOptions(),
	)
}
