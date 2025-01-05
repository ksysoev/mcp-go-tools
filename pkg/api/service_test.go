package api

import (
	"context"
	"errors"
	"testing"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockToolHandler is a mock implementation of ToolHandler for testing
type mockToolHandler struct {
	rules []core.Rule
	err   error
}

func (m *mockToolHandler) GetCodeStyle(_ context.Context, _ []string, _ string) ([]core.Rule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rules, nil
}

func TestNew(t *testing.T) {
	// Arrange
	cfg := &Config{}
	handler := &mockToolHandler{}

	// Act
	svc := New(cfg, handler)

	// Assert
	assert.NotNil(t, svc)
	assert.Equal(t, cfg, svc.config)
	assert.Equal(t, handler, svc.handler)
}

func TestService_setupTools(t *testing.T) {
	// This test verifies that the codestyle tool is properly registered
	tests := []struct {
		name    string
		handler *mockToolHandler
		wantErr bool
	}{
		{
			name: "successful registration",
			handler: &mockToolHandler{
				rules: []core.Rule{
					{
						Name:        "test_rule",
						Category:    "testing",
						Description: "Test rule",
						Language:    "go",
						Examples: []core.Example{
							{
								Description: "Example",
								Code:        "test code",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "handler error",
			handler: &mockToolHandler{
				err: assert.AnError,
			},
			wantErr: false, // Registration should succeed even if handler has errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			svc := New(&Config{}, tt.handler)
			server := mcp.NewServer(stdio.NewStdioServerTransport())

			// Act
			err := svc.setupTools(server)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestService_Run(t *testing.T) {
	tests := []struct {
		name    string
		handler *mockToolHandler
		wantErr bool
	}{
		{
			name: "successful run",
			handler: &mockToolHandler{
				rules: []core.Rule{
					{
						Name:        "test_rule",
						Category:    "testing",
						Description: "Test rule",
						Language:    "go",
						Examples: []core.Example{
							{
								Description: "Example",
								Code:        "test code",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "handler error",
			handler: &mockToolHandler{
				err: assert.AnError,
			},
			wantErr: false, // Service should start even if handler has errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			svc := New(&Config{}, tt.handler)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Act
			errCh := make(chan error)
			go func() {
				errCh <- svc.Run(ctx)
			}()

			// Cancel context after a short delay to stop the service
			cancel()

			// Assert
			err := <-errCh
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCodeStyleArgs_Validation(t *testing.T) {
	tests := []struct {
		name    string
		args    CodeStyleArgs
		wantErr bool
	}{
		{
			name: "valid args",
			args: CodeStyleArgs{
				Categories: "testing",
				Language:   "go",
			},
			wantErr: false,
		},
		{
			name: "multiple categories",
			args: CodeStyleArgs{
				Categories: "testing,documentation",
				Language:   "go",
			},
			wantErr: false,
		},
		{
			name: "empty categories",
			args: CodeStyleArgs{
				Categories: "",
				Language:   "go",
			},
			wantErr: true,
		},
		{
			name: "empty language",
			args: CodeStyleArgs{
				Categories: "testing",
				Language:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.args.Validate()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// Helper function to add validation to CodeStyleArgs
func (a *CodeStyleArgs) Validate() error {
	if a.Categories == "" {
		return errors.New("categories is required")
	}
	if a.Language == "" {
		return errors.New("language is required")
	}
	return nil
}
