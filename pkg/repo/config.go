package repo

import (
	"fmt"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
	"github.com/ksysoev/mcp-go-tools/pkg/repo/static"
	"github.com/ksysoev/mcp-go-tools/pkg/repo/vector"
)

// Type represents the type of repository to use
type Type string

const (
	// Static represents the static file-based repository
	Static Type = "static"
	// Vector represents the vector database repository
	Vector Type = "vector"
)

// Config represents the repository configuration
type Config struct {
	// Type specifies which repository implementation to use
	Type Type `mapstructure:"type"`
	// Rules defines the code generation rules and patterns
	Rules []static.Rule `mapstructure:"rules"`
}

// New creates a new repository instance based on the configuration
func New(cfg *Config) (core.ResourceRepo, error) {
	switch cfg.Type {
	case Static, "":
		return static.New(&cfg.Rules), nil
	case Vector:
		repo, err := vector.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create vector repository: %w", err)
		}

		// Convert static rules to core rules for initialization
		var rules []core.Rule
		for _, r := range cfg.Rules {
			rules = append(rules, static.ConvertRule(r))
		}

		if err := repo.InitializeFromConfig(rules); err != nil {
			return nil, fmt.Errorf("failed to initialize vector repository: %w", err)
		}

		return repo, nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.Type)
	}
}
