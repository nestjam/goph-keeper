package client

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	serverAddress = "serveraddress"
)

type Config struct {
	ServerAddress string
}

type ConfigOption func(*viper.Viper) error

func NewConfig(opts ...ConfigOption) (*Config, error) {
	const op = "new config"

	v := viper.New()

	for _, opt := range opts {
		err := opt(v)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
	}

	c := Config{}
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &c, nil
}

func FromArgs(args []string) ConfigOption {
	return func(v *viper.Viper) error {
		const op = "from args"

		flagSet := setupFlagSet()

		err := flagSet.Parse(args)
		if err != nil {
			return errors.Wrap(err, op)
		}

		err = v.BindPFlags(flagSet)
		if err != nil {
			return errors.Wrap(err, op)
		}

		return nil
	}
}

func setupFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	flagSet.StringP(serverAddress, "s", "", "server address")
	return flagSet
}
