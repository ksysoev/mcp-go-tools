package static

import (
	"testing"
)

func TestGetRulesByCategory(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{
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
		},
	}

	svc := New(cfg)
	rules := svc.GetRulesByCategory("testing")

	if len(rules) != 2 {
		t.Errorf("Expected 2 testing rules, got %d", len(rules))
	}

	for _, rule := range rules {
		if rule.Category != "testing" {
			t.Errorf("Expected testing category, got %s", rule.Category)
		}
	}
}

func TestValidateCode(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{
			{
				Name:     "constructor_pattern",
				Category: "code_pattern",
				Type:     "template",
				Pattern: RulePattern{
					Validation: "^func New[A-Z][a-zA-Z0-9]*\\(",
				},
				AppliesTo: []string{"struct"},
				Constraints: []Constraint{
					{
						Type:    "max",
						Value:   50,
						Message: "Constructor too long",
					},
				},
			},
		},
	}

	svc := New(cfg)

	tests := []struct {
		name    string
		code    string
		context string
		want    bool
	}{
		{
			name: "valid constructor",
			code: `func NewUser(name string) *User {
				return &User{name: name}
			}`,
			context: "struct",
			want:    true,
		},
		{
			name:    "invalid constructor name",
			code:    "func CreateUser(name string) *User {}",
			context: "struct",
			want:    false,
		},
		{
			name:    "wrong context",
			code:    "func NewUser(name string) *User {}",
			context: "interface",
			want:    true, // No applicable rules for interface context
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.ValidateCode(tt.code, tt.context)
			if result.Valid != tt.want {
				t.Errorf("ValidateCode() = %v, want %v", result.Valid, tt.want)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	expectedTemplate := "func New{{.TypeName}}() *{{.TypeName}} {}"
	cfg := &Config{
		Rules: []Rule{
			{
				Name: "constructor",
				Pattern: RulePattern{
					Template: expectedTemplate,
				},
			},
		},
	}

	svc := New(cfg)

	template, err := svc.GetTemplate("constructor")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if template != expectedTemplate {
		t.Errorf("Expected template %q, got %q", expectedTemplate, template)
	}

	_, err = svc.GetTemplate("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent template, got nil")
	}
}

func TestGetExamples(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{
			{
				Name: "test_rule",
				Examples: []Example{
					{
						Description: "Example 1",
						Code:        "example code",
					},
				},
			},
		},
	}

	svc := New(cfg)

	examples, err := svc.GetExamples("test_rule")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(examples) != 1 {
		t.Errorf("Expected 1 example, got %d", len(examples))
	}

	if examples[0].Description != "Example 1" {
		t.Errorf("Expected description 'Example 1', got %q", examples[0].Description)
	}

	_, err = svc.GetExamples("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent rule examples, got nil")
	}
}

func TestValidateConstraint(t *testing.T) {
	svc := New(&Config{})

	tests := []struct {
		name       string
		code       string
		constraint Constraint
		want       bool
	}{
		{
			name: "max constraint valid",
			code: "short code",
			constraint: Constraint{
				Type:  "max",
				Value: 20,
			},
			want: true,
		},
		{
			name: "max constraint invalid",
			code: "this is a very long piece of code",
			constraint: Constraint{
				Type:  "max",
				Value: 10,
			},
			want: false,
		},
		{
			name: "regex constraint valid",
			code: "func NewUser()",
			constraint: Constraint{
				Type:  "regex",
				Value: "^func New",
			},
			want: true,
		},
		{
			name: "regex constraint invalid",
			code: "func CreateUser()",
			constraint: Constraint{
				Type:  "regex",
				Value: "^func New",
			},
			want: false,
		},
		{
			name: "forbidden constraint valid",
			code: "good code",
			constraint: Constraint{
				Type:  "forbidden",
				Value: []string{"bad", "wrong"},
			},
			want: true,
		},
		{
			name: "forbidden constraint invalid",
			code: "bad code",
			constraint: Constraint{
				Type:  "forbidden",
				Value: []string{"bad", "wrong"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.validateConstraint(tt.code, tt.constraint)
			if got != tt.want {
				t.Errorf("validateConstraint() = %v, want %v", got, tt.want)
			}
		})
	}
}
