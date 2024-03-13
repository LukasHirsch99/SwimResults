package database

import (
  "testing"
)

func TestZeroPackage(t *testing.T) {
  var c Supabase
  var data []Heat
  data = append(data, Heat{})
  c.InsertInto("heat", data)
}
