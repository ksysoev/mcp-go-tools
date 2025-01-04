package static

import (
	"fmt"
	"regexp"
)

// Config represents the main configuration structure for code generation guidelines
type Config struct {
	Resources []Resource `mapstructure:"resources"`
	Rules     []Rule     `mapstructure:"rules"`
}

// Resource represents a static resource with name and data
type Resource struct {
	Name string `mapstructure:"name"`
	Data string `mapstructure:"data"`
}

// Rule defines a universal structure for all types of code generation rules
type Rule struct {
	// Name of the rule
	Name string `mapstructure:"name"`

	// Category of the rule (e.g., "testing", "function", "error_handling", "style", etc.)
	Category string `mapstructure:"category"`

	// Type of the rule (e.g., "pattern", "constraint", "template", "naming")
	Type string `mapstructure:"type"`

	// Description explains the purpose and usage of the rule
	Description string `mapstructure:"description"`

	// Pattern contains the actual rule definition
	Pattern RulePattern `mapstructure:"pattern"`

	// Constraints define any limitations or requirements
	Constraints []Constraint `mapstructure:"constraints"`

	// Examples provide usage examples for the rule
	Examples []Example `mapstructure:"examples"`

	// AppliesTo defines where this rule should be applied (e.g., ["functions", "methods", "interfaces"])
	AppliesTo []string `mapstructure:"applies_to"`

	// Priority defines the importance of the rule (higher number means higher priority)
	Priority int `mapstructure:"priority"`

	// IsRequired indicates if this rule must be followed
	IsRequired bool `mapstructure:"required"`
}

// RulePattern defines how the rule should be implemented
type RulePattern struct {
	// Template for code generation
	Template string `mapstructure:"template"`

	// Replacements define variables that can be used in the template
	Replacements map[string]string `mapstructure:"replacements"`

	// Validation is a regex pattern to validate if code matches this rule
	Validation string `mapstructure:"validation"`

	// Format defines how the rule should be formatted (e.g., "go", "yaml", "markdown")
	Format string `mapstructure:"format"`
}

// Constraint defines limitations or requirements for a rule
type Constraint struct {
	// Type of constraint (e.g., "max", "min", "regex", "forbidden")
	Type string `mapstructure:"type"`

	// Value of the constraint
	Value interface{} `mapstructure:"value"`

	// Message to show when constraint is violated
	Message string `mapstructure:"message"`
}

// Example provides a usage example for a rule
type Example struct {
	// Description of what the example demonstrates
	Description string `mapstructure:"description"`

	// Code showing the example
	Code string `mapstructure:"code"`

	// Context provides additional information about when to use this example
	Context string `mapstructure:"context"`
}

// ValidationResult represents the result of validating code against rules
type ValidationResult struct {
	Valid    bool
	Messages []string
}

// Service provides functionality to work with static resources and code rules
type Service struct {
	config *Config
}

// New creates a new instance of the Service
func New(cfg *Config) *Service {
	return &Service{
		config: cfg,
	}
}

// GetRulesByCategory returns all rules for a given category
func (s *Service) GetRulesByCategory(category string) []Rule {
	var rules []Rule
	for _, rule := range s.config.GetRules() {
		if rule.Category == category {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetRulesByType returns all rules of a given type
func (s *Service) GetRulesByType(ruleType string) []Rule {
	var rules []Rule
	for _, rule := range s.config.GetRules() {
		if rule.Type == ruleType {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetApplicableRules returns all rules that apply to a given context
func (s *Service) GetApplicableRules(context string) []Rule {
	var rules []Rule
	for _, rule := range s.config.GetRules() {
		for _, applies := range rule.AppliesTo {
			if applies == context {
				rules = append(rules, rule)
				break
			}
		}
	}
	return rules
}

// ValidateCode validates the provided code against applicable rules
func (s *Service) ValidateCode(code, context string) ValidationResult {
	result := ValidationResult{
		Valid:    true,
		Messages: make([]string, 0),
	}

	rules := s.GetApplicableRules(context)
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

	return result
}

// validateConstraint checks if the code satisfies a given constraint
func (s *Service) validateConstraint(code string, constraint Constraint) bool {
	switch constraint.Type {
	case "max":
		// Handle maximum value constraints (e.g., line length, complexity)
		if val, ok := constraint.Value.(int); ok {
			return len(code) <= val
		}
	case "min":
		// Handle minimum value constraints
		if val, ok := constraint.Value.(int); ok {
			return len(code) >= val
		}
	case "regex":
		// Handle regex pattern constraints
		if pattern, ok := constraint.Value.(string); ok {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return false
			}
			return re.MatchString(code)
		}
	case "forbidden":
		// Handle forbidden patterns/words
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

// GetTemplate returns the template for a given rule name
func (s *Service) GetTemplate(ruleName string) (string, error) {
	for _, rule := range s.config.GetRules() {
		if rule.Name == ruleName {
			return rule.Pattern.Template, nil
		}
	}
	return "", fmt.Errorf("template not found for rule: %s", ruleName)
}

// GetExamples returns examples for a given rule name
func (s *Service) GetExamples(ruleName string) ([]Example, error) {
	for _, rule := range s.config.GetRules() {
		if rule.Name == ruleName {
			return rule.Examples, nil
		}
	}
	return nil, fmt.Errorf("examples not found for rule: %s", ruleName)
}

// GetRules returns all rules from the config
func (c *Config) GetRules() []Rule {
	if c == nil {
		return nil
	}
	return c.Rules
}
