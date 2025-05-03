package app

import (
	"contacts-list/pkg/sl"
	"context"
)

type Runner interface {
	Run() error
}

type App struct {
	runners []Runner
	logger  sl.Logger
}

func New(logger sl.Logger, runners ...Runner) *App {
	return &App{
		runners: runners,
		logger:  logger,
	}
}

func (a *App) Run(ctx context.Context) {
	errs := make(chan error)

	for _, runner := range a.runners {
		go func() {
			if err := runner.Run(); err != nil {
				errs <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
	case err := <-errs:
		a.logger.Error(ctx, "error while running component: %w", sl.Err(err))
	}
}
