package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/danielpacak/myevents-contracts"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/danielpacak/myevents-events-service/domain"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type eventsHandler struct {
	repository persistence.EventsRepository
	emitter    msgqueue.EventEmitter
}

type ResponseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"error"`
}

func (h *eventsHandler) getById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		_, _ = fmt.Fprint(w, `{"error"": "No id found in the path vars"}`)
		return
	}
	eventId, err := hex.DecodeString(id)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		_, _ = fmt.Fprint(w, `{"error"": "Internal server error"}`)
		return
	}
	event, err := h.repository.FindById(eventId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		_, _ = fmt.Fprint(w, `{"error": "event not found"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *eventsHandler) getByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		_, _ = fmt.Fprint(w, `{"error"": "No name found in the path vars"}`)
		return
	}
	event, _ := h.repository.FindByName(name)

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *eventsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.repository.FindAll()
	if err != nil {
		h.writeInternalServerErrorResponse(w)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		h.writeInternalServerErrorResponse(w)
		return
	}
}

func (h *eventsHandler) create(w http.ResponseWriter, r *http.Request) {
	event := domain.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "Error decoding event data", "error": "%s"}`, err)
		return
	}
	eventId, err := h.repository.Create(event)
	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "Error persisting event", "error": "%s"}`, err)
		return
	}
	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(eventId),
		Name:       event.Name,
		LocationID: event.Location.ID.String(),
		Start:      event.StartDate,
		End:        event.EndDate,
	}
	_ = h.emitter.Emit(&msg)
}

func (h *eventsHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(&ResponseError{
		StatusCode: statusCode,
		Message:    message,
	})
	if err != nil {
		panic(err)
	}
}

func (h *eventsHandler) writeInternalServerErrorResponse(w http.ResponseWriter) {
	h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}

func NewAPIServer(repository persistence.EventsRepository, emitter msgqueue.EventEmitter) http.Handler {
	router := mux.NewRouter()
	eventsRouter := router.PathPrefix("/events").Subrouter()

	eventsHandler := &eventsHandler{
		repository: repository,
		emitter:    emitter,
	}

	eventsRouter.Methods("GET").Path("/name/{name}").HandlerFunc(eventsHandler.getByName)
	eventsRouter.Methods("GET").Path("/id/{id}").HandlerFunc(eventsHandler.getById)
	eventsRouter.Methods("GET").Path("/").HandlerFunc(eventsHandler.getAll)
	eventsRouter.Methods("POST").Path("").HandlerFunc(eventsHandler.create)

	return handlers.CORS()(router)
}
