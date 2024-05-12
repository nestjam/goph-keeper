package main

import (
	"log"
	"os"

	"github.com/nestjam/goph-keeper/internal/app/client"
)

var (
	BuildVersion string
	BuildDate    string
)

func main() {
	if err := client.NewApp().Run(BuildVersion, BuildDate, os.Args); err != nil {
		log.Fatal(err)
	}
}
