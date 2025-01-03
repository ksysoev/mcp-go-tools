package server

import "encoding/json"

// ServerInfo represents basic information about the MCP server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// JSONSchema represents a JSON Schema definition
type JSONSchema struct {
	Type                 string                `json:"type"`
	Description          string                `json:"description,omitempty"`
	Properties           map[string]JSONSchema `json:"properties,omitempty"`
	Required             []string              `json:"required,omitempty"`
	AdditionalProperties *JSONSchema           `json:"additionalProperties,omitempty"`
}

// ToolInfo represents information about a tool provided by the server
type ToolInfo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema JSONSchema `json:"input_schema"`
}

// Content represents a response content with type and text
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ErrorCode represents standard MCP error codes
type ErrorCode string

const (
	ErrorCodeInvalidParams  ErrorCode = "invalid_params"
	ErrorCodeMethodNotFound ErrorCode = "method_not_found"
	ErrorCodeInternalError  ErrorCode = "internal_error"
)

// Error represents an MCP protocol error
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// Request types
type ListToolsRequest struct{}

type ListToolsResponse struct {
	Tools []ToolInfo `json:"tools"`
}

type CallToolRequest struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type CallToolResponse struct {
	Content []Content `json:"content"`
}
