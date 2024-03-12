package main

import "github.com/oapi-codegen/nullable"

var maxSessionId uint = 0
var maxEventId uint = 0

type Club struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

type Swimmer struct {
	Id        uint   `json:"id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	ClubId    uint   `json:"clubid"`
	BirthYear uint   `json:"birthyear"`
	Gender    string `json:"gender"`
	IsRelay   bool   `json:"isrelay"`
}

type Heat struct {
	Id      uint `json:"id"`
	EventId uint `json:"eventid"`
	HeatNr  uint `json:"heatnr"`
}

type HeatWithStarts struct {
	Heat
	Starts []Start `json:"start"`
}

type Start struct {
	HeatId    uint                      `json:"heatid"`
	SwimmerId uint                      `json:"swimmerid"`
	Lane      uint                      `json:"lane"`
	Time      nullable.Nullable[string] `json:"time"`
}

type StartWithHeat struct {
	Start
	Heat Heat `json:"heat"`
}

type SessionInfo struct {
	day          string
	displaynr    uint
	warmupstart  string
	sessionstart string
}

type EventInfo struct {
	displaynr uint
	name      string
}

type Session struct {
	Id           uint   `json:"id"`
	Meetid       uint   `json:"meetid"`
	Day          string `json:"day"`
	Displaynr    uint   `json:"displaynr"`
	Warmupstart  string `json:"warmupstart"`
	Sessionstart string `json:"sessionstart"`
}

type SessionWithEvents struct {
	Session
	Events []Event `json:"event"`
}

func (si SessionInfo) toSessionIncMaxId(meetId uint) Session {
	maxSessionId++
	return Session{
		Id:           maxSessionId,
		Meetid:       meetId,
		Day:          si.day,
		Displaynr:    si.displaynr,
		Warmupstart:  si.warmupstart,
		Sessionstart: si.sessionstart,
	}
}

type Event struct {
	Id        uint   `json:"id"`
	SessionId uint   `json:"sessionid"`
	DisplayNr uint   `json:"displaynr"`
	Name      string `json:"name"`
}

func (ei EventInfo) toEventIncMaxId(sessionId uint) Event {
	maxEventId++
	return Event{
		Id:        maxEventId,
		SessionId: sessionId,
		DisplayNr: ei.displaynr,
		Name:      ei.name,
	}
}

type EventWithSession struct {
	Event
	Session Session `json:"session"`
}

type Result struct {
	Id             uint                      `json:"id"`
	EventId        uint                      `json:"eventid"`
	SwimmerId      uint                      `json:"swimmerid"`
	Time           nullable.Nullable[string] `json:"time"`
	Splits         nullable.Nullable[string] `json:"splits"`
	FinaPoints     nullable.Nullable[uint]   `json:"finapoints"`
	AdditionalInfo nullable.Nullable[string] `json:"additionalinfo"`
	Penalty        bool                      `json:"penalty"`
}

type AgeClass struct {
	Name        string                    `json:"name"`
	ResultId    uint                      `json:"resultid"`
	Position    nullable.Nullable[uint]   `json:"position"`
	TimeToFirst nullable.Nullable[string] `json:"timetofirst"`
}
