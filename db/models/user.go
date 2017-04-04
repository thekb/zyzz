package models

import (
	"github.com/thekb/zyzz/db"
	"github.com/jmoiron/sqlx"
	"fmt"
	"time"
)

const (
	CREATE_USER = `INSERT INTO users (short_id, name, description, email, nickname, avatarurl, fbid, access_token)
			VALUES (:short_id, :name, :description, :email, :nickname, :avatarurl, :fbid, :access_token)
			returning id;`
	UPDATE_USER = `UPDATE users set name=:name, description=:description, email=:email, nickname=:nickname,
	 		avatarurl=:avatarurl, access_token=:access_token, username=:username, language=:language
	 		WHERE users.short_id=:short_id;`
	GET_USER_SHORT_ID = `SELECT A.* FROM users A
			WHERE A.short_id=$1;`
	GET_USER_USER_NAME = `SELECT A.* FROM users A
			WHERE A.username=$1;`
	GET_USER_ID = `SELECT A.* FROM users A
			WHERE A.id=$1;`
	GET_DEFAULT_USER = `SELECT A.* FROM users A
	 		ORDER BY A.id
	 		ASC LIMIT 1;`
	GET_USER_FBID = `SELECT A.* FROM users A WHERE A.fbid=$1`
)

type User struct {
	Id          int `db:"id" json:"-"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ShortId     string `db:"short_id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Published   int `db:"published" json:"published"`
	Subscribed  int `db:"subscribed" json:"subscribed"`
	Email       string `db:"email" json:"_"`
	NickName    string `db:"nickname" json:"nickname"`
	AvatarURL   string `db:"avatarurl" json:"avatarurl"`
	FBId        string `db:"fbid" json:"fbid"`
	AccessToken string `db:"access_token" json:"_"`
	UserName    string `db:"username" json:"username"`
	Language    string `db:"language" json:"language"`
}

func CreateUser(d *sqlx.DB, user *User) (int64, error) {
	id, err := db.InsertStruct(d, CREATE_USER, user)
	if err != nil {
		fmt.Println("unable to create user:", err)
	}
	return id, err
}

func UpdateUser(d *sqlx.DB, user *User) (error) {
	err := db.UpdateObj(d, UPDATE_USER, user)
	if err != nil {
		fmt.Println("unable to update user:", err)
	}
	return err
}

func GetUserForShortId(d *sqlx.DB, short_id string) (User, error) {
	var user User
	err := db.Get(d, GET_USER_SHORT_ID, &user, short_id)
	if err != nil {
		fmt.Println("unable to fetch user:", err)
	}
	return user, err
}

func GetUserForName(d *sqlx.DB, user_name string) (User, error) {
	var user User
	err := db.Get(d, GET_USER_USER_NAME, &user, user_name)
	if err != nil {
		fmt.Println("unable to fetch user:", err)
	}
	return user, err
}


func GetUserForFBId(d *sqlx.DB, fbid string) (User, error) {
	var user User
	err := db.Get(d, GET_USER_FBID, &user, fbid)
	if err != nil {
		fmt.Println("unable to fetch user:", err)
	}
	fmt.Println("got user")
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