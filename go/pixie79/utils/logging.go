package utils

import (
	"log/slog"
	"os"
	"strings"
)

func SetupLogger() *slog.Logger {
	logLevelEnv := os.Getenv("LOG_LEVEL")
	logLevel := mapLogLevel(logLevelEnv)

	// Configure the logger
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func mapLogLevel(logLevelEnv string) slog.Level {
	switch strings.ToUpper(logLevelEnv) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
