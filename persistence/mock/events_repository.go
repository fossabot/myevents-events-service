package mock

import (
	"github.com/danielpacak/myevents-events-service/domain"
	"github.com/stretchr/testify/mock"
)

type EventsRepository struct {
	mock.Mock
}

func (m *EventsRepository) Create(event domain.Event) ([]byte, error) {
	return nil, nil
}

func (m *EventsRepository) FindById(id []byte) (domain.Event, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Event), args.Error(1)
}

func (m *EventsRepository) FindByName(name string) (domain.Event, error) {
	args := m.Called(name)
	return args.Get(0).(domain.Event), args.Error(1)
}

func (m *EventsRepository) FindAll() ([]domain.Event, error) {
	return nil, nil
}
