package main

import (
	"contacts-list/internal/app"
	loggerpkg "contacts-list/internal/common/logger"
	"contacts-list/internal/config"
	"contacts-list/internal/domain/usecases"
	"contacts-list/internal/primary/rest/controllers"
	"contacts-list/internal/primary/rest/endpoints"
	"contacts-list/internal/primary/rest/fiber"
	"contacts-list/internal/secondary/postgres"
	"contacts-list/pkg/sl"
	"context"
	"os"
	"os/signal"
)

//todo postman collection with tests and do it with go handler tests
//todo the main todo is a validation after i master grpc

//todo slog lint very good

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	logger, err := loggerpkg.NewLogger(conf.Logger())
	if err != nil {
		panic(err)
	}

	postgresPool, postgresStop, err := postgres.New(ctx, logger, conf.Postgres())
	if err != nil {
		logger.Error(ctx, "failed to initialize postgres connections pool", sl.Err(err))
		return
	}
	defer postgresStop()

	srv, srvStop, err := fiber.New(ctx, logger, conf.FiberServer())
	if err != nil {
		logger.Error(ctx, "failed to initialize http server", sl.Err(err))
		return
	}
	defer srvStop()

	restV1, err := endpoints.RegisterV1(srv, logger, conf.EndpointsV1())
	if err != nil {
		logger.Error(ctx, "cannot create endpoints.V1", sl.Err(err))
		return
	}

	contactsRepo := postgres.NewContacts(postgresPool)

	contactsUseCases := usecases.NewContacts(contactsRepo, logger)

	contactsControllers := controllers.NewContacts(contactsUseCases)

	endpoints.RegisterContacts(contactsControllers, restV1)

	app.New(logger, srv).Run(ctx)
	<-ctx.Done()
}
