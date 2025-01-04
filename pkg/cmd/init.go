package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type args struct {
	build      string
	version    string
	LogLevel   string `mapstructure:"loglevel"`
	ConfigPath string `mapstructure:"config"`
	TextFormat bool   `mapstructure:"logtext"`
	LogFile    string `mapstructure:"logfile"`
}

// InitCommands initializes and returns the root command for the Backend for Frontend (BFF) service.
// It sets up the command structure and adds subcommands, including setting up persistent flags.
// It returns a pointer to a cobra.Command which represents the root command.
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

// ServerCommand creates a new cobra.Command to start the BFF server for Deriv API.
// It takes cfgPath of type *string which is the path to the configuration file.
// It returns a pointer to a cobra.Command which can be executed to start the server.
// It returns an error if the logger initialization fails, the configuration cannot be loaded, or the server fails to run.
func startCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start BFF server",
		Long:  "Start BFF server for Deriv API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := initLogger(arg); err != nil {
				return err
			}

			slog.Info("Starting Deriv API BFF server", slog.String("version", arg.version), slog.String("build", arg.build))

			cfg, err := initConfig(arg)
			if err != nil {
				return err
			}

			return runStart(cmd.Context(), cfg)
		},
	}
}
