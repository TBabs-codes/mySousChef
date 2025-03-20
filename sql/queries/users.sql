-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, premium, email, hashed_password)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: ReturnUserFromEmail :one
SELECT * FROM users WHERE email = ?;

-- name: ReturnUserFromID :one
SELECT * FROM users WHERE id = ?;

-- name: UpdateUserInfo :exec
UPDATE users
SET email = ?, hashed_password = ?, updated_at = ?
WHERE id = ?;

-- name: UpgradeUserMembership :exec
UPDATE users
SET premium = true
WHERE id = ?;

-- name: DowngradeUserMembership :exec
UPDATE users
SET premium = false
WHERE id = ?;