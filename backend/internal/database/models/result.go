package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Splits map[int]time.Time

func (a Splits) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Splits) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Result struct {
	Id             int
	Swimmerid      int
	Ageclassid     int
	Time           sql.NullTime
	Splits         Splits
	Finapoints     sql.NullInt16
	Additionalinfo sql.NullString
	Penalty        bool
	Reactiontime   sql.NullFloat64
}
