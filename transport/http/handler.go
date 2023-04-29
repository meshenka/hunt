package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Heartbeat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}
