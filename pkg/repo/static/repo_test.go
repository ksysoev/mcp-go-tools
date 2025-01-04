package static

import (
	"context"
	"testing"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
)

func TestGetRulesByCategory(t *testing.T) {
	config := Config{
		{
			Name:     "test_rule1",
			Category: "testing",
		},
		{
			Name:     "test_rule2",
			Category: "testing",
		},
		{
			Name:     "style_rule",
			Category: "style",
		},
	}
	cfg := &config

	svc := New(cfg)
	ctx := context.Background()

	var coreRules []core.Rule
	coreRules, err := svc.GetRulesByCategory(ctx, "testing")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(coreRules) != 2 {
		t.Errorf("Expected 2 testing rules, got %d", len(coreRules))
	}

	for _, rule := range coreRules {
		if rule.Category != "testing" {
			t.Errorf("Expected testing category, got %s", rule.Category)
		}
	}
}

func TestGetTemplate(t *testing.T) {
	expectedTemplate := "func New{{.TypeName}}() *{{.TypeName}} {}"
	config := Config{
		{
			Name: "constructor",
			Pattern: RulePattern{
				Template: expectedTemplate,
			},
		},
	}
	cfg := &config

	svc := New(cfg)

	ctx := context.Background()
	var template string
	var err error
	template, err = svc.GetTemplate(ctx, "constructor")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if template != expectedTemplate {
		t.Errorf("Expected template %q, got %q", expectedTemplate, template)
	}

	_, err = svc.GetTemplate(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent template, got nil")
	}
}

func TestGetExamples(t *testing.T) {
	config := Config{
		{
			Name: "test_rule",
			Examples: []Example{
				{
					Description: "Example 1",
					Code:        "example code",
				},
			},
		},
	}
	cfg := &config

	svc := New(cfg)

	ctx := context.Background()
	var coreExamples []core.Example
	var err error
	coreExamples, err = svc.GetExamples(ctx, "test_rule")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(coreExamples) != 1 {
		t.Errorf("Expected 1 example, got %d", len(coreExamples))
	}

	if coreExamples[0].Description != "Example 1" {
		t.Errorf("Expected description 'Example 1', got %q", coreExamples[0].Description)
	}

	_, err = svc.GetExamples(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent template, got nil")
	}
}
