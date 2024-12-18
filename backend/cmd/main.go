package main

import (
	"log"
	"net/http"
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/database/repositories"
	updateschedule "swimresults-backend/updateSchedule"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var repos *repositories.Repositories

func getHeatsWithStartForEvent(c *gin.Context) {
	heats, _ := repos.HeatRepository.GetHeatsWithStartsForEvent(19898)

	c.IndentedJSON(http.StatusOK, heats)
}

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

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repos = cfg.InitializeRepositories(db)

	router := gin.Default()
	router.GET("/swimmers", getHeatsWithStartForEvent)

	done := make(chan bool)
	todaysMeetsTicker := time.NewTicker(5 * time.Minute)
	upcomingMeetsTicker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-done:
				upcomingMeetsTicker.Stop()
				todaysMeetsTicker.Stop()
        log.Println("Done")
				return
			case <-todaysMeetsTicker.C:
				meets, err := repos.MeetRepository.GetTodaysMeets()
				if err != nil {
					log.Println(err)
				}
				for _, m := range meets {
					updateschedule.UpdateSchedule(m.Id, repos)
				}
			case <-upcomingMeetsTicker.C:
				log.Println("Updating Upcoming Meets")
			}
		}
	}()

	if err := router.Run("localhost:8080"); err != nil {
		log.Println(err)
	}
  done <- true
}
