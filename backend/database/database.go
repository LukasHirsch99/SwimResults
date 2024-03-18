package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"swimresults-backend/entities"
	sharedlist "swimresults-backend/sharedList"
	"sync"
	"time"

	"github.com/supabase-community/supabase-go"
	"github.com/supabase/postgrest-go"
)

const API_URL = "https://qeudknoyuvjztxvgbmou.supabase.co"
const API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU"

type Supabase struct {
	client *supabase.Client
}

var singleInstance *Supabase
var lock = &sync.Mutex{}

func GetClient() (*Supabase, error) {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			var err error
			var s Supabase
			s.client, err = supabase.NewClient(API_URL, API_KEY, nil)
			if err != nil {
				return nil, err
			}
			singleInstance = &s
		} else {
			fmt.Println("Already initialized")
		}
	}
	return singleInstance, nil
}

func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func executeAndParse[T any](f *postgrest.FilterBuilder) (T, int64, error) {
	var err error
	var r T
	data, cnt, err := f.Execute()
	if err != nil {
		return r, 0, err
	}
	err = json.Unmarshal(data, &r)
	return r, cnt, err
}

func (s Supabase) UpsertInto(table string, value interface{}) error {
	_, _, err := s.client.From(table).Upsert(value, "", "", "planned").Execute()
	return err
}

func (s Supabase) InsertInto(table string, value interface{}) error {
	_, _, err := s.client.From(table).Insert(value, false, "", "", "planned").Execute()
	return err
}

func (s Supabase) Upsert(entity sharedlist.Entity) error {
  if entity.GetItemCnt() == 0 {
    return nil
  }
  fmt.Printf("Upserting %d %v\n", entity.GetItemCnt(), reflect.TypeOf(entity))
	_, _, err := s.client.From(entity.GetTableName()).Upsert(entity.GetItems(), "", "", "planned").Execute()
	return err
}

func (s Supabase) Insert(entity sharedlist.Entity) error {
  if entity.GetItemCnt() == 0 {
    return nil
  }
  fmt.Printf("Inserting %d %v\n", entity.GetItemCnt(), reflect.TypeOf(entity))
	_, _, err := s.client.From(entity.GetTableName()).Insert(entity.GetItems(), false, "", "", "planned").Execute()
	return err
}

func (s Supabase) GetUpcomingMeets() ([]*entities.Meet) {
	today := time.Now().Format("2006-01-02")
	meets, _, err := executeAndParse[[]*entities.Meet](s.client.From("meet").Select("*", "planned", false).Gte("startdate", today))
  if err != nil {
    panic(err)
  }
	return meets
}

func (s Supabase) GetIds(tablename string) []uint {
	idMap, _, err := executeAndParse[[]map[string]uint](s.client.From(tablename).Select("id", "planned", false))
	var ids []uint
	if err != nil {
		panic(err)
	}
	for _, v := range idMap {
		ids = append(ids, v["id"])
	}
	return ids
}

func (s Supabase) GetMaxId(tablename string) sharedlist.MaxId {
	var maxId uint
	err := json.Unmarshal([]byte(s.client.Rpc("maxid", "exact", map[string]string{"tablename": tablename})), &maxId)
	if err != nil {
		panic(err)
	}
	return sharedlist.MaxId{Id: maxId}
}

func (s Supabase) GetHeatsWithStartsByEventid(eventId uint) ([]entities.HeatWithStarts, int64, error) {
	return executeAndParse[[]entities.HeatWithStarts](s.client.
		From("heat").
		Select("*, start(*)", "exact", false).
		Eq("eventid", UintToString(eventId)))
}

func (s Supabase) GetSessionsWithEventsByMeetId(meetId uint) ([]entities.SessionWithEvents, int64, error) {
	return executeAndParse[[]entities.SessionWithEvents](s.client.
		From("session").
		Select("*, event(*)", "exact", false).
		Eq("meetid", UintToString(meetId)))
}

func (s Supabase) GetAgeclassCntByEventId(eventId uint) (int64, error) {
	_, dbResultCnt, err := s.client.From("ageclass").Select("*, result!inner(*)", "exact", false).Eq("result.eventid", UintToString(eventId)).Execute()
	return dbResultCnt, err
}

func (s Supabase) DeleteHeatsByEventId(eventId uint) {
	s.client.From("heat").Delete("*", "planned").Eq("eventid", UintToString(eventId)).Execute()
}

func (s Supabase) DeleteResultsByEventId(eventId uint) {
	s.client.From("result").Delete("*", "planned").Eq("eventid", UintToString(eventId)).Execute()
}

func (s Supabase) DeleteSessionsByMeetId(meetId uint) {
	s.client.From("session").Delete("*", "planned").Eq("meetid", UintToString(meetId)).Execute()
}

func (s Supabase) GetTodaysMeets() ([]entities.Meet, int64, error) {
	today := time.Now().Format("2006-01-02")
	return executeAndParse[[]entities.Meet](s.client.From("meet").Select("*", "planned", false).Gte("enddate", today).Lte("startdate", today))
}
