-- name: CreateUser :one
INSERT INTO users (email, hashed_password) 
VALUES ($1, $2) 
RETURNING id, created_at, updated_at, email;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email 
FROM users WHERE id = $1 
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email 
FROM users WHERE email = $1 
LIMIT 1;

-- name: UpdateUserEmailAndPassword :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1 
RETURNING id, created_at, updated_at, email, hashed_password;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;

--auth-only
-- name: GetUserByEmailForAuth :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

