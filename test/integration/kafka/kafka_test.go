package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue/kafka"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (e *TestEvent) PartitionKey() string {
	return e.ID
}

func (e *TestEvent) EventName() string {
	return "event.test"
}

func TestEventEmitterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test")
	}
	time.Sleep(30 * time.Second)
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	saramaConfig.ClientID = "integration-test"
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Net.MaxOpenRequests = 1
	saramaConfig.Producer.Idempotent = true
	client, err := sarama.NewClient([]string{"kafka:9092"}, saramaConfig)
	require.NoError(t, err)

	emitter, err := kafka.NewKafkaEventEmitter(client)
	require.NoError(t, err)

	err = emitter.Emit(&TestEvent{
		ID:   "test123",
		Name: "Test event 123",
	})
	require.NoError(t, err)

	adminConfig := sarama.NewConfig()
	adminConfig.Version= sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin([]string{"kafka:9092"}, adminConfig)
	require.NoError(t, err)

	topics, err := admin.ListTopics()
	require.NoError(t, err)
	t.Logf("Admin:topics: %v", topics)
}
