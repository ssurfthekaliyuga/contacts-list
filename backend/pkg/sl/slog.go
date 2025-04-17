package sl

import (
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
	Level       slog.Level
	HandlerType HandlerType
}

var defaultOptions = &Options{
	AddSource:   true,
	Level:       slog.LevelInfo,
	HandlerType: HandlerTypeJSON,
}

func NewLogger(opts *Options) (*slog.Logger, error) {
	if opts == nil {
		opts = defaultOptions
	}

	o := slog.HandlerOptions{
		AddSource: opts.AddSource,
		Level:     slog.Level(opts.Level),
	}

	var handler slog.Handler

	switch opts.HandlerType {
	case HandlerTypeText:
		handler = slog.NewTextHandler(os.Stdout, &o)
	case HandlerTypeJSON:
		handler = slog.NewJSONHandler(os.Stdout, &o)
	default:
		return nil, fmt.Errorf("unexpected HandlerType: %q", opts.HandlerType)
	}

	contextHandler := &handlerContext{
		Handler: handler,
	}

	logger := slog.New(contextHandler)

	return logger, nil
}
