-- name: GetSessionCntForMeet :one
SELECT count(*) FROM session WHERE meetid = $1;

-- name: DeleteSessionsForMeet :exec
DELETE FROM session WHERE meetid = $1;

-- name: GetSessionByPk :one
SELECT * FROM session WHERE meetid = $1 AND displaynr = $2;

-- name: CreateSession :one
INSERT INTO session (meetid, displaynr, day, warmupstart, sessionstart)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
