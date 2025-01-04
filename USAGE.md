# MCP Code Tools Usage Guide

This guide provides detailed examples and usage patterns for the MCP Code Tools CLI.

## Table of Contents
- [Configuration](#configuration)
- [Code Pattern Rules](#code-pattern-rules)
- [Command Usage Examples](#command-usage-examples)
- [Advanced Usage](#advanced-usage)

## Configuration

### Basic Configuration File
Create a `config.yaml` file to define your code patterns and rules:

```yaml
# Basic constructor pattern rule
- name: "constructor_pattern"
  category: "code_pattern"
  type: "template"
  description: "Standard constructor pattern for Go types"
  pattern:
    template: |
      func New{{.TypeName}}({{.Params}}) *{{.TypeName}} {
          return &{{.TypeName}}{
              {{.Fields}}
          }
      }
```

### Environment Variables
All configuration options can be set via environment variables:

```bash
# Set log level
export MCP_LOGLEVEL=debug

# Set config path
export MCP_CONFIG=/path/to/config.yaml

# Enable text logging
export MCP_LOGTEXT=true
```

## Code Pattern Rules

### 1. Constructor Pattern
```yaml
# Example usage in config.yaml
- name: "constructor_pattern"
  pattern:
    template: |
      func NewService(repo Repository, logger *slog.Logger) *Service {
          return &Service{
              repo: repo,
              logger: logger,
          }
      }
```

### 2. Error Handling Pattern
```yaml
- name: "error_handling"
  pattern:
    template: |
      if err != nil {
          return fmt.Errorf("{{.Operation}}: %w", err)
      }
```

### 3. Test Pattern
```yaml
- name: "test_pattern"
  pattern:
    template: |
      func TestFunction(t *testing.T) {
          tests := []struct{
              name string
              input string
              want string
              wantErr bool
          }{
              {
                  name: "valid input",
                  input: "test",
                  want: "TEST",
                  wantErr: false,
              },
          }
          for _, tt := range tests {
              t.Run(tt.name, func(t *testing.T) {
                  got, err := Function(tt.input)
                  if (err != nil) != tt.wantErr {
                      t.Errorf("unexpected error: %v", err)
                  }
                  if got != tt.want {
                      t.Errorf("got %v, want %v", got, tt.want)
                  }
              })
          }
      }
```

## Command Usage Examples

### Starting the Server

1. Basic start:
```bash
mcp start --config config.yaml
```

2. With debug logging:
```bash
mcp start --config config.yaml --loglevel debug
```

3. With text logging format:
```bash
mcp start --config config.yaml --logtext
```

## Advanced Usage

### Custom Pattern Rules

You can define custom patterns in your config file:

```yaml
- name: "custom_logger"
  category: "logging"
  type: "template"
  description: "Standardized logging pattern"
  pattern:
    template: |
      slog.{{.Level}}("{{.Message}}",
          {{range .Fields}}
          "{{.Name}}", {{.Value}},
          {{end}}
      )
  applies_to: ["function"]
  priority: 1
  required: true
```

### Integration with Development Workflow

1. Git Pre-commit Hook:
```bash
#!/bin/bash
mcp start --config .mcp/config.yaml --loglevel error
```

2. CI/CD Pipeline:
```yaml
validate_code:
  script:
    - mcp start --config ci/mcp-config.yaml --loglevel error
```

### Logging Configuration

Configure structured logging for better observability:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
slog.SetDefault(logger)
```

## Best Practices

1. **Configuration Organization**
   - Keep patterns organized by category
   - Use descriptive names for patterns
   - Include examples in pattern definitions

2. **Pattern Development**
   - Start with essential patterns
   - Iterate based on team feedback
   - Document pattern rationale

3. **Integration Tips**
   - Use version control for pattern configs
   - Implement gradual pattern adoption
   - Monitor pattern effectiveness

## Troubleshooting

Common issues and solutions:

1. **Configuration Issues**
   - Ensure YAML syntax is correct
   - Verify file paths
   - Check environment variables

2. **Performance Optimization**
   - Use specific pattern categories
   - Configure appropriate log levels
