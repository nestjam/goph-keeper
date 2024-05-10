package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
)

type PGSQLRepositoryTestHelper struct {
	DataSourceName string
}

func (h *PGSQLRepositoryTestHelper) Run(m *testing.M) {
	const (
		userName = "user"
		dbName   = "goph-keeper"
		hostPort = "5432/tcp"
		password = "secret"
	)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + userName,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort(hostPort)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", userName, password, hostAndPort, dbName)

	log.Println("connecting to database on url: ", dsn)

	const seconds = 120
	_ = resource.Expire(seconds) // tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = seconds * time.Second
	if err = pool.Retry(func() error {
		const op = "retry"
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return errors.Wrap(err, op)
		}

		err = db.Ping()
		if err != nil {
			return errors.Wrap(err, op)
		}

		return nil
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	h.DataSourceName = dsn

	code := m.Run()

	// you can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}
