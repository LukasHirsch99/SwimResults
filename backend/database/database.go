package database

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"swimresults-backend/entities"
	sharedlist "swimresults-backend/sharedList"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

const API_URL = "https://qeudknoyuvjztxvgbmou.supabase.co"
const API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU"

type DBConnection struct {
	connection *pgx.Conn
	clientCnt  uint
}

var singleInstance *DBConnection
var lock = &sync.Mutex{}

func GetClient() *DBConnection {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			dbUrl := "postgres://admin:admin@localhost:5432/swim-results"
			var err error
			var c DBConnection
			c.connection, err = pgx.Connect(context.Background(), dbUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
				os.Exit(1)
			}

			singleInstance = &c
		} else {
			fmt.Println("Already initialized")
		}
	}
	singleInstance.clientCnt++
	return singleInstance
}

func (c DBConnection) CloseClient() {
	c.clientCnt--
	if c.clientCnt == 0 {
		singleInstance.connection.Close(context.Background())
	}
}

func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

// func executeAndParse[T any](f *postgrest.FilterBuilder) (T, int64, error) {
// 	var err error
// 	var r T
// 	data, cnt, err := f.Execute()
// 	if err != nil {
// 		return r, 0, err
// 	}
// 	err = json.Unmarshal(data, &r)
// 	return r, cnt, err
// }
//
// func (s Supabase) UpsertInto(table string, value interface{}) error {
// 	_, _, err := s.client.From(table).Upsert(value, "", "", "planned").Execute()
// 	return err
// }
//
// func (s Supabase) InsertInto(table string, value interface{}) error {
// 	_, _, err := s.client.From(table).Insert(value, false, "", "", "planned").Execute()
// 	return err
// }

// func (c DBConnection) Upsert(entity sharedlist.Entity) error {
//   fmt.Printf("Upserting %d %v\n", entity.GetItemCnt(), reflect.TypeOf(entity))
// 	_, _, err := s.client.From(entity.GetTableName()).Upsert(entity.GetItems(), "", "", "planned").Execute()
//   if err != nil {
//     panic(err)
//   }
// 	return err
// }

func (c DBConnection) BulkInsert(entity sharedlist.Entity) error {
	fmt.Printf("Inserting %d %v\n", entity.GetItemCnt(), reflect.TypeOf(entity))
	if entity.GetItemCnt() == 0 {
		return nil
	}
	fmt.Println(entity.GetRows())

	_, err := c.connection.CopyFrom(
		context.Background(),
		pgx.Identifier{entity.GetTableName()},
		entity.GetColumnNames(),
		pgx.CopyFromRows(entity.GetRows()),
	)

	if err != nil {
		panic(err)
	}
	return nil
}

func (c DBConnection) GetUpcomingMeets() []entities.Meet {
	rows, _ := c.connection.Query(context.Background(), "SELECT * FROM meet where enddate >= $1", time.Now())
	meets, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Meet])
	if err != nil {
		panic(err)
	}
	return meets
}

func (c DBConnection) GetMaxId(tablename string) sharedlist.MaxId {
	row, _ := c.connection.Query(context.Background(), fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", tablename))
	maxId, err := pgx.CollectExactlyOneRow(row, pgx.RowTo[uint])
	if err != nil {
		panic(err)
	}
	return sharedlist.MaxId{Id: maxId}
}

func (c DBConnection) GetIds(tablename string) []uint {
	rows, _ := c.connection.Query(context.Background(), fmt.Sprintf("SELECT id FROM %s", tablename))
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uint])
	if err != nil {
		panic(err)
	}
	return ids
}

// func (s Supabase) GetHeatsWithStartsByEventid(eventId uint) ([]entities.HeatWithStarts, int64, error) {
// 	return executeAndParse[[]entities.HeatWithStarts](s.client.
// 		From("heat").
// 		Select("*, start(*)", "exact", false).
// 		Eq("eventid", UintToString(eventId)))
// }
//
// func (s Supabase) GetSessionsWithEventsByMeetId(meetId uint) ([]entities.SessionWithEvents, int64, error) {
// 	return executeAndParse[[]entities.SessionWithEvents](s.client.
// 		From("session").
// 		Select("*, event(*)", "exact", false).
// 		Eq("meetid", UintToString(meetId)))
// }
//
// func (s Supabase) GetAgeclassCntByEventId(eventId uint) (int64, error) {
// 	_, dbResultCnt, err := s.client.From("ageclass").Select("*, result!inner(*)", "exact", false).Eq("result.eventid", UintToString(eventId)).Execute()
// 	return dbResultCnt, err
// }
//
// func (s Supabase) DeleteHeatsByEventId(eventId uint) {
// 	s.client.From("heat").Delete("*", "planned").Eq("eventid", UintToString(eventId)).Execute()
// }
//
// func (s Supabase) DeleteResultsByEventId(eventId uint) {
// 	s.client.From("result").Delete("*", "planned").Eq("eventid", UintToString(eventId)).Execute()
// }
//
// func (s Supabase) DeleteSessionsByMeetId(meetId uint) {
// 	s.client.From("session").Delete("*", "planned").Eq("meetid", UintToString(meetId)).Execute()
// }
//
// func (s Supabase) GetTodaysMeets() ([]entities.Meet, int64, error) {
// 	today := time.Now().Format("2006-01-02")
// 	return executeAndParse[[]entities.Meet](s.client.From("meet").Select("*", "planned", false).Gte("enddate", today).Lte("startdate", today))
// }
