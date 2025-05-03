package endpoints

import (
	loggermw "contacts-list/internal/primary/rest/middlewares/logger"
	requestmw "contacts-list/internal/primary/rest/middlewares/request"
	"contacts-list/pkg/sl"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

type V1Config struct {
	Cors cors.Config
}

func RegisterV1(router fiber.Router, logger sl.Logger, conf V1Config) (fiber.Router, error) {
	router.Use(requestmw.New( //todo real id and other mw from video about grafana
		requestmw.WithHeaders("X-Request-ContactID", "xRequestID"),
		requestmw.WithLoggerKey("request_id"),
		requestmw.WithGenerator(uuid.NewString),
	))
	router.Use(recoverer.New(recoverer.Config{ //todo think about logging stack trace and my own mw
		EnableStackTrace: false,
	}))
	router.Use(loggermw.New(
		loggermw.WithLevel(sl.LevelInfo),
		loggermw.WithLogger(logger),
		loggermw.WithMessage("receive request"),
	))
	router.Use(cors.New(conf.Cors))

	return router.Group("/v1"), nil
}
