package models

import (
	"github.com/thekb/zyzz/db"
	"github.com/jmoiron/sqlx"
	"fmt"
	"time"
)

const (
	CREATE_USER = `INSERT INTO user (short_id, name, description, email, nickname, avatarurl, fbid)
			VALUES (:short_id, :name, :description, :email, :nickname, :avatarurl, :fbid);`
	GET_USER_SHORT_ID = `SELECT A.* FROM user A
			WHERE A.short_id=$1;`
	GET_USER_ID = `SELECT A.* FROM user A
			WHERE A.id=$1;`
	GET_DEFAULT_USER = `SELECT A.* FROM user A
	 		ORDER BY A.id
	 		ASC LIMIT 1;`

)

type User struct {
	Id int `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ShortId string `db:"short_id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Published int `db:"published" json:"published"`
	Subscribed int `db:"subscribed" json:"subscribed"`
	Email string `db:"email" json:"email"`
	NickName string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatarurl" json:"avatarurl"`
	FBId string `db:"fbid" json:"fbid"`
}

func CreateUser(d *sqlx.DB, user *User) (int64, error) {
	id, err := db.InsertStruct(d, CREATE_USER, user)
	if err != nil {
		fmt.Println("unable to create user:", err)
	}
	return id, err
}

func GetUserForShortId(d *sqlx.DB, short_id string) (User, error) {
	var user User
	err := db.Get(d, GET_USER_SHORT_ID, &user, short_id)
	if err != nil {
		fmt.Println("unable to fetch user:", err)
	}
	return user, err
}

func GetUserForId(d *sqlx.DB, id int64) (User, error) {
	var user User
	err := db.Get(d, GET_USER_ID, &user, id)
	if err != nil {
		fmt.Println("unable to fetch user:", err)
	}
	return user, err
}

func GetDefaultUser(d *sqlx.DB) User {
	var user User
	err := db.Get(d, GET_DEFAULT_USER, &user)
	if err != nil {
		fmt.Println("unable to get default user:", err)
	}
	return user
}