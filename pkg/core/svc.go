// Package core provides the core business logic for code generation rule management.
//
// It defines the domain types and interfaces for managing code generation rules,
// their templates, and examples. The package implements a service layer that
// coordinates rule retrieval and management operations.
package core

import (
	"context"
	"fmt"
	"strings"
)

// ResourceRepo defines the interface for managing code generation rules and resources.
// It provides methods to retrieve rules by categories and language.
type ResourceRepo interface {
	// GetCodeStyle returns all rules that match the specified categories and keywords
	// If keywords is empty, all rules matching categories are returned
	// If a rule has no keywords defined, it is considered a general rule and is always returned
	GetCodeStyle(ctx context.Context, categories []string, keywords []string) ([]Rule, error)
}

// Rule defines a universal structure for all types of code generation rules.
// It encapsulates the complete definition of a code generation rule including
// its metadata and examples.
type Rule struct {
	Name        string    `json:"name"`
	Category    string    `json:"category"` // One of: "documentation", "testing", "code"
	Description string    `json:"description"`
	Examples    []Example `json:"examples"`
}

// FormatForLLM returns a concise, token-optimized string representation of the rule
// that is easy for Language Models to parse and understand.
func (r *Rule) FormatForLLM() string {
	var parts []string

	// Always include name and description as they're essential
	if r.Description != "" {
		parts = append(parts, fmt.Sprintf("Description: %s", r.Description))
	}

	// Include examples if present
	if len(r.Examples) > 0 {
		examples := make([]string, 0, len(r.Examples))

		for _, ex := range r.Examples {
			if ex.Description != "" && ex.Code != "" {
				examples = append(examples, fmt.Sprintf("Example (%s):\n```\n%s```", ex.Description, ex.Code))
			}
		}

		parts = append(parts, strings.Join(examples, "\n"))
	}

	return strings.Join(parts, "\n")
}

// Example provides a usage example for a rule.
// It includes a description of what the example demonstrates,
// the actual code snippet, and the context in which it applies.
type Example struct {
	Description string `json:"description"`
	Code        string `json:"code"`
}

// Service implements the core business logic for rule management.
// This is safe for concurrent use as it delegates operations to the underlying repository.
type Service struct {
	resource ResourceRepo
}

// New creates a new Service instance with the provided resource repository.
// The repository must be properly initialized before being passed to this constructor.
func New(resource ResourceRepo) *Service {
	return &Service{
		resource: resource,
	}
}

// GetCodeStyle retrieves rules that match the specified categories.
// It returns a slice of rules and any error encountered during the retrieval.
// Returns error if the repository access fails.
func (s *Service) GetCodeStyle(ctx context.Context, categories []string) ([]Rule, error) {
	var keywords []string
	return s.resource.GetCodeStyle(ctx, categories, keywords)
}

// String implements the Stringer interface for Rule.
// It uses FormatForLLM to provide a string representation optimized for LLMs.
func (r *Rule) String() string {
	return r.FormatForLLM()
}
