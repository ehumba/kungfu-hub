-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;