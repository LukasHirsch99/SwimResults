package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type AgeclassRepositoryInterface interface {
	Create(*models.Ageclass) error
	Update(*models.Ageclass) error
	Delete(id int) error
	DeleteForEvent(eventid int) error
	CountForEvent(eventid int) (int, error)
}

type AgeclassRepository struct {
	db *sqlx.DB
}

func NewAgeclassRepository(db *sqlx.DB) *AgeclassRepository {
	return &AgeclassRepository{
		db: db,
	}
}

func (ar *AgeclassRepository) Create(ageclass *models.Ageclass) error {
	r, err := ar.db.NamedQuery(`INSERT INTO ageclass (eventid, name)
        VALUES (:eventid, :name) RETURNING id`, ageclass)
	if err != nil {
		return err
	}

	for r.Next() {
		err := r.StructScan(&ageclass)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ar *AgeclassRepository) Update(ageclass *models.Ageclass) error {
	_, err := ar.db.NamedExec("UPDATE ageclass SET eventid = :eventid, name = :name) WHERE id = :id", ageclass)
	return err
}

func (ar *AgeclassRepository) Delete(id int) error {
	_, err := ar.db.Exec("DELETE FROM ageclass WHERE id = $1", id)
	return err
}

func (ar *AgeclassRepository) CountForEvent(eventid int) (int, error) {
	var cnt int
	query := "SELECT COUNT(*) FROM ageclass WHERE eventid = $1"
	err := ar.db.Get(&cnt, query, eventid)
	return cnt, err
}

func (ar *AgeclassRepository) DeleteForEvent(eventid int) error {
	_, err := ar.db.Exec("DELETE FROM ageclass WHERE eventid = $1", eventid)
	return err
}
