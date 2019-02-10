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
	"strings"
	"time"
)

type eventServiceHandler struct {
	dbhandler    persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error: No search criteria found, you can either search by id via /id/4
		to search by name via /name/coldplayconcert}`)
		return
	}
	searchkey, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error: No search keys found, you can either search by id via /id/4
		to search by name via /name/coldplayconcert}`)
		return
	}
	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
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

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: error occured while decoding event data %s}", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: error occured while persisting event %d %s", id, err)
		return
	}
	msg := contracts.EventCreatedEvent{
		ID:   hex.EncodeToString(id),
		Name: event.Name,
		// LocationID: event.Location.ID,
		Start: time.Unix(event.StartDate, 0),
		End:   time.Unix(event.EndDate, 0),
	}
	eh.eventEmitter.Emit(&msg)
}

func newEventHandler(databaseHandler persistence.DatabaseHandler, emitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler:    databaseHandler,
		eventEmitter: emitter,
	}
}

func ServeAPI(httpHandler string, databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) chan error {
	handler := newEventHandler(databaseHandler, eventEmitter)
	router := mux.NewRouter()
	eventsRouter := router.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrChan := make(chan error)

	server := handlers.CORS()(router)

	go func() { httpErrChan <- http.ListenAndServe(httpHandler, server) }()
	return httpErrChan
}
