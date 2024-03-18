package main

import (
	"swimresults-backend/database"
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
		updateschedule.UpdateSchedule(m.Id, &wg)
	}
	wg.Wait()
}
