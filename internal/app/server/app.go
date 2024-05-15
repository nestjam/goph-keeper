package server

import (
	"bytes"
	"context"
	"os"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/server"
	"github.com/nestjam/goph-keeper/migration"
)

type app struct {
}

func NewApp() *app {
	return &app{}
}

func (a *app) Run(ctx context.Context, args []string) error {
	const op = "run app"

	conf, err := getConfig(args)
	if err != nil {
		return errors.Wrap(err, op)
	}

	migrator := migration.NewDatabaseMigrator(conf.Postgres.DataSourceName)
	if err := migrator.Up(); err != nil {
		return errors.Wrap(err, op)
	}

	s := server.New(conf)

	if err := s.Run(ctx); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func getConfig(args []string) (*config.Config, error) {
	const op = "run app"

	confFile, err := config.NewConfigFile(args)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	opts := []config.ConfigOption{
		config.FromArgs(args),
	}

	if confFile != nil {
		b, err := os.ReadFile(confFile.Name)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
		opts = append(opts, config.FromYaml(bytes.NewReader(b)))
	}

	conf, err := config.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return conf, nil
}
