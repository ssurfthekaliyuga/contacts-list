//go:build wireinject

package wire

import (
	"contacts-list/internal/app"
	"contacts-list/internal/common/logger"
	"contacts-list/internal/config"
	"contacts-list/internal/domain/usecases"
	"contacts-list/internal/primary/rest/controllers"
	"contacts-list/internal/primary/rest/endpoints"
	"contacts-list/internal/primary/rest/fiber"
	"contacts-list/internal/secondary/postgres"
	"contacts-list/pkg/sl"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/wire"
)

func Wire(_ context.Context, _ config.Config) (*app.App, func(), error) {
	wire.Build(
		provideLoggerConfig,
		providePostgresConfig,
		provideFiberConfig,
		provideEndpointsConfig,

		logger.NewLogger,
		postgres.New,

		provideContacts,

		provideFiber,
		provideRunners,
		app.New,
	)

	return nil, nil, errors.New("wire injection should not be called directly")
}

func provideLoggerConfig(cfg config.Config) logger.Config {
	return cfg.Logger()
}

func providePostgresConfig(cfg config.Config) postgres.Config {
	return cfg.Postgres()
}

func provideFiberConfig(cfg config.Config) fiber.Config {
	return cfg.FiberServer()
}

func provideEndpointsConfig(cfg config.Config) endpoints.V1Config {
	return cfg.EndpointsV1()
}

func provideRunners(fiberSrv *fiber.Server) []app.Runner {
	return []app.Runner{fiberSrv}
}

func provideFiber(ctx context.Context, logger sl.Logger, srvConf fiber.Config, v1conf endpoints.V1Config, controllers *controllers.Contacts) (*fiber.Server, func(), error) {
	srv, stop, err := fiber.New(ctx, logger, srvConf)
	if err != nil {
		return nil, nil, err
	}

	router, err := endpoints.RegisterV1(srv, logger, v1conf)
	if err != nil {
		return nil, nil, err
	}

	endpoints.RegisterContacts(controllers, router)

	return srv, stop, nil
}

func provideContacts(l sl.Logger, pool *pgxpool.Pool) *controllers.Contacts {
	repo := postgres.NewContacts(pool)
	use := usecases.NewContacts(repo, l)
	return controllers.NewContacts(use)
}
