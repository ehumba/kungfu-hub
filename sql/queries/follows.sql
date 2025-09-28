-- name: CreateFollow :one
WITH p AS(
    SELECT
    $1::UUID AS follower_id,
    $2::UUID AS followed_id
)
INSERT INTO follows (follower_id, followed_id)
SELECT follower_id, followed_id FROM p
ON CONFLICT (follower_id, followed_id) DO NOTHING
RETURNING follower_id, followed_id;

-- name: DeleteFollow :exec
DELETE FROM follows WHERE follower_id = $1 AND followed_id = $2;