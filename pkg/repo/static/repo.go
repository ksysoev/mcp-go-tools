// Package static provides a static file-based implementation of the code generation rule repository.
//
// It implements the core.ResourceRepo interface by managing rules through configuration
// files. The package handles rule storage, retrieval, and conversion between internal
// and core domain types.
package static

import (
	"context"
	"strings"

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
	Keywords    []string  `mapstructure:"keywords,omitempty"`
}

// Example provides a usage example for a rule.
// It includes a description of what the example demonstrates,
// the actual code snippet, and optional keywords for categorization.
type Example struct {
	Description string   `mapstructure:"description"`
	Code        string   `mapstructure:"code"`
	Keywords    []string `mapstructure:"keywords,omitempty"`
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
func (r *Repository) convertRule(rule *Rule) core.Rule {
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
// GetCodeStyle returns rules filtered by categories and keywords.
// If keywords is empty, all rules matching categories are returned.
// If a rule has no keywords defined, it is considered a general rule and is always returned.
func (r *Repository) GetCodeStyle(ctx context.Context, categories, keywords []string) ([]core.Rule, error) {
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
			// Skip if category doesn't match
			if len(categories) > 0 && !categoryMap[rule.Category] {
				continue
			}

			// If no keywords specified or rule has no keywords, include the rule
			if len(keywords) == 0 || len(rule.Keywords) == 0 {
				rules = append(rules, r.convertRule(&rule))
				continue
			}

			// Check if any of the requested keywords match rule's keywords
			for _, keyword := range keywords {
				for _, ruleKeyword := range rule.Keywords {
					if strings.EqualFold(keyword, ruleKeyword) {
						rules = append(rules, r.convertRule(&rule))
						goto nextRule
					}
				}
			}
		nextRule:
		}

		return rules, nil
	}
}
