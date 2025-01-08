// Package cmd implements the command-line interface for the MCP code tools server.
package cmd

import (
	"context"

	"github.com/ksysoev/mcp-go-tools/pkg/api"
	"github.com/ksysoev/mcp-go-tools/pkg/core"
	"github.com/ksysoev/mcp-go-tools/pkg/repo"
)

// runStart initializes and runs the MCP code tools server with the provided configuration.
// It sets up the component chain in the following order:
// 1. Repository (static or vector) for rule storage
// 2. Core service for business logic
// 3. MCP API service for handling tool requests
//
// The function runs until the context is cancelled or an error occurs.
// Returns error if any component initialization fails or the server encounters an error.
func runStart(ctx context.Context, cfg *Config) error {
	repository, err := repo.New(&cfg.Repository)
	if err != nil {
		return err
	}

	toolHandler := core.New(repository)

	mcpAPI := api.New(&cfg.API, toolHandler)

	return mcpAPI.Run(ctx)
}
