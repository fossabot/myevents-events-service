package persistence

type EventsRepository interface {
	Create(event Event) ([]byte, error)
	FindById(id []byte) (Event, error)
	FindByName(name string) (Event, error)
	FindAll() ([]Event, error)
}
