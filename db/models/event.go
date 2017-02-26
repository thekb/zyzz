package models

import (
	"time"
	"github.com/thekb/zyzz/db"
	"github.com/jmoiron/sqlx"
	"fmt"
)

const (
	CREATE_EVENT = `INSERT INTO event (name, description, short_id, starttime, endtime, running_now)
			VALUES (:name, :description, :short_id, :starttime, :endtime, :running_now);`
	GET_EVENT_ID = `SELECT E.* FROM event E
			WHERE E.id=$1;`
	GET_EVENT_SHORT_ID = `SELECT E.* FROM event E
			WHERE E.short_id=$1;`
)

type Event struct {
	Id          int `db:"id" json:"-"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ShortId     string `db:"short_id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	StartTime   time.Time `db:"starttime" json:"starttime"`
	EndTime     time.Time `db:"endtime" json:"endtime"`
	running_now bool `db:"running_now" json:"running_now"`
}

type Event_Stream struct {
	Id	    int `db:"id" json:"-"`
	ShortId     string `db:"short_id" json:"id"`
	EventId	    int `db:"event_id" json:"event_id"`
	StreamId    int `db:"stream_id" json:"stream_id"`
}

func CreateEvent(d *sqlx.DB, event *Event) (int64, error) {
	id, err := db.InsertStruct(d, CREATE_EVENT, event)
	if err != nil {
		fmt.Println("unable to create event:", err)
	}
	return id, err
}

func GetEventForId(d *sqlx.DB, id int64) (Event, error) {
	var event Event
	err := db.Get(d, GET_EVENT_ID, &event, id)
	if err != nil {
		fmt.Println("unable to fetch event:", err)
	}
	return event, err
}

func GetEventForShortId(d *sqlx.DB, short_id string) (Event, error) {
	var event Event
	err := db.Get(d, GET_EVENT_SHORT_ID, &event, short_id)
	if err != nil {
		fmt.Println("unable to fetch event:", err)
	}
	return event, err
}