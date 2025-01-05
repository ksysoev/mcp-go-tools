package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/ksysoev/mcp-code-tools/pkg/api"
	"github.com/ksysoev/mcp-code-tools/pkg/repo/static"
	"github.com/stretchr/testify/assert"
)

func TestStartServerWithError(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		runErr    error
		wantError bool
	}{
		{
			name: "successful startup",
			config: &Config{
				API: api.Config{},
				Rules: static.Config{
					{
						Name:        "test_rule",
						Category:    "testing",
						Description: "test rule",
					},
				},
			},
			wantError: false,
		},
		{
			name: "empty rules",
			config: &Config{
				API:   api.Config{},
				Rules: static.Config{},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Start the server
			errCh := make(chan error, 1)
			go func() {
				errCh <- runStart(ctx, tt.config)
			}()

			// Wait for either error or timeout
			err := <-errCh
			if tt.wantError {
				assert.Error(t, err)
				if tt.runErr != nil {
					assert.ErrorIs(t, err, tt.runErr)
				}
			} else {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
			}

		})
	}
}

func TestStartServerContextCancellation(t *testing.T) {
	// Arrange
	config := &Config{
		API: api.Config{},
		Rules: static.Config{
			{
				Name:     "test_rule",
				Category: "testing",
			},
		},
	}

	// Create context that we'll cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start server
	errCh := make(chan error)
	go func() {
		errCh <- runStart(ctx, config)
	}()

	// Wait a bit to ensure server is running
	time.Sleep(50 * time.Millisecond)

	// Cancel context
	cancel()

	// Verify server stops with context cancellation
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "server did not stop after context cancellation")
	}
}
