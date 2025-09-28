-- +gooseUp
CREATE TABLE events(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    martial_art_id INT NOT NULL REFERENCES martial_arts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    location TEXT NOT NULL,
    description TEXT
);

-- +gooseDown
DROP TABLE events;