package logger

import (
	"contacts-list/pkg/sl"
)

type Options struct { //todo make all Options required refuse functional options
	Level   sl.Level
	Logger  sl.Logger
	Message string
}

type Option func(*Options)

func WithLevel(level sl.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithLogger(logger sl.Logger) Option {
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
	res := &Options{
		Level:   sl.LevelInfo,
		Logger:  sl.Logger{}, //todo
		Message: "handle request",
	}

	for _, fn := range opts {
		fn(res)
	}

	return res
}
