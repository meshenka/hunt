package http

import (
	"context"
	"fmt"
	"net/http"

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

	log.Debug().Str("dsn", cfg.databaseDSN).Msg("database string")
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
	api.GET("/__internal__/ping", Heartbeat())
    // TODO add proper routing.
	return http.ListenAndServe(cfg.host, router)
}

