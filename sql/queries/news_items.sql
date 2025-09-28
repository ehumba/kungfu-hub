-- name: GetNewsByMartialArt :many
SELECT id, event_id, martial_art_id, title, content, published_at FROM news_items
WHERE martial_art_id = $1 AND published_at >= NOW() - INTERVAL '7 days'
ORDER BY published_at DESC;

-- name: GetNewsByEvent :many
SELECT id, event_id, martial_art_id, title, content, published_at FROM news_items
WHERE event_id = $1 AND published_at >= NOW() - INTERVAL '7 days'
ORDER BY published_at DESC;

-- name: GetAllNews :many
SELECT id, event_id, martial_art_id, title, content, published_at FROM news_items
WHERE published_at >= NOW() - INTERVAL '7 days'
ORDER BY published_at DESC;