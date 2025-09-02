-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id) VALUES ($1)
RETURNING token;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP 
WHERE token = $1
RETURNING token;

-- name: GetUserIDFromRefreshToken :one
SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1;