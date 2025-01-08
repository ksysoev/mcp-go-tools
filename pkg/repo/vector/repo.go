package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
	"github.com/philippgille/chromem-go"
)

// Repository implements core.ResourceRepo interface using chromem-go vector database
type Repository struct {
	db          *chromem.DB
	collections map[string]*chromem.Collection
	mu          sync.RWMutex
}

// Document represents a rule document in the vector database
type Document struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Rule     core.Rule `json:"rule"`
}

// New creates a new Repository instance
func New() (*Repository, error) {
	db := chromem.NewDB()

	return &Repository{
		db:          db,
		collections: make(map[string]*chromem.Collection),
	}, nil
}

// GetCodeStyle implements core.ResourceRepo interface
// Returns all rules that match the specified categories
func (r *Repository) GetCodeStyle(ctx context.Context, categories []string) ([]core.Rule, error) {
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
		collection, ok := r.collections[category]
		if !ok {
			continue
		}

		// For now, we'll get all documents from the collection using a broad query
		// In the future, this could be enhanced with similarity search
		results, err := collection.Query(ctx, "", 100, nil, nil) // Get all documents with empty query
		if err != nil {
			return nil, fmt.Errorf("failed to get documents from collection %s: %w", category, err)
		}

		for _, result := range results {
			var document Document
			if err := json.Unmarshal([]byte(result.Content), &document); err != nil {
				return nil, fmt.Errorf("failed to unmarshal document: %w", err)
			}

			rules = append(rules, document.Rule)
		}
	}

	return rules, nil
}

// createCollection creates a new collection for a category if it doesn't exist
func (r *Repository) createCollection(_ context.Context, category string) (*chromem.Collection, error) {
	collection, ok := r.collections[category]
	if !ok {
		var err error
		collection, err = r.db.CreateCollection(fmt.Sprintf("rules_%s", category), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}

		r.collections[category] = collection
	}

	return collection, nil
}

// AddRule adds a new rule to the appropriate category collection
func (r *Repository) AddRule(ctx context.Context, rule core.Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	collection, err := r.createCollection(ctx, rule.Category)
	if err != nil {
		return err
	}

	// Create document
	document := Document{
		ID:       fmt.Sprintf("%s_%s", rule.Category, rule.Name),
		Name:     rule.Name,
		Category: rule.Category,
		Rule:     rule,
	}

	data, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Add document to collection
	err = collection.AddDocuments(ctx, []chromem.Document{
		{
			ID:      document.ID,
			Content: string(data),
		},
	}, runtime.NumCPU())
	if err != nil {
		return fmt.Errorf("failed to add document to collection: %w", err)
	}

	return nil
}

// InitializeFromConfig initializes collections from existing config
func (r *Repository) InitializeFromConfig(cfg []core.Rule) error {
	ctx := context.Background()

	for _, rule := range cfg {
		if err := r.AddRule(ctx, rule); err != nil {
			return fmt.Errorf("failed to add rule: %w", err)
		}
	}

	return nil
}

// SearchSimilar finds similar rules using vector similarity
func (r *Repository) SearchSimilar(ctx context.Context, query string, limit int) ([]core.Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allRules []core.Rule

	// Search in each collection
	for _, collection := range r.collections {
		results, err := collection.Query(ctx, query, limit, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to query collection: %w", err)
		}

		for _, result := range results {
			var document Document
			if err := json.Unmarshal([]byte(result.Content), &document); err != nil {
				return nil, fmt.Errorf("failed to unmarshal document: %w", err)
			}

			allRules = append(allRules, document.Rule)
		}
	}

	// Limit total results
	if len(allRules) > limit {
		allRules = allRules[:limit]
	}

	return allRules, nil
}
