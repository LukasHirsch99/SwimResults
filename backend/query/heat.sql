-- name: GetHeatCntForEvent :one
SELECT count(*) FROM heat
WHERE eventid = $1;

-- name: DeleteHeatsForEvent :exec
DELETE FROM heat
WHERE eventid = $1;

-- name: CreateHeat :one
INSERT INTO heat (eventid, heatnr)
VALUES ($1, $2) RETURNING id;
