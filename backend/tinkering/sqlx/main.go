package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Swimmer struct {
	Id        int
	Clubid    int
	Firstname string
	Lastname  string
	Birthyear sql.NullInt16
	Gender    string
	Isrelay   bool
}

type Club struct {
	Id          int
	Name        string
	Nationality sql.NullString
}

func main() {
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	dsn := "host=localhost user=admin password=admin dbname=swim-results port=5432 sslmode=disable TimeZone=Europe/Vienna"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

  duration, err := time.ParseDuration("15m")
	if err != nil {
    panic(err)
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
    panic(err)
	}

  defer db.Close()

	for i := range 500 {
		swimmer := Swimmer{
			Id:        i,
			Clubid:    1,
			Firstname: "Insert",
			Lastname:  "Test",
			Gender:    "M",
			Isrelay:   false,
		}
		fmt.Printf("Inserting %d\n", i)
    _, err = db.NamedExec(`INSERT INTO swimmer (id, clubid, firstname, lastname, gender, isrelay) VALUES (:id, :clubid, :firstname, :lastname, :gender, :isrelay)`, swimmer)
		if err != nil {
			panic(err)
		}
	}
}
