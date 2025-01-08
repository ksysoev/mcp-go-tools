package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfigWithInvalidContent(t *testing.T) {
	// Create a temporary config file with invalid YAML
	configContent := `
api: {
  invalid: yaml: content:
    missing: quotes
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	// Test config loading with invalid content
	args := &args{
		ConfigPath: configPath,
	}

	cfg, err := initConfig(args)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "did not find expected ',' or '}'")
}

func TestInitConfigWithDifferentFileTypes(t *testing.T) {
	tests := []struct {
		name         string
		fileContent  string
		fileExt      string
		errorMessage string
		wantError    bool
	}{
		{
			name: "valid yaml",
			fileContent: `
api: {}
repository:
  type: "static"
  rules:
    - name: "test_rule"
      category: "testing"
`,
			fileExt:   ".yaml",
			wantError: false,
		},
		{
			name: "valid json",
			fileContent: `{
				"api": {},
				"repository": {
					"type": "static",
					"rules": [{
						"name": "test_rule",
						"category": "testing"
					}]
				}
			}`,
			fileExt:   ".json",
			wantError: false,
		},
		{
			name: "invalid extension",
			fileContent: `
api: {}
repository:
  type: "static"
  rules:
    - name: "test_rule"
`,
			fileExt:      ".invalid",
			wantError:    true,
			errorMessage: "failed to read config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config"+tt.fileExt)
			err := os.WriteFile(configPath, []byte(tt.fileContent), 0o600)
			require.NoError(t, err)

			args := &args{
				ConfigPath: configPath,
			}

			cfg, err := initConfig(args)

			if tt.wantError {
				assert.Error(t, err)

				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}

				assert.Nil(t, cfg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.NotEmpty(t, cfg.Repository.Rules)
			assert.Equal(t, "test_rule", cfg.Repository.Rules[0].Name)
			assert.Equal(t, "static", string(cfg.Repository.Type))
		})
	}
}
