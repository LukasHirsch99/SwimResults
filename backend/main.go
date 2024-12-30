package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"swimresults-backend/api"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/repository"
	updateschedule "swimresults-backend/updateSchedule"
	updateupcomingmeets "swimresults-backend/updateUpcomingMeets"
	"time"

	"github.com/joho/godotenv"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.Connect(ctx, migrations)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	repo := repository.New(db)

	api := api.New(repo, logger)

	todaysMeetsTicker := time.NewTicker(5 * time.Minute)
	upcomingMeetsTicker := time.NewTicker(24 * time.Hour)

	go func() {
		updateupcomingmeets.UpdateUpcomingMeets(repo, logger)
		for {
			select {
			case <-ctx.Done():
				upcomingMeetsTicker.Stop()
				todaysMeetsTicker.Stop()
				logger.Info("Stopping")
				return
			case <-todaysMeetsTicker.C:
				meets, err := repo.GetTodaysMeets(context.Background())
				if err != nil {
					logger.Error("failed getting todays meets", slog.Any("error", err))
				}
				for _, m := range meets {
					updateschedule.UpdateSchedule(m.ID, repo)
				}
			case <-upcomingMeetsTicker.C:
				logger.Info("Updating Upcoming Meets")
				updateupcomingmeets.UpdateUpcomingMeets(repo, logger)
			}
		}
	}()

	if err := api.Start(ctx); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
	}
}
