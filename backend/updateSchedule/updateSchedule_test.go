package updateschedule

import (
	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"testing"

	"github.com/joho/godotenv"
)

func TestUpdateSchedule(*testing.T) {
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

	repos := cfg.InitializeRepositories(db)

	UpdateSchedule(2134, repos)
	// for i := range 500 {
	// 	swimmer := models.Swimmer{
	// 		Id:        i,
	// 		Clubid:    1,
	// 		Firstname: "Insert",
	// 		Lastname:  "Test",
	// 		Gender:    "M",
	// 		Isrelay:   false,
	// 	}
	// 	fmt.Printf("Inserting %d\n", i)
	// 	// _, err = db.NamedExec(`INSERT INTO swimmer (id, clubid, firstname, lastname, gender, isrelay) VALUES (:id, :clubid, :firstname, :lastname, :gender, :isrelay)`, swimmer)
	// 	err = repos.SwimmerRepository.Create(&swimmer)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}
