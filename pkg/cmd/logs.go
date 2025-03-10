// Package cmd implements the command-line interface for the MCP code tools server.
//
// This file provides logging configuration and initialization using slog.
// Logging features include:
// - JSON and text output formats
// - Configurable log levels (debug, info, warn, error)
// - File output support with automatic file creation
// - Version and application tagging for all log entries
package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// initLogger initializes the default logger for the application using slog.
// It configures the logger based on command-line arguments:
//   - LogLevel: Sets the minimum log level (debug, info, warn, error)
//   - TextFormat: Uses human-readable format instead of JSON
//   - LogFile: Writes logs to specified file
//
// The logger adds version and application tags to all log entries.
// Returns error if log level is invalid or file access fails.
func initLogger(arg *args) error {
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(arg.LogLevel))

	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Set up writer based on logfile flag
	var writer = io.Discard

	// Open log file if specified
	if arg.LogFile != "" {
		writer, err = os.OpenFile(arg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)

		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	// Create handler based on format
	var logHandler slog.Handler
	if arg.TextFormat {
		logHandler = slog.NewTextHandler(writer, options)
	} else {
		logHandler = slog.NewJSONHandler(writer, options)
	}

	logger := slog.New(logHandler).With(
		slog.String("ver", arg.version),
		slog.String("app", "mcp-go-tools"),
	)

	slog.SetDefault(logger)

	return nil
}
