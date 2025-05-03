package logging

import (
	"log/slog"
	"os"

	"github.com/rahul/api-gateway/pkg/config"
)

func ConfigureLogger(logLevel slog.Level) *slog.Logger {
	levelVar := new(slog.LevelVar)
	levelVar.Set(logLevel)

	handlerOptions := &slog.HandlerOptions{
		Level: levelVar,
	}

	return slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
}

func UpdateLogLevel(logger *slog.Logger, logLevelStr config.LogLevel) error {
	logLevel := new(slog.LevelVar)

	// Convert string log level to slog.Level
	switch logLevelStr {
	case config.Debug:
		logLevel.Set(slog.LevelDebug)
	case config.Info:
		logLevel.Set(slog.LevelInfo)
	case config.Warn:
		logLevel.Set(slog.LevelWarn)
	case config.Error:
		logLevel.Set(slog.LevelError)
	default:
		return nil
	}

	// Replace the handler with a new one using the updated level
	newHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	slog.SetDefault(slog.New(newHandler))
	*logger = *slog.Default()

	return nil
}
