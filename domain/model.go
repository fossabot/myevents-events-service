package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string        `json:"name"`
	Duration  int           `json:"duration"`
	StartDate time.Time     `json:"start_date"`
	EndDate   time.Time     `json:"end_date"`
	Location  Location      `json:"location"`
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
