package dblayer

import (
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/danielpacak/myevents-events-service/mongolayer"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/streadway/amqp"
	msgqueue_amqp "github.com/danielpacak/myevents-contracts/lib/msgqueue/amqp"
	"log"
)

type DBTYPE string

const (
	MONGODB DBTYPE = "mongodb"
)

func NewEventsRepository(options DBTYPE, connection string) (persistence.EventsRepository, error) {
	log.Printf("Connecting to database at %s\n", connection)

	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}

func NewEventEmitter(brokerUrl string) (msgqueue.EventEmitter, error) {
	log.Printf("Connecting to message broker at %s\n", brokerUrl)

	conn, err := amqp.Dial(brokerUrl)
	if err != nil {
		return nil, err
	}
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		return nil, err
	}
	return emitter, nil
}
