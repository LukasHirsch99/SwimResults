package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Meet struct {
	Id             int
	Name           string
	Image          sql.NullString
	Invitations    pq.StringArray
	Deadline       sql.NullTime
	Address        string
	Startdate      time.Time
	Enddate        time.Time
	Googlemapslink sql.NullString
	Msecmid        sql.NullInt16
}
