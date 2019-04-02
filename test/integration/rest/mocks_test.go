package rest

import (
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/stretchr/testify/mock"
)

type mockEventEmitter struct {
	mock.Mock
}

func (m *mockEventEmitter) Emit(event msgqueue.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
