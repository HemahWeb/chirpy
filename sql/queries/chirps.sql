-- name: CreateChirp :one
INSERT INTO chirps (body, user_id) VALUES ($1, $2) RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1 LIMIT 1;

-- name: UpdateChirp :one
UPDATE chirps SET body = $2 WHERE id = $1 RETURNING *;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

