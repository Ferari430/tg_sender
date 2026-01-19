package main

import (
	"context"
	"log"

	"github.com/Ferari430/tg_sender/cmd/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Println(err)
	}
	app.RunApp(context.Background())

}
