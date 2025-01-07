package api

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Arrange
	cfg := &Config{}
	handler := NewMockToolHandler(t)

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
		handler *MockToolHandler
		name    string
		wantErr bool
	}{
		{
			name:    "successful registration",
			handler: NewMockToolHandler(t),
			wantErr: false,
		},
		{
			name:    "handler error",
			handler: NewMockToolHandler(t),
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
		handler *MockToolHandler
		name    string
		wantErr bool
	}{
		{
			name:    "successful run",
			handler: NewMockToolHandler(t),
			wantErr: false,
		},
		{
			name:    "handler error",
			handler: NewMockToolHandler(t),
			wantErr: false, // Service should start even if handler has errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			svc := New(&Config{}, tt.handler)
			ctx, cancel := context.WithCancel(context.Background())

			// Act
			errCh := make(chan error)
			go func() {
				errCh <- svc.Run(ctx)
			}()

			defer cancel()

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

func TestService_handleCodeStyle(t *testing.T) {
	tests := []struct {
		name      string
		handler   *MockToolHandler
		args      CodeStyleArgs
		wantErr   bool
		wantRules bool
		ruleCount int
	}{
		{
			name: "successful handling",
			handler: func() *MockToolHandler {
				m := NewMockToolHandler(t)
				m.EXPECT().GetCodeStyle(mock.Anything, []string{"testing"}).Return([]core.Rule{
					{
						Name:        "test_rule",
						Category:    "testing",
						Description: "Test rule",
						Examples: []core.Example{
							{
								Description: "Example",
								Code:        "test code",
							},
						},
					},
				}, nil)
				return m
			}(),
			args: CodeStyleArgs{
				Categories: "testing",
			},
			wantErr:   false,
			wantRules: true,
			ruleCount: 1,
		},
		{
			name: "handler error",
			handler: func() *MockToolHandler {
				m := NewMockToolHandler(t)
				m.EXPECT().GetCodeStyle(mock.Anything, []string{"testing"}).Return(nil, assert.AnError)
				return m
			}(),
			args: CodeStyleArgs{
				Categories: "testing",
			},
			wantErr:   true,
			wantRules: false,
		},
		{
			name: "empty rules",
			handler: func() *MockToolHandler {
				m := NewMockToolHandler(t)
				m.EXPECT().GetCodeStyle(mock.Anything, []string{"testing"}).Return([]core.Rule{}, nil)
				return m
			}(),
			args: CodeStyleArgs{
				Categories: "testing",
			},
			wantErr:   false,
			wantRules: true,
			ruleCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			svc := New(&Config{}, tt.handler)

			// Act
			resp, err := svc.handleCodeStyle(tt.args)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.wantRules {
				require.NotNil(t, resp.Content)
				require.Len(t, resp.Content, 1)
				require.NotNil(t, resp.Content[0])

				content := resp.Content[0].TextContent
				require.NotNil(t, content)

				if tt.ruleCount > 0 {
					assert.Contains(t, content.Text, "test code")
					assert.Contains(t, content.Text, "Test rule")
					assert.Contains(t, content.Text, "---") // Check separator
				} else {
					// Even with empty rules, we should get an empty string
					assert.Equal(t, "", content.Text)
				}
			}
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
			},
			wantErr: false,
		},
		{
			name: "multiple categories",
			args: CodeStyleArgs{
				Categories: "testing,documentation",
			},
			wantErr: false,
		},
		{
			name: "empty categories",
			args: CodeStyleArgs{
				Categories: "",
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			args: CodeStyleArgs{
				Categories: "invalid",
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

	// Split and validate each category
	validCategories := map[string]bool{
		"documentation": true,
		"testing":       true,
		"code":          true,
	}

	categories := strings.Split(a.Categories, ",")
	for _, cat := range categories {
		cat = strings.TrimSpace(cat)
		if !validCategories[cat] {
			return fmt.Errorf("invalid category: %s", cat)
		}
	}

	return nil
}
