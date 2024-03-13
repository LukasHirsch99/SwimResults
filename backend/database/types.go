package database

import (
	"github.com/guregu/null/v5"
)

type Meet struct {
	Id             uint        `json:"id"`
	Name           string      `json:"name"`
	Image          null.String `json:"image"`
	Invitations    []string    `json:"invitations"`
	Deadline       string      `json:"deadline"`
	Address        string      `json:"address"`
	StartDate      string      `json:"startdate"`
	EndDate        string      `json:"enddate"`
	GoogleMapsLink null.String `json:"googlemapslink"`
}

type Session struct {
	Id           uint        `json:"id"`
	Meetid       uint        `json:"meetid"`
	Day          string      `json:"day"`
	Warmupstart  null.String `json:"warmupstart"`
	Sessionstart null.String `json:"sessionstart"`
	Displaynr    uint        `json:"displaynr"`
}

type SessionInfo struct {
	Day          string
	DisplayNr    uint
	WarmupStart  string
	SessionStart string
}

type SessionWithEvents struct {
	Session
	Events []Event `json:"event"`
}

func (si SessionInfo) ToSession(meetId uint, sessionId uint) Session {
	return Session{
		Id:        sessionId,
		Meetid:    meetId,
		Day:       si.Day,
		Displaynr: si.DisplayNr,
	}
}

type Event struct {
	Id        uint   `json:"id"`
	SessionId uint   `json:"sessionid"`
	DisplayNr uint   `json:"displaynr"`
	Name      string `json:"name"`
}

type EventInfo struct {
	DisplayNr uint
	Name      string
}

type EventWithSession struct {
	Event
	Session Session `json:"session"`
}

func (ei EventInfo) ToEvent(sessionId uint, eventId uint) Event {
	return Event{
		Id:        eventId,
		SessionId: sessionId,
		DisplayNr: ei.DisplayNr,
		Name:      ei.Name,
	}
}

type Club struct {
	Id          uint        `json:"id"`
	Name        string      `json:"name"`
	Nationality null.String `json:"nationality"`
}

type Swimmer struct {
	Id        uint     `json:"id"`
	Lastname  string   `json:"lastname"`
	Firstname string   `json:"firstname"`
	ClubId    uint     `json:"clubid"`
	BirthYear null.Int `json:"birthyear"`
	Gender    string   `json:"gender"`
	IsRelay   bool     `json:"isrelay"`
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
	HeatId    uint        `json:"heatid"`
	SwimmerId uint        `json:"swimmerid"`
	Lane      uint        `json:"lane"`
	Time      null.String `json:"time"`
}

type StartWithHeat struct {
	Start
	Heat Heat `json:"heat"`
}

type Result struct {
	Id             uint        `json:"id"`
	EventId        uint        `json:"eventid"`
	SwimmerId      uint        `json:"swimmerid"`
	Time           null.String `json:"time"`
	Splits         null.String `json:"splits"`
	FinaPoints     null.Int    `json:"finapoints"`
	AdditionalInfo null.String `json:"additionalinfo"`
	Penalty        bool        `json:"penalty"`
}

type AgeClass struct {
	Name        string      `json:"name"`
	ResultId    uint        `json:"resultid"`
	Position    null.Int    `json:"position"`
	TimeToFirst null.String `json:"timetofirst"`
}
