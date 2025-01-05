package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksysoev/mcp-go-tools/pkg/cmd"
)

// version is the version of the application. It should be set at build time.
var version = "dev"
var build = "dev"

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	rootCmd, err := cmd.InitCommands(build, version)
	if err != nil {
		slog.Error("failed to initialize commands", slog.Any("error", err))
		os.Exit(1)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("failed to execute command", slog.Any("error", err))
		os.Exit(1)
	}
}
