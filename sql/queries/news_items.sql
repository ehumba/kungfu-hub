-- name: GetNewsByMartialArt :many
SELECT id, event_id, martial_art_id, title, content, published_at FROM news_items
WHERE martial_art_id = $1
ORDER BY published_at DESC;

-- name: GetAllNews :many
SELECT id, event_id, martial_art_id, title, content, published_at FROM news_items
ORDER BY published_at DESC;

-- name: InsertNewsItem :one
INSERT INTO news_items (event_id, martial_art_id, title, content, published_at, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, event_id, martial_art_id, title, content, published_at;