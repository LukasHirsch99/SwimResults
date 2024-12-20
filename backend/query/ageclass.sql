-- name: CreateAgeclass :one
INSERT INTO ageclass (name)
VALUES ($1)
RETURNING *;

-- name: GetAgeclassByName :one
SELECT * FROM ageclass WHERE name = $1;

-- name: DeleteAgeclassesByEvent :exec
DELETE FROM ageclass
  WHERE id in (SELECT ageclassid FROM ageclass_to_result WHERE eventid = $1);
