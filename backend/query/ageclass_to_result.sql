-- name: GetAgeclassCntByEvent :one
SELECT count(distinct ageclassid) FROM ageclass_to_result WHERE eventid = $1;

-- name: CreateAgeclassToResults :copyfrom
INSERT INTO ageclass_to_result (eventid, ageclassid, resultid) 
VALUES ($1, $2, $3);

-- name: DeleteAgeclass_to_Results_ByEvent :exec
DELETE FROM ageclass_to_result WHERE eventid = $1;
