package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockResourceRepo is a mock implementation of ResourceRepo for testing
type mockResourceRepo struct {
	rules []Rule
	err   error
}

func (m *mockResourceRepo) GetCodeStyle(_ context.Context, categories []string) ([]Rule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rules, nil
}

func TestNew(t *testing.T) {
	// Arrange
	repo := &mockResourceRepo{}

	// Act
	svc := New(repo)

	// Assert
	assert.NotNil(t, svc)
	assert.Equal(t, repo, svc.resource)
}

func TestGetCodeStyle(t *testing.T) {
	// Define test cases using table-driven test pattern
	tests := []struct {
		name       string
		repo       *mockResourceRepo
		categories []string
		want       []Rule
		wantErr    bool
	}{
		{
			name: "successful retrieval",
			repo: &mockResourceRepo{
				rules: []Rule{
					{
						Name:        "test_rule",
						Category:    "testing",
						Description: "Test rule description",
						Examples: []Example{
							{
								Description: "Example description",
								Code:        "example code",
							},
						},
					},
				},
			},
			categories: []string{"testing"},
			want: []Rule{
				{
					Name:        "test_rule",
					Category:    "testing",
					Description: "Test rule description",
					Examples: []Example{
						{
							Description: "Example description",
							Code:        "example code",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockResourceRepo{
				err: assert.AnError,
			},
			categories: []string{"testing"},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "empty result",
			repo: &mockResourceRepo{
				rules: []Rule{},
			},
			categories: []string{"nonexistent"},
			want:       []Rule{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			svc := New(tt.repo)
			ctx := context.Background()

			// Act
			got, err := svc.GetCodeStyle(ctx, tt.categories)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRule_String(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
		want string
	}{
		{
			name: "full rule",
			rule: Rule{
				Name:        "test_rule",
				Category:    "testing",
				Description: "Test description",
				Examples: []Example{
					{
						Description: "Example 1",
						Code:        "test code",
					},
				},
			},
			want: "Rule: test_rule\nDescription: Test description\nCategory: testing\nExample (Example 1):\ntest code",
		},
		{
			name: "minimal rule",
			rule: Rule{
				Name: "minimal_rule",
			},
			want: "Rule: minimal_rule",
		},
		{
			name: "rule with no examples",
			rule: Rule{
				Name:        "no_examples",
				Category:    "testing",
				Description: "No examples here",
			},
			want: "Rule: no_examples\nDescription: No examples here\nCategory: testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got := tt.rule.String()

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
