package cmd

import (
	"context"

	"github.com/ksysoev/mcp-code-tools/pkg/api"
	"github.com/ksysoev/mcp-code-tools/pkg/core"
	"github.com/ksysoev/mcp-code-tools/pkg/repo/static"
)

func runStart(ctx context.Context, cfg *Config) error {

	staticRepo := static.New(&cfg.Rules)

	toolHandler := core.New(staticRepo)

	mcpApi := api.New(&cfg.API, toolHandler)

	return mcpApi.Run(ctx)
}
