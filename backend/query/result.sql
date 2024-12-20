-- name: CreateResults :copyfrom
INSERT INTO result
  (swimmerid, time, splits, finapoints, additionalinfo, penalty, reactiontime) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateResult :one
INSERT INTO result
  (swimmerid, time, splits, finapoints, additionalinfo, penalty, reactiontime) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: DeleteResultsByEvent :exec
DELETE FROM result
  WHERE id in (SELECT resultid FROM ageclass_to_result WHERE eventid = $1);
