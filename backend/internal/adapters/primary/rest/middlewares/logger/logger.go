package logger

import (
	"contacts-list/pkg/sl"
	"context"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"reflect"
	"time"
)

func New(logger *slog.Logger, extractors ...AttrExtractor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			start = time.Now()
			ctx   = c.UserContext()
			attrs = make([]slog.Attr, 0)
		)

		for _, extractor := range extractors {
			attr := extractor(ctx)
			if attr != nil {
				attrs = append(attrs, *attr)
			}
		}

		c.SetUserContext(sl.ContextWithAttrs(ctx, attrs...))

		defer func() {
			logger.InfoContext(c.UserContext(), "complete request",
				slog.String("method", c.Method()),
				slog.String("endpoint", c.Path()),
				slog.Duration("duration", time.Since(start)), //todo как в fiber logger mw сделать типа через горутину
				slog.Int("status code", c.Response().StatusCode()),
			)
		}()

		return c.Next()
	}
}

type ValueExtractor func(ctx context.Context) any

type AttrExtractor func(ctx context.Context) *slog.Attr

type Constructor func(extractor ValueExtractor, slogKey string) AttrExtractor

func NewAttrExtractor(extractor ValueExtractor, slogKey string) AttrExtractor {
	return func(ctx context.Context) *slog.Attr {
		value := extractor(ctx)
		if reflect.ValueOf(value).IsZero() {
			return nil
		}

		attr := slog.Any(slogKey, value)

		return &attr
	}
}
