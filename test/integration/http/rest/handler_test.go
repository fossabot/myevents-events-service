package rest

import (
	"encoding/hex"
	"github.com/danielpacak/myevents-events-service/pkg/domain"
	"github.com/danielpacak/myevents-events-service/pkg/http/rest"
	"github.com/danielpacak/myevents-events-service/pkg/persistence/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test")
	}

	t.Run("should handle GET /events/ request", func(t *testing.T) {
		var server *http.Server
		baseURL, err := url.Parse("http://localhost:8080/")
		require.NoError(t, err)

		repository := new(mock.EventsRepository)
		emitter := new(mockEventEmitter)

		go func() {
			handler := rest.NewHandler(repository, emitter)
			server = &http.Server{Addr: "localhost:8080", Handler: handler}
			_ = server.ListenAndServe()
		}()

		client := NewClient(baseURL)

		time.Sleep(3 * time.Second)

		e1 := domain.Event{Name: "e1", Duration: 3}
		e2 := domain.Event{Name: "e2", Duration: 4}
		repository.On("FindAll").Return([]domain.Event{e1, e2}, nil)

		events, err := client.ListEvents()
		require.NoError(t, err)

		assert.Equal(t, []domain.Event{e1, e2}, events)

		err = server.Close()
		require.NoError(t, err)
	})

	t.Run("should handle GET /events/{id} request", func(t *testing.T) {
		var server *http.Server
		baseURL, err := url.Parse("http://localhost:8080/")
		require.NoError(t, err)

		repository := new(mock.EventsRepository)
		emitter := new(mockEventEmitter)

		go func() {
			handler := rest.NewHandler(repository, emitter)
			server = &http.Server{Addr: "localhost:8080", Handler: handler}
			_ = server.ListenAndServe()
		}()

		client := NewClient(baseURL)

		time.Sleep(3 * time.Second)

		eventId, err := hex.DecodeString("abcd")
		require.NoError(t, err)

		e1 := domain.Event{Name: "e1"}

		repository.On("FindById", eventId).Return(e1, nil)

		event, err := client.GetById(eventId)
		require.NoError(t, err)

		assert.Equal(t, &e1, event)

		err = server.Close()
		require.NoError(t, err)
	})

	t.Run("should handle GET /events/{name} request", func(t *testing.T) {
		var server *http.Server
		baseURL, err := url.Parse("http://localhost:8080/")
		require.NoError(t, err)

		repository := new(mock.EventsRepository)
		emitter := new(mockEventEmitter)

		go func() {
			handler := rest.NewHandler(repository, emitter)
			server = &http.Server{Addr: "localhost:8080", Handler: handler}
			_ = server.ListenAndServe()
		}()

		client := NewClient(baseURL)

		time.Sleep(3 * time.Second)

		e1 := domain.Event{Name: "e1"}

		repository.On("FindByName", "peep show").Return(e1, nil)

		event, err := client.GetByName("peep show")
		require.NoError(t, err)

		assert.Equal(t, &e1, event)

		err = server.Close()
		require.NoError(t, err)
	})
}
