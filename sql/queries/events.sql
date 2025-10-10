-- name: GetEventsByMartialArtID :many
SELECT id, martial_art_id, title, date, location, description FROM events
WHERE martial_art_id = $1
ORDER BY date ASC;

-- name: GetAllEvents :many
SELECT id, martial_art_id, title, date, location, description FROM events
ORDER BY date ASC;

-- name: InsertEvent :one
INSERT INTO events (martial_art_id, title, date, location, description, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, martial_art_id, title, date, location, description;