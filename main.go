package main

import (
	"github.com/danielpacak/myevents-events-service/configuration"
	"github.com/danielpacak/myevents-events-service/dblayer"
	"github.com/danielpacak/myevents-events-service/rest"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

func main() {
	config, _ := configuration.ExtractConfiguration()

	repository, err := dblayer.NewEventsRepository(config.DatabaseType, config.DBConnection)
	if err != nil {
		panic(err)
	}

	emitter, err := dblayer.NewEventEmitter(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Printf("Serving metrics at `%s`", config.MetricsAddr)
		h := http.NewServeMux()
		h.Handle("/metrics", prometheus.Handler())
		err := http.ListenAndServe(config.MetricsAddr, h)
		if err != nil {
			panic(err)
		}
	}()

	httpErrChan := rest.ServeAPI(config.RestApiAddr, repository, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	}
}
