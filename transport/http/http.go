package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/meshenka/hunt/storage"
	"github.com/rs/zerolog/log"
)

func logRoute(method, path, h string, n int) {
	log.Debug().Str("method", method).Str("path", path).Str("handler", h).Int("n", n).Msg("route setup")
}

// init setup logger properly for all cases.
func init() {
	gin.DebugPrintRouteFunc = logRoute
	gin.DefaultWriter = log.Logger
	gin.DefaultErrorWriter = log.Logger
}

func Listen(ctx context.Context, opts ...Option) error {
	cfg := new(config)

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}

	if cfg.databaseDSN == "" {
		return fmt.Errorf("database DSN not configured")
	}

	// TODO connect to db and build a handler that actually do something.
	// conn, err := storage.Connect(cfg.databaseDSN)
	// if err != nil {
	//     return fmt.Errorf("cannot connect to db: %w", err)
	// }

	router := gin.New()
	api := router.Group("/api/v1")

	// TODO add middleware for authenticate.

	opportunities := api.Group("/opportunites")
	opportunities.GET("/", Heartbeat())    // list/search your opportunities.
	opportunities.GET("/:id", Heartbeat()) // show single opportinity.

	internal := router.Group("/__internal__")
	internal.GET("/ping", Heartbeat())
	internal.GET("/metrics", Heartbeat()) // TODO plug in prometheus.

	admin := router.Group("/admin")
	admin.GET("/login", Heartbeat()) // TODO setup real handlers

	return serve(ctx, cfg.host, router)
}

// Serve routes HTTP requests to handler.
func serve(ctx context.Context, addr string, handler http.Handler) error {
	log.Debug().Str("address", addr).Msg("starting HTTP server")

	srv := new(http.Server)
	srv.Addr = addr
	srv.Handler = handler

	sink := make(chan error, 1)

	go func() {
		defer close(sink)
		sink <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return shutdown(srv)
	case err := <-sink:
		return err
	}
}

func shutdown(srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(),  time.Second * 10)
	defer cancel()
	return srv.Shutdown(ctx)
}

