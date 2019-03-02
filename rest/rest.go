package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/danielpacak/myevents-contracts"
	"github.com/danielpacak/myevents-contracts/lib/msgqueue"
	"github.com/danielpacak/myevents-events-service/persistence"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type eventsHandler struct {
	repository persistence.EventsRepository
	emitter    msgqueue.EventEmitter
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
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while trying to find all available events %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while trying to encode events to JSON %s}", err)
	}
}

func (h *eventsHandler) create(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, `{"message": "Error decoding event data", "error": "%s"}`, err)
		return
	}
	eventId, err := h.repository.Create(event)
	if nil != err {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, `{"message": "Error persisting event", "error": "%s"}`, err)
		return
	}
	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(eventId),
		Name:       event.Name,
		LocationID: event.Location.ID.String(),
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
	}
	_ = h.emitter.Emit(&msg)
}

func newEventsHandler(repository persistence.EventsRepository, emitter msgqueue.EventEmitter) *eventsHandler {
	return &eventsHandler{
		repository: repository,
		emitter:    emitter,
	}
}

func NewAPIServer(repository persistence.EventsRepository, emitter msgqueue.EventEmitter) http.Handler {
	eventsHandler := newEventsHandler(repository, emitter)
	router := mux.NewRouter()
	eventsRouter := router.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("/name/{name}").HandlerFunc(eventsHandler.getByName)
	eventsRouter.Methods("GET").Path("/id/{id}").HandlerFunc(eventsHandler.getById)
	eventsRouter.Methods("GET").Path("/").HandlerFunc(eventsHandler.getAll)
	eventsRouter.Methods("POST").Path("").HandlerFunc(eventsHandler.create)

	return handlers.CORS()(router)
}
