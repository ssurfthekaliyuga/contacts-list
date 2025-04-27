package sl

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type HandlerType string

const (
	HandlerTypeText = "text"
	HandlerTypeJSON = "JSON"
)

type Options struct {
	AddSource   bool
	Level       Level
	HandlerType HandlerType
}

var defaultOptions = &Options{
	AddSource:   false,
	Level:       LevelInfo,
	HandlerType: HandlerTypeJSON,
}

type Logger struct {
	logger *slog.Logger
}

func NewDefaultLogger() Logger {
	logger, _ := NewLogger(nil)
	return logger
}

func NewLogger(opts *Options) (Logger, error) {
	if opts == nil {
		opts = defaultOptions
	}

	o := slog.HandlerOptions{
		AddSource: opts.AddSource,
		Level:     opts.Level,
	}

	var handler slog.Handler

	switch opts.HandlerType {
	case HandlerTypeText:
		handler = slog.NewTextHandler(os.Stdout, &o)
	case HandlerTypeJSON:
		handler = slog.NewJSONHandler(os.Stdout, &o)
	default:
		return Logger{}, fmt.Errorf("unexpected HandlerType: %q", opts.HandlerType)
	}

	contextHandler := &handlerContext{
		Handler: handler,
	}

	logger := Logger{
		logger: slog.New(contextHandler),
	}

	return logger, nil
}

func (l Logger) Log(ctx context.Context, level Level, msg string, args ...any) { //todo accept only attrs
	l.logger.Log(ctx, level, msg, args...)
}

func (l Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelDebug, msg, args...)
}

func (l Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelInfo, msg, args...)
}

func (l Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelWarn, msg, args...)
}

func (l Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelError, msg, args...)
}

func (l Logger) With(args ...any) Logger {
	return Logger{
		logger: l.logger.With(args...),
	}
}

func (l Logger) WithGroup(name string) Logger {
	return Logger{
		logger: l.logger.WithGroup(name),
	}
}
