-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOrdersByUser :one
SELECT count(*) FROM orders
WHERE user_id = $1;

-- name: ListOrdersByUser :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListAllOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllOrders :one
SELECT count(*) FROM orders;

-- name: ListOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateOrder :one
INSERT INTO orders (user_id, address_id, total)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $1, updated_at = now()
WHERE id = $2
RETURNING *;
