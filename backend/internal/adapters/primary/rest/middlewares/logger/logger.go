package logger

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

func New(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		now := time.Now()

		defer func() {
			logger.InfoContext(c.UserContext(), "complete request",
				slog.String("method", c.Method()),
				slog.String("endpoint", c.Path()),
				slog.Duration("duration", time.Since(now)), //todo как в fiber logger mw сделать типа через горутину
				slog.Int("status_code", c.Response().StatusCode()),
			)
		}()

		return c.Next()
	}
}
