package client

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	config "github.com/nestjam/goph-keeper/internal/config/client"
	"github.com/nestjam/goph-keeper/internal/tui/auth"
)

type app struct {
}

func NewApp() *app {
	return &app{}
}

func (a *app) Run(buildVersion, buildDate string, args []string) error {
	const op = "run app"

	conf, err := getConfig(args)
	if err != nil {
		return errors.Wrap(err, op)
	}

	m := auth.NewLoginModel(conf.ServerAddress)
	m.BuildDate = buildDate
	m.BuildVersion = buildVersion

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func getConfig(args []string) (*config.Config, error) {
	const op = "run app"

	conf, err := config.NewConfig(config.FromArgs(args))
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return conf, nil
}
