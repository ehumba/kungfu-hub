-- name: PostComment :one
INSERT INTO comments (id, news_item_id, user_id, parent_comment_id, content)
VALUES (DEFAULT, $1, $2, $3, $4)
RETURNING id, news_item_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: EditComment :one
UPDATE comments
SET content = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, news_item_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;

-- name: GetCommentsByPostID :many
SELECT id, news_item_id, user_id, parent_comment_id, content, created_at, updated_at FROM comments
WHERE news_item_id = $1
ORDER BY created_at ASC;

-- name: GetCommentsByUserID :many
SELECT id, news_item_id, user_id, parent_comment_id, content, created_at, updated_at FROM comments
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetCommentByID :one
SELECT id, news_item_id, user_id, parent_comment_id, content, created_at, updated_at FROM comments
WHERE id = $1;