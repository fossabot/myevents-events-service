package mongolayer

import (
	"github.com/danielpacak/myevents-events-service/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	EVENTS = "events"
)

type MongoEventsRepository struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (*MongoEventsRepository, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoEventsRepository{
		session: s,
	}, err
}

func (m *MongoEventsRepository) getFreshSession() *mgo.Session {
	return m.session.Copy()
}

func (m *MongoEventsRepository) Create(e persistence.Event) ([] byte, error) {
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

func (m *MongoEventsRepository) FindById(id []byte) (persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (m *MongoEventsRepository) FindByName(name string) (persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (m *MongoEventsRepository) FindAll() ([] persistence.Event, error) {
	s := m.getFreshSession()
	defer s.Close()
	var events []persistence.Event
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
