// Package api implements the MCP (Model Context Protocol) server functionality.
//
// It provides a Service that registers and handles MCP tools for code generation rule management.
// The package uses stdio transport for MCP communication and supports concurrent operations
// through error groups. Each tool is registered with debug logging for request tracking.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"golang.org/x/sync/errgroup"
)

const codeStyleDescription = `Retrieve coding style guidelines and best practices for generating idiomatic code.

This tool helps Language Models understand and apply consistent coding standards when generating or modifying code. It provides language-specific rules, patterns, and examples for writing high-quality, maintainable code.

Use this tool when you need to:
1. Generate new code that follows language idioms
2. Understand naming conventions (e.g., how to name interfaces, variables)
3. Apply proper code organization (e.g., file structure, package layout)
4. Implement language-specific patterns and practices
5. Format code according to language standards

Input Parameters:
- categories: Comma separated list of rule categories to filter by
  * "naming" - conventions for naming variables, functions, types
  * "formatting" - code formatting and style rules
  * "organization" - code structure and layout guidelines
  * "patterns" - common design patterns and implementations
  * "documentation" - rules for comments and documentation
  * "interfaces" - interface design principles (Go-specific)
  * "packages" - package organization rules (Go-specific)
  * "errors" - error handling conventions (Go-specific)
  * "concurrency" - concurrent programming patterns (Go-specific)

- language: Target programming language (lowercase)
  Examples: "go", "typescript"

Returns:
- Array of matching style rules, each containing:
  * Name and description
  * Code templates and examples
  * Priority level
  * Whether the rule is required

Example Usage:
1. Get Go interface naming rules:
   {
     "categories": ["naming", "interfaces"],
     "language": "go"
   }

2. Get Python documentation standards:
   {
     "categories": ["documentation"],
     "language": "python"
   }
`

// ToolHandler defines the interface for handling code generation rule operations.
// Implementations must be safe for concurrent use as methods may be called
// simultaneously by different MCP tool handlers.
type ToolHandler interface {
	GetCodeStyle(ctx context.Context, categories []string, language string) ([]core.Rule, error)
}

// Config holds the service configuration parameters.
// Currently empty but maintained for future configuration options.
type Config struct {
}

// Service implements the MCP server functionality for code generation rules.
// It registers tools for rule management and handles their execution through
// the provided ToolHandler. The service is safe for concurrent use.
type Service struct {
	config  *Config
	handler ToolHandler
}

// New creates a new Service instance with the provided configuration and handler.
// The handler must be properly initialized and safe for concurrent use.
func New(cfg *Config, handler ToolHandler) *Service {
	return &Service{
		config:  cfg,
		handler: handler,
	}
}

// Run starts the MCP server and begins handling tool requests.
// It sets up all available tools and starts the server with stdio transport.
// The server runs until the context is cancelled or an error occurs.
// Returns error if tool setup fails or server encounters an error.
func (s *Service) Run(ctx context.Context) error {
	server := mcp.NewServer(stdio.NewStdioServerTransport())

	if err := s.setupTools(server); err != nil {
		return fmt.Errorf("failed to setup tools: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(server.Serve)

	eg.Go(func() error {
		<-ctx.Done()

		// TODO: Implement graceful shutdown, when it'll be supported by the mcp library.

		return ctx.Err()
	})

	err := eg.Wait()
	if errors.Is(err, context.Canceled) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to run service: %w", err)
	}

	return nil
}

// Tool argument types define the expected input parameters for each tool.
// These types are used for JSON unmarshaling of tool arguments.

// CategoryArgs holds the category parameter for rule filtering.
// Used to specify the category of code generation rules to retrieve.
type CodeStyleArgs struct {
	// Category name for filtering rules
	// Example: "controller", "repository", "service"
	Categories string `json:"category" jsonschema:"required,description=The category name for filtering code generation rules, it can contation comma separate list. Examples: 'testing', 'documentation', 'code', 'project'"`
	Language   string `json:"language" jsonschema:"required,description=The language name for filtering code generation rules. Should be lowercase and consistent. Examples: 'go', 'typescript', 'python', 'java', 'csharp', 'rust'"`
}

// mustMarshal marshals the value to JSON and panics on error.
// This is an internal helper used for response formatting where JSON
// marshaling errors indicate a programming error rather than runtime condition.
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}

// setupTools registers all available tools with the MCP server.
// Each tool is registered with debug logging and proper error handling.
// Returns error if any tool registration fails.
func (s *Service) setupTools(server *mcp.Server) error {
	// Register get rules by category tool
	err := server.RegisterTool("codestyle", codeStyleDescription, func(args CodeStyleArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_code_guidelines request", "category", args.Categories, "language", args.Language)

		// Split categories by comma
		categories := strings.Split(args.Categories, ",")
		for i, cat := range categories {
			categories[i] = strings.TrimSpace(cat)
		}

		rules, err := s.handler.GetCodeStyle(context.Background(), categories, args.Language)
		if err != nil {
			slog.Debug("get_rules_by_category failed", "error", err)
			return nil, fmt.Errorf("get rules by category: %w", err)
		}

		slog.Debug("get_rules_by_category completed", "rules_count", len(rules))

		// Format rules in an LLM-friendly way
		var formattedRules []string
		for _, rule := range rules {
			// Include both the rule format and its LLM-friendly representation
			formattedRules = append(formattedRules,
				rule.FormatForLLM(),
				"---") // Separator between rules
		}

		return mcp.NewToolResponse(mcp.NewTextContent(strings.Join(formattedRules, "\n"))), nil
	})
	if err != nil {
		return fmt.Errorf("register get rules by category tool: %w", err)
	}

	return nil
}
