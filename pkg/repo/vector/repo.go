package vector

import (
	"context"
	"fmt"
	"sync"

	"github.com/ksysoev/mcp-go-tools/pkg/core"
)

// Repository implements core.ResourceRepo interface using vector storage
type Repository struct {
	collections map[string][]Document
	mu          sync.RWMutex
}

// Document represents a rule document in the vector database
type Document struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Rule     core.Rule `json:"rule"`
	Vector   []float32 `json:"vector"`
}

// New creates a new Repository instance
func New() (*Repository, error) {
	return &Repository{
		collections: make(map[string][]Document),
	}, nil
}

// generateVector creates a simple vector representation of a rule
func generateVector(rule core.Rule) []float32 {
	// Create a simple 384-dimensional vector for testing
	vector := make([]float32, 384)
	content := rule.Name + rule.Description

	// Create a simple hash of the rule content
	var hash int

	for i := 0; i < len(content); i++ {
		hash = hash*31 + int(content[i])
	}

	// Use the hash to generate vector values
	for i := range vector {
		vector[i] = float32(hash%100) / 100.0
		hash /= 100
	}

	return vector
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
			if docs, ok := r.collections[category]; ok {
				for i := range docs {
					rules = append(rules, docs[i].Rule)
				}
			}
		}

		return rules, nil
	}
}

// AddRule adds a new rule to the appropriate category collection
func (r *Repository) AddRule(ctx context.Context, rule core.Rule) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		// Create document
		document := Document{
			ID:       fmt.Sprintf("%s_%s", rule.Category, rule.Name),
			Name:     rule.Name,
			Category: rule.Category,
			Rule:     rule,
			Vector:   generateVector(rule),
		}

		// Add to collection
		r.collections[rule.Category] = append(r.collections[rule.Category], document)

		return nil
	}
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

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32

	var normA float32

	var normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(float64(normA) * float64(normB)))
}

// SearchSimilar finds similar rules using vector similarity
func (r *Repository) SearchSimilar(ctx context.Context, query string, limit int) ([]core.Rule, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		// Create a query vector
		queryVector := generateVector(core.Rule{
			Name:        query,
			Description: query,
		})

		type docWithScore struct {
			doc   Document
			score float32
		}

		var allDocs []docWithScore

		// Calculate similarity scores for all documents
		for _, docs := range r.collections {
			for i := range docs {
				score := cosineSimilarity(queryVector, docs[i].Vector)
				allDocs = append(allDocs, docWithScore{doc: docs[i], score: score})
			}
		}

		// Sort by similarity score (simple bubble sort for now)
		for i := 0; i < len(allDocs)-1; i++ {
			for j := 0; j < len(allDocs)-i-1; j++ {
				if allDocs[j].score < allDocs[j+1].score {
					allDocs[j], allDocs[j+1] = allDocs[j+1], allDocs[j]
				}
			}
		}

		// Get top N results
		var rules []core.Rule
		for i := 0; i < len(allDocs) && i < limit; i++ {
			rules = append(rules, allDocs[i].doc.Rule)
		}

		return rules, nil
	}
}
