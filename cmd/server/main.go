package main

import (
	"context"
	"log"

	"github.com/nestjam/goph-keeper/internal/app/server"
)

func main() {
	ctx := context.Background()
	if err := server.NewApp().Run(ctx); err != nil {
		log.Fatal(err)
	}
}
