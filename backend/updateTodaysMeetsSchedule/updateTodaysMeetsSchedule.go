package main

import (
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	updateschedule "swimresults-backend/updateSchedule"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig()

	err = cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repos := cfg.InitializeRepositories(db)

	meets, err := repos.MeetRepository.GetTodaysMeets()
	if err != nil {
		panic(err)
	}

	// var wg sync.WaitGroup
	for _, m := range meets {
		// wg.Add(1)
		// go updateschedule.UpdateSchedule(m.Id, repos)
		updateschedule.UpdateSchedule(m.Id, repos)
	}
	// wg.Wait()
}
