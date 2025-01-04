package core

import (
	"context"
	"fmt"
	"regexp"
)

// ResoureRepo defines the interface for managing code generation rules and resources
type ResoureRepo interface {
	// GetRulesByCategory returns all rules for a given category
	GetRulesByCategory(ctx context.Context, category string) ([]Rule, error)

	// GetRulesByType returns all rules of a given type
	GetRulesByType(ctx context.Context, ruleType string) ([]Rule, error)

	// GetApplicableRules returns all rules that apply to a given context
	GetApplicableRules(ctx context.Context, context string) ([]Rule, error)

	// ValidateCode validates the provided code against applicable rules
	ValidateCode(ctx context.Context, code, context string) (*ValidationResult, error)

	// GetTemplate returns the template for a given rule name
	GetTemplate(ctx context.Context, ruleName string) (string, error)

	// GetExamples returns examples for a given rule name
	GetExamples(ctx context.Context, ruleName string) ([]Example, error)
}

// Rule defines a universal structure for all types of code generation rules
type Rule struct {
	Name        string       `json:"name"`
	Category    string       `json:"category"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Pattern     RulePattern  `json:"pattern"`
	Constraints []Constraint `json:"constraints"`
	Examples    []Example    `json:"examples"`
	AppliesTo   []string     `json:"applies_to"`
	Priority    int          `json:"priority"`
	IsRequired  bool         `json:"required"`
}

// RulePattern defines how the rule should be implemented
type RulePattern struct {
	Template     string            `json:"template"`
	Replacements map[string]string `json:"replacements"`
	Validation   string            `json:"validation"`
	Format       string            `json:"format"`
}

// Constraint defines limitations or requirements for a rule
type Constraint struct {
	Type    string      `json:"type"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// Example provides a usage example for a rule
type Example struct {
	Description string `json:"description"`
	Code        string `json:"code"`
	Context     string `json:"context"`
}

// ValidationResult represents the result of validating code against rules
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Messages []string `json:"messages"`
}

type Service struct {
	resource ResoureRepo
}

func New(resource ResoureRepo) *Service {
	return &Service{
		resource: resource,
	}
}

// ValidateCode validates the provided code against applicable rules
func (s *Service) ValidateCode(ctx context.Context, code, context string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Messages: make([]string, 0),
	}

	rules, err := s.GetApplicableRules(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("get applicable rules: %w", err)
	}

	for _, rule := range rules {
		if rule.Pattern.Validation != "" {
			re, err := regexp.Compile(rule.Pattern.Validation)
			if err != nil {
				result.Messages = append(result.Messages,
					fmt.Sprintf("Invalid validation pattern in rule %s: %v", rule.Name, err))
				continue
			}

			if !re.MatchString(code) {
				result.Valid = false
				result.Messages = append(result.Messages,
					fmt.Sprintf("Code does not match pattern for rule %s", rule.Name))
			}
		}

		for _, constraint := range rule.Constraints {
			if !s.validateConstraint(code, constraint) {
				result.Valid = false
				result.Messages = append(result.Messages, constraint.Message)
			}
		}
	}

	return result, nil
}

// validateConstraint checks if the code satisfies a given constraint
func (s *Service) validateConstraint(code string, constraint Constraint) bool {
	switch constraint.Type {
	case "max":
		if val, ok := constraint.Value.(int); ok {
			return len(code) <= val
		}
	case "min":
		if val, ok := constraint.Value.(int); ok {
			return len(code) >= val
		}
	case "regex":
		if pattern, ok := constraint.Value.(string); ok {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return false
			}
			return re.MatchString(code)
		}
	case "forbidden":
		if patterns, ok := constraint.Value.([]string); ok {
			for _, pattern := range patterns {
				if regexp.MustCompile(pattern).MatchString(code) {
					return false
				}
			}
			return true
		}
	}
	return false
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
