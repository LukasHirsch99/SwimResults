-- name: GetStartCntForEvent :one
SELECT count(*) FROM start
JOIN heat on heat.id = start.heatid
WHERE heat.eventid = $1;

-- name: CreateStarts :copyfrom
INSERT INTO start (heatid, swimmerid, lane, time)
VALUES ($1, $2, $3, $4);
