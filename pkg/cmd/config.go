package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ksysoev/mcp-code-tools/pkg/api"
	"github.com/spf13/viper"
)

type Config struct {
	API api.Config `mapstructure:"api"`
}

// initConfig initializes the configuration by reading from the specified config file.
// It takes configPath of type string which is the path to the configuration file.
// It returns a pointer to a config struct and an error.
// It returns an error if the configuration file cannot be read or if the configuration cannot be unmarshaled.
func initConfig(arg *args) (*Config, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	if arg.ConfigPath != "" {
		v.SetConfigFile(arg.ConfigPath)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
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
