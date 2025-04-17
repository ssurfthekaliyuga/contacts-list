package logger

import (
	"log/slog"
	"os"
)

type Options struct {
	Level   slog.Level
	Logger  *slog.Logger
	Message string
}

type Option func(*Options)

func WithLevel(level slog.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithMessage(msg string) Option {
	return func(o *Options) {
		o.Message = msg
	}
}

func config(opts []Option) *Options {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	res := &Options{
		Level:   slog.LevelInfo,
		Logger:  logger,
		Message: "complete request",
	}

	for _, fn := range opts {
		fn(res)
	}

	return res
}
