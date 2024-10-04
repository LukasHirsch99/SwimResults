package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type ClubRepositoryInterface interface {
	Create(*models.Club) error
	Update(*models.Club) error
	Delete(id int) error
	CheckId(id int) bool
	GetIds() ([]int, error)
}

type ClubRepository struct {
	db *sqlx.DB
}

func NewClubRepository(db *sqlx.DB) *ClubRepository {
	return &ClubRepository{
		db: db,
	}
}

func (cr *ClubRepository) Create(club *models.Club) error {
	_, err := cr.db.NamedExec(`INSERT INTO club (id, name, nationality)
                              VALUES (:id, :name, :nationality) ON CONFLICT(id) DO NOTHING;`, club)
	return err
}

func (cr *ClubRepository) Update(club *models.Club) error {
	_, err := cr.db.NamedExec("UPDATE club SET name = :name, nationality = :nationality) WHERE id = :id", club)
	return err
}

func (cr *ClubRepository) Delete(id int) error {
	_, err := cr.db.Exec("DELETE FROM club WHERE id = $1", id)
	return err
}

func (cr *ClubRepository) CheckId(id int) bool {
	c := models.Club{}
	err := cr.db.Get(&c, "SELECT * FROM club WHERE id = $1", id)
	return err == nil
}

func (cr *ClubRepository) GetIds() ([]int, error) {
	ids := []int{}
	err := cr.db.Select(&ids, "SELECT id FROM club")
	return ids, err
}
