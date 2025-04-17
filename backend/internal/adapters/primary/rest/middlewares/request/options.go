package request

import (
	"math/rand/v2"
	"strconv"
)

type Options struct {
	Headers   []string
	Generator func() string
}

type Option func(*Options)

func WithGenerator(generator func() string) Option {
	return func(o *Options) {
		o.Generator = generator
	}
}

func WithHeader(header string) Option {
	return func(o *Options) {
		o.Headers = []string{header}
	}
}

func WithHeadersList(headers []string) Option {
	return func(o *Options) {
		o.Headers = headers
	}
}

func generator() string {
	return strconv.FormatInt(rand.Int64(), 10)
}
