package service

import (
	"context"

	"github.com/kirill/mcp-code-guidelines/pkg/core"
)

// GoProvider implements GuidelineProvider for Go language
type GoProvider struct{}

// NewGoProvider creates a new instance of GoProvider
func NewGoProvider() *GoProvider {
	return &GoProvider{}
}

// supported project types
const (
	ProjectTypeAPI     = "api"
	ProjectTypeCLI     = "cli"
	ProjectTypeLibrary = "library"
)

// SupportsProjectType implements GuidelineProvider
func (p *GoProvider) SupportsProjectType(projectType string) bool {
	switch projectType {
	case ProjectTypeAPI, ProjectTypeCLI, ProjectTypeLibrary:
		return true
	default:
		return false
	}
}

// GetGuidelines implements GuidelineProvider
func (p *GoProvider) GetGuidelines(ctx context.Context, projectType string) ([]core.Guideline, error) {
	guidelines := []core.Guideline{
		{
			Category: "Project Structure",
			Rules: []core.Rule{
				{
					Title:       "Standard Layout",
					Description: "Use standard Go project layout with cmd/, pkg/, and internal/ directories",
					Priority:    1,
				},
				{
					Title:       "Package Organization",
					Description: "Organize packages by feature, keeping them focused and cohesive",
					Priority:    1,
				},
			},
			Examples: []string{
				`project/
├── cmd/                    # Main applications
│   └── app/               # Application-specific code
│       └── main.go        # Application entry point
├── pkg/                   # Public library code
│   ├── api/              # API handlers and routes
│   ├── core/             # Core types and interfaces
│   ├── service/          # Business logic
│   └── repo/             # Data access layer`,
			},
		},
		{
			Category: "Code Style",
			Rules: []core.Rule{
				{
					Title:       "Package Names",
					Description: "Use single, lowercase words for package names. For multi-word packages, use no underscores or mixedCaps",
					Priority:    1,
				},
				{
					Title:       "Interface Names",
					Description: "Use -er suffix for single-method interfaces describing actions",
					Priority:    2,
				},
			},
			Examples: []string{
				`package user // Good
package imageutil // Good for multi-word
package UserService // Bad - don't use mixed caps`,
				`type Reader interface { // Good - single method
    Read(p []byte) (n int, error)
}`,
			},
		},
		{
			Category: "Error Handling",
			Rules: []core.Rule{
				{
					Title:       "Error Wrapping",
					Description: "Wrap errors with context using fmt.Errorf and %w verb",
					Priority:    1,
				},
				{
					Title:       "Custom Errors",
					Description: "Define custom errors for specific error cases",
					Priority:    2,
				},
			},
			Examples: []string{
				`if err != nil {
    return fmt.Errorf("validate user: %w", err)
}`,
				`var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
)`,
			},
		},
	}

	// Add project-specific guidelines
	switch projectType {
	case ProjectTypeAPI:
		guidelines = append(guidelines, p.getAPIGuidelines()...)
	case ProjectTypeCLI:
		guidelines = append(guidelines, p.getCLIGuidelines()...)
	case ProjectTypeLibrary:
		guidelines = append(guidelines, p.getLibraryGuidelines()...)
	}

	return guidelines, nil
}

func (p *GoProvider) getAPIGuidelines() []core.Guideline {
	return []core.Guideline{
		{
			Category: "API Design",
			Rules: []core.Rule{
				{
					Title:       "Handler Structure",
					Description: "Use consistent handler structure with dependency injection",
					Priority:    1,
				},
				{
					Title:       "Error Responses",
					Description: "Use consistent error response format and appropriate HTTP status codes",
					Priority:    1,
				},
			},
			Examples: []string{
				`type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    result, err := h.service.Get(ctx)
    if err != nil {
        slog.Error("get failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(result)
}`,
			},
		},
	}
}

func (p *GoProvider) getCLIGuidelines() []core.Guideline {
	return []core.Guideline{
		{
			Category: "CLI Design",
			Rules: []core.Rule{
				{
					Title:       "Command Structure",
					Description: "Use cobra for CLI applications with clear command hierarchy",
					Priority:    1,
				},
				{
					Title:       "Flag Handling",
					Description: "Use consistent flag naming and provide clear descriptions",
					Priority:    2,
				},
			},
			Examples: []string{
				`var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "A brief description",
    Long:  "A longer description",
}

func Execute() error {
    return rootCmd.Execute()
}`,
			},
		},
	}
}

func (p *GoProvider) getLibraryGuidelines() []core.Guideline {
	return []core.Guideline{
		{
			Category: "Library Design",
			Rules: []core.Rule{
				{
					Title:       "API Design",
					Description: "Design clear, consistent APIs with good documentation",
					Priority:    1,
				},
				{
					Title:       "Versioning",
					Description: "Follow semantic versioning and maintain backwards compatibility",
					Priority:    1,
				},
			},
			Examples: []string{
				`// Package example provides a clear example of good library design
package example

// Client handles all library operations
type Client struct {
    config Config
}

// NewClient creates a new client with the provided configuration
func NewClient(config Config) *Client {
    return &Client{config: config}
}`,
			},
		},
	}
}
