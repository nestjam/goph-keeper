package main

import (
	"bytes"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

func main() {
	const (
		certFile = "servercert.crt"
		keyfile  = "servercert.key"
	)

	if !exists(certFile) || !exists(keyfile) {
		if err := generate(certFile, keyfile); err != nil {
			log.Fatal(err)
		}
	}
}

func generate(certFile, keyfile string) error {
	const op = "generate certificate and key"
	cert, key, err := utils.GenerateCert()
	if err != nil {
		return errors.Wrap(err, op)
	}

	if err = writeFile(certFile, cert); err != nil {
		return errors.Wrap(err, op)
	}

	if err = writeFile(keyfile, key); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func writeFile(name string, data bytes.Buffer) error {
	const perm os.FileMode = 0600
	err := os.WriteFile(name, data.Bytes(), perm)
	if err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
