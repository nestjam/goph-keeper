package server

import (
	"context"

	"github.com/nestjam/goph-keeper/internal/server"
	"github.com/pkg/errors"
)

type app struct {
}

func NewApp() *app {
	return &app{}
}

func (a *app) Run(ctx context.Context) error {
	const op = "run app"
	const baseURL = "localhost:8080"

	s := server.New(baseURL)

	if err := s.Run(ctx); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
