package database

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestZeroPackage(t *testing.T) {
  var m Meet
  date := MeetDate{"1", "2"}
  m.StartDate = "B"
  m.MeetDate = date
  d, _ := json.Marshal(m)
  fmt.Println(string(d))
}
