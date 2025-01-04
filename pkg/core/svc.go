package core

import (
	"context"
)

// ResoureRepo defines the interface for managing code generation rules and resources
type ResoureRepo interface {
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

// Rule defines a universal structure for all types of code generation rules
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

// RulePattern defines how the rule should be implemented
type RulePattern struct {
	Template     string            `json:"template"`
	Replacements map[string]string `json:"replacements"`
	Format       string            `json:"format"`
}

// Example provides a usage example for a rule
type Example struct {
	Description string `json:"description"`
	Code        string `json:"code"`
	Context     string `json:"context"`
}

type Service struct {
	resource ResoureRepo
}

func New(resource ResoureRepo) *Service {
	return &Service{
		resource: resource,
	}
}

// GetRulesByCategory returns all rules for a given category
func (s *Service) GetRulesByCategory(ctx context.Context, category string) ([]Rule, error) {
	return s.resource.GetRulesByCategory(ctx, category)
}

// GetRulesByType returns all rules of a given type
func (s *Service) GetRulesByType(ctx context.Context, ruleType string) ([]Rule, error) {
	return s.resource.GetRulesByType(ctx, ruleType)
}

// GetApplicableRules returns all rules that apply to a given context
func (s *Service) GetApplicableRules(ctx context.Context, context string) ([]Rule, error) {
	return s.resource.GetApplicableRules(ctx, context)
}

// GetTemplate returns the template for a given rule name
func (s *Service) GetTemplate(ctx context.Context, ruleName string) (string, error) {
	return s.resource.GetTemplate(ctx, ruleName)
}

// GetExamples returns examples for a given rule name
func (s *Service) GetExamples(ctx context.Context, ruleName string) ([]Example, error) {
	return s.resource.GetExamples(ctx, ruleName)
}
