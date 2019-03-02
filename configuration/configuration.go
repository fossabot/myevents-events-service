package configuration

import (
	"github.com/danielpacak/myevents-events-service/dblayer"
	"os"
)

var (
	DBTypeDefault            = dblayer.DBTYPE("mongodb")
	DBConnectionDefault      = "mongodb://127.0.0.1"
	RestApiAddrDefault       = ":8181"
	MetricsAddrDefault       = ":9100"
	AMQPMessageBrokerDefault = "amqp://localhost:5672"
)

type ServiceConfig struct {
	DatabaseType      dblayer.DBTYPE
	DBConnection      string
	RestApiAddr       string
	MetricsAddr       string
	AMQPMessageBroker string
}

func ExtractConfiguration() ServiceConfig {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestApiAddrDefault,
		MetricsAddrDefault,
		AMQPMessageBrokerDefault,
	}
	if brokerUrl := os.Getenv("AMQP_BROKER_URL"); brokerUrl != "" {
		conf.AMQPMessageBroker = brokerUrl
	}
	if mongoUrl := os.Getenv("MONGO_URL"); mongoUrl != "" {
		conf.DBConnection = mongoUrl
	}
	if listenUrl := os.Getenv("LISTEN_URL"); listenUrl != "" {
		conf.RestApiAddr = listenUrl
	}

	return conf
}
