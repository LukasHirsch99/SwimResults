package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type StartRepositoryInterface interface {
	Create(*models.Start) error
	CreateMany(starts []models.Start) error
	Update(*models.Start) error
	Delete(id int) error
	CountForEvent(eventid int) (int, error)
}

type StartRepository struct {
	db *sqlx.DB
}

func NewStartRepository(db *sqlx.DB) *StartRepository {
	return &StartRepository{
		db: db,
	}
}

func (mr *StartRepository) Create(start *models.Start) error {
  query := "INSERT INTO start (heatid, swimmerid, lane, time) VALUES (:heatid, :swimmerid, :lane, :time)"
	_, err := mr.db.NamedExec(query, start)
	return err
}

func (mr *StartRepository) CreateMany(starts []models.Start) error {
	query := "INSERT INTO start (heatid, swimmerid, lane, time) VALUES (:heatid, :swimmerid, :lane, :time)"
	_, err := mr.db.NamedExec(query, starts)
	return err
}

func (mr *StartRepository) Update(start *models.Start) error {
	query := "UPDATE start SET lane = $1, time = $2) WHERE heatid = $3 AND swimmerid = $4"
	_, err := mr.db.Exec(query,
		start.Lane,
		start.Time,
		start.Heatid,
		start.Swimmerid)

	return err
}

func (mr *StartRepository) Delete(id int) error {
	query := "DELETE FROM start WHERE id = id"
	_, err := mr.db.Exec(query, id)
	return err
}

func (sr *StartRepository) CountForEvent(eventid int) (int, error) {
	var cnt int
	query := "SELECT COUNT(*) FROM start s JOIN heat h on s.heatid = h.id WHERE h.eventid = $1"
	err := sr.db.Get(&cnt, query, eventid)
	return cnt, err
}
