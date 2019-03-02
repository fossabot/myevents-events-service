package persistence

import (
	"github.com/stretchr/testify/mock"
)

type MockEventsRepository struct {
	mock.Mock
}

func (m *MockEventsRepository) Create(event Event) ([]byte, error) {
	return nil, nil
}

func (m *MockEventsRepository) FindById(id []byte) (Event, error) {
	args := m.Called(id)
	return args.Get(0).(Event), args.Error(1)
}

func (m *MockEventsRepository) FindByName(name string) (Event, error) {
	return *new(Event), nil
}

func (m *MockEventsRepository) FindAll() ([]Event, error) {
	return nil, nil
}
