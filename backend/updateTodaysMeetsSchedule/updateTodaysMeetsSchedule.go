package main

import (
	"fmt"
	"swimresults-backend/database"
	updateschedule "swimresults-backend/updateSchedule"
	"sync"
	"time"
)

func main() {
  st := time.Now()
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
  fmt.Printf("Took: %v\n", time.Now().Sub(st))
}
