package globalMutex

import (
	"encoding/json"
	"fmt"
	"swimresults-backend/database"
	"testing"

	"github.com/guregu/null/v5"
)

func TestSharedList(t *testing.T) {
	sl := CreateSharedList[database.Start]()
	sl.Add(database.Start{
		HeatId: 1,
		Lane:   2,
		Time:   null.NewString("27.01", true),
	})
	j, _ := json.Marshal(sl)
	fmt.Println(string(j))
}

func TestSharedListWithUniqueId(t *testing.T) {
	sl := CreateSharedListWithUniqueId()
	sl.Add(&Swimmer{
		Firstname: "Lukas",
		Lastname:  "Hirsch",
		Id:        Id{2},
	})
	sl.Add(&Swimmer{
		Firstname: "Lukas",
		Lastname:  "Hirsch",
		Id:        Id{2},
	})
	j, _ := json.Marshal(sl)
	fmt.Println(string(j))
}

func TestSharedListWithMaxId(t *testing.T) {
	sl := CreateSharedListWithMaxId()
	sl.Add(&Heat{
		EventId: 1,
		HeatNr:  2,
	})
	sl.Add(&Heat{
		EventId: 1,
		HeatNr:  2,
	})
	j, _ := json.Marshal(sl)
	fmt.Println(string(j))
}
