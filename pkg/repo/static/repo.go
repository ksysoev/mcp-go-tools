// Package static provides a static file-based implementation of the code generation rule repository.
//
// It implements the core.ResourceRepo interface by managing rules through configuration
// files. The package handles rule storage, retrieval, and conversion between internal
// and core domain types.
package static

import (
	"context"

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
		Examples:  convertExamples(rule.Examples),
		AppliesTo: rule.AppliesTo,
		Priority:  rule.Priority,
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
		}
	}
	return result
}

// GetCodeStyle returns all rules that match the specified categories and language.
// It filters the configuration rules by categories and language, converting matches to core.Rule format.
// Returns error if the context is cancelled.
func (r *Repository) GetCodeStyle(ctx context.Context, categories []string, language string) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var rules []core.Rule

		// Create a map for faster category lookup
		categoryMap := make(map[string]bool)
		for _, cat := range categories {
			categoryMap[cat] = true
		}

		for _, rule := range *r.config {
			// Check if rule matches requested language
			if rule.Pattern.Format != language {
				continue
			}

			// Check if rule category matches any requested category
			if categoryMap[rule.Category] {
				rules = append(rules, r.convertRule(rule))
			}
		}
		return rules, nil
	}
}
