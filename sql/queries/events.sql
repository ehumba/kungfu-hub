-- name: GetEventsByMartialArtID :many
SELECT id, martial_art_id, title, date, location, description FROM events
WHERE martial_art_id = $1 AND date >= NOW()
ORDER BY date ASC;

-- name: GetAllEvents :many
SELECT id, martial_art_id, title, date, location, description FROM events
WHERE date >= NOW()
ORDER BY date ASC;

