// Package static provides a static file-based implementation of the code generation rule repository.
//
// It implements the core.ResourceRepo interface by managing rules through configuration
// files. The package handles rule storage, retrieval, and conversion between internal
// and core domain types.
package static

import (
	"context"
	"fmt"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
)

// Config represents the main configuration structure for code generation guidelines.
// It is a slice of Rule that can be loaded from configuration files.
type Config = []Rule

// Rule defines a universal structure for all types of code generation rules.
// It mirrors core.Rule but uses mapstructure tags for configuration file parsing.
// Each rule contains metadata, pattern definition, examples, and applicability criteria.
type Rule struct {
	Name        string      `mapstructure:"name"`
	Category    string      `mapstructure:"category"`
	Type        string      `mapstructure:"type"`
	Description string      `mapstructure:"description"`
	Pattern     RulePattern `mapstructure:"pattern"`
	Examples    []Example   `mapstructure:"examples"`
	AppliesTo   []string    `mapstructure:"applies_to"`
	Priority    int         `mapstructure:"priority"`
	IsRequired  bool        `mapstructure:"required"`
}

// RulePattern defines how the rule should be implemented.
// It contains the template to be used for code generation, variable replacements,
// and the format specification for the generated code.
type RulePattern struct {
	Template     string            `mapstructure:"template"`
	Replacements map[string]string `mapstructure:"replacements"`
	Format       string            `mapstructure:"format"`
}

// Example provides a usage example for a rule.
// It includes a description of what the example demonstrates,
// the actual code snippet, and the context in which it applies.
type Example struct {
	Description string `mapstructure:"description"`
	Code        string `mapstructure:"code"`
	Context     string `mapstructure:"context"`
}

// Repository provides functionality to work with static resources and code rules.
// It implements core.ResourceRepo interface and is safe for concurrent use
// as it operates on immutable configuration data.
type Repository struct {
	config *Config
}

// New creates a new instance of the Repository.
// The provided configuration must be properly initialized and will be used
// as the source of all rule data.
func New(cfg *Config) *Repository {
	return &Repository{
		config: cfg,
	}
}

// convertRule converts internal Rule to core.Rule.
// This is an internal helper method that maps between the configuration
// and domain representations of a rule.
func (r *Repository) convertRule(rule Rule) core.Rule {
	return core.Rule{
		Name:        rule.Name,
		Category:    rule.Category,
		Type:        rule.Type,
		Description: rule.Description,
		Pattern: core.RulePattern{
			Template:     rule.Pattern.Template,
			Replacements: rule.Pattern.Replacements,
			Format:       rule.Pattern.Format,
		},
		Examples:   convertExamples(rule.Examples),
		AppliesTo:  rule.AppliesTo,
		Priority:   rule.Priority,
		IsRequired: rule.IsRequired,
	}
}

// convertExamples converts internal Examples to core.Examples.
// This is an internal helper method that maps between the configuration
// and domain representations of examples.
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

// GetRulesByCategory returns all rules for a given category.
// It filters the configuration rules by category and converts them to core.Rule format.
// Returns error if the context is cancelled.
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

// GetRulesByType returns all rules of a given type.
// It filters the configuration rules by type and converts them to core.Rule format.
// Returns error if the context is cancelled.
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

// GetApplicableRules returns all rules that apply to a given context.
// It filters the configuration rules by their AppliesTo field and converts matches to core.Rule format.
// Returns error if the context is cancelled.
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

// GetTemplate returns the template for a given rule name.
// It searches for a rule by name and returns its pattern template.
// Returns error if the rule is not found or the context is cancelled.
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

// GetExamples returns examples for a given rule name.
// It searches for a rule by name and returns its converted examples.
// Returns error if the rule is not found or the context is cancelled.
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
