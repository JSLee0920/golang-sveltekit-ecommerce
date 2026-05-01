-- name: GetBrandByID :one
SELECT * FROM brands
WHERE id = $1;

-- name: GetBrandBySlug :one
SELECT * FROM brands
WHERE slug = $1;

-- name: ListBrands :many
SELECT * FROM brands
ORDER BY name ASC;

-- name: ListBrandsBySupplier :many
SELECT * FROM brands
WHERE supplier_id = $1
ORDER BY name ASC;

-- name: CreateBrand :one
INSERT INTO brands (supplier_id, name, slug, logo_url, description)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateBrand :one
UPDATE brands
SET name = $1, slug = $2, logo_url = $3, description = $4, updated_at = NOW()
WHERE id = $5
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brands
WHERE id = $1;
