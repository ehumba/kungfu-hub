-- +gooseUp
CREATE TABLE news_items(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    martial_art_id INT NOT NULL REFERENCES martial_arts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    url TEXT NOT NULL
);

-- +gooseDown
DROP TABLE news_items;