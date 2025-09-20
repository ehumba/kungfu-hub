-- name: AddSubscription :one
INSERT INTO subscriptions(user_id, martial_art_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveSubscription :exec
DELETE FROM subscriptions
WHERE user_id = $1 AND martial_art_id = $2;

-- name: GetUserSubscriptions :many
SELECT ma.* FROM martial_arts ma
JOIN subscriptions s ON ma.id = s.martial_art_id
WHERE s.user_id = $1;