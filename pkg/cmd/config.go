package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ksysoev/mcp-go-tools/pkg/api"
	"github.com/ksysoev/mcp-go-tools/pkg/repo"
	"github.com/spf13/viper"
)

// Config represents the complete application configuration structure.
// It combines API service configuration and repository configuration loaded from
// configuration files and environment variables.
type Config struct {
	// API holds the MCP server configuration
	API api.Config `mapstructure:"api"`
	// Repository defines the repository configuration including type and rules
	Repository repo.Config `mapstructure:"repository"`
}

// initConfig initializes the configuration from the specified file and environment.
// It supports both YAML/JSON configuration files and environment variables,
// where environment variables override file settings. Environment variables
// use underscore (_) as separator for nested fields (e.g., "api_port").
//
// The function logs the final configuration at debug level for troubleshooting.
// Returns error if the configuration file cannot be read or parsed.
func initConfig(arg *args) (*Config, error) {
	v := viper.NewWithOptions()

	v.SetConfigFile(arg.ConfigPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	slog.Debug("Config loaded", slog.Any("config", cfg))

	return &cfg, nil
}
