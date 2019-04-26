package app

import (
	"github.com/Shopify/sarama"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue/kafka"
	"github.com/danielpacak/myevents-events-service/pkg/config"
	"github.com/danielpacak/myevents-events-service/pkg/metrics"
	"github.com/danielpacak/myevents-events-service/pkg/persistence"
	"github.com/danielpacak/myevents-events-service/pkg/persistence/mongo"
	"github.com/danielpacak/myevents-events-service/pkg/http/rest"
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
		server := rest.NewHandler(a.repository, a.emitter)
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
	return mongo.NewMongoEventsRepository(cfg.MongoDBConfig)
}

func newEventEmitter(cfg config.AppConfig) (msgqueue.EventEmitter, error) {
	return newKafkaEventEmitter(cfg.KafkaConfig)
}

func newKafkaEventEmitter(cfg config.KafkaConfig) (msgqueue.EventEmitter, error) {
	log.Printf("Connecting to Kafka message broker at %s", cfg.Brokers)
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = "myevents-events-service"
	client, err := sarama.NewClient([]string{cfg.Brokers}, saramaConfig)
	emitter, err := kafka.NewKafkaEventEmitter(client)
	if err != nil {
		return nil, err
	}

	return emitter, nil
}
