package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// initLogger initializes the default logger for the application using slog.
// It supports writing logs to both stdout and an optional log file.
// The log format can be either JSON (default) or text based on the TextFormat flag.
// Returns an error if logger initialization fails, including file access errors.
func initLogger(arg *args) error {
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(arg.LogLevel)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Set up writers
	writers := []io.Writer{os.Stdout}
	var logFile *os.File
	if arg.LogFile != "" {
		var err error
		logFile, err = os.OpenFile(arg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, logFile)
	}

	// Create multi-writer
	writer := io.MultiWriter(writers...)

	// Create handler based on format
	var logHandler slog.Handler
	if arg.TextFormat {
		logHandler = slog.NewTextHandler(writer, options)
	} else {
		logHandler = slog.NewJSONHandler(writer, options)
	}

	logger := slog.New(logHandler).With(
		slog.String("ver", arg.version),
		slog.String("app", "mcp-code-tools"),
	)

	slog.SetDefault(logger)

	return nil
}
