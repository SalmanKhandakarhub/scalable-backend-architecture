package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		ctx.Next()

		latency := time.Since(startTime)

		statusCode := ctx.Writer.Status()

		log.Printf(
			"[%s] %s %s | Status: %d | Latency: %v | IP: %s | User-Agent: %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Request.Proto,
			statusCode,
			latency,
			ctx.ClientIP(),
			ctx.Request.UserAgent(),
		)
	}
}
