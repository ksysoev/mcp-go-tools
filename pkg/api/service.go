package api

import (
	"context"
	"errors"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"golang.org/x/sync/errgroup"
)

type Config struct {
}

type Service struct {
	config *Config
}

func New(cfg *Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) Run(ctx context.Context) error {
	server := mcp.NewServer(stdio.NewStdioServerTransport())

	if err := s.setupTools(server); err != nil {
		return fmt.Errorf("failed to setup tools: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(server.Serve)

	eg.Go(func() error {
		<-ctx.Done()

		// TODO: Implement graceful shutdown, when it'll be supported by the mcp library.

		return ctx.Err()
	})

	err := eg.Wait()
	if errors.Is(err, context.Canceled) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to run service: %w", err)
	}

	return nil
}

func (s *Service) setupTools(*mcp.Server) error {
	//TODO: Implement setup tools for mcp server

	return nil
}
