package repositories

import (
	"swimresults-backend/internal/database/models"

	"github.com/jmoiron/sqlx"
)

type SessionRepositoryInterface interface {
	GetByPK(*models.Session) error
	Create(*models.Session) error
	Update(*models.Session) error
	Delete(int) error
	DeleteForMeet(int) error
	CountForMeet() (int, error)
}

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) GetByPK(session *models.Session) error {
	return sr.db.Get(&session, "SELECT * FROM session WHERE meetid = $1 AND displaynr = $2", session.Meetid, session.Displaynr)
}

func (sr *SessionRepository) Create(session *models.Session) error {
	query := "INSERT INTO session (meetid, displaynr, warmupstart, sessionstart) VALUES ($1, $2, $3, $4) RETURNING id"
	return sr.db.QueryRow(query,
		session.Meetid,
		session.Displaynr,
		session.Warmupstart,
		session.Sessionstart).Scan(&session.Id)
}

func (sr *SessionRepository) Update(session *models.Session) error {
	query := "UPDATE session SET meetid = $1, displaynr = $2, warmupstart = $3, sessionstart = $4) WHERE id = $5"
	_, err := sr.db.Exec(query,
		session.Meetid,
		session.Displaynr,
		session.Warmupstart,
		session.Sessionstart,
		session.Id)

	return err
}

func (sr *SessionRepository) Delete(id int) error {
	query := "DELETE FROM session WHERE id = id"
	_, err := sr.db.Exec(query, id)
	return err
}

func (sr *SessionRepository) CountForMeet(meetid int) (int, error) {
	query := "SELECT COUNT(*) FROM session WHERE meetid = $1"
	var cnt int
	err := sr.db.QueryRow(query, meetid).Scan(&cnt)
	return cnt, err
}

func (sr *SessionRepository) DeleteForMeet(meetid int) error {
	query := "DELETE FROM session WHERE meetid = $1"
	_, err := sr.db.Exec(query, meetid)
	return err
}
