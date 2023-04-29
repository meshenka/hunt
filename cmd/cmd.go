// Package cmd regroups reusable components to build CLIs.
package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Env returns the value of the given environment variable or uses the provided
// fallback value instead.
func Env(name, fallback string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return fallback
}

// Run runs the given function with a context that is closed as soon as an OS
// signal is caught.
func Run(f func(context.Context) error) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	return f(ctx)
}

func getLogLevel(name string) zerolog.Level {
	switch strings.ToLower(name) {
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	case "panic":
		return zerolog.PanicLevel
	}
	return zerolog.ErrorLevel
}

// init setup logger properly for all cases.
func init() {
	zerolog.LevelFieldName = "severity"
	zerolog.DefaultContextLogger = &log.Logger
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorMarshalFunc = multiErrorSupport // proper formatting of errors.Join-ed errors
	zerolog.SetGlobalLevel(getLogLevel(Env("LOGLEVEL", "debug")))
}

const newline = "\n"

func multiErrorSupport(err error) any {
	if err == nil {
		return nil
	}
	if msg := err.Error(); strings.Contains(msg, newline) {
		return strings.Split(msg, newline)
	}
	return err
}

