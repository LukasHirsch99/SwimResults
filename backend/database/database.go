package database

import (
	"encoding/json"
	"reflect"
	"strconv"

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

func InitClient() (*Supabase, error) {
	var err error
  var s Supabase
	s.client, err = supabase.NewClient(API_URL, API_KEY, nil)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func (s Supabase) InsertInto(table string, value interface{}) error {
	sl, ok := takeArg(value, reflect.Slice)
	if ok && sl.Len() == 0 {
		return nil
	}
	_, _, err := s.client.From(table).Insert(value, false, "", "", "").Execute()
	return err
}

func (s Supabase) GetClubIdSet() ([]uint, error) {
	cSet, _, err := executeAndParse[[]map[string]uint](s.client.From("club").Select("id", "exact", false))
	var clubIdSet []uint
	if err != nil {
		return clubIdSet, err
	}
	for _, v := range cSet {
		clubIdSet = append(clubIdSet, v["id"])
	}
	return clubIdSet, err
}

func (s Supabase) GetSwimmerIdSet() ([]uint, error) {
	sSet, _, err := executeAndParse[[]map[string]uint](s.client.From("swimmer").Select("id", "exact", false))
	var swimmerIdSeet []uint
	if err != nil {
		return swimmerIdSeet, err
	}
	for _, v := range sSet {
		swimmerIdSeet = append(swimmerIdSeet, v["id"])
	}
	return swimmerIdSeet, err
}

func (s Supabase) GetMaxIds() (MaxIds, error) {
	var maxIdMap map[string]uint
	var maxIds MaxIds
	err := json.Unmarshal([]byte(s.client.Rpc("maxids", "exact", "")), &maxIds)
	if err != nil {
		return maxIds, err
	}
	return MaxIds{
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
	s.client.From("heat").Delete("*", "exact").Eq("eventid", UintToString(eventId)).Execute()
}

func (s Supabase) DeleteResultsByEventId(eventId uint) {
	s.client.From("result").Delete("*", "exact").Eq("eventid", UintToString(eventId)).Execute()
}
func (s Supabase) DeleteSessionsByMeetId(meetId uint) {
	s.client.From("session").Delete("*", "exact").Eq("meetid", UintToString(meetId)).Execute()
}
