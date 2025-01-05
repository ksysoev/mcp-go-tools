# MCP Code Tools

A Go-focused Model Context Protocol (MCP) server that provides idiomatic Go code generation, style guidelines, and best practices. This tool helps Language Models understand and generate high-quality Go code following established patterns and conventions.

> **Note**: This project is under active development. The following features are currently being implemented:
> - Go-specific code pattern recognition and generation
> - Idiomatic Go code style enforcement
> - Go project structure templates
> 
> The project follows clean architecture principles and is designed for extensibility.

## About MCP Protocol

The Model Context Protocol (MCP) is implemented using the `github.com/metoro-io/mcp-golang` library. The server:
- Uses stdio transport for communication
- Supports graceful shutdown
- Handles concurrent operations using errgroups
- Provides Go-specific code generation tools

## Installation

```bash
go install github.com/ksysoev/mcp-code-tools/cmd/mcp@latest
```

## Features

- Go-specific code generation and style guidelines
- Command-line interface built with Cobra
- Flexible configuration using YAML/JSON files
- Structured logging with slog
  - File output support with --log-file flag (writes to file instead of stdout)
  - JSON and text formats
  - Configurable log levels
  - Debug logging for request tracking
- Server management commands
- Signal handling for graceful shutdown

## Quick Start

### Basic Command Structure

```bash
mcp [command] [flags]
```

For detailed usage examples and patterns, see [USAGE.md](USAGE.md).

### Common Commands

#### Start Server
Starts the MCP server with the specified configuration:
```bash
mcp start --config config.yaml --log-level debug
```

#### Run with JSON Logging
Run the server with structured JSON logging (default):
```bash
mcp start --config config.yaml
```

#### Run with Text Logging
Run the server with human-readable text logging:
```bash
mcp start --config config.yaml --log-text
```

#### Run with File Logging
Run the server with logs written to a file instead of stdout:
```bash
# JSON format (default)
mcp start --config config.yaml --log-file=server.log

# Text format with debug level for request tracking
mcp start --config config.yaml --log-file=server.log --log-text --log-level=debug
```

Note: When --log-file is provided, logs will be written only to the specified file, not to stdout.

## Architecture

The application follows a clean, layered architecture typical of Go projects:

1. **API Layer** (`pkg/api`)
   - Handles MCP protocol communication via stdio transport
   - Manages server lifecycle with graceful shutdown
   - Implements Go code generation tools
   - Uses errgroups for concurrent operations

2. **Core Layer** (`pkg/core`)
   - Implements Go code pattern recognition
   - Manages code style rules
   - Processes MCP requests
   - Designed for extensibility with interface-based components

3. **Repository Layer** (`pkg/repo`)
   - Manages Go code patterns and templates
   - Supports named resource definitions
   - Implements simple data persistence
   - Uses Viper for resource configuration mapping

4. **Command Layer** (`pkg/cmd`)
   - Implements CLI commands
   - Handles configuration and logging setup

### Global Flags

```bash
--config string      Config file path
--log-level string   Log level (debug, info, warn, error) (default "info")
--log-text          Log in text format, otherwise JSON
--log-file string   Log file path (if set, logs to stdout)
```

### Configuration File

The tool supports configuration via a JSON/YAML file. Specify the config file path using the `--config` flag. See example.config.yaml for Go-specific patterns and rules.

## Project Structure

```
.
├── cmd/
│   └── mcp/              # Main application entry point
├── pkg/
│   ├── api/             # API service implementation
│   ├── cmd/             # Command implementations
│   ├── core/            # Core business logic
│   └── repo/            # Data repositories
```

## Dependencies

- Go 1.23.4 or higher
- github.com/metoro-io/mcp-golang - MCP protocol implementation
- github.com/spf13/cobra - CLI framework
- github.com/spf13/viper - Configuration management
- golang.org/x/sync - Synchronization primitives

## Development

### Project Status

The project is in active development with the following components:
- ✅ CLI framework and command structure
- ✅ Configuration management
- ✅ Enhanced logging system
- ✅ MCP protocol integration
- ✅ Go code pattern recognition
- ✅ Idiomatic Go code generation
- 🚧 Go project templates
- 🚧 Advanced code analysis

### Building from Source

```bash
go build -o mcp ./cmd/mcp
```

### Running Tests

```bash
go test ./...
```

## Using with Cline

To use this MCP server with Cline, add it to Cline's MCP settings file located at:
- VSCode: `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`
- Claude Desktop App: `~/Library/Application Support/Claude/claude_desktop_config.json`

Add the following configuration to the `mcpServers` object in the settings file:

```json
{
  "mcpServers": {
    "code-tools": {
      "command": "mcp",
      "args": ["start"],
      "env": {}
    }
  }
}
```

After adding the configuration, Cline will have access to the following tool:

### codestyle

Retrieves Go-specific coding style guidelines and best practices for generating idiomatic Go code. This tool helps Language Models understand and apply consistent coding standards when writing or modifying Go code.

Parameters:
- `category`: Comma-separated list of rule categories to filter by
  * "code" - Go code patterns, structures, and best practices
  * "documentation" - Go documentation conventions and standards
  * "testing" - Go testing patterns and practices

Returns:
- Go-specific formatting guidelines
- Code style rules with examples and templates
- Priority levels and requirement status

Example usage in Cline:
```
You: Show me Go interface patterns
Cline: Let me get those rules for you...
[Uses codestyle tool with categories="code"]

You: How should I document my Go package?
Cline: I'll check the documentation guidelines...
[Uses codestyle tool with categories="documentation"]

You: What's the best way to write table-driven tests?
Cline: I'll show you the testing patterns...
[Uses codestyle tool with categories="testing"]
```

The tool returns rules in an LLM-optimized format that includes:
- Rule name and description
- Go code templates and examples
- Go-specific formatting context
- Priority and requirement status

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
