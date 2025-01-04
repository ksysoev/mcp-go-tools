// Package cmd implements the command-line interface for the MCP code tools server.
//
// It provides command initialization, configuration management, and logging setup.
// The package uses cobra for CLI implementation and viper for configuration handling,
// supporting both command-line flags and environment variables.
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// args holds all command-line arguments and configuration options.
// It supports both command-line flags and environment variables through
// viper's mapstructure tags.
type args struct {
	build      string
	version    string
	LogLevel   string `mapstructure:"loglevel"`
	ConfigPath string `mapstructure:"config"`
	TextFormat bool   `mapstructure:"logtext"`
	LogFile    string `mapstructure:"logfile"`
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
		Use:   "mcp",
		Short: "",
		Long:  "",
	}

	cmd.AddCommand(startCommand(args))

	cmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&args.LogLevel, "loglevel", "info", "log level (debug, info, warn, error)")
	cmd.PersistentFlags().BoolVar(&args.TextFormat, "logtext", false, "log in text format, otherwise JSON")
	cmd.PersistentFlags().StringVar(&args.LogFile, "logfile", "", "log file path (if not set, logs to stdout)")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, fmt.Errorf("failed to parse env args: %w", err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(args); err != nil {
		return nil, fmt.Errorf("failed to unmarshal args: %w", err)
	}

	return cmd, nil
}

// startCommand creates a new cobra.Command to start the MCP code tools server.
// It initializes logging, loads configuration, and starts the server.
// Returns error if logger initialization fails, configuration loading fails,
// or the server encounters an error during startup.
func startCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start MCP code tools server",
		Long:  "Start the Model Context Protocol server for code generation tools",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := initLogger(arg); err != nil {
				return err
			}

			slog.Info("Starting MCP code tools server", slog.String("version", arg.version), slog.String("build", arg.build))

			cfg, err := initConfig(arg)
			if err != nil {
				return err
			}

			return runStart(cmd.Context(), cfg)
		},
	}
}
