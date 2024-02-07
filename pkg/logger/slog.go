package logger

import (
	"log/slog"
	"os"
)

type Log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type SlogAdapter struct {
	Logger Log
}

// getLogLevelFromEnv returns the log level from the LOG_LEVEL environment variable.
func getLogLevelFromEnv() slog.Level {
	level := os.Getenv("LOG_LEVEL")

	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLogger returns a new Logger.
func NewLogger() *SlogAdapter {
	logFormat := os.Getenv("LOG_FORMAT")

	var handler slog.Handler
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevelFromEnv(),
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevelFromEnv(),
		})
	}

	logger := slog.New(handler)

	return &SlogAdapter{
		Logger: logger,
	}
}
