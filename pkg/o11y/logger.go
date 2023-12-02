package o11y

import (
	"log/slog"
	"os"
)

type LoggerInterface interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type LoggerOptions struct {
	EnableJSONHandler bool
	EnableStdError    bool
}

type LogImpl struct {
	impl *slog.Logger
}

func (l *LogImpl) Info(msg string, args ...any) {
	if len(args) == 0 {
		l.impl.Info(msg)
		return
	}

	l.impl.Info(msg, args...)
}

func (l *LogImpl) Warn(msg string, args ...any) {
	if len(args) == 0 {
		l.impl.Warn(msg)
		return
	}

	l.impl.Warn(msg, args...)
}

func (l *LogImpl) Error(msg string, args ...any) {
	if len(args) == 0 {
		l.impl.Error(msg)
		return
	}

	l.impl.Error(msg, args...)
}

func (l *LogImpl) Debug(msg string, args ...any) {
	if len(args) == 0 {
		l.impl.Debug(msg)
		return
	}

	l.impl.Debug(msg, args...)
}

func NewLogger(options LoggerOptions) LoggerInterface {
	var logger *slog.Logger
	var jsonHandler *slog.JSONHandler
	var textHandler *slog.TextHandler

	if options.EnableJSONHandler {
		if options.EnableStdError {
			jsonHandler = slog.NewJSONHandler(os.Stderr, nil)
		} else {
			jsonHandler = slog.NewJSONHandler(os.Stdout, nil)
		}

		logger = slog.New(jsonHandler)
	} else {
		if options.EnableStdError {
			textHandler = slog.NewTextHandler(os.Stderr, nil)
		} else {
			textHandler = slog.NewTextHandler(os.Stdout, nil)
		}

		logger = slog.New(textHandler)
	}

	return &LogImpl{
		impl: logger,
	}
}
