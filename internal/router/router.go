package router

import (
	"net/http"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/config"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/middleware"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *user.Handler, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	app := gin.New()

	// add global middleware
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(middleware.CORS())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.ErrorHandler())

	// checking endpoint
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Server is running",
			"app":     cfg.AppName,
			//"version": cfg.AppVersion,
		})
	})

	v1 := app.Group("api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		users := v1.Group("/users")
		{
			// Public endpoints
			users.GET("", userHandler.GetAll)         // GET /api/v1/users
			users.GET("/:id", userHandler.GetById)    // GET /api/v1/users/:id
			users.GET("/search", userHandler.Search)  // GET /api/v1/users/search
			users.GET("/stats", userHandler.GetStats) // GET /api/v1/users/stats
			users.POST("/users", userHandler.Create)  // POST /api/v1/users

			// protected endpoints (require authentication)
			protected := users.Group("")

			protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protected.PUT("/:id", userHandler.Update)
				protected.DELETE("/:id", userHandler.Delete)
				protected.PUT("/:id/password", userHandler.ChangePassword)
			}
		}
	}

	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Route not found",
			"error":   "The requested endpoint does not exist",
		})
	})

	return app
}
