package main

import (
	"contacts-list/internal/config"
	"contacts-list/internal/wire"
	"context"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	app, cleanup, err := wire.Wire(ctx, *conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	app.Run(ctx)
}
