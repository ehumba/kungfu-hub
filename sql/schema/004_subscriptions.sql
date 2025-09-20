-- +goose Up
CREATE TABLE subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    martial_art_id UUID NOT NULL REFERENCES martial_arts(id) ON DELETE CASCADE,
    UNIQUE(user_id, martial_art_id));

-- +goose Down
DROP TABLE subscriptions;