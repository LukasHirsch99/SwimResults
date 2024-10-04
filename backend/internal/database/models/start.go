package models

import (
	"database/sql"
)

type Start struct {
	Heatid    int
	Swimmerid int
	Lane      int
	Time      sql.NullTime
}
