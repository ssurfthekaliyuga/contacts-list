package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os-lab-3-1/internal/config"
	"os-lab-3-1/internal/controllers"
	"os-lab-3-1/internal/repositories"
	"os/signal"
	"time"
)

//todo panics are not logged via fiber/logger

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.MustLoad()

	postgresURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	postgresConnPool, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("cannot connect to postgreSQL: %s", err)
	}
	defer postgresConnPool.Close()

	if err = postgresConnPool.Ping(context.Background()); err != nil {
		log.Fatalf("cannot ping postgreSQL: %s", err)
	}

	contactsRepo := repositories.NewContactsRepository(postgresConnPool)

	errorHandler := func(c *fiber.Ctx, inError error) error {
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

	server.Use(recoverer.New())
	server.Use(logger.New())

	v1Group := server.Group("/v1")

	contactsGroup := v1Group.Group("/contact")
	contactsGroup.Get("/", controllers.NewGetContacts(contactsRepo))
	contactsGroup.Post("/", controllers.NewCreateContact(contactsRepo))
	contactsGroup.Put("/", controllers.NewUpdateContact(contactsRepo))
	contactsGroup.Delete("/", controllers.NewDeleteContact(contactsRepo))

	addr := fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("starting server panic: %v", r)
			}
		}()

		if err = server.Listen(addr); err != nil {
			log.Fatalf("cannot start http server: %s", err)
		}
	}()

	<-ctx.Done()

	if err = server.ShutdownWithTimeout(time.Minute); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}
}
