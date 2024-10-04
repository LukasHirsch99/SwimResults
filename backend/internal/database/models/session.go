package models

import (
	"database/sql"
)

type Session struct {
	Id           int
	Meetid       int
	Displaynr    int
	Warmupstart  sql.NullTime
	Sessionstart sql.NullTime
}
