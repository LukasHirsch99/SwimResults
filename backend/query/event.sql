-- name: GetEventCntForMeet :one
SELECT count(*) FROM event
JOIN session ON event.sessionid = session.id AND session.meetid = $1;

-- name: CreateEvent :one
INSERT INTO event (sessionid, displaynr, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEventByPk :one
SELECT * FROM event
WHERE sessionid = $1 AND displaynr = $2 AND name = $3;
