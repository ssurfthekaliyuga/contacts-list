package request

import (
	"contacts-list/pkg/sl"
	"context"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func New(options ...Option) fiber.Handler {
	opts := defaultOptions()
	for _, fn := range options {
		fn(opts)
	}

	return func(c *fiber.Ctx) error {
		var id string
		for _, header := range opts.Headers {
			id = c.Get(header)
			if id != "" {
				break
			}
		}

		if id == "" {
			id = opts.Generator()
		}

		ctx := c.UserContext()
		ctx = context.WithValue(c.UserContext(), key{}, id)
		ctx = sl.ContextWithAttrs(ctx, slog.String(opts.LoggerKey, id))

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func Extract(ctx context.Context) string {
	value := ctx.Value(key{})

	if value == nil {
		return ""
	}

	return value.(string)
}

type key struct{}
