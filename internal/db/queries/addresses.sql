-- name: GetAddressByID :one
SELECT * FROM addresses
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAddressesByUser :many
SELECT * FROM addresses
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetDefaultAddress :one
SELECT * FROM addresses 
WHERE user_id = $1 AND is_default = TRUE AND deleted_at IS NULL;

-- name: CreatedAddress :one
INSERT INTO addresses (user_id, full_name, phone, address_line_1, address_line_2, city, state, postal_code, country, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: SetDefaultAddress :exec
UPDATE addresses SET is_default = FALSE
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: UpdateAddress :one
UPDATE addresses
SET full_name = $1, phone = $2, address_line_1 = $3, address_line_2 = $4,
    city = $5, state = $6, postal_code = $7, country = $8
WHERE id = $9 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAddress :exec
UPDATE addresses SET deleted_at = NOW()
WHERE id = $1;
