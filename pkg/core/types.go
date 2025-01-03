package core

import "context"

// GuidelineRequest represents a request for code guidelines
type GuidelineRequest struct {
	Language    string            `json:"language"`
	ProjectType string            `json:"project_type"`
	Options     map[string]string `json:"options,omitempty"`
}

// Guideline represents a code guideline with specific rules and examples
type Guideline struct {
	Category   string   `json:"category"`
	Rules      []Rule   `json:"rules"`
	Examples   []string `json:"examples"`
	References []string `json:"references,omitempty"`
}

// Rule represents a specific coding rule or best practice
type Rule struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// GuidelineService defines the interface for retrieving code guidelines
type GuidelineService interface {
	// GetGuidelines returns a set of guidelines based on the provided request
	GetGuidelines(ctx context.Context, req GuidelineRequest) ([]Guideline, error)
}
