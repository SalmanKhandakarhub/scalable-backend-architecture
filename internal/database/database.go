package database

import (
	"fmt"
	"log"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/config"
	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dns := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully.")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migration...")

	err := db.AutoMigrate(
		&user.User{},
		// Add other models here
	)

	if err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}
	log.Println("Database migations completed")
	return nil
}
