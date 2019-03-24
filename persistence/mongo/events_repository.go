package mongo

import (
	"github.com/danielpacak/myevents-events-service/config"
	"github.com/danielpacak/myevents-events-service/domain"
	"github.com/danielpacak/myevents-events-service/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	eventsCollection = "events"
)

type mongoEventsRepository struct {
	config  *config.MongoDBConfig
	session *mgo.Session
}

func NewMongoEventsRepository(config *config.MongoDBConfig) (persistence.EventsRepository, error) {
	log.Printf("Connecting to database at %s", config.ConnectionURL)
	session, err := mgo.Dial(config.ConnectionURL)
	if err != nil {
		return nil, err
	}
	return &mongoEventsRepository{
		config:  config,
		session: session,
	}, err
}

func (m *mongoEventsRepository) Create(e domain.Event) ([] byte, error) {
	s := m.getFreshSession()
	defer s.Close()
	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), m.withEventsCollection(s).Insert(e)
}

func (m *mongoEventsRepository) FindById(id []byte) (domain.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := domain.Event{}
	err := m.withEventsCollection(s).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (m *mongoEventsRepository) FindByName(name string) (domain.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := domain.Event{}
	err := m.withEventsCollection(s).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (m *mongoEventsRepository) FindAll() ([] domain.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	var events []domain.Event
	err := m.withEventsCollection(s).Find(nil).All(&events)
	return events, err
}

func (m *mongoEventsRepository) getFreshSession() *mgo.Session {
	return m.session.Copy()
}

func (m *mongoEventsRepository) withEventsCollection(s *mgo.Session) *mgo.Collection {
	return s.DB(m.config.DatabaseName).C(eventsCollection)
}
