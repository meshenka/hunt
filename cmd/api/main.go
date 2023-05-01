package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/meshenka/hunt/cmd"
	"github.com/meshenka/hunt/transport/http"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.Run(run); err != nil {
		log.Error().Err(err).Msg("exiting with error")
		os.Exit(1)
	}
	os.Exit(0)
}

func run(parent context.Context) error {
	ctx, cancel := signal.NotifyContext(
		parent,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	return http.Listen(
		ctx,
		http.WithDatabaseDSN(cmd.Env("HUNT_DATABASE_DSN", "postgres://root:root@localhost:5432/hunt")),
		http.WithHost(fmt.Sprintf("0.0.0.0:%s", cmd.Env("HUNT_PORT", "4000"))),
	)
}
