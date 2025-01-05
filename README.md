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
go install github.com/ksysoev/mcp-go-tools/cmd/mcp-go-tools@latest
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
mcp start --config config.yaml
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
- ✅ Go project templates
- 🚧 Add integration with mockery
- 🚧 Add integration linters

### Building from Source

```bash
go build -o mcp-go-tools ./cmd/mcp-go-tools
```

### Running Tests

```bash
go test ./...
```

## Using with Cline

To use this MCP server with Cline, add it to Cline's MCP settings

Add the following configuration to the `mcpServers` object in the settings file:

```json
{
  "mcpServers": {
    "code-tools": {
      "command": "mcp-go-tools",
      "args": ["server", "--config=/Users/user/mcp-go-tools/example.config.yaml"],
      "env": {}
    }
  }
}

```

Custom instructions example:

```
Use project template to initialize new applications for GoLang, it's available in MCP server code-tools `codestyle ` with category `template`

Every time you need to generate code use MCP server code-tools `codestyle ` with required category `code`, `documentation`, `testing`
```


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
