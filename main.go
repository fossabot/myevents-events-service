package main

import (
	"github.com/danielpacak/myevents-events-service/configuration"
	"github.com/danielpacak/myevents-events-service/dblayer"
	"github.com/danielpacak/myevents-events-service/metrics"
	"github.com/danielpacak/myevents-events-service/rest"
	"log"
	"net/http"
	"sync"
)

func main() {
	config := configuration.ExtractConfiguration()

	repository, err := dblayer.NewEventsRepository(config.DatabaseType, config.DBConnection)
	if err != nil {
		panic(err)
	}

	emitter, err := dblayer.NewEventEmitter(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		server := metrics.NewMetricsServer()
		log.Printf("Serving metrics at `%s`", config.MetricsAddr)
		err := http.ListenAndServe(config.MetricsAddr, server)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	go func() {
		server := rest.NewAPIServer(repository, emitter)
		log.Printf("Serving API at `%s`", config.RestApiAddr)
		err := http.ListenAndServe(config.RestApiAddr, server)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
