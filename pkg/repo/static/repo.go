// Package static provides a static file-based implementation of the code generation rule repository.
//
// It implements the core.ResourceRepo interface by managing rules through configuration
// files. The package handles rule storage, retrieval, and conversion between internal
// and core domain types.
package static

import (
	"context"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
)

// Config represents the main configuration structure for code generation guidelines.
// It is a slice of Rule that can be loaded from configuration files.
type Config = []Rule

// Rule defines a universal structure for all types of code generation rules.
// It mirrors core.Rule but uses mapstructure tags for configuration file parsing.
type Rule struct {
	Name        string    `mapstructure:"name"`
	Category    string    `mapstructure:"category"` // One of: "documentation", "testing", "code"
	Description string    `mapstructure:"description"`
	Examples    []Example `mapstructure:"examples"`
}

// Example provides a usage example for a rule.
// It includes a description of what the example demonstrates
// and the actual code snippet.
type Example struct {
	Description string `mapstructure:"description"`
	Code        string `mapstructure:"code"`
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

// ConvertRule converts internal Rule to core.Rule.
// This is a helper function that maps between the configuration
// and domain representations of a rule.
func ConvertRule(rule Rule) core.Rule {
	return core.Rule{
		Name:        rule.Name,
		Category:    rule.Category,
		Description: rule.Description,
		Examples:    convertExamples(rule.Examples),
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

// GetCodeStyle returns all rules that match the specified categories.
// It filters the configuration rules by categories, converting matches to core.Rule format.
// Returns error if the context is cancelled.
func (r *Repository) GetCodeStyle(ctx context.Context, categories []string) ([]core.Rule, error) {
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
			// Check if rule matches requested category
			if categoryMap[rule.Category] {
				rules = append(rules, ConvertRule(rule))
			}
		}

		return rules, nil
	}
}
