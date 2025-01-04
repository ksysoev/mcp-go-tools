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
	// GetCodeStyle returns all rules that match the specified categories and language
	GetCodeStyle(ctx context.Context, categories []string, language string) ([]Rule, error)
}

// Rule defines a universal structure for all types of code generation rules.
// It encapsulates the complete definition of a code generation rule including
// its metadata, pattern definition, examples, and applicability criteria.
type Rule struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Pattern     RulePattern `json:"pattern"`
	Examples    []Example   `json:"examples"`
	AppliesTo   []string    `json:"applies_to"`
	Priority    int         `json:"priority"`
	IsRequired  bool        `json:"required"`
}

// FormatForLLM returns a concise, token-optimized string representation of the rule
// that is easy for Language Models to parse and understand. It omits empty fields
// and unnecessary metadata to reduce token usage while preserving essential information.
func (r *Rule) FormatForLLM() string {
	var parts []string

	// Always include name and description as they're essential
	parts = append(parts, fmt.Sprintf("Rule: %s", r.Name))
	if r.Description != "" {
		parts = append(parts, fmt.Sprintf("Description: %s", r.Description))
	}

	// Include category and type if present
	if r.Category != "" {
		parts = append(parts, fmt.Sprintf("Category: %s", r.Category))
	}
	if r.Type != "" {
		parts = append(parts, fmt.Sprintf("Type: %s", r.Type))
	}

	// Add pattern information if present
	if r.Pattern.Template != "" {
		parts = append(parts, fmt.Sprintf("Template:\n%s", r.Pattern.Template))
	}
	if len(r.Pattern.Replacements) > 0 {
		replacements := make([]string, 0, len(r.Pattern.Replacements))
		for k, v := range r.Pattern.Replacements {
			replacements = append(replacements, fmt.Sprintf("%s -> %s", k, v))
		}
		parts = append(parts, fmt.Sprintf("Replacements: %s", strings.Join(replacements, ", ")))
	}

	// Include examples if present
	if len(r.Examples) > 0 {
		examples := make([]string, 0, len(r.Examples))
		for _, ex := range r.Examples {
			if ex.Description != "" && ex.Code != "" {
				examples = append(examples, fmt.Sprintf("Example (%s):\n%s", ex.Description, ex.Code))
			}
		}
		if len(examples) > 0 {
			parts = append(parts, strings.Join(examples, "\n"))
		}
	}

	// Add applicability and priority information
	if len(r.AppliesTo) > 0 {
		parts = append(parts, fmt.Sprintf("Applies to: %s", strings.Join(r.AppliesTo, ", ")))
	}
	if r.Priority > 0 {
		parts = append(parts, fmt.Sprintf("Priority: %d", r.Priority))
	}
	if r.IsRequired {
		parts = append(parts, "Required: yes")
	}

	return strings.Join(parts, "\n")
}

// RulePattern defines how the rule should be implemented.
// It contains the template to be used, any variable replacements,
// and the format specification for the generated code.
type RulePattern struct {
	Template     string            `json:"template"`
	Replacements map[string]string `json:"replacements"`
	Format       string            `json:"format"`
}

// Example provides a usage example for a rule.
// It includes a description of what the example demonstrates,
// the actual code snippet, and the context in which it applies.
type Example struct {
	Description string `json:"description"`
	Code        string `json:"code"`
	Context     string `json:"context"`
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

// GetCodeStyle retrieves rules that match the specified categories and language.
// It returns a slice of rules and any error encountered during the retrieval.
// Returns error if the repository access fails.
func (s *Service) GetCodeStyle(ctx context.Context, categories []string, language string) ([]Rule, error) {
	return s.resource.GetCodeStyle(ctx, categories, language)
}

// String implements the Stringer interface for Rule.
// It uses FormatForLLM to provide a string representation optimized for LLMs.
func (r *Rule) String() string {
	return r.FormatForLLM()
}
