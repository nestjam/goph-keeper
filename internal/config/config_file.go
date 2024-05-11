package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	configTag = "config"

	defaultConfigFilePath = "config.yml"
)

type ConfigFile struct {
	Name string `mapstructure:"config"`
}

func NewConfigFile(args []string) (*ConfigFile, error) {
	const op = "new config file"

	v := viper.New()

	v.SetDefault(configTag, defaultConfigFilePath)

	flagSet := setupFlagSet()

	err := flagSet.Parse(args)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	err = v.BindPFlags(flagSet)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	c := ConfigFile{}
	err = v.Unmarshal(&c)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return &c, nil
}
