-- name: GetClubIds :many
SELECT id FROM club;

-- name: GetClubs :many
SELECT * FROM club;

-- name: CreateClub :exec
INSERT INTO club (id, name, nationality)
VALUES ($1, $2, $3)
ON CONFLICT(id) DO NOTHING;

-- name: CheckClubId :one
SELECT CASE WHEN EXISTS (
    SELECT *
    FROM club
    WHERE id = $1
)
THEN true
ELSE false END;
