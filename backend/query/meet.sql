-- name: GetMeetById :one
SELECT * FROM meet where id = $1;

-- name: GetMeetByMsecmId :one
SELECT * FROM meet where msecmid = $1;

-- name: UpsertMeet :exec
INSERT INTO meet (id, name, image, invitations, deadline, address, startdate, enddate, googlemapslink, msecmid) 
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  ON CONFLICT (id)
  DO UPDATE SET name = $2, image = $3, invitations = $4, deadline = $5, address = $6,
  startdate = $7, enddate = $8, googlemapslink = $9, msecmid = $10;

-- name: GetTodaysMeets :many
SELECT * FROM meet WHERE startdate <= now() AND enddate >= now();
