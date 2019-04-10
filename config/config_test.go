package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractConfig(t *testing.T) {

	t.Run("should return default values", func(t *testing.T) {
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, AppConfig{
			DatabaseType: MongoDB,
			RestApiAddr:  ":8181",
			MetricsAddr:  ":9100",
			BrokerType:   AMQP,
			MongoDBConfig: &MongoDBConfig{
				ConnectionURL: "mongodb://127.0.0.1",
				DatabaseName:  "myevents",
			},
			AMQPConfig: &AMQPConfig{
				ConnectionURI: "amqp://localhost:5672",
			},
			KafkaConfig: &KafkaConfig{
				Brokers: "localhost:9092",
			},
		}, config)
	})

	t.Run("should override broker url with env variable", func(t *testing.T) {
		// given
		_ = os.Setenv("AMQP_CONNECTION_URI", "amqp://somewhere:1234")
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, "amqp://somewhere:1234", config.AMQPConfig.ConnectionURI)
	})

	t.Run("should override mongo url with env variable", func(t *testing.T) {
		// given
		_ = os.Setenv("MONGODB_CONNECTION_URL", "mongodb://somewhere:1234")
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, "mongodb://somewhere:1234", config.MongoDBConfig.ConnectionURL)
	})

	t.Run("should override REST API listen address with env variable", func(t *testing.T) {
		// given
		_ = os.Setenv("LISTEN_URL", ":32000")
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, ":32000", config.RestApiAddr)
	})

}
