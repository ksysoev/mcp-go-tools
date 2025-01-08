package vector

import (
	"context"
	"fmt"
	"sync"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
)

// Repository implements core.ResourceRepo interface using vector storage
type Repository struct {
	config  map[string][]core.Rule // Cache for rules by category
	vectors map[string][]float32   // map[ruleID]vector
	mu      sync.RWMutex
}

// New creates a new Repository instance
func New() (*Repository, error) {
	return &Repository{
		config:  make(map[string][]core.Rule),
		vectors: make(map[string][]float32),
	}, nil
}

// GetCodeStyle implements core.ResourceRepo interface
// Returns all rules that match the specified categories
func (r *Repository) GetCodeStyle(ctx context.Context, categories []string) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		var rules []core.Rule

		// Create a map for faster category lookup
		categoryMap := make(map[string]bool)
		for _, cat := range categories {
			categoryMap[cat] = true
		}

		// Get rules from each requested category
		for category := range categoryMap {
			if categoryRules, ok := r.config[category]; ok {
				rules = append(rules, categoryRules...)
			}
		}

		return rules, nil
	}
}

// generateVector creates a simple vector representation of a rule
// This is a placeholder implementation - in the future, this should use proper embeddings
func (r *Repository) generateVector(rule core.Rule) []float32 {
	// For now, we'll create a simple 384-dimensional vector
	// This should be replaced with proper embeddings in the future
	vector := make([]float32, 384)

	// Simple hash-based vector generation
	// This is just for demonstration - should be replaced with proper embeddings
	hash := 0
	for _, c := range rule.Name + rule.Description {
		hash = hash*31 + int(c)
	}

	// Distribute the hash across the vector
	for i := range vector {
		vector[i] = float32(hash % 100)
		hash /= 100
	}

	return vector
}

// AddRule adds a new rule to the appropriate category collection
func (r *Repository) AddRule(ctx context.Context, rule core.Rule) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		// Store in memory cache
		r.config[rule.Category] = append(r.config[rule.Category], rule)

		// Generate and store vector representation
		ruleID := fmt.Sprintf("%s_%s", rule.Category, rule.Name)
		r.vectors[ruleID] = r.generateVector(rule)

		return nil
	}
}

// InitializeFromConfig initializes repository from existing config
func (r *Repository) InitializeFromConfig(cfg []core.Rule) error {
	for _, rule := range cfg {
		if err := r.AddRule(context.Background(), rule); err != nil {
			return fmt.Errorf("failed to add rule: %w", err)
		}
	}

	return nil
}

// SearchSimilar finds similar rules using vector similarity
// This is a placeholder implementation that will be enhanced in the future
func (r *Repository) SearchSimilar(_ context.Context, _ string, limit int) ([]core.Rule, error) {
	// This is where we would implement vector similarity search
	// For now, we just return all rules as we're focusing on the infrastructure
	var allRules []core.Rule
	for _, rules := range r.config {
		allRules = append(allRules, rules...)
	}

	if len(allRules) > limit {
		allRules = allRules[:limit]
	}

	return allRules, nil
}
