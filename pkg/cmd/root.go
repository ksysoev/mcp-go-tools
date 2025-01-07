// Package cmd implements the command-line interface for the MCP code tools server.
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// args holds all command-line arguments and configuration options.
type args struct {
	build      string
	version    string
	LogLevel   string
	ConfigPath string
	LogFile    string
	TextFormat bool
}

// InitCommands initializes and returns the root command for the MCP code tools server.
// It sets up the command structure, persistent flags, and environment variable bindings.
// The build and version parameters are used for logging and version information.
// Returns error if flag binding or configuration unmarshaling fails.
func InitCommands(build, version string) (*cobra.Command, error) {
	args := &args{
		build:   build,
		version: version,
	}

	cmd := &cobra.Command{
		Use:     "mcp-go-tools",
		Short:   "MCP code tools server",
		Long:    "Model Context Protocol server for code generation tools",
		Version: fmt.Sprintf("%s (Build: %s)", version, build),
	}

	// Add server subcommand
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start MCP code tools server",
		Long:  "Start the Model Context Protocol server for code generation tools",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := initLogger(args); err != nil {
				return fmt.Errorf("init logger: %w", err)
			}

			slog.Info("Starting MCP code tools server",
				slog.String("version", args.version),
				slog.String("build", args.build))

			cfg, err := initConfig(args)
			if err != nil {
				return fmt.Errorf("init config: %w", err)
			}

			return runStart(cmd.Context(), cfg)
		},
	}

	// Add persistent flags
	serverCmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "", "config file path")
	serverCmd.PersistentFlags().StringVar(&args.LogLevel, "log-level", "info", "log level (debug, info, warn, error)")
	serverCmd.PersistentFlags().BoolVar(&args.TextFormat, "log-text", false, "log in text format, otherwise JSON")
	serverCmd.PersistentFlags().StringVar(&args.LogFile, "log-file", "", "log file path (if not set, logs to stdout)")

	cmd.AddCommand(serverCmd)

	return cmd, nil
}
