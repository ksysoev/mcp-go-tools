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
// It provides methods to retrieve rules by category, type, and context, as well as
// accessing templates and examples for specific rules.
type ResourceRepo interface {
	// GetRulesByCategory returns all rules for a given category
	GetRulesByCategory(ctx context.Context, category string) ([]Rule, error)

	// GetRulesByType returns all rules of a given type
	GetRulesByType(ctx context.Context, ruleType string) ([]Rule, error)

	// GetApplicableRules returns all rules that apply to a given context
	GetApplicableRules(ctx context.Context, context string) ([]Rule, error)

	// GetTemplate returns the template for a given rule name
	GetTemplate(ctx context.Context, ruleName string) (string, error)

	// GetExamples returns examples for a given rule name
	GetExamples(ctx context.Context, ruleName string) ([]Example, error)
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

// GetRulesByCategory retrieves all rules that belong to the specified category.
// It returns a slice of rules and any error encountered during the retrieval.
// Returns error if the repository access fails.
func (s *Service) GetRulesByCategory(ctx context.Context, category string) ([]Rule, error) {
	return s.resource.GetRulesByCategory(ctx, category)
}

// GetRulesByType retrieves all rules of the specified type.
// It returns a slice of rules and any error encountered during the retrieval.
// Returns error if the repository access fails.
func (s *Service) GetRulesByType(ctx context.Context, ruleType string) ([]Rule, error) {
	return s.resource.GetRulesByType(ctx, ruleType)
}

// GetApplicableRules retrieves all rules that are applicable to the specified context.
// It returns a slice of rules that can be applied in the given context.
// Returns error if the repository access fails.
func (s *Service) GetApplicableRules(ctx context.Context, context string) ([]Rule, error) {
	return s.resource.GetApplicableRules(ctx, context)
}

// GetTemplate retrieves the template associated with the specified rule name.
// It returns the template string that can be used for code generation.
// Returns error if the rule is not found or repository access fails.
func (s *Service) GetTemplate(ctx context.Context, ruleName string) (string, error) {
	return s.resource.GetTemplate(ctx, ruleName)
}

// GetExamples retrieves all examples associated with the specified rule name.
// It returns a slice of examples that demonstrate the rule's usage.
// Returns error if the rule is not found or repository access fails.
func (s *Service) GetExamples(ctx context.Context, ruleName string) ([]Example, error) {
	return s.resource.GetExamples(ctx, ruleName)
}
