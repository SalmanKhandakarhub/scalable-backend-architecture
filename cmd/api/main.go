package main

import (
	"log"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/config"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/database"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/router"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/user"
	"github.com/joho/godotenv"
)

func main() {
	// Load invironment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	// load configaretion
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories, services, handelers
	userRepo := user.NewRepository(db)

	userService := user.NewService(userRepo)

	userHandeler := user.NewHandler(userService)

	// Setup router
	r := router.SetupRouter(userHandeler, cfg)

	// Start server
	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("Server starting on http://%s", addr)
	log.Printf("API Documentation: http://%s/health", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
