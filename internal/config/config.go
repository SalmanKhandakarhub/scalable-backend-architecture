package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	ServerHost string
	ServerPort string
	GinMode    string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret          string
	JWTExpirationHours int

	// App
	AppName string
}

func LoadConfig() *Config {
	JwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

	return &Config{
		ServerHost:         getEnv("SERVER_HOST", "localhost"),
		ServerPort:         getEnv("SERVER_PORT", "8026"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "ows1234"),
		DBName:             getEnv("DB_NAME", "myapp_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-key"),
		JWTExpirationHours: JwtExpHours,
		AppName:            getEnv("APP_NAME", "MyApp API"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
