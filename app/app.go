package app

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	msgqueue_amqp "github.com/danielpacak/myevents-contracts/lib/msgqueue/amqp"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue/kafka"
	"github.com/danielpacak/myevents-events-service/config"
	"github.com/danielpacak/myevents-events-service/metrics"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/danielpacak/myevents-events-service/persistence/mongo"
	"github.com/danielpacak/myevents-events-service/rest"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"sync"
)

type App struct {
	cfg        config.AppConfig
	repository persistence.EventsRepository
	emitter    msgqueue.EventEmitter
}

func NewApp(cfg config.AppConfig) (*App, error) {
	repository, err := newEventsRepository(cfg)
	if err != nil {
		panic(err)
	}

	emitter, err := newEventEmitter(cfg)
	if err != nil {
		panic(err)
	}

	return &App{cfg: cfg, repository: repository, emitter: emitter}, nil
}

func (a *App) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		server := metrics.NewMetricsServer()
		log.Printf("Serving metrics at `%s`", a.cfg.MetricsAddr)
		err := http.ListenAndServe(a.cfg.MetricsAddr, server)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	go func() {
		server := rest.NewAPIServer(a.repository, a.emitter)
		log.Printf("Serving API at `%s`", a.cfg.RestApiAddr)
		err := http.ListenAndServe(a.cfg.RestApiAddr, server)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func newEventsRepository(cfg config.AppConfig) (persistence.EventsRepository, error) {
	switch cfg.DatabaseType {
	case config.MongoDB:
		return mongo.NewMongoEventsRepository(cfg.MongoDBConfig)
	}
	return nil, nil
}

func newEventEmitter(cfg config.AppConfig) (msgqueue.EventEmitter, error) {
	switch cfg.BrokerType {
	case config.AMQP:
		return newAMQPEventEmitter(cfg.AMQPConfig)
	case config.Kafka:
		return newKafkaEventEmitter(cfg.KafkaConfig)
	}
	return nil, fmt.Errorf("unrecognized broker type %s", cfg.BrokerType)
}

func newAMQPEventEmitter(cfg *config.AMQPConfig) (msgqueue.EventEmitter, error) {
	log.Printf("Connecting to AMQP message broker at %s", cfg.ConnectionURI)
	conn, err := amqp.Dial(cfg.ConnectionURI)
	if err != nil {
		return nil, err
	}
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		return nil, err
	}
	return emitter, nil
}

func newKafkaEventEmitter(cfg *config.KafkaConfig) (msgqueue.EventEmitter, error) {
	log.Printf("Connecting to Kafka message broker at %s", cfg.Brokers)
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	saramaConfig.ClientID = "myevents-events-service"
	client, err := sarama.NewClient([]string{cfg.Brokers}, saramaConfig)
	emitter, err := kafka.NewKafkaEventEmitter(client)
	if err != nil {
		return nil, err
	}

	return emitter, nil
}