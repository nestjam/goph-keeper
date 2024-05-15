package config

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	serverAddress  = "server.address"
	serverCertFile = "server.certfile"
	serverKeyFile  = "server.keyfile"
	vaultMasterKey = "vault.masterkey"

	defaultServerAddress = "localhost:8080"
	defaultCertFile      = "servercert.crt"
	defaultKeyFile       = "servercert.key"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Vault    VaultConfig
	JWTAuth  JWTAuthConfig
}

type ServerConfig struct {
	Address  string
	CertFile string
	KeyFile  string
}

type PostgresConfig struct {
	DataSourceName string
}

type JWTAuthConfig struct {
	SignKey       string
	TokenExpiryIn time.Duration
}

type VaultConfig struct {
	MasterKey string
}

type ConfigOption func(*viper.Viper) error

func New(opts ...ConfigOption) (*Config, error) {
	const op = "new config"

	v := viper.New()
	setDefaults(v)

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

func setDefaults(v *viper.Viper) {
	v.SetDefault(serverAddress, defaultServerAddress)
	v.SetDefault(serverCertFile, defaultCertFile)
	v.SetDefault(serverKeyFile, defaultKeyFile)
}

func FromYaml(in io.Reader) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigType("yaml")

		err := v.ReadConfig(in)
		if err != nil {
			return errors.Wrap(err, "from yaml")
		}

		return nil
	}
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
	flagSet.StringP(configTag, "c", "", "config file path")
	flagSet.StringP(serverAddress, "a", "", "server address")
	flagSet.StringP(vaultMasterKey, "k", "", "vault master key")
	return flagSet
}
