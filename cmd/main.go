package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rxtech-lab/i18n-mcp/internal/mcp"
	"github.com/rxtech-lab/i18n-mcp/internal/service"
)

func main() {
	// Default to no log output
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	// Create PoService (will be passed to tools as needed)
	poService := &service.PoService{}

	// Create and initialize MCP server
	mcpServer := mcp.NewMCPServer(poService)

	// Start the MCP server in a goroutine
	go func() {
		if err := mcpServer.Start(); err != nil {
			log.Fatal("Failed to start MCP server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
}
