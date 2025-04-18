package main

import (
	"contacts-list/internal/adapters/primary/rest/controllers"
	loggermw "contacts-list/internal/adapters/primary/rest/middlewares/logger"
	"contacts-list/internal/adapters/primary/rest/middlewares/request"
	"contacts-list/internal/adapters/secondary/postgres"
	"contacts-list/internal/config"
	"contacts-list/pkg/sl"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

//todo panics are not logged via fiber/loggermw
//todo if there is no .env file godotenv panics!!
//todo you can run like this POSTGRES_USER=admin go run main.go

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	logger, err := sl.NewLogger(&sl.Options{
		AddSource:   conf.Logger.AddSource,
		Level:       slog.Level(conf.Logger.Level),
		HandlerType: sl.HandlerType(conf.Logger.HandlerType),
	})
	if err != nil {
		panic(err)
	}

	logger.Info("logger was initialized successfully",
		slog.String("level", slog.Level(conf.Logger.Level).String()),
		slog.String("handler_type", conf.Logger.HandlerType),
	)

	postgresURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.DB,
	)

	postgresConnPool, err := pgxpool.New(context.Background(), postgresURL) // TODO: посмотреть на работе как передается конфиг
	if err != nil {
		logger.Error("cannot connect to postgreSQL", sl.Error(err))
	}
	defer postgresConnPool.Close()

	if err = postgresConnPool.Ping(context.Background()); err != nil {
		logger.Error("cannot ping postgreSQL", sl.Error(err))
	}

	contactsRepo := postgres.NewContacts(postgresConnPool)

	errorHandler := func(c *fiber.Ctx, inError error) error {
		logger.Error("could not process request", sl.Error(inError))

		defaultError := fiber.ErrInternalServerError
		errors.As(inError, &defaultError)

		body := fiber.Map{
			"error": map[string]any{
				"msg": defaultError.Error(),
			},
		}

		return c.
			Status(defaultError.Code).
			JSON(body)
	}

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		CaseSensitive:         true,
		ErrorHandler:          errorHandler,
	})

	generator := func() string {
		return uuid.New().String()
	}

	server.Use(request.New(
		request.WithHeaders("X-Request-ID", "xRequestID"),
		request.WithLoggerKey("request_id"),
		request.WithGenerator(generator),
	))
	server.Use(recoverer.New(recoverer.Config{
		EnableStackTrace: false,
	}))
	server.Use(loggermw.New(
		loggermw.WithLevel(slog.LevelInfo),
		loggermw.WithLogger(logger),
		loggermw.WithMessage("receive request"),
	))
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	v1Group := server.Group("/v1")

	contactsGroup := v1Group.Group("/contact")
	contactsGroup.Get("/", controllers.NewGetContacts(contactsRepo))
	contactsGroup.Post("/", controllers.NewCreateContact(contactsRepo))
	contactsGroup.Put("/", controllers.NewUpdateContact(contactsRepo)) //todo patch
	contactsGroup.Delete("/", controllers.NewDeleteContact(contactsRepo))

	addr := fmt.Sprintf("%s:%s", conf.HTTPServer.Host, conf.HTTPServer.Port)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic occurred while http server running", sl.Panic(r))
			}
		}()

		if err = server.Listen(addr); err != nil {
			logger.Error("cannot start http server", sl.Error(err))
		}
	}()

	<-ctx.Done()

	if err = server.ShutdownWithTimeout(time.Minute); err != nil {
		logger.Error("cannot shutting down the server", sl.Error(err))
	}
}
