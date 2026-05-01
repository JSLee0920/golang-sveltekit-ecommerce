-- name: GetSupplierByID :one
SELECT * FROM suppliers
WHERE id = $1;

-- name: GetSupplierByUserID: one
SELECT * FROM suppliers
WHERE user_id = $1;

-- name: ListSuppliers :many
SELECT * FROM suppliers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSuppliers :one
SELECT COUNT(*) FROM suppliers;

-- name: CreateSupplier :one
INSERT INTO suppliers (user_id, business_name, email, phone, address)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSupplierStatus :one
UPDATE suppliers SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateSupplier :one
UPDATE suppliers
SET business_name = $1, email = $2, phone = $3, address = $4, updated_at = NOW()
WHERE id = $5
RETURNING *;
