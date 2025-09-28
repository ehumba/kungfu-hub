-- name: FetchFeed :many
-- $1 = current user's UUID
SELECT n.id, n.event_id, n.martial_art_id, n.title, n.content, n.published_at
FROM news_items AS n
JOIN subscriptions AS s
  ON s.martial_art_id = n.martial_art_id
WHERE s.user_id = $1
ORDER BY n.published_at DESC
LIMIT 50;
