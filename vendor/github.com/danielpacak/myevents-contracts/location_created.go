package contracts

type LocationCreatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (e *LocationCreatedEvent) PartitionKey() string {
	return e.ID
}

func (e *LocationCreatedEvent) EventName() string {
	return "location.created"
}
