package persistence

import "github.com/danielpacak/myevents-events-service/domain"

// EventsRepository defines method for managing Event entities.
type EventsRepository interface {
	Create(event domain.Event) ([]byte, error)
	FindById(id []byte) (domain.Event, error)
	FindByName(name string) (domain.Event, error)
	FindAll() ([]domain.Event, error)
}
