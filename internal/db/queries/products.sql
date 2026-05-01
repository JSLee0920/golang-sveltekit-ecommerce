-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProductBySlug :one
SELECT * FROM products
WHERE slug = $1;

-- name: CountProducts :one
SELECT count(*) FROM products
WHERE active = true;

-- name: ListProducts :many
SELECT * FROM products 
WHERE active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProductsBySupplier :many
SELECT * FROM products
WHERE supplier_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountProductsBySupplier :one
SELECT count(*) FROM products
WHERE supplier_id = $1;

-- name: ListProductsByCategory :many
SELECT * FROM products
WHERE category_id = $1 AND active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProductsByBrand :many
SELECT * FROM products
WHERE brand_id = $1 AND active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateProduct :one
INSERT INTO products (supplier_id, category_id, brand_id, name, slug, description, price, stock, image_url)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $1, slug = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = now()
WHERE id = $7
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products SET stock = stock - $1, updated_at = now()
WHERE id = $2 AND stock >= $1
RETURNING *;

-- name: RestoreProductStock :exec
UPDATE products SET stock = stock + $1, updated_at = now()
WHERE id = $2;

-- name: SetProductActive :one
UPDATE products SET active = $1, updated_at = now()
WHERE id = $2
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
