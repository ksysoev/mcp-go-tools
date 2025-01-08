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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := runStart(ctx, cfg)
	assert.NoError(t, err)
}

func TestRunStart_ContextCanceled(t *testing.T) {
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

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := runStart(ctx, cfg)
	assert.NoError(t, err)
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
