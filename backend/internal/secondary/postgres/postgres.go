package postgres

import (
	"contacts-list/pkg/sl"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnString string
}

func New(ctx context.Context, logger sl.Logger, conf Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.New(ctx, conf.ConnString)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping postgres: %w", err)
	}

	logger.Info(ctx, "successfully connected to postgres")

	stop := func() {
		logger.Info(ctx, "closing connections to postgres")
		pool.Close()
		logger.Info(ctx, "successfully close connections to postgres")
	}

	return pool, stop, nil
}
