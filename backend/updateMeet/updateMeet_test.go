package updatemeet

import (
	"context"
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/repository"
	"testing"

	"github.com/joho/godotenv"
)

func TestUpdateMeet(*testing.T) {
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

	UpdateMeet(2134, repo, nil)
}
