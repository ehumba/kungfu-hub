-- name: FetchFeed :many
-- $1 = current user's UUID
(
  SELECT
    'news' AS item_type,
    n.id          AS item_id,
    n.title       AS title,
    n.content     AS content,
    n.published_at AS created_at
  FROM news_items n
  JOIN subscriptions s
    ON s.martial_art_id = n.martial_art_id
  WHERE s.user_id = $1
)
UNION ALL
(
  SELECT
    'comment' AS item_type,
    c.id       AS item_id,
    NULL       AS title,         -- comments don’t have titles
    c.content  AS content,
    c.created_at
  FROM comments c
  JOIN follows f
    ON f.followed_id = c.user_id
  WHERE f.follower_id = $1
)
ORDER BY created_at DESC
LIMIT 20;
