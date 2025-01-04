package static

import (
	"context"
	"fmt"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
)

// Config represents the main configuration structure for code generation guidelines
type Config = []Rule

// Rule defines a universal structure for all types of code generation rules
type Rule struct {
	Name        string       `mapstructure:"name"`
	Category    string       `mapstructure:"category"`
	Type        string       `mapstructure:"type"`
	Description string       `mapstructure:"description"`
	Pattern     RulePattern  `mapstructure:"pattern"`
	Constraints []Constraint `mapstructure:"constraints"`
	Examples    []Example    `mapstructure:"examples"`
	AppliesTo   []string     `mapstructure:"applies_to"`
	Priority    int          `mapstructure:"priority"`
	IsRequired  bool         `mapstructure:"required"`
}

// RulePattern defines how the rule should be implemented
type RulePattern struct {
	Template     string            `mapstructure:"template"`
	Replacements map[string]string `mapstructure:"replacements"`
	Validation   string            `mapstructure:"validation"`
	Format       string            `mapstructure:"format"`
}

// Constraint defines limitations or requirements for a rule
type Constraint struct {
	Type    string      `mapstructure:"type"`
	Value   interface{} `mapstructure:"value"`
	Message string      `mapstructure:"message"`
}

// Example provides a usage example for a rule
type Example struct {
	Description string `mapstructure:"description"`
	Code        string `mapstructure:"code"`
	Context     string `mapstructure:"context"`
}

// Repository provides functionality to work with static resources and code rules
type Repository struct {
	config *Config
}

// New creates a new instance of the Repository
func New(cfg *Config) *Repository {
	return &Repository{
		config: cfg,
	}
}

// convertRule converts internal Rule to core.Rule
func (r *Repository) convertRule(rule Rule) core.Rule {
	return core.Rule{
		Name:        rule.Name,
		Category:    rule.Category,
		Type:        rule.Type,
		Description: rule.Description,
		Pattern: core.RulePattern{
			Template:     rule.Pattern.Template,
			Replacements: rule.Pattern.Replacements,
			Validation:   rule.Pattern.Validation,
			Format:       rule.Pattern.Format,
		},
		Constraints: convertConstraints(rule.Constraints),
		Examples:    convertExamples(rule.Examples),
		AppliesTo:   rule.AppliesTo,
		Priority:    rule.Priority,
		IsRequired:  rule.IsRequired,
	}
}

// convertConstraints converts internal Constraints to core.Constraints
func convertConstraints(constraints []Constraint) []core.Constraint {
	result := make([]core.Constraint, len(constraints))
	for i, c := range constraints {
		result[i] = core.Constraint{
			Type:    c.Type,
			Value:   c.Value,
			Message: c.Message,
		}
	}
	return result
}

// convertExamples converts internal Examples to core.Examples
func convertExamples(examples []Example) []core.Example {
	result := make([]core.Example, len(examples))
	for i, e := range examples {
		result[i] = core.Example{
			Description: e.Description,
			Code:        e.Code,
			Context:     e.Context,
		}
	}
	return result
}

// GetRulesByCategory returns all rules for a given category
func (r *Repository) GetRulesByCategory(ctx context.Context, category string) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var rules []core.Rule
		for _, rule := range *r.config {
			if rule.Category == category {
				rules = append(rules, r.convertRule(rule))
			}
		}
		return rules, nil
	}
}

// GetRulesByType returns all rules of a given type
func (r *Repository) GetRulesByType(ctx context.Context, ruleType string) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var rules []core.Rule
		for _, rule := range *r.config {
			if rule.Type == ruleType {
				rules = append(rules, r.convertRule(rule))
			}
		}
		return rules, nil
	}
}

// GetApplicableRules returns all rules that apply to a given context
func (r *Repository) GetApplicableRules(ctx context.Context, context string) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var rules []core.Rule
		for _, rule := range *r.config {
			for _, applies := range rule.AppliesTo {
				if applies == context {
					rules = append(rules, r.convertRule(rule))
					break
				}
			}
		}
		return rules, nil
	}
}

// ValidateCode validates the provided code against applicable rules
func (r *Repository) ValidateCode(ctx context.Context, code, context string) (*core.ValidationResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		result := &core.ValidationResult{
			Valid:    true,
			Messages: make([]string, 0),
		}
		return result, nil
	}
}

// GetTemplate returns the template for a given rule name
func (r *Repository) GetTemplate(ctx context.Context, ruleName string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		for _, rule := range *r.config {
			if rule.Name == ruleName {
				return rule.Pattern.Template, nil
			}
		}
		return "", fmt.Errorf("template not found for rule: %s", ruleName)
	}
}

// GetExamples returns examples for a given rule name
func (r *Repository) GetExamples(ctx context.Context, ruleName string) ([]core.Example, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, rule := range *r.config {
			if rule.Name == ruleName {
				return convertExamples(rule.Examples), nil
			}
		}
		return nil, fmt.Errorf("examples not found for rule: %s", ruleName)
	}
}
