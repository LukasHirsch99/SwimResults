package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type EventRepositoryInterface interface {
	GetByPK(*models.Event) error
	Create(*models.Event) error
	Update(*models.Event) error
	Delete(id int) error
}

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (er *EventRepository) GetByPK(event *models.Event) error {
	return er.db.Get(&event, "SELECT * FROM event WHERE sessionid = $1 AND displaynr = $2 AND name = $3", event.Sessionid, event.Displaynr, event.Name)
}

func (er *EventRepository) Create(event *models.Event) error {
	r, err := er.db.NamedQuery(`INSERT INTO event (sessionid, displaynr, name)
    VALUES (:sessionid, :displaynr, :name) RETURNING id`, event)
	if err != nil {
		return err
	}

	for r.Next() {
		err := r.StructScan(&event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (er *EventRepository) Update(event *models.Event) error {
	query := "UPDATE event SET sessionid = :sessionid, displaynr = :displaynr, name = :name) WHERE id = :id"
	_, err := er.db.NamedExec(query, event)
	return err
}

func (er *EventRepository) Delete(id int) error {
	_, err := er.db.Exec("DELETE FROM event WHERE id = $1", id)
	return err
}

func (er *EventRepository) CountForMeet(meetid int) (int, error) {
	var cnt int
	query := "SELECT COUNT(*) FROM event e JOIN session s on e.sessionid = s.id WHERE s.meetid = $1"
	err := er.db.Get(&cnt, query, meetid)
	return cnt, err
}
