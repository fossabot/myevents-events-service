package main

import (
	"flag"
	"fmt"
	"github.com/danielpacak/go-rest-api-seed/configuration"
	"github.com/danielpacak/go-rest-api-seed/dblayer"
	"github.com/danielpacak/go-rest-api-seed/rest"
	msgqueue_amqp "github.com/danielpacak/myevents-contracts/lib/msgqueue/amqp"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	fmt.Println("Connecting to message broker")
	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTlsendpoint, dbhandler, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
