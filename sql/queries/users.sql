-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserData :exec
UPDATE users
SET username = $2,
email = $3,
password_hash = $4,
updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;