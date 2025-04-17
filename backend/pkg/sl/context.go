package sl

import (
	"context"
	"log/slog"
)

type attrsKey struct{}

func injectAttrs(ctx context.Context, attrs []slog.Attr) context.Context {
	ctx = context.WithValue(ctx, attrsKey{}, attrs)
	return ctx
}

func extractAttrs(ctx context.Context) []slog.Attr {
	value := ctx.Value(attrsKey{})

	if value == nil {
		return nil
	}

	attrs := value.([]slog.Attr)
	return attrs
}

func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	attrs = append(extractAttrs(ctx), attrs...)

	if len(attrs) == 0 {
		return ctx
	}

	ctx = injectAttrs(ctx, attrs)
	return ctx
}
