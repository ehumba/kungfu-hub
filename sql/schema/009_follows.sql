-- +gooseUp
CREATE TABLE follows(
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    followed_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id, followed_id)
);

-- +gooseDown
DROP TABLE follows;