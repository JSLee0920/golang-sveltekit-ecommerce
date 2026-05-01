-- name: GetOrderDetailByID :one
SELECT * FROM order_details
WHERE id = $1;

-- name: ListOrderDetailsByOrder :many
SELECT * FROM order_details
WHERE order_id = $1;

-- name: ListOrderDetailsBySupplier :many
SELECT * FROM order_details
WHERE supplier_id = $1;

-- name: CreateOrderDetail :one
INSERT INTO order_details (order_id, product_id, supplier_id, quantity, price, subtotal)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
