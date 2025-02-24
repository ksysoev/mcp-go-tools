package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule_FormatForLLM(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		rule     Rule
	}{
		{
			name: "full rule with examples",
			rule: Rule{
				Name:        "TestRule",
				Category:    "testing",
				Description: "Test description",
				Examples: []Example{
					{
						Description: "Example 1",
						Code:        "code1",
					},
				},
			},
			expected: "Description: Test description\nExample (Example 1):\n```\ncode1```",
		},
		{
			name: "rule without examples",
			rule: Rule{
				Name:        "TestRule",
				Category:    "testing",
				Description: "Test description",
			},
			expected: "Description: Test description",
		},
		{
			name: "rule with empty description",
			rule: Rule{
				Name:     "TestRule",
				Category: "testing",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule.FormatForLLM()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRule_String(t *testing.T) {
	rule := Rule{
		Name:        "TestRule",
		Category:    "testing",
		Description: "Test description",
		Examples: []Example{
			{
				Description: "Example 1",
				Code:        "code1",
			},
		},
	}

	expected := rule.FormatForLLM()
	assert.Equal(t, expected, rule.String())
}

func TestNew(t *testing.T) {
	mockRepo := NewMockResourceRepo(t)
	svc := New(mockRepo)

	assert.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.resource)
}

func TestService_GetCodeStyle(t *testing.T) {
	ctx := context.Background()
	categories := []string{"testing", "code"}

	expectedRules := []Rule{
		{
			Name:        "Rule1",
			Category:    "testing",
			Description: "Test rule",
		},
		{
			Name:        "Rule2",
			Category:    "code",
			Description: "Code rule",
		},
	}

	mockRepo := NewMockResourceRepo(t)

	var keywords []string

	mockRepo.EXPECT().
		GetCodeStyle(ctx, categories, keywords).
		Return(expectedRules, nil)

	svc := New(mockRepo)
	rules, err := svc.GetCodeStyle(ctx, categories)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}
