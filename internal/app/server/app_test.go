package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/server"
)

func TestRun(t *testing.T) {
	const baseURL = "localhost:8080"
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sut := server.New(baseURL)

	go func() {
		err := sut.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	c := resty.New()
	c.BaseURL = "http://" + baseURL
	r := c.R()
	r.Body = httpAuth.RegisterUserRequest{
		Email:    "user@email.com",
		Password: "1234",
	}
	resp, err := r.Post("/register")

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}
