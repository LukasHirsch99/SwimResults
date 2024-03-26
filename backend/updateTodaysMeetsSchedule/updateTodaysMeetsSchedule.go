package main

import (
	"swimresults-backend/database"
	"swimresults-backend/store"
	updateschedule "swimresults-backend/updateSchedule"
	"sync"
)

func main() {
	supabase, err := database.GetClient()
	if err != nil {
		panic(err)
	}
	meets, _, err := supabase.GetTodaysMeets()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, m := range meets {
		wg.Add(1)
		go updateschedule.UpdateSchedule(m.Id, &wg)
	}
	wg.Wait()

	supabase.Insert(store.Clubs)
	supabase.Insert(store.Swimmers)
	supabase.Insert(store.Sessions)
	supabase.Insert(store.Events)
	supabase.Insert(store.Heats)
	supabase.Insert(store.Results)
	supabase.Insert(store.Starts)
	supabase.Insert(store.Ageclasses)
}
