package sl

import (
	"context"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(handler slog.Handler) (Logger, error) {
	contextHandler := &HandlerContext{
		Handler: handler,
	}

	logger := Logger{
		logger: slog.New(contextHandler),
	}

	return logger, nil
}

func (l Logger) Log(ctx context.Context, level Level, msg string, attrs ...Attr) {
	l.logger.Log(ctx, level, msg, l.attrs(attrs)...)
}

func (l Logger) Debug(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.Log(ctx, LevelDebug, msg, l.attrs(attrs)...)
}

func (l Logger) Info(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.Log(ctx, LevelInfo, msg, l.attrs(attrs)...)
}

func (l Logger) Warn(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.Log(ctx, LevelWarn, msg, l.attrs(attrs)...)
}

func (l Logger) Error(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.Log(ctx, LevelError, msg, l.attrs(attrs)...)
}

func (l Logger) With(attrs ...Attr) Logger {
	return Logger{
		logger: l.logger.With(l.attrs(attrs)...),
	}
}

func (l Logger) WithGroup(name string) Logger {
	return Logger{
		logger: l.logger.WithGroup(name),
	}
}

func (l Logger) attrs(attrs []Attr) []any {
	slice := make([]any, len(attrs))
	for i, a := range attrs {
		slice[i] = a
	}
	return slice
}
