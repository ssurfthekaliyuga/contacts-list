package main

import (
	"contacts-list/internal/adapters/primary/rest"
	"contacts-list/internal/adapters/primary/rest/controllers"
	loggermv "contacts-list/internal/adapters/primary/rest/middlewares/logger"
	"contacts-list/internal/adapters/primary/rest/middlewares/request"
	"contacts-list/internal/adapters/secondary/postgres"
	"contacts-list/internal/config"
	"contacts-list/internal/domain/usecases"
	"contacts-list/pkg/sl"
	"context"
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
//todo postman collection with tests and do it with go handler tests

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	logger, err := sl.NewLogger(&sl.Options{
		AddSource:   conf.Logger.AddSource,
		Level:       sl.Level(conf.Logger.Level),
		HandlerType: sl.HandlerType(conf.Logger.HandlerType),
	})
	if err != nil {
		panic(err)
	}

	logger.Info(ctx, "logger was initialized successfully",
		sl.String("level", sl.Level(conf.Logger.Level).String()),
		sl.String("handler_type", conf.Logger.HandlerType),
	)

	postgresConnPool, err := pgxpool.New(ctx, conf.Postgres)
	if err != nil {
		logger.Error(ctx, "cannot connect to postgres", sl.Err(err))
		return
	}
	defer postgresConnPool.Close()

	if err = postgresConnPool.Ping(ctx); err != nil {
		logger.Error(ctx, "cannot ping postgres", sl.Err(err))
		return
	}

	contactsRepo := postgres.NewContacts(postgresConnPool)

	contactsUseCases := usecases.NewContacts(contactsRepo, logger)

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		CaseSensitive:         true,
		ErrorHandler:          rest.NewErrorHandler(logger),
	})

	generator := func() string {
		return uuid.New().String()
	}

	server.Use(request.New(
		request.WithHeaders("X-Request-ContactID", "xRequestID"),
		request.WithLoggerKey("request_id"),
		request.WithGenerator(generator),
	))
	server.Use(recoverer.New(recoverer.Config{
		EnableStackTrace: false,
	}))
	server.Use(loggermv.New(
		loggermv.WithLevel(slog.LevelInfo),
		loggermv.WithLogger(logger),
		loggermv.WithMessage("receive request"),
	))
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	v1Group := server.Group("/v1")

	contactsGroup := v1Group.Group("/contacts")
	contactsGroup.Get("/", controllers.NewGetContacts(contactsUseCases))
	contactsGroup.Post("/", controllers.NewCreateContact(contactsUseCases))
	contactsGroup.Patch("/:contactID", controllers.NewUpdateContact(contactsUseCases)) //вынести :contactID куда то
	contactsGroup.Delete("/:contactID", controllers.NewDeleteContact(contactsUseCases))

	addr := fmt.Sprintf("%s:%s", conf.HTTPServer.Host, conf.HTTPServer.Port)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "panic occurred while http server running", sl.Panic(r))
			}
		}()

		if err = server.Listen(addr); err != nil {
			logger.Error(ctx, "cannot start http server", sl.Err(err))
		}
	}()

	<-ctx.Done()

	if err = server.ShutdownWithTimeout(time.Minute); err != nil {
		logger.Error(ctx, "cannot shutting down the server", sl.Err(err))
	}
}
