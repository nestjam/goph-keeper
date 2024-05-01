package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type Server struct {
	rootKey *model.MasterKey
	baseURL string
}

func New(baseURL, rootKey string) *Server {
	return &Server{
		baseURL: baseURL,
		rootKey: model.NewMasterKey([]byte(rootKey)),
	}
}

func (s *Server) Run(ctx context.Context) error {
	const op = "run server"

	h := s.mapHandlers()

	if err := http.ListenAndServe(s.baseURL, h); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
