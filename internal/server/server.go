package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type Server struct {
	conf    *config.Config
	rootKey *model.MasterKey
}

func New(conf *config.Config) *Server {
	return &Server{
		conf:    conf,
		rootKey: model.NewMasterKey([]byte(conf.Vault.MasterKey)),
	}
}

func (s *Server) Run(ctx context.Context) error {
	const op = "run server"

	h, err := s.mapHandlers(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}

	c := s.conf.Server
	if err := http.ListenAndServeTLS(c.Address, c.CertFile, c.KeyFile, h); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
