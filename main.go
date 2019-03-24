package main

import (
	"github.com/danielpacak/myevents-events-service/app"
	"github.com/danielpacak/myevents-events-service/config"
)

func main() {
	appConfig := config.ExtractConfig()

	application, err := app.NewApp(appConfig)
	if err != nil {
		panic(err)
	}
	application.Start()
}
