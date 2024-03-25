package globalMutex

import (
	"github.com/guregu/null"
)

type AutoId struct {
	Id uint `json:"id"`
}

type Id struct {
	Id uint `json:"id"`
}

type EntityWithAutoId interface {
	SetId(AutoId)
}

type EntityWithId interface {
	GetId() Id
}

type Result struct {
	AutoId
	EventId uint   `json:"eventid"`
	Time    string `json:"time"`
}

func (r *Result) SetId(id AutoId) {
	r.AutoId = id
}

type Heat struct {
	AutoId
	EventId uint `json:"eventid"`
	HeatNr  uint `json:"heatnr"`
}

func (h *Heat) SetId(id AutoId) {
	h.AutoId = id
}

type Swimmer struct {
	Id
	Lastname  string   `json:"lastname"`
	Firstname string   `json:"firstname"`
	ClubId    uint     `json:"clubid"`
	BirthYear null.Int `json:"birthyear"`
	Gender    string   `json:"gender"`
	IsRelay   bool     `json:"isrelay"`
}

func (s *Swimmer) GetId() Id {
	return s.Id
}
