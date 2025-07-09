package main

import (
	"log"
	"os"

	"github.com/your-org/go-sqlpp-mcp-proxy/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generate_config.go <output_file>")
	}

	filename := os.Args[1]
	if err := config.GenerateExampleConfig(filename); err != nil {
		log.Fatalf("Failed to generate config file: %v", err)
	}

	log.Printf("Generated example config file: %s", filename)
}
