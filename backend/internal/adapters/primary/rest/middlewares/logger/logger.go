package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log/slog"
	"time"
)

func New(opts ...Option) fiber.Handler {
	conf := config(opts)

	return func(c *fiber.Ctx) error { //todo разделить на два сообщения получили/обработали короче конец обработки нужно перенести в errorHandler или вызывать его тут
		now := time.Now()

		err := c.Next()
		logger.New()
		conf.Logger.Log(c.UserContext(), conf.Level, conf.Message,
			slog.String("method", c.Method()),
			slog.String("endpoint", c.Path()),
			slog.Duration("duration", time.Since(now)), //todo как в fiber logger mw сделать типа через горутину
			slog.Int("status_code", c.Response().StatusCode()),
		)

		return err
	}
}
