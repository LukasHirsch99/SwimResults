package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	AgeclassRepository *AgeclassRepository
	ClubRepository     *ClubRepository
	EventRepository    *EventRepository
	HeatRepository     *HeatRepository
	MeetRepository     *MeetRepository
	ResultRepository   *ResultRepository
	SessionRepository  *SessionRepository
	StartRepository    *StartRepository
	SwimmerRepository  *SwimmerRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		AgeclassRepository: NewAgeclassRepository(db),
		ClubRepository:     NewClubRepository(db),
		EventRepository:    NewEventRepository(db),
		HeatRepository:     NewHeatRepository(db),
		MeetRepository:     NewMeetRepository(db),
		ResultRepository:   NewResultRepository(db),
		SessionRepository:  NewSessionRepository(db),
		StartRepository:    NewStartRepository(db),
		SwimmerRepository:  NewSwimmerRepository(db),
	}
}
