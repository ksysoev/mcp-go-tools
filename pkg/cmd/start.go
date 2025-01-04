package cmd

import (
	"context"

	"github.com/ksysoev/mcp-code-tools/pkg/api"
)

func runStart(ctx context.Context, cfg *Config) error {

	mcpApi := api.New(&cfg.API)

	return mcpApi.Run(ctx)
}
