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
	const (
		baseURL = "localhost:8080"
		rootKey = "N3SaEN8k2z3?DCf@4&8j+Yc92pTrFt6W"
	)

	s := server.New(baseURL, rootKey)

	if err := s.Run(ctx); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
