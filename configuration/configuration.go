package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/danielpacak/myevents-event-service/dblayer"
	"os"
)

var (
	DBTypeDefault            = dblayer.DBTYPE("mongodb")
	DBConnectionDefault      = "mongodb://127.0.0.1"
	RestfulEPDefault         = "localhost:8181"
	RestfulTlsEPDefault      = "localhost:4343"
	AMQPMessageBrokerDefault = "amqp://localhost:5672"
)

type ServiceConfig struct {
	Databasetype       dblayer.DBTYPE `json:"databasetype"`
	DBConnection       string         `json:"dbconnection"`
	RestfulEndpoint    string         `json:"restfulapi_endpoint"`
	RestfulTlsendpoint string         `json:"restfulapi_tlsendpoint"`
	AMQPMessageBroker  string         `json:"amqp_message_broker"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTlsEPDefault,
		AMQPMessageBrokerDefault,
	}
	if broker := os.Getenv("AMQP_URL"); broker != "" {
		conf.AMQPMessageBroker = broker
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}
	err = json.NewDecoder(file).Decode(&conf)

	return conf, err
}
