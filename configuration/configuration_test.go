package configuration

import (
	"github.com/danielpacak/myevents-events-service/dblayer"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractConfiguration(t *testing.T) {

	t.Run("should return default values", func(t *testing.T) {
		config, err := ExtractConfiguration()

		assert.NoError(t, err)
		assert.Equal(t, ServiceConfig{DatabaseType: dblayer.MONGODB,
			DBConnection:      "mongodb://127.0.0.1",
			RestApiAddr:       ":8181",
			MetricsAddr:       ":9100",
			AMQPMessageBroker: "amqp://localhost:5672",
		}, config)
	})

	t.Run("should override broker url with env variable", func(t *testing.T) {
		_ = os.Setenv("AMQP_BROKER_URL", "amqp://somewhere:1234")
		config, err := ExtractConfiguration()
		assert.NoError(t, err)
		assert.Equal(t, "amqp://somewhere:1234", config.AMQPMessageBroker)
	})

	t.Run("should override mongo url with env variable", func(t *testing.T) {
		_ = os.Setenv("MONGO_URL", "mongodb://somewhere:1234")
		config, err := ExtractConfiguration()
		assert.NoError(t, err)
		assert.Equal(t, "mongodb://somewhere:1234", config.DBConnection)
	})

}