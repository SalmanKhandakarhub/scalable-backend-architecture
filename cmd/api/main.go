package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load invironment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}
}
