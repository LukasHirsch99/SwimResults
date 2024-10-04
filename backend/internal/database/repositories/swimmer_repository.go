package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type SwimmerRepositoryInterface interface {
	GetIds() ([]int, error)
	Create(*models.Swimmer) error
	CreateMany(swimmer []models.Swimmer) error
	Update(*models.Swimmer) error
	Delete(id int) error
}

type SwimmerRepository struct {
	db *sqlx.DB
}

func NewSwimmerRepository(db *sqlx.DB) *SwimmerRepository {
	return &SwimmerRepository{
		db: db,
	}
}

func (sr *SwimmerRepository) Create(swimmer *models.Swimmer) error {
	_, err := sr.db.NamedExec(`INSERT INTO swimmer (id, clubid, firstname, lastname, birthyear, gender)
    VALUES (:id, :clubid, :firstname, :lastname, :birthyear, :gender) ON CONFLICT DO NOTHING;`, swimmer)
	return err
}

func (sr *SwimmerRepository) CreateMany(swimmer []models.Swimmer) error {
	_, err := sr.db.NamedExec(`INSERT INTO swimmer (id, clubid, firstname, lastname, birthyear, gender)
    VALUES (:id, :clubid, :firstname, :lastname, :birthyear, :gender) ON CONFLICT(id) DO NOTHING;`, swimmer)
	return err
}

func (sr *SwimmerRepository) Update(swimmer *models.Swimmer) error {
	query := "UPDATE swimmer SET clubid = :clubid, firstname = :firstname, lastname = :lastname, birthyear = :birthyear, gender = :gender) WHERE id = :id"
	_, err := sr.db.NamedExec(query, swimmer)

	return err
}

func (sr *SwimmerRepository) Delete(id int) error {
	query := "DELETE FROM swimmer WHERE id = id"
	_, err := sr.db.Exec(query, id)
	return err
}

func (sr *SwimmerRepository) CheckId(id int) bool {
	s := models.Swimmer{}
	err := sr.db.Get(&s, "SELECT * FROM swimmer WHERE id = $1", id)
	return err == nil
}

func (sr *SwimmerRepository) GetIds() ([]int, error) {
	ids := []int{}
	err := sr.db.Select(&ids, "SELECT id FROM swimmer")
	return ids, err
}
