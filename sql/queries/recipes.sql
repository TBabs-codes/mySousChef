-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, ingredients, instructions, notes, img_link, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllRecipes :exec
DELETE FROM recipes;

-- name: ReturnAllRecipes :many
SELECT * FROM recipes ORDER BY created_at;

-- name: ReturnAllRecipesByAuthor :many
SELECT * FROM recipes WHERE user_id = ? ORDER BY created_at;

-- name: ReturnAllRecipesDESC :many
SELECT * FROM recipes ORDER BY created_at DESC;

-- name: ReturnAllRecipesByAuthorDESC :many
SELECT * FROM recipes WHERE user_id = ? ORDER BY created_at DESC;

-- name: ReturnRecipe :one
SELECT * FROM recipes WHERE id = ?;

-- name: DeteleRecipe :exec
DELETE FROM recipes WHERE id = ?;
