package rest

import (
	"errors"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/danielpacak/myevents-events-service/persistence/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeAPI(t *testing.T) {

	t.Run("getById", func(t *testing.T) {

		t.Run("should return 200", func(t *testing.T) {
			// setup
			repository := new(mock.MockEventsRepository)
			// given
			handler := eventsHandler{
				repository: repository,
			}
			// and
			request := httptest.NewRequest("GET", "/id/1234", nil)
			response := httptest.NewRecorder()
			// and
			event := persistence.Event{}
			eventId := []byte{18, 52}
			// and
			repository.On("FindById", eventId).Return(event, nil)
			// when
			router := mux.NewRouter()
			router.HandleFunc("/id/{id}", handler.getById)
			router.ServeHTTP(response, request)
			// then
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "application/json;charset=utf8", response.Header().Get("Content-Type"))
			// and
			expectedResponse := `{"ID":"","Name":"","Duration":0,"StartDate":0,"EndDate":0,"Location":{"ID":"","Name":"","Address":"","Country":"","OpenTime":0,"CloseTime":0,"Halls":null}}`
			assert.JSONEq(t, expectedResponse, response.Body.String())
			// finally
			repository.AssertExpectations(t)
		})

		t.Run("should return 400", func(t *testing.T) {
			// setup
			repository := new(mock.MockEventsRepository)
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
			repository.On("FindById", eventId).Return(persistence.Event{}, errors.New("not found"))
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
}
