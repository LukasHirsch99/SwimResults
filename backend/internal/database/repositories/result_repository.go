package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type ResultRepositoryInterface interface {
	Create(*models.Result) error
	CreateMany([]models.Result) error
	Update(*models.Result) error
	Delete(id int) error
}

type ResultRepository struct {
	db *sqlx.DB
}

func NewResultRepository(db *sqlx.DB) *ResultRepository {
	return &ResultRepository{
		db: db,
	}
}

func (rr *ResultRepository) Create(result *models.Result) error {
	query := `INSERT INTO result (swimmerid, ageclassid, time, splits, finapoints, additionalinfo, penalty, reactiontime) 
  VALUES (:swimmerid, :ageclassid, :time, :splits, :finapoints, :additionalinfo, :penalty, :reactiontime) RETURNING id`
	_, err := rr.db.NamedExec(query, result)
	return err
}

func (rr *ResultRepository) CreateMany(results []models.Result) error {
	query := `INSERT INTO result (swimmerid, ageclassid, time, splits, finapoints, additionalinfo, penalty, reactiontime) VALUES
                               (:swimmerid, :ageclassid, :time, :splits, :finapoints, :additionalinfo, :penalty, :reactiontime)`
	_, err := rr.db.NamedExec(query, results)
	return err
}

func (rr *ResultRepository) Update(result *models.Result) error {
  query := "UPDATE result SET swimmerid = :swimmerid, time = :time, splits = :splits, finapoints = :finapoints, additionalinfo = :additionalinfo, penalty = :penalty, reactiontime = :reactiontime) WHERE id = :id RETURNING id"
	_, err := rr.db.NamedExec(query, result)
	return err
}

func (rr *ResultRepository) Delete(id int) error {
	query := "DELETE FROM result WHERE id = id"
	_, err := rr.db.Exec(query, id)
	return err
}

func (rr *ResultRepository) CountForEvent(eventid int) (int, error) {
	var cnt int
	query := "SELECT COUNT(*) FROM result r JOIN ageclass a on r.ageclassid = a.id WHERE a.eventid = $1"
	err := rr.db.Get(&cnt, query, eventid)
	return cnt, err
}
