-- name: GetCategoryByID :one
SELECT * FROM categories 
WHERE id = $1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name ASC;

-- name: CreateCategory :one
INSERT INTO categories (name, slug)
VALUES($1, $2)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories SET name = $1, slug = $2
WHERE id = $3
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;
