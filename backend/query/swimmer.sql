-- name: GetAllSwimmers :many
SELECT * FROM swimmer WHERE lastname = $1;

-- name: GetClubWithSwimmers :many
SELECT sqlc.embed(c), sqlc.embed(s) FROM club c
JOIN swimmer s ON c.id = s.clubid;

-- name: CreateSwimmer :exec
INSERT INTO swimmer (id, clubid, firstname, lastname, birthyear, gender)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO NOTHING;

-- name: GetSwimmerIds :many
SELECT id FROM swimmer;

-- name: CheckSwimmerId :one
SELECT CASE WHEN EXISTS (
    SELECT *
    FROM swimmer
    WHERE id = $1
)
THEN true
ELSE false END;
