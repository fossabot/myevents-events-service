package config

import (
	"os"
)

const (
	RestApiAddrDefault = ":8181"
	MetricsAddrDefault = ":9100"

	DefaultMongoDBConnectionURL = "mongodb://127.0.0.1"
	DefaultMongoDBDatabaseName  = "events-svc"

	DefaultKafkaBrokers = "localhost:9092"
)

type AppConfig struct {
	RestApiAddr string
	MetricsAddr string

	MongoDBConfig MongoDBConfig
	KafkaConfig   KafkaConfig
}

type MongoDBConfig struct {
	ConnectionURL string
	DatabaseName  string
}

type KafkaConfig struct {
	Brokers string
}

func ExtractConfig() AppConfig {
	conf := AppConfig{
		RestApiAddr: RestApiAddrDefault,
		MetricsAddr: MetricsAddrDefault,
		MongoDBConfig: MongoDBConfig{
			ConnectionURL: DefaultMongoDBConnectionURL,
			DatabaseName:  DefaultMongoDBDatabaseName,
		},
		KafkaConfig: KafkaConfig{
			Brokers: DefaultKafkaBrokers,
		},
	}

	if connectionURL := os.Getenv("EVENTS_SVC_MONGODB_CONNECTION_URL"); connectionURL != "" {
		conf.MongoDBConfig.ConnectionURL = connectionURL
	}
	if dbName := os.Getenv("EVENTS_SVC_MONGODB_DATABASE_NAME"); dbName != "" {
		conf.MongoDBConfig.DatabaseName = dbName
	}
	if brokerUrl := os.Getenv("EVENTS_SVC_KAFKA_BROKERS"); brokerUrl != "" {
		conf.KafkaConfig.Brokers = brokerUrl
	}
	if listenUrl := os.Getenv("EVENTS_SVC_REST_API_TCP_ADDRESS"); listenUrl != "" {
		conf.RestApiAddr = listenUrl
	}
	if metricsUrl := os.Getenv("EVENTS_SVC_METRICS_TCP_ADDRESS"); metricsUrl != "" {
		conf.MetricsAddr = metricsUrl
	}

	return conf
}
