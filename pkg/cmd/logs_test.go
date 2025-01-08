package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		args      *args
		name      string
		errMsg    string
		wantError bool
	}{
		{
			name: "valid log level - info",
			args: &args{
				LogLevel:   "info",
				TextFormat: true,
			},
			wantError: false,
		},
		{
			name: "valid log level - debug",
			args: &args{
				LogLevel:   "debug",
				TextFormat: true,
			},
			wantError: false,
		},
		{
			name: "invalid log level",
			args: &args{
				LogLevel:   "invalid",
				TextFormat: true,
			},
			wantError: true,
			errMsg:    "invalid log level",
		},
		{
			name: "append to existing file",
			args: &args{
				LogLevel:   "info",
				TextFormat: true,
				LogFile:    filepath.Join(t.TempDir(), "test.log"),
			},
			wantError: false,
		},
		{
			name: "invalid file path",
			args: &args{
				LogLevel:   "info",
				TextFormat: true,
				LogFile:    "/invalid/path/test.log",
			},
			wantError: true,
			errMsg:    "failed to open log file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initLogger(tt.args)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
