package fiber

import (
	"contacts-list/pkg/sl"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	address string
	*fiber.App
}

func New(ctx context.Context, logger sl.Logger, conf Config) (*Server, func(), error) {
	fiberApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		CaseSensitive:         true,
		ErrorHandler:          newErrorHandler(logger),
	})

	stop := func() {
		logger.Info(ctx, "closing fiber fiberApp")

		if err := fiberApp.Shutdown(); err != nil {
			logger.Error(context.Background(), "cannot graceful shutdown fiber fiberApp", sl.Err(err))
			return
		}

		logger.Info(context.Background(), "successfully close fiber fiberApp")
	}

	server := &Server{
		address: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		App:     fiberApp,
	}

	return server, stop, nil
}

func (s *Server) Run() error {
	return s.Listen(s.address)
}
