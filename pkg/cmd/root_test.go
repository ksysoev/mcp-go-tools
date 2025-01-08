package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCommands(t *testing.T) {
	tests := []struct {
		name      string
		build     string
		version   string
		wantError bool
	}{
		{
			name:      "successful initialization",
			build:     "test-build",
			version:   "1.0.0",
			wantError: false,
		},
		{
			name:      "empty build and version",
			build:     "",
			version:   "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			cmd, err := InitCommands(tt.build, tt.version)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cmd)

			// Verify command structure
			assert.Equal(t, "mcp-go-tools", cmd.Use)
			assert.Equal(t, "MCP code tools server", cmd.Short)
			assert.Equal(t, "Model Context Protocol server for code generation tools", cmd.Long)
			assert.Equal(t, tt.version+" (Build: "+tt.build+")", cmd.Version)

			// Verify subcommands
			subCmds := cmd.Commands()
			require.Len(t, subCmds, 1)
			serverCmd := subCmds[0]
			assert.Equal(t, "server", serverCmd.Use)
			assert.Equal(t, "Start MCP code tools server", serverCmd.Short)

			// Verify flags
			flags := serverCmd.PersistentFlags()

			configFlag := flags.Lookup("config")
			require.NotNil(t, configFlag)
			assert.Equal(t, "", configFlag.DefValue)

			logLevelFlag := flags.Lookup("log-level")
			require.NotNil(t, logLevelFlag)
			assert.Equal(t, "info", logLevelFlag.DefValue)

			logTextFlag := flags.Lookup("log-text")
			require.NotNil(t, logTextFlag)
			assert.Equal(t, "false", logTextFlag.DefValue)

			logFileFlag := flags.Lookup("log-file")
			require.NotNil(t, logFileFlag)
			assert.Equal(t, "", logFileFlag.DefValue)
		})
	}
}

func TestServerCommandExecution(t *testing.T) {
	// Create a test config file
	configContent := `
api: {}
repository:
  type: "static"
  rules:
    - name: "test_rule"
      category: "testing"
      description: "test rule"
      language: "go"
      examples:
        - description: "Example"
          code: "test code"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	// Create a test log file directory
	logDir := filepath.Join(tmpDir, "logs")
	err = os.MkdirAll(logDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "missing config flag",
			args:      []string{"server"},
			wantError: true,
		},
		{
			name: "invalid log level",
			args: []string{
				"server",
				"--config", configPath,
				"--log-level", "invalid",
			},
			wantError: true,
		},
		{
			name: "invalid log file path",
			args: []string{
				"server",
				"--config", configPath,
				"--log-file", "/invalid/path/test.log",
			},
			wantError: true,
		},
		{
			name: "invalid config file",
			args: []string{
				"server",
				"--config", filepath.Join(tmpDir, "nonexistent.yaml"),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cmd, err := InitCommands("test", "1.0.0")
			require.NoError(t, err)

			// Set args
			cmd.SetArgs(tt.args)

			// Act
			err = cmd.Execute()

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify log file was created if specified
			if logFile := cmd.Flags().Lookup("log-file").Value.String(); logFile != "" {
				_, err := os.Stat(logFile)
				assert.NoError(t, err)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		name        string
		build       string
		version     string
		wantVersion string
	}{
		{
			name:        "full version info",
			build:       "abc123",
			version:     "1.0.0",
			wantVersion: "1.0.0 (Build: abc123)",
		},
		{
			name:        "empty build",
			build:       "",
			version:     "1.0.0",
			wantVersion: "1.0.0 (Build: )",
		},
		{
			name:        "empty version",
			build:       "abc123",
			version:     "",
			wantVersion: " (Build: abc123)",
		},
		{
			name:        "both empty",
			build:       "",
			version:     "",
			wantVersion: " (Build: )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := InitCommands(tt.build, tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.wantVersion, cmd.Version)
		})
	}
}
