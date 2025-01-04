// Package core provides the core business logic for code generation rule management.
//
// It defines the domain types and interfaces for managing code generation rules,
// their templates, and examples. The package implements a service layer that
// coordinates rule retrieval and management operations.
package core

import (
	"context"
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
