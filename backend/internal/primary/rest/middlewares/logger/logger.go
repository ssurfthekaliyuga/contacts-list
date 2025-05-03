package logger

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func New(opts ...Option) fiber.Handler {
	conf := config(opts)

	return func(c *fiber.Ctx) error {
		conf.Logger.Log(c.UserContext(), conf.Level, conf.Message,
			slog.String("method", c.Method()),
			slog.String("endpoint", c.Path()),
			slog.String("path", c.Path()),
			slog.String("url", c.OriginalURL()),
		)

		return c.Next()
	}
}
