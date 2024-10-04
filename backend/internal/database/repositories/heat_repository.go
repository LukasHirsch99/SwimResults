package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type HeatRepositoryInterface interface {
	Create(*models.Heat) error
	Update(*models.Heat) error
	Delete(id int) error
	DeleteForEvent(eventid int) error
	CountForEvent(eventid int) (int, error)
}

type HeatRepository struct {
	db *sqlx.DB
}

func NewHeatRepository(db *sqlx.DB) *HeatRepository {
	return &HeatRepository{
		db: db,
	}
}

func (hr *HeatRepository) Create(heat *models.Heat) error {
	query := "INSERT INTO heat (eventid, heatnr) VALUES (:eventid, :heatnr) RETURNING id"
	r, err := hr.db.NamedQuery(query, heat)
	if err != nil {
		return err
	}

	for r.Next() {
		err := r.StructScan(&heat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hr *HeatRepository) Update(heat *models.Heat) error {
	query := "UPDATE heat SET eventid = :eventid, heatnr = :heatnr) WHERE id = :id"
	_, err := hr.db.NamedExec(query, heat)
	return err
}

func (hr *HeatRepository) Delete(id int) error {
	_, err := hr.db.Exec("DELETE FROM heat WHERE id = $1", id)
	return err
}

func (hr *HeatRepository) DeleteForEvent(eventid int) error {
	_, err := hr.db.Exec("DELETE FROM heat WHERE eventid = $1", eventid)
	return err
}

func (hr *HeatRepository) CountForEvent(eventid int) (int, error) {
	query := "SELECT COUNT(*) FROM heat WHERE eventid = $1"
	var cnt int
	err := hr.db.QueryRow(query, eventid).Scan(&cnt)
	return cnt, err
}
