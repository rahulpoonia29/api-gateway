package prettylog

import (
	"log/slog"
	"os"

	"github.com/rahul/api-gateway/pkg/config"
)

func ConfigureLogger(logLevel slog.Level) *slog.Logger {
	handler := NewPrettyHandler(os.Stdout, logLevel)

	return slog.New(handler)
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

	// Create a new pretty handler with the updated level
	newHandler := NewPrettyHandler(os.Stdout, logLevel.Level())

	// Replace the logger with the updated handler
	*logger = *slog.New(newHandler)

	return nil
}
