package updatemeet

import (
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"testing"

	"github.com/joho/godotenv"
)

func TestUpdateMeet(*testing.T) {
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

	UpdateMeet(2134, repos, nil)
}
