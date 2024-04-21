package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type Server struct {
	baseURL string
}

func New(baseURL string) *Server {
	return &Server{
		baseURL: baseURL,
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
