package repositories

import (
	"swimresults-backend/internal/database/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MeetRepositoryInterface interface {
	Get(int) (models.Meet, error)
	GetByMsecmId(int) (models.Meet, error)
	GetTodaysMeets() ([]models.Meet, error)
	Create(*models.Meet) error
	Update(*models.Meet) error
	Upsert(meet *models.Meet) error
	Delete(int) error
}

type MeetRepository struct {
	db *sqlx.DB
}

func NewMeetRepository(db *sqlx.DB) *MeetRepository {
	return &MeetRepository{
		db: db,
	}
}

func (mr *MeetRepository) Create(meet *models.Meet) error {
	query := `INSERT INTO meet (id, name, image, invitations, deadline, address, startdate, enddate, googlemapslink, msecmid) 
  VALUES (:id, :name, :image, :invitations, :deadline, :address, :startdate, :enddate, :googlemapslink, :msecmid)`
	_, err := mr.db.NamedExec(query, meet)
	return err
}

func (mr *MeetRepository) Update(meet *models.Meet) error {
	query := `UPDATE meet SET name = :name, image = :image, invitations = :invitations, deadline = :deadline, address = :address,
            startdate = :startdate, enddate = :enddate, googlemapslink = :googlemapslink, msecmid = :msecmid) WHERE id = :id`
	_, err := mr.db.Exec(query, meet)
	return err
}

func (mr *MeetRepository) Upsert(meet *models.Meet) error {
	query := `INSERT INTO meet (id, name, image, invitations, deadline, address, startdate, enddate, googlemapslink, msecmid) 
            VALUES (:id, :name, :image, :invitations, :deadline, :address, :startdate, :enddate, :googlemapslink, :msecmid)
            ON CONFLICT (id)
            DO UPDATE SET name = :name, image = :image, invitations = :invitations, deadline = :deadline, address = :address,
            startdate = :startdate, enddate = :enddate, googlemapslink = :googlemapslink, msecmid = :msecmid`
	_, err := mr.db.NamedExec(query, meet)
	return err
}

func (mr *MeetRepository) Delete(id int) error {
	query := "DELETE FROM meet WHERE id = id"
	_, err := mr.db.Exec(query, id)
	return err
}

func (mr *MeetRepository) GetById(meetid int) (*models.Meet, error) {
	meet := models.Meet{}
	err := mr.db.Get(&meet, "SELECT * FROM meet WHERE id = $1", meetid)
	return &meet, err
}

func (mr *MeetRepository) GetByMsecmId(msecmid int) (*models.Meet, error) {
	meet := models.Meet{}
	err := mr.db.Get(&meet, "SELECT * FROM meet WHERE msecmid = $1", msecmid)
	return &meet, err
}

func (mr *MeetRepository) GetTodaysMeets() ([]models.Meet, error) {
	meets := []models.Meet{}
	err := mr.db.Select(&meets, "SELECT * FROM meet WHERE startdate <= $1 AND enddate >= $1", time.Now())
  return meets, err
}
