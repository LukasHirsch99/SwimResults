package database

import (
	"swimresults-backend/entities"
	sharedlist "swimresults-backend/sharedList"
	"testing"
	"time"

	"github.com/jackc/pgtype"
)

func TestBulkInsert(t *testing.T) {
	c := GetClient()
	defer c.CloseClient()
	meet := entities.Meet{
		Name:     "Test Meet 2",
		Address:  "Test Address",
		Deadline: pgtype.Timestamp{Time: time.Now(), Status: pgtype.Present, InfinityModifier: pgtype.None},
		MeetDate: entities.MeetDate{
			StartDate: pgtype.Date{Time: time.Now(), Status: pgtype.Present, InfinityModifier: pgtype.None},
			EndDate:   pgtype.Date{Time: time.Now(), Status: pgtype.Present, InfinityModifier: pgtype.None},
		},
	}
  var meetList = make([]entities.Meet, 1)
  meetList[0] = meet
	var meets = sharedlist.GetSharedListWithUniqueIdWithItems("meet", meetList)
	c.BulkInsert(meets)
}
