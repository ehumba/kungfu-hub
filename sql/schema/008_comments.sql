-- +gooseUp
CREATE TABLE comments(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    news_item_id INT REFERENCES news_items(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +gooseDown
DROP TABLE comments;