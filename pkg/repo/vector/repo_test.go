package vector

import (
	"context"
	"testing"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetCodeStyle(t *testing.T) {
	repo, err := New()
	assert.NoError(t, err)

	// Test data
	rules := []core.Rule{
		{
			Name:        "Test Rule 1",
			Category:    "testing",
			Description: "Test description 1",
			Examples: []core.Example{
				{
					Description: "Example 1",
					Code:        "test code 1",
				},
			},
		},
		{
			Name:        "Test Rule 2",
			Category:    "code",
			Description: "Test description 2",
			Examples: []core.Example{
				{
					Description: "Example 2",
					Code:        "test code 2",
				},
			},
		},
	}

	// Initialize repository with test data
	err = repo.InitializeFromConfig(rules)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		categories []string
		want       int
	}{
		{
			name:       "Get testing category",
			categories: []string{"testing"},
			want:       1,
		},
		{
			name:       "Get code category",
			categories: []string{"code"},
			want:       1,
		},
		{
			name:       "Get both categories",
			categories: []string{"testing", "code"},
			want:       2,
		},
		{
			name:       "Get non-existent category",
			categories: []string{"nonexistent"},
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetCodeStyle(context.Background(), tt.categories)
			assert.NoError(t, err)
			assert.Len(t, got, tt.want)
		})
	}
}

func TestRepository_SearchSimilar(t *testing.T) {
	repo, err := New()
	assert.NoError(t, err)

	// Test data
	rules := []core.Rule{
		{
			Name:        "Test Rule 1",
			Category:    "testing",
			Description: "Test description 1",
		},
		{
			Name:        "Test Rule 2",
			Category:    "code",
			Description: "Test description 2",
		},
		{
			Name:        "Test Rule 3",
			Category:    "documentation",
			Description: "Test description 3",
		},
	}

	// Initialize repository with test data
	err = repo.InitializeFromConfig(rules)
	assert.NoError(t, err)

	// Test search with different limits
	tests := []struct {
		name  string
		query string
		limit int
		want  int
	}{
		{
			name:  "Search with limit 1",
			query: "test",
			limit: 1,
			want:  1,
		},
		{
			name:  "Search with limit 2",
			query: "test",
			limit: 2,
			want:  2,
		},
		{
			name:  "Search with limit exceeding total rules",
			query: "test",
			limit: 5,
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.SearchSimilar(context.Background(), tt.query, tt.limit)
			assert.NoError(t, err)
			assert.Len(t, got, tt.want)
		})
	}
}

func TestRepository_AddRule(t *testing.T) {
	repo, err := New()
	assert.NoError(t, err)

	rule := core.Rule{
		Name:        "Test Rule",
		Category:    "testing",
		Description: "Test description",
		Examples: []core.Example{
			{
				Description: "Example",
				Code:        "test code",
			},
		},
	}

	// Add rule
	err = repo.AddRule(context.Background(), rule)
	assert.NoError(t, err)

	// Verify rule was added
	rules, err := repo.GetCodeStyle(context.Background(), []string{"testing"})
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, rule.Name, rules[0].Name)
	assert.Equal(t, rule.Category, rules[0].Category)
	assert.Equal(t, rule.Description, rules[0].Description)
	assert.Equal(t, rule.Examples, rules[0].Examples)
}
