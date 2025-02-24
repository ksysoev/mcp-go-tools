package static

import (
	"context"
	"testing"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
)

func TestGetCodeStyle(t *testing.T) {
	config := Config{
		{
			Name:        "test_rule1",
			Category:    "testing",
			Description: "Test rule 1",
			Keywords:    []string{"unit", "basic"},
			Examples: []Example{
				{
					Description: "Example 1",
					Code:        "func TestExample() {}",
				},
			},
		},
		{
			Name:        "test_rule2",
			Category:    "testing",
			Description: "Test rule 2",
			Keywords:    []string{"integration", "advanced"},
			Examples: []Example{
				{
					Description: "Example 2",
					Code:        "func TestExample2() {}",
				},
			},
		},
		{
			Name:        "code_rule",
			Category:    "code",
			Description: "Code style rule",
			Examples: []Example{
				{
					Description: "Style example",
					Code:        "var myVar = 42",
				},
			},
		},
	}
	cfg := &config

	svc := New(cfg)
	ctx := context.Background()

	tests := []struct {
		name       string
		categories []string
		keywords   []string
		want       int
	}{
		{
			name:       "single category rules",
			categories: []string{"testing"},
			keywords:   nil,
			want:       2,
		},
		{
			name:       "multiple categories rules",
			categories: []string{"testing", "code"},
			keywords:   nil,
			want:       3,
		},
		{
			name:       "no matching rules",
			categories: []string{"nonexistent"},
			keywords:   nil,
			want:       0,
		},
		{
			name:       "filter by keyword unit",
			categories: []string{"testing"},
			keywords:   []string{"unit"},
			want:       1,
		},
		{
			name:       "filter by keyword advanced",
			categories: []string{"testing"},
			keywords:   []string{"advanced"},
			want:       1,
		},
		{
			name:       "filter by multiple keywords",
			categories: []string{"testing"},
			keywords:   []string{"unit", "integration"},
			want:       2,
		},
		{
			name:       "no keyword match",
			categories: []string{"testing"},
			keywords:   []string{"nonexistent"},
			want:       0,
		},
		{
			name:       "general rule without keywords",
			categories: []string{"code"},
			keywords:   []string{"any"},
			want:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rules []core.Rule
			rules, err := svc.GetCodeStyle(ctx, tt.categories, tt.keywords)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(rules) != tt.want {
				t.Errorf("Expected %d rules, got %d", tt.want, len(rules))
			}

			for _, rule := range rules {
				found := false

				for _, cat := range tt.categories {
					if rule.Category == cat {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Rule category %s not in expected categories %v", rule.Category, tt.categories)
				}
			}
		})
	}
}
