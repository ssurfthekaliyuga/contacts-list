package logger

import (
	"contacts-list/pkg/sl"
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

type Config struct {
	AddSource   bool
	Level       sl.Level
	HandlerType HandlerType
}

func NewLogger(conf Config) (sl.Logger, error) {
	opts := slog.HandlerOptions{
		AddSource:   conf.AddSource,
		Level:       conf.Level,
		ReplaceAttr: nil,
	}

	var handler slog.Handler
	switch conf.HandlerType {
	case HandlerTypeText:
		handler = slog.NewTextHandler(os.Stderr, &opts)
	case HandlerTypeJSON:
		handler = slog.NewJSONHandler(os.Stderr, &opts)
	default:
		return sl.Logger{}, fmt.Errorf("invalid handler type: %s", conf.HandlerType)
	}

	logger, err := sl.NewLogger(handler)
	if err != nil {
		return sl.Logger{}, fmt.Errorf("cannot create new logger: %s", err)
	}

	logger.Info(context.Background(), "logger was initialized successfully",
		slog.String("level", conf.Level.String()),
		slog.String("handler_type", string(conf.HandlerType)),
	)

	return logger, nil
}
