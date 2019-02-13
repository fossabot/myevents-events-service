package main

import (
	"fmt"
	msgqueue_amqp "github.com/danielpacak/myevents-contracts/lib/msgqueue/amqp"
	"github.com/danielpacak/myevents-events-service/configuration"
	"github.com/danielpacak/myevents-events-service/dblayer"
	"github.com/danielpacak/myevents-events-service/rest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	config, _ := configuration.ExtractConfiguration()
	fmt.Printf("Connecting to database at %s\n", config.DBConnection)
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	fmt.Printf("Connecting to message broker at %s\n", config.AMQPMessageBroker)
	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Serving metrics API")
		h := http.NewServeMux()
		h.Handle("/metrics", prometheus.Handler())
		err := http.ListenAndServe(":9100", h)
		if err != nil {
			panic(err)
		}
	}()

	httpErrChan := rest.ServeAPI(config.RestfulEndpoint, dbhandler, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	}
}
