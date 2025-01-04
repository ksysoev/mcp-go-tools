# MCP Code Tools

A command-line interface (CLI) tool for managing Model Context Protocol (MCP) servers. This tool provides a robust set of commands for initializing, configuring, and managing MCP servers with support for structured logging and flexible configuration options.

> **Note**: This project is under active development. The following features are currently being implemented:
> - MCP tool setup and registration
> - Resource repository integration
> - Core service functionality
> 
> The project follows clean architecture principles and is designed for extensibility.

## About MCP Protocol

The Model Context Protocol (MCP) is implemented using the `github.com/metoro-io/mcp-golang` library. The server:
- Uses stdio transport for communication
- Supports graceful shutdown
- Handles concurrent operations using errgroups
- Provides a flexible tool registration system

## Installation

```bash
go install github.com/ksysoev/mcp-code-tools/cmd/mcp@latest
```

## Features

- Command-line interface built with Cobra
- Flexible configuration using Viper (supports file-based and environment variable configuration)
- Structured logging with slog
  - File output support with --logfile flag
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
mcp start --config config.yaml --loglevel debug
```

#### Run with JSON Logging
Run the server with structured JSON logging (default):
```bash
mcp start --config config.yaml
```

#### Run with Text Logging
Run the server with human-readable text logging:
```bash
mcp start --config config.yaml --logtext
```

#### Run with File Logging
Run the server with logs written to both stdout and file:
```bash
# JSON format (default)
mcp start --config config.yaml --logfile=server.log

# Text format with debug level for request tracking
mcp start --config config.yaml --logfile=server.log --logtext --loglevel=debug
```

## Architecture

The application follows a clean, layered architecture:

1. **API Layer** (`pkg/api`)
   - Handles MCP protocol communication via stdio transport
   - Manages server lifecycle with graceful shutdown
   - Implements tool registration and setup
   - Uses errgroups for concurrent operations

2. **Core Layer** (`pkg/core`)
   - Implements tool handling logic through dependency injection
   - Manages resource repositories
   - Processes MCP requests
   - Designed for extensibility with interface-based components

3. **Repository Layer** (`pkg/repo`)
   - Manages static resources through configuration
   - Supports named resource definitions
   - Implements simple data persistence
   - Uses Viper for resource configuration mapping

4. **Command Layer** (`pkg/cmd`)
   - Implements CLI commands
   - Handles configuration and logging setup

### Global Flags

```bash
--config string     Config file path
--loglevel string   Log level (debug, info, warn, error) (default "info")
--logtext           Log in text format, otherwise JSON
```

### Environment Variables

All configuration options can also be set via environment variables. Environment variables are automatically mapped from the configuration structure (dots are replaced with underscores).

Example environment variables:
```bash
# Log level configuration
MCP_LOGLEVEL=debug
```

### Configuration File

The tool supports configuration via a JSON/YAML file. Specify the config file path using the `--config` flag.

Example configuration file (config.yaml):
```yaml
api:
rules:
  resources:
    - name: "example-resource"
      data: "resource-data"
    - name: "another-resource"
      data: "more-data"
```


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
  - File output support
  - Structured JSON/text formats
  - Configurable log levels
  - Comprehensive debug logging
- ✅ MCP protocol integration
  - Stdio transport implementation
  - Request handler debug logging
  - Error tracking and reporting
- ✅ Tool registration system
  - Category-based rule management
  - Template handling
  - Example management
- 🚧 Resource repository
- 🚧 Core service implementation

### Building from Source

```bash
go build -o mcp ./cmd/mcp
```

### Running Tests

```bash
go test ./...
```

## Version Information

The application includes version and build information that can be set at build time. This information is displayed in logs and can be useful for debugging.

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

After adding the configuration, Cline will have access to the following tools:
- `get_rules_by_category` - Get all rules for a given category
- `get_rules_by_type` - Get all rules of a given type
- `get_applicable_rules` - Get rules that apply to a given context
- `get_template` - Get template for a given rule
- `get_examples` - Get examples for a given rule

Example usage in Cline:
```
You: Show me Go interface naming rules
Cline: Let me get those rules for you...
[Uses get_rules_by_category with "golang" and "interface_naming"]
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
