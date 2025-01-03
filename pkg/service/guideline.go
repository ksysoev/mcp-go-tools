package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kirill/mcp-code-guidelines/pkg/core"
)

// GuidelineProvider defines the interface for specific language guideline providers
type GuidelineProvider interface {
	// GetGuidelines returns language-specific guidelines
	GetGuidelines(ctx context.Context, projectType string) ([]core.Guideline, error)
	// SupportsProjectType checks if the provider supports the given project type
	SupportsProjectType(projectType string) bool
}

// GuidelineService implements core.GuidelineService
type GuidelineService struct {
	providers map[string]GuidelineProvider
}

// NewGuidelineService creates a new instance of GuidelineService
func NewGuidelineService() *GuidelineService {
	return &GuidelineService{
		providers: make(map[string]GuidelineProvider),
	}
}

// RegisterProvider registers a new language-specific guideline provider
func (s *GuidelineService) RegisterProvider(language string, provider GuidelineProvider) {
	s.providers[language] = provider
}

// GetGuidelines implements core.GuidelineService
func (s *GuidelineService) GetGuidelines(ctx context.Context, req core.GuidelineRequest) ([]core.Guideline, error) {
	// Validate request
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	// Get provider for the requested language
	provider, ok := s.providers[req.Language]
	if !ok {
		return nil, core.ErrLanguageNotSupported
	}

	// Check if provider supports the project type
	if !provider.SupportsProjectType(req.ProjectType) {
		return nil, core.ErrProjectTypeNotSupported
	}

	// Get guidelines from the provider
	guidelines, err := provider.GetGuidelines(ctx, req.ProjectType)
	if err != nil {
		slog.Error("failed to get guidelines",
			"language", req.Language,
			"project_type", req.ProjectType,
			"error", err)
		return nil, fmt.Errorf("get guidelines from provider: %w", err)
	}

	return guidelines, nil
}

func (s *GuidelineService) validateRequest(req core.GuidelineRequest) error {
	if req.Language == "" {
		return core.ErrInvalidRequest
	}
	if req.ProjectType == "" {
		return core.ErrInvalidRequest
	}
	return nil
}
