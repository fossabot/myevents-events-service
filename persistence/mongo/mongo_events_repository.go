package mongo

import (
	"github.com/danielpacak/myevents-events-service/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	EVENTS = "events"
)

type mongoEventsRepository struct {
	session *mgo.Session
}

func NewMongoEventsRepository(connection string) (persistence.EventsRepository, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &mongoEventsRepository{
		session: s,
	}, err
}

func (m *mongoEventsRepository) getFreshSession() *mgo.Session {
	return m.session.Copy()
}

func (m *mongoEventsRepository) Create(e persistence.Event) ([] byte, error) {
	s := m.getFreshSession()
	defer s.Close()
	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}

func (m *mongoEventsRepository) FindById(id []byte) (persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (m *mongoEventsRepository) FindByName(name string) (persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (m *mongoEventsRepository) FindAll() ([] persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	var events []persistence.Event
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
