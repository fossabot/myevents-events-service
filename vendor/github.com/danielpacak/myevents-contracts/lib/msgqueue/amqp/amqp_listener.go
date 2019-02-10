package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/danielpacak/myevents-contracts"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: conn,
		queue:      queue,
	}
	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil
	}
	defer channel.Close()
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}
	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}
	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("msg did not contain x-event-name header")
				msg.Nack(false, true)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name header is not string, but %t", rawEventName)
				msg.Nack(false, true)
				continue
			}
			var event msgqueue.Event
			switch eventName {
			case "event.created":
				event = new(contracts.EventCreatedEvent)
			default:
				errors <- fmt.Errorf("event type %s is unknown", eventName)
				continue
			}
			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- err
				continue
			}
			events <- event
		}
	}()
	return events, errors, nil
}
