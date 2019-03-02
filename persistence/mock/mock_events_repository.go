package mock

import (
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/stretchr/testify/mock"
)

type MockEventsRepository struct {
	mock.Mock
}

func (m *MockEventsRepository) Create(event persistence.Event) ([]byte, error) {
	return nil, nil
}

func (m *MockEventsRepository) FindById(id []byte) (persistence.Event, error) {
	args := m.Called(id)
	return args.Get(0).(persistence.Event), args.Error(1)
}

func (m *MockEventsRepository) FindByName(name string) (persistence.Event, error) {
	args := m.Called(name)
	return args.Get(0).(persistence.Event), args.Error(1)
}

func (m *MockEventsRepository) FindAll() ([]persistence.Event, error) {
	return nil, nil
}
