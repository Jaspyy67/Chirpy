-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;


-- name: GetChirpsByID :one
SELECT id, created_at, updated_at, user_id, body
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: GetChirp :one
SELECT id, body, user_id, created_at, updated_at
FROM chirps
WHERE id = $1;