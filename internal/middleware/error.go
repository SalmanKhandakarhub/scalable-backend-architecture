package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		}
	}
}
