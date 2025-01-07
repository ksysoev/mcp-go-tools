package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerFileHandling(t *testing.T) {
	tests := []struct {
		setupFile func(t *testing.T) string
		name      string
		wantError bool
	}{
		{
			name: "append to existing file",
			setupFile: func(t *testing.T) string {
				t.Helper()
				file := filepath.Join(t.TempDir(), "existing.log")
				err := os.WriteFile(file, []byte("existing content\n"), 0o600)
				require.NoError(t, err)
				return file
			},
			wantError: false,
		},
		{
			name: "invalid path",
			setupFile: func(t *testing.T) string {
				t.Helper()
				return filepath.Join(t.TempDir(), "subdir", "test.log")
			},
			wantError: true,
		},
		{
			name: "invalid permissions",
			setupFile: func(t *testing.T) string {
				t.Helper()
				dir := filepath.Join(t.TempDir(), "readonly")
				err := os.Mkdir(dir, 0o500) // Read-only directory
				require.NoError(t, err)
				return filepath.Join(dir, "test.log")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file
			logFile := tt.setupFile(t)

			// Initialize logger
			args := &args{
				LogLevel:   "info",
				TextFormat: true,
				LogFile:    logFile,
			}

			err := initLogger(args)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file exists
			_, err = os.Stat(logFile)
			assert.NoError(t, err)

			if strings.Contains(tt.name, "append") {
				// Verify content was appended
				content, err := os.ReadFile(logFile)
				require.NoError(t, err)
				assert.Contains(t, string(content), "existing content")
			}
		})
	}
}
