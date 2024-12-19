// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: event.sql

package repository

import (
	"context"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO event (sessionid, displaynr, name)
VALUES ($1, $2, $3)
RETURNING id, sessionid, displaynr, name
`

type CreateEventParams struct {
	Sessionid int32
	Displaynr int32
	Name      string
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createEvent, arg.Sessionid, arg.Displaynr, arg.Name)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Sessionid,
		&i.Displaynr,
		&i.Name,
	)
	return i, err
}

const getEventByPk = `-- name: GetEventByPk :one
SELECT id, sessionid, displaynr, name FROM event
WHERE sessionid = $1 AND displaynr = $2 AND name = $3
`

type GetEventByPkParams struct {
	Sessionid int32
	Displaynr int32
	Name      string
}

func (q *Queries) GetEventByPk(ctx context.Context, arg GetEventByPkParams) (Event, error) {
	row := q.db.QueryRow(ctx, getEventByPk, arg.Sessionid, arg.Displaynr, arg.Name)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Sessionid,
		&i.Displaynr,
		&i.Name,
	)
	return i, err
}

const getEventCntForMeet = `-- name: GetEventCntForMeet :one
SELECT count(*) FROM event
JOIN session ON event.sessionid = session.id AND session.meetid = $1
`

func (q *Queries) GetEventCntForMeet(ctx context.Context, meetid int32) (int64, error) {
	row := q.db.QueryRow(ctx, getEventCntForMeet, meetid)
	var count int64
	err := row.Scan(&count)
	return count, err
}