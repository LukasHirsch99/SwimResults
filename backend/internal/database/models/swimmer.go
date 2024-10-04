package models

import "database/sql"

type Swimmer struct {
	Id        int
	Clubid    int
	Firstname string
	Lastname  string
	Birthyear sql.NullInt16
	Gender    string
	Isrelay   bool
}
