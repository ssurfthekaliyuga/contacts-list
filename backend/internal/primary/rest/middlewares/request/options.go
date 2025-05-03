package request

import (
	"math/rand/v2"
	"strconv"
)

type Options struct {
	Headers   []string
	LoggerKey string
	Generator func() string
}

type Option func(*Options)

func WithGenerator(generator func() string) Option {
	return func(o *Options) {
		o.Generator = generator
	}
}

func WithHeaders(headers ...string) Option {
	return func(o *Options) {
		o.Headers = headers
	}
}

func WithLoggerKey(key string) Option {
	return func(o *Options) {
		o.LoggerKey = key
	}
}

func config(opts []Option) *Options {
	res := &Options{
		Generator: generator,
		Headers:   []string{"X-Request-ContactID"},
		LoggerKey: "X-Request-ContactID",
	}

	for _, fn := range opts {
		fn(res)
	}

	return res
}

func generator() string {
	return strconv.FormatInt(rand.Int64(), 10)
}
