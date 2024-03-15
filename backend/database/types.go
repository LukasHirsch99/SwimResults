package database

import (
	"strconv"

	"github.com/guregu/null/v5"
)

func StringToUint(s string) uint {
	u, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return uint(u)
}

func StringToInt(s string) int64 {
	u, _ := strconv.Atoi(s)
	return int64(u)
}

type MeetDate struct {
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

type Meet struct {
	Id             uint        `json:"id"`
	Name           string      `json:"name"`
	Image          null.String `json:"image"`
	Invitations    []string    `json:"invitations"`
	Deadline       string      `json:"deadline"`
	Address        string      `json:"address"`
	GoogleMapsLink null.String `json:"googlemapslink"`
	MsecmId        null.Int    `json:"msecmid"`
	MeetDate
}

type Session struct {
	Id     uint `json:"id"`
	Meetid uint `json:"meetid"`
	SessionInfo
}

type SessionInfo struct {
	Day          string      `json:"day"`
	WarmupStart  null.String `json:"warmupstart"`
	SessionStart null.String `json:"sessionstart"`
	DisplayNr    uint        `json:"displaynr"`
}

type SessionWithEvents struct {
	Session
	Events []Event `json:"event"`
}

type Event struct {
	Id        uint `json:"id"`
	SessionId uint `json:"sessionid"`
	EventInfo
}

type EventInfo struct {
	DisplayNr uint   `json:"displaynr"`
	Name      string `json:"name"`
}

type EventWithSession struct {
	Event
	Session Session `json:"session"`
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
