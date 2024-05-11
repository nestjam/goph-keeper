package main

import (
	"context"
	"log"
	"os"

	"github.com/nestjam/goph-keeper/internal/app/server"
)

func main() {
	ctx := context.Background()
	if err := server.NewApp().Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
