package msgqueue

type Event interface {
	PartitionKey() string
	EventName() string
}
