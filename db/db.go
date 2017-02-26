package db

import (
	"github.com/rubenv/sql-migrate"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	DRIVER_NAME = "sqlite3"
	DB_NAME = "/opt/zyzz/db/zyzz.sqlite"
)

var sqlxDB *sqlx.DB

func GetDB() (*sqlx.DB, error) {
	if sqlxDB == nil {
		var err error
		err = RunMigrations()
		if err != nil {
			return nil, err
		}
		sqlxDB, err = sqlx.Connect(DRIVER_NAME, DB_NAME)
		if err != nil {
			return nil, err
		}
		//set open connections to 1 to prevent db errors when accessing multiple goroutines
		sqlxDB.SetMaxIdleConns(1)
		sqlxDB.SetMaxOpenConns(1)
		return sqlxDB, err
	}
	return sqlxDB, nil
}

// runs migrations
func RunMigrations() error {
	migrations := &migrate.AssetMigrationSource{
		Asset: Asset,
		AssetDir: AssetDir,
		Dir: "db/migrations",
	}

	db, err := sql.Open(DRIVER_NAME, DB_NAME)
	if err != nil {
		fmt.Println("unable to open database:", err)
		return err
	}
	defer db.Close()
	n, err := migrate.Exec(db, DRIVER_NAME, migrations, migrate.Up)
	if err != nil {
		fmt.Println("unable to run migrations:", err)
		return err
	}
	fmt.Println("applied migrations:", n)
	return nil
}

// db wrappers

//wrapper for sql select
func Select(db *sqlx.DB, query string, destination interface{}, args ...interface{}) error {
	err := db.Select(destination, query, args...)
	if err != nil {
		fmt.Println("unable to select:", err)
		return err
	}
	return nil
}

func Get(db *sqlx.DB, query string, destination interface{}, args ...interface{}) error {
	err := db.Get(destination, query, args...)
	if err != nil {
		fmt.Println("unable to get:", err)
		return err
	}
	return nil
}

//wrapper for count
func Count(db *sqlx.DB, query string, args ...interface{}) (int, error) {
	var count int
	err := db.Get(&count, query, args...)
	if err != nil {
		fmt.Println("unable to count:", err)
	}
	return count, nil
}

//wrapper for inserting multiple structs
func InsertStructs(db *sqlx.DB, query string, objects []interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("unable to begin transaction:", err)
		return err
	}
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("unable to prepare statement:", err)
		return err
	}
	for _, object := range objects {
		_, err := stmt.Exec(object)
		if err != nil {
			fmt.Println("unable to execute", query, object)
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("unable to commit transaction:", err)
		return err
	}
	return nil

}

//wrapper for insert struct
func InsertStruct(db *sqlx.DB, query string, object interface{}) (int64, error) {
	result, err := db.NamedExec(query, object)
	if err != nil {
		fmt.Println("unable to insert struct:", err)
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

//wrapper for update
func Update(db *sqlx.DB, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println("unable to update:", err)
	}
	return err
}

//wrapper for updating object
func UpdateObj(db *sqlx.DB, query string, object interface{}) error {
	_, err := db.NamedExec(query, object)
	if err != nil {
		fmt.Println("unable to update:", err)
	}
	return err
}