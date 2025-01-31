package main

import (
	"fmt"
	"handbook-scraper/server"
	"handbook-scraper/utils"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := utils.LoadEnv(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	server.StartServer()
}
