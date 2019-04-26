package rest

import (
	"errors"
	"github.com/danielpacak/myevents-events-service/pkg/domain"
	"github.com/danielpacak/myevents-events-service/pkg/persistence/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerAPI(t *testing.T) {

	t.Run("GetById", func(t *testing.T) {

		t.Run("should return 200", func(t *testing.T) {
			// setup
			repository := new(mock.EventsRepository)
			router := NewAPIServer(repository, nil)
			// given
			event := domain.Event{}
			eventId := []byte{18, 52}
			repository.On("FindById", eventId).Return(event, nil)
			// and
			request := httptest.NewRequest("GET", "/events/id/1234", nil)
			response := httptest.NewRecorder()
			// when
			router.ServeHTTP(response, request)
			// then
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "application/json;charset=utf8", response.Header().Get("Content-Type"))
			// and
			expectedResponse := `{"ID":"","name":"","duration":0,"start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","location":{"ID":"","Name":"","Address":"","Country":"","OpenTime":0,"CloseTime":0,"Halls":null}}`
			assert.JSONEq(t, expectedResponse, response.Body.String())
			// finally
			repository.AssertExpectations(t)
		})

		t.Run("should return 400", func(t *testing.T) {
			// setup
			repository := new(mock.EventsRepository)
			// given
			handler := eventsHandler{
				repository: repository,
			}
			// and
			request := httptest.NewRequest("GET", "/id/1234", nil)
			response := httptest.NewRecorder()
			// and
			eventId := []byte{18, 52}
			// and
			repository.On("FindById", eventId).Return(domain.Event{}, errors.New("not found"))
			// when
			router := mux.NewRouter()
			router.HandleFunc("/id/{id}", handler.getById)
			router.ServeHTTP(response, request)
			// then
			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, "application/json;charset=utf8", response.Header().Get("Content-Type"))
			// and
			expectedResponse := `{"error": "event not found"}`
			assert.JSONEq(t, expectedResponse, response.Body.String())
			// finally
			repository.AssertExpectations(t)
		})

	})

	t.Run("GetAll", func(t *testing.T) {

		t.Run("should return 200", func(t *testing.T) {
			repository := new(mock.EventsRepository)
			repository.On("FindAll").Return([]domain.Event{
				{Name: "e1"},
				{Name: "e2"},
			}, nil)

			server := NewAPIServer(repository, nil)

			request := httptest.NewRequest("GET", "/events/", nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assert.Equal(t, http.StatusOK, response.Code)
			repository.AssertExpectations(t)
		})

		t.Run("should return 500 when accessing data fails", func(t *testing.T) {
			repository := new(mock.EventsRepository)
			repository.On("FindAll").
				Return(nil, errors.New("backend error"))

			server := NewAPIServer(repository, nil)

			request := httptest.NewRequest("GET", "/events/", nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.JSONEq(t, `{"error": "Internal server error", "statusCode": 500}`, response.Body.String())
			repository.AssertExpectations(t)
		})
	})
}
