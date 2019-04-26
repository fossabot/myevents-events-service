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
			RestApiAddr: ":8181",
			MetricsAddr: ":9100",
			MongoDBConfig: MongoDBConfig{
				ConnectionURL: "mongodb://127.0.0.1",
				DatabaseName:  "events-svc",
			},
			KafkaConfig: KafkaConfig{
				Brokers: "localhost:9092",
			},
		}, config)
	})

	t.Run("should override mongo url with env variable", func(t *testing.T) {
		// given
		_ = os.Setenv("EVENTS_SVC_MONGODB_CONNECTION_URL", "mongodb://somewhere:1234")
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, "mongodb://somewhere:1234", config.MongoDBConfig.ConnectionURL)
	})

	t.Run("should override REST API listen address with env variable", func(t *testing.T) {
		// given
		_ = os.Setenv("EVENTS_SVC_REST_API_TCP_ADDRESS", ":32000")
		// when
		config := ExtractConfig()
		// then
		assert.Equal(t, ":32000", config.RestApiAddr)
	})

}
