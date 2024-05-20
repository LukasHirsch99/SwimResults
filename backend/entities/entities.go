package entities

import (
	"strconv"
	"swimresults-backend/sharedList"

	"github.com/guregu/null/v5"
	"github.com/jackc/pgtype"
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
	StartDate pgtype.Date `json:"startdate"`
	EndDate   pgtype.Date `json:"enddate"`
}

type Meet struct {
	Id             uint                `json:"id" db:"id"`
	Name           string              `json:"name"`
	Image          pgtype.Varchar      `json:"image"`
	Invitations    pgtype.VarcharArray `json:"invitations"`
	Deadline       pgtype.Timestamp    `json:"deadline"`
	Address        string              `json:"address"`
	GoogleMapsLink pgtype.Varchar      `json:"googlemapslink"`
	MsecmId        pgtype.Int4         `json:"msecmid"`
	MeetDate
}

func (m Meet) GetId() uint {
	return m.Id
}

type Session struct {
	sharedlist.MaxId
	Meetid uint `json:"meetid"`
	SessionInfo
}

func (s *Session) SetId(id sharedlist.MaxId) {
	s.MaxId = id
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
	sharedlist.MaxId
	SessionId uint `json:"sessionid"`
	EventInfo
}

func (e *Event) SetId(id sharedlist.MaxId) {
	e.MaxId = id
}

type EventInfo struct {
	DisplayNr uint   `json:"displaynr"`
	Name      string `json:"name"`
}

type Club struct {
	Id          uint        `json:"id"`
	Name        string      `json:"name"`
	Nationality null.String `json:"nationality"`
}

func (c *Club) GetId() uint {
	return c.Id
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

func (s *Swimmer) GetId() uint {
	return s.Id
}

type Heat struct {
	sharedlist.MaxId
	EventId uint `json:"eventid"`
	HeatNr  uint `json:"heatnr"`
}

func (h *Heat) SetId(id sharedlist.MaxId) {
	h.MaxId = id
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

type Result struct {
	sharedlist.MaxId
	EventId        uint        `json:"eventid"`
	SwimmerId      uint        `json:"swimmerid"`
	Time           null.String `json:"time"`
	Splits         null.String `json:"splits"`
	FinaPoints     null.Int    `json:"finapoints"`
	AdditionalInfo null.String `json:"additionalinfo"`
	Penalty        bool        `json:"penalty"`
	ReactionTime   null.Float  `json:"reactiontime"`
}

func (r *Result) SetId(id sharedlist.MaxId) {
	r.MaxId = id
}

type AgeClass struct {
	Name        string      `json:"name"`
	ResultId    uint        `json:"resultid"`
	Position    null.Int    `json:"position"`
	TimeToFirst null.String `json:"timetofirst"`
}
