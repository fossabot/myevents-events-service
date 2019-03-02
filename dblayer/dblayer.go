package dblayer

import (
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	msgqueue_amqp "github.com/danielpacak/myevents-contracts/lib/msgqueue/amqp"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/danielpacak/myevents-events-service/persistence/mongo"
	"github.com/streadway/amqp"
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
		return mongo.NewMongoEventsRepository(connection)
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
