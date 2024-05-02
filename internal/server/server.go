package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type Server struct {
	rootKey  *model.MasterKey
	baseURL  string
	certFile string
	keyFile  string
}

func New(baseURL, rootKey, certFile, keyFile string) *Server {
	return &Server{
		baseURL:  baseURL,
		rootKey:  model.NewMasterKey([]byte(rootKey)),
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *Server) Run(ctx context.Context) error {
	const op = "run server"

	h := s.mapHandlers()

	if err := http.ListenAndServeTLS(s.baseURL, s.certFile, s.keyFile, h); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
