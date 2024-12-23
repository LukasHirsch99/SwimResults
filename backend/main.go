package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"swimresults-backend/api"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/repository"
	updateschedule "swimresults-backend/updateSchedule"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.Connect(ctx)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	repo := repository.New(db)

	api := api.New(repo, logger)

	todaysMeetsTicker := time.NewTicker(5 * time.Minute)
	upcomingMeetsTicker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-ctx.Done():
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

	if err := api.Start(ctx); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
	}
}
