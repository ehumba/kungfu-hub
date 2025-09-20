-- +goose Up
CREATE TABLE martial_arts(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL UNIQUE);

-- +goose Down
DROP TABLE martial_arts;