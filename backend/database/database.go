package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/supabase-community/supabase-go"
	"github.com/supabase/postgrest-go"
)

const API_URL = "https://qeudknoyuvjztxvgbmou.supabase.co"
const API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU"

type MaxIds struct {
	MaxSessionId uint
	MaxEventId   uint
	MaxHeatId    uint
	MaxResultId  uint
}

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

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func (s Supabase) UpsertInto(table string, value interface{}) error {
	sl, ok := takeArg(value, reflect.Slice)
	if ok && sl.Len() == 0 {
		return nil
	}
	_, _, err := s.client.From(table).Upsert(value, "", "", "").Execute()
	return err
}

func (s Supabase) InsertInto(table string, value interface{}) error {
	sl, ok := takeArg(value, reflect.Slice)
	if ok && sl.Len() == 0 {
		return nil
	}
	_, _, err := s.client.From(table).Insert(value, false, "", "", "").Execute()
	return err
}

func (s Supabase) GetClubIds() ([]uint, error) {
	cSet, _, err := executeAndParse[[]map[string]uint](s.client.From("club").Select("id", "planned", false))
	var clubIdSet []uint
	if err != nil {
		return clubIdSet, err
	}
	for _, v := range cSet {
		clubIdSet = append(clubIdSet, v["id"])
	}
	return clubIdSet, err
}

func (s Supabase) GetSwimmerIds() ([]uint, error) {
	sSet, _, err := executeAndParse[[]map[string]uint](s.client.From("swimmer").Select("id", "planned", false))
	var swimmerIdSeet []uint
	if err != nil {
		return swimmerIdSeet, err
	}
	for _, v := range sSet {
		swimmerIdSeet = append(swimmerIdSeet, v["id"])
	}
	return swimmerIdSeet, err
}

func (s Supabase) GetUpcomingMeets() ([]Meet, error) {
	today := time.Now().Format("2006-01-02")
	meets, _, err := executeAndParse[[]Meet](s.client.From("meet").Select("*", "planned", false).Gte("startdate", today))
	return meets, err
}

func (s Supabase) GetMaxIds() (*MaxIds, error) {
	var maxIdMap map[string]uint
	err := json.Unmarshal([]byte(s.client.Rpc("maxids", "exact", "")), &maxIdMap)
	if err != nil {
		return nil, err
	}
	return &MaxIds{
		maxIdMap["maxsessionid"],
		maxIdMap["maxeventid"],
		maxIdMap["maxheatid"],
		maxIdMap["maxresultid"],
	}, nil
}

func (s Supabase) GetHeatsWithStartsByEventid(eventId uint) ([]HeatWithStarts, int64, error) {
	return executeAndParse[[]HeatWithStarts](s.client.
		From("heat").
		Select("*, start!inner(*)", "exact", false).
		Eq("eventid", UintToString(eventId)))
}

func (s Supabase) GetSessionsWithEventsByMeetId(meetId uint) ([]SessionWithEvents, int64, error) {
	return executeAndParse[[]SessionWithEvents](s.client.
		From("session").
		Select("*, event!inner(*)", "exact", false).
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

func (s Supabase) GetTodaysMeets() ([]Meet, int64, error) {
	today := time.Now().Format("2006-01-02")
	return executeAndParse[[]Meet](s.client.From("meet").Select("*", "planned", false).Gte("enddate", today).Lte("startdate", today))
}
