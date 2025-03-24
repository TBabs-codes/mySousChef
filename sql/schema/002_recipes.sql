-- +goose Up
CREATE TABLE recipes (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    ingredients TEXT NOT NULL,
    instructions TEXT,
    notes TEXT,
    img_link TEXT,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
