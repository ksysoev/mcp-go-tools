package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/kirill/mcp-code-guidelines/pkg/core"
)

// Server represents the MCP server for code guidelines
type Server struct {
	info    ServerInfo
	service core.GuidelineService
	reader  *bufio.Reader
	writer  *bufio.Writer
}

// NewServer creates a new MCP server instance
func NewServer(guidelineService core.GuidelineService) *Server {
	return &Server{
		info: ServerInfo{
			Name:    "code-guidelines",
			Version: "0.1.0",
		},
		service: guidelineService,
		reader:  bufio.NewReader(os.Stdin),
		writer:  bufio.NewWriter(os.Stdout),
	}
}

// Run starts the MCP server
func (s *Server) Run() error {
	slog.Info("Code guidelines MCP server started")

	for {
		// Read request
		line, err := s.reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read request: %w", err)
		}

		// Parse request
		var request struct {
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.writeError(&Error{
				Code:    ErrorCodeInvalidParams,
				Message: "invalid JSON request",
			})
			continue
		}

		// Handle request
		var response interface{}
		var handleErr *Error

		switch request.Method {
		case "list_tools":
			response, handleErr = s.handleListTools(context.Background(), &ListToolsRequest{})
		case "call_tool":
			var callReq CallToolRequest
			if err := json.Unmarshal(request.Params, &callReq); err != nil {
				handleErr = &Error{
					Code:    ErrorCodeInvalidParams,
					Message: fmt.Sprintf("invalid call_tool params: %v", err),
				}
				break
			}
			response, handleErr = s.handleCallTool(context.Background(), &callReq)
		default:
			handleErr = &Error{
				Code:    ErrorCodeMethodNotFound,
				Message: fmt.Sprintf("unknown method: %s", request.Method),
			}
		}

		if handleErr != nil {
			s.writeError(handleErr)
			continue
		}

		// Write response
		if err := s.writeResponse(response); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
	}
}

func (s *Server) handleListTools(ctx context.Context, req *ListToolsRequest) (*ListToolsResponse, *Error) {
	return &ListToolsResponse{
		Tools: []ToolInfo{
			{
				Name:        "get_guidelines",
				Description: "Get code guidelines for a specific programming language and project type",
				InputSchema: JSONSchema{
					Type: "object",
					Properties: map[string]JSONSchema{
						"language": {
							Type:        "string",
							Description: "Programming language (e.g., 'go', 'python')",
						},
						"project_type": {
							Type:        "string",
							Description: "Type of project (e.g., 'api', 'cli', 'library')",
						},
						"options": {
							Type: "object",
							AdditionalProperties: &JSONSchema{
								Type: "string",
							},
							Description: "Additional options for customizing guidelines",
						},
					},
					Required: []string{"language", "project_type"},
				},
			},
		},
	}, nil
}

func (s *Server) handleCallTool(ctx context.Context, req *CallToolRequest) (*CallToolResponse, *Error) {
	switch req.Name {
	case "get_guidelines":
		return s.handleGetGuidelines(ctx, req.Arguments)
	default:
		return nil, &Error{
			Code:    ErrorCodeMethodNotFound,
			Message: fmt.Sprintf("unknown tool: %s", req.Name),
		}
	}
}

func (s *Server) handleGetGuidelines(ctx context.Context, args json.RawMessage) (*CallToolResponse, *Error) {
	var request core.GuidelineRequest
	if err := json.Unmarshal(args, &request); err != nil {
		return nil, &Error{
			Code:    ErrorCodeInvalidParams,
			Message: fmt.Sprintf("invalid request format: %v", err),
		}
	}

	guidelines, err := s.service.GetGuidelines(ctx, request)
	if err != nil {
		switch {
		case core.IsNotSupported(err):
			return nil, &Error{
				Code:    ErrorCodeInvalidParams,
				Message: err.Error(),
			}
		case core.IsInvalidRequest(err):
			return nil, &Error{
				Code:    ErrorCodeInvalidParams,
				Message: err.Error(),
			}
		default:
			slog.Error("failed to get guidelines", "error", err)
			return nil, &Error{
				Code:    ErrorCodeInternalError,
				Message: "internal server error",
			}
		}
	}

	// Format guidelines as markdown for better readability
	markdown := formatGuidelinesMarkdown(guidelines)

	return &CallToolResponse{
		Content: []Content{
			{
				Type: "markdown",
				Text: markdown,
			},
		},
	}, nil
}

func (s *Server) writeError(err *Error) error {
	response := struct {
		Error *Error `json:"error"`
	}{
		Error: err,
	}
	return s.writeResponse(response)
}

func (s *Server) writeResponse(response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}

	if _, err := s.writer.Write(data); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	if err := s.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("write newline: %w", err)
	}
	return s.writer.Flush()
}

func formatGuidelinesMarkdown(guidelines []core.Guideline) string {
	var result string
	result = "# Code Guidelines\n\n"

	for _, g := range guidelines {
		result += fmt.Sprintf("## %s\n\n", g.Category)

		for _, r := range g.Rules {
			result += fmt.Sprintf("### %s\n", r.Title)
			result += fmt.Sprintf("Priority: %d\n\n", r.Priority)
			result += fmt.Sprintf("%s\n\n", r.Description)
		}

		if len(g.Examples) > 0 {
			result += "### Examples\n\n"
			for _, example := range g.Examples {
				result += "```go\n" + example + "\n```\n\n"
			}
		}

		if len(g.References) > 0 {
			result += "### References\n\n"
			for _, ref := range g.References {
				result += fmt.Sprintf("- %s\n", ref)
			}
			result += "\n"
		}
	}

	return result
}
