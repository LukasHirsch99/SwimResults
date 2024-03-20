package updateschedule

import (
	"swimresults-backend/database"
	"swimresults-backend/store"
	"sync"
	"testing"
)

func TestUpdateSchedule(*testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go UpdateSchedule(2123, &wg)
	wg.Wait()

	supabase, err := database.GetClient()
	if err != nil {
		panic(err)
	}
	var sessions = store.Sessions
	var events = store.Events
	var heats = store.Heats
	var results = store.Results
	var clubs = store.Clubs
	var swimmers = store.Swimmers
	var starts = store.Starts
	var ageclasses = store.Ageclasses

  err = supabase.Insert(swimmers)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(clubs)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(sessions)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(events)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(heats)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(starts)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(results)
  if err != nil {
    panic(err)
  }
  err = supabase.Insert(ageclasses)
  if err != nil {
    panic(err)
  }
}
