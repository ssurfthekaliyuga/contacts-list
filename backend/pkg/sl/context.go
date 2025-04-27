package sl

import (
	"context"
)

type attrsKey struct{}

func ContextWithAttrs(ctx context.Context, attrs ...Attr) context.Context {
	attrs = append(extractAttrs(ctx), attrs...)

	if len(attrs) == 0 {
		return ctx
	}

	ctx = injectAttrs(ctx, attrs)
	return ctx
}

func injectAttrs(ctx context.Context, attrs []Attr) context.Context {
	ctx = context.WithValue(ctx, attrsKey{}, attrs)
	return ctx
}

func extractAttrs(ctx context.Context) []Attr {
	value := ctx.Value(attrsKey{})

	if value == nil {
		return nil
	}

	attrs := value.([]Attr)
	return attrs
}
