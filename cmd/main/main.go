package main

import (
	"log"

	"github.com/sikehish/OneCompile/internal/server"
)

func main() {
	addr := "localhost:8080"
	err := server.RunServer(addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("Server listening on %s", addr)
}
