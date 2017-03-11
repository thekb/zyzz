package models

import (
	"time"
	"github.com/thekb/zyzz/db"
	"github.com/jmoiron/sqlx"
	"fmt"
)

const (
	CREATE_EVENT = `INSERT INTO event (name, description, short_id, starttime, endtime, running_now, matchid, matchurl)
			VALUES (:name, :description, :short_id, :starttime, :endtime, :running_now, :matchid, :matchurl);`
	UPDATE_EVENT = `UPDATE event set name=:name, description=:description, starttime=:starttime,
			endtime=:endtime, running_now=:running_now, matchid=:matchid, matchurl=:matchurl
			WHERE event.short_id=:short_id`
	GET_EVENT_ID = `SELECT E.* FROM event E
			WHERE E.id=$1;`
	GET_EVENT_SHORT_ID = `SELECT E.* FROM event E
			WHERE E.short_id=$1;`
	GET_EVENTS = `SELECT A.* FROM event A
			WHERE A.running_now = 1
				ORDER By A.id ASC;`
)

type Event struct {
	Id          int64 `db:"id" json:"-"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ShortId     string `db:"short_id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	StartTime   time.Time `db:"starttime" json:"starttime"`
	EndTime     time.Time `db:"endtime" json:"endtime"`
	RunningNow  int `db:"running_now" json:"running_now"`
	MatchId     int `db:"matchid" json:"matchid"`
	MatchUrl    string `db:"matchurl" json:"matchurl"`
}


func CreateEvent(d *sqlx.DB, event *Event) (int64, error) {
	id, err := db.InsertStruct(d, CREATE_EVENT, event)
	if err != nil {
		fmt.Println("unable to create event:", err)
	}
	return id, err
}

func UpdateEvent(d *sqlx.DB, event *Event) (error) {
	err := db.UpdateObj(d, UPDATE_EVENT, event)
	if err != nil {
		fmt.Println("Unable to update event", err)
	}
	return err
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

func GetEvents(d *sqlx.DB) ([]Event, error) {
	var events []Event
	err := db.Select(d, GET_EVENTS, &events)
	if err != nil {
		fmt.Println("unable to get streams:", err)
	}
	return events, err
}

