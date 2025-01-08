package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/ksysoev/mcp-go-tools/pkg/api"
	"github.com/ksysoev/mcp-go-tools/pkg/repo"
	"github.com/ksysoev/mcp-go-tools/pkg/repo/static"
	"github.com/stretchr/testify/assert"
)

func TestRunStart(t *testing.T) {
	cfg := &Config{
		API: api.Config{},
		Repository: repo.Config{
			Type: repo.Static,
			Rules: []static.Rule{
				{
					Name:     "test",
					Category: "test",
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	err := runStart(ctx, cfg)
	assert.Error(t, err) // Service should return error when context times out
}

func TestRunStart_InvalidConfig(t *testing.T) {
	cfg := &Config{
		API: api.Config{},
		Repository: repo.Config{
			Type: "invalid",
			Rules: []static.Rule{
				{
					Name:     "test",
					Category: "test",
				},
			},
		},
	}

	ctx := context.Background()
	err := runStart(ctx, cfg)
	assert.Error(t, err)
}
