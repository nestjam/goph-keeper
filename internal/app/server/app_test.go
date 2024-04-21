package server_test

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/server"
)

func TestRun(t *testing.T) {
	const baseURL = "localhost:8080"
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sut := server.New(baseURL)

	go func() {
		_ = sut.Run(ctx)
	}()

	c := resty.New()
	c.BaseURL = "http://" + baseURL
	_, err := c.R().Post("/register")

	require.NoError(t, err)
}
