package main

import (
	"log"
	"os"

	"connectionInfo/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting connectionInfo server on port %s", port)
	if err := server.Run(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
