package main

import (
	"context"
	"fmt"
	"log"
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/repository"
	updateschedule "swimresults-backend/updateSchedule"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := config.NewConfig()

	err = cfg.ParseFlags()
	if err != nil {
		panic("Failed to parse command-line flags")
	}

	ctx := context.Background()
	db, err := database.Connect(cfg, ctx)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repo := repository.New(db)

	done := make(chan bool)
	todaysMeetsTicker := time.NewTicker(5 * time.Minute)
	upcomingMeetsTicker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-done:
				upcomingMeetsTicker.Stop()
				todaysMeetsTicker.Stop()
        fmt.Println("Stopping")
				return
			case <-todaysMeetsTicker.C:
				meets, err := repo.GetTodaysMeets(context.Background())
				if err != nil {
					log.Println(err)
				}
				for _, m := range meets {
					updateschedule.UpdateSchedule(m.ID, repo)
				}
			case <-upcomingMeetsTicker.C:
				fmt.Println("Updating Upcoming Meets")
			}
		}
	}()

  for ;; {}
  // time.Sleep(10 * time.Second)
  // done <- true
}
