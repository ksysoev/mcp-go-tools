package main

import (
	"log/slog"
	"os"

	"github.com/kirill/mcp-code-guidelines/pkg/server"
	"github.com/kirill/mcp-code-guidelines/pkg/service"
)

func main() {
	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create guideline service with Go provider
	guidelineService := service.NewGuidelineService()
	guidelineService.RegisterProvider("go", service.NewGoProvider())

	// Create and run MCP server
	mcpServer := server.NewServer(guidelineService)
	if err := mcpServer.Run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
