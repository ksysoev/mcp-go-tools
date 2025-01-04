package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ksysoev/mcp-code-tools/pkg/core"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"golang.org/x/sync/errgroup"
)

type ToolHandler interface {
	GetRulesByCategory(ctx context.Context, category string) ([]core.Rule, error)
	GetRulesByType(ctx context.Context, ruleType string) ([]core.Rule, error)
	GetApplicableRules(ctx context.Context, context string) ([]core.Rule, error)
	GetTemplate(ctx context.Context, ruleName string) (string, error)
	GetExamples(ctx context.Context, ruleName string) ([]core.Example, error)
}

type Config struct {
}

type Service struct {
	config  *Config
	handler ToolHandler
}

func New(cfg *Config, handler ToolHandler) *Service {
	return &Service{
		config:  cfg,
		handler: handler,
	}
}

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

// Tool argument types
type CategoryArgs struct {
	Category string `json:"category"`
}

type TypeArgs struct {
	Type string `json:"type"`
}

type ContextArgs struct {
	Context string `json:"context"`
}

type RuleNameArgs struct {
	RuleName string `json:"rule_name"`
}

// mustMarshal marshals the value to JSON and panics on error
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}

func (s *Service) setupTools(server *mcp.Server) error {
	// Register get rules by category tool
	if err := server.RegisterTool("get_rules_by_category", "Get all rules for a given category", func(args CategoryArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_rules_by_category request", "category", args.Category)

		rules, err := s.handler.GetRulesByCategory(context.Background(), args.Category)
		if err != nil {
			slog.Debug("get_rules_by_category failed", "error", err)
			return nil, fmt.Errorf("get rules by category: %w", err)
		}

		slog.Debug("get_rules_by_category completed", "rules_count", len(rules))
		return mcp.NewToolResponse(mcp.NewTextContent(string(mustMarshal(rules)))), nil
	}); err != nil {
		return fmt.Errorf("register get rules by category tool: %w", err)
	}

	// Register get rules by type tool
	if err := server.RegisterTool("get_rules_by_type", "Get all rules of a given type", func(args TypeArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_rules_by_type request", "type", args.Type)

		rules, err := s.handler.GetRulesByType(context.Background(), args.Type)
		if err != nil {
			slog.Debug("get_rules_by_type failed", "error", err)
			return nil, fmt.Errorf("get rules by type: %w", err)
		}

		slog.Debug("get_rules_by_type completed", "rules_count", len(rules))
		return mcp.NewToolResponse(mcp.NewTextContent(string(mustMarshal(rules)))), nil
	}); err != nil {
		return fmt.Errorf("register get rules by type tool: %w", err)
	}

	// Register get applicable rules tool
	if err := server.RegisterTool("get_applicable_rules", "Get all rules that apply to a given context", func(args ContextArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_applicable_rules request", "context", args.Context)

		rules, err := s.handler.GetApplicableRules(context.Background(), args.Context)
		if err != nil {
			slog.Debug("get_applicable_rules failed", "error", err)
			return nil, fmt.Errorf("get applicable rules: %w", err)
		}

		slog.Debug("get_applicable_rules completed", "rules_count", len(rules))
		return mcp.NewToolResponse(mcp.NewTextContent(string(mustMarshal(rules)))), nil
	}); err != nil {
		return fmt.Errorf("register get applicable rules tool: %w", err)
	}

	// Register get template tool
	if err := server.RegisterTool("get_template", "Get template for a given rule", func(args RuleNameArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_template request", "rule_name", args.RuleName)

		template, err := s.handler.GetTemplate(context.Background(), args.RuleName)
		if err != nil {
			slog.Debug("get_template failed", "error", err)
			return nil, fmt.Errorf("get template: %w", err)
		}

		slog.Debug("get_template completed", "template_length", len(template))
		return mcp.NewToolResponse(mcp.NewTextContent(template)), nil
	}); err != nil {
		return fmt.Errorf("register get template tool: %w", err)
	}

	// Register get examples tool
	if err := server.RegisterTool("get_examples", "Get examples for a given rule", func(args RuleNameArgs) (*mcp.ToolResponse, error) {
		slog.Debug("handling get_examples request", "rule_name", args.RuleName)

		examples, err := s.handler.GetExamples(context.Background(), args.RuleName)
		if err != nil {
			slog.Debug("get_examples failed", "error", err)
			return nil, fmt.Errorf("get examples: %w", err)
		}

		slog.Debug("get_examples completed", "examples_count", len(examples))
		return mcp.NewToolResponse(mcp.NewTextContent(string(mustMarshal(examples)))), nil
	}); err != nil {
		return fmt.Errorf("register get examples tool: %w", err)
	}

	return nil
}
