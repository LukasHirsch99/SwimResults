package updatemeet

import (
	// "context"
	// "swimresults-backend/internal/config"
	// "swimresults-backend/internal/database"
	// "swimresults-backend/internal/repository"
	"strings"
	"swimresults-backend/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "github.com/joho/godotenv"
)

func TestExtractDate(t *testing.T) {
	examples := []string{
		"01.-05.08.2020",
		"03.10.2020",
		"29.02.-01.03.2020",
	}

	correct := []string{
		"2020-08-01", "2020-08-05",
		"2020-10-03", "2020-10-03",
		"2020-02-29", "2020-03-01",
	}

	var m = repository.Meet{}

	for i, example := range examples {
		err := extractMeetDate(strings.TrimSpace(example), &m)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		} else {
			assert.Equal(t, correct[2*i], m.Startdate.Time.Format(time.DateOnly))
			assert.Equal(t, correct[2*i+1], m.Enddate.Time.Format(time.DateOnly))
			// fmt.Printf("Start: %s, End: %s\n", m.Startdate.Time.Format("2006-01-02"), m.Enddate.Time.Format("2006-01-02"))
		}
	}
}

// func TestUpdateMeet(*testing.T) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		panic("Error loading .env file")
// 	}
//
// 	cfg := config.NewConfig()
//
// 	err = cfg.ParseFlags()
// 	if err != nil {
// 		panic("Failed to parse command-line flags")
// 	}
//
// 	ctx := context.Background()
// 	db, err := database.Connect(cfg, ctx)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	defer db.Close()
//
// 	repo := repository.New(db)
//
// 	UpdateMeet(2134, repo, nil, nil)
// }
