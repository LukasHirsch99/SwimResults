package models

import "database/sql"

type Club struct {
	Id          int
	Name        string
	Nationality sql.NullString
}
