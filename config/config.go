package config

import (
	"os"
)

type DatabaseType string
type BrokerType string

const (
	MongoDB DatabaseType = "mongodb"
	AMQP    BrokerType   = "AMQP"
	Kafka   BrokerType   = "Kafka"

	DatabaseTypeDefault = MongoDB
	BrokerTypeDefault   = AMQP

	RestApiAddrDefault = ":8181"
	MetricsAddrDefault = ":9100"

	DefaultMongoDBConnectionURL = "mongodb://127.0.0.1"
	DefaultMongoDBDatabaseName  = "myevents"

	DefaultAMQPConnectionURI = "amqp://localhost:5672"

	DefaultKafkaBrokers = "localhost:9092"
)

type AppConfig struct {
	DatabaseType DatabaseType
	RestApiAddr  string
	MetricsAddr  string
	BrokerType   BrokerType

	MongoDBConfig *MongoDBConfig
	AMQPConfig    *AMQPConfig
	KafkaConfig   *KafkaConfig
}

type MongoDBConfig struct {
	ConnectionURL string
	DatabaseName  string
}

type AMQPConfig struct {
	ConnectionURI string
}

type KafkaConfig struct {
	Brokers string
}

func ExtractConfig() AppConfig {
	conf := AppConfig{
		RestApiAddr:  RestApiAddrDefault,
		MetricsAddr:  MetricsAddrDefault,
		DatabaseType: DatabaseTypeDefault,
		BrokerType:   BrokerTypeDefault,
		MongoDBConfig: &MongoDBConfig{
			ConnectionURL: DefaultMongoDBConnectionURL,
			DatabaseName:  DefaultMongoDBDatabaseName,
		},
		AMQPConfig: &AMQPConfig{
			ConnectionURI: DefaultAMQPConnectionURI,
		},
		KafkaConfig: &KafkaConfig{
			Brokers: DefaultKafkaBrokers,
		},
	}

	if connectionURL := os.Getenv("MONGODB_CONNECTION_URL"); connectionURL != "" {
		conf.MongoDBConfig.ConnectionURL = connectionURL
	}
	if dbName := os.Getenv("MONGODB_DATABASE_NAME"); dbName != "" {
		conf.MongoDBConfig.DatabaseName = dbName
	}
	if brokerUrl := os.Getenv("AMQP_CONNECTION_URI"); brokerUrl != "" {
		conf.AMQPConfig.ConnectionURI = brokerUrl
	}
	if brokerUrl := os.Getenv("KAFKA_BROKERS"); brokerUrl != "" {
		conf.KafkaConfig.Brokers = brokerUrl
	}
	if listenUrl := os.Getenv("LISTEN_URL"); listenUrl != "" {
		conf.RestApiAddr = listenUrl
	}

	return conf
}
