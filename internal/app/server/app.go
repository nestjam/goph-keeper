package server

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/server"
)

type app struct {
}

func NewApp() *app {
	return &app{}
}

func (a *app) Run(ctx context.Context) error {
	const op = "run app"

	const (
		baseURL  = "localhost:8080"
		rootKey  = "N3SaEN8k2z3?DCf@4&8j+Yc92pTrFt6W"
		certFile = "servercert.crt"
		keyfile  = "servercert.key"
	)

	s := server.New(baseURL, rootKey, certFile, keyfile)

	if err := s.Run(ctx); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
