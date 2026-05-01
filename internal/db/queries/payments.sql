-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1;

-- name: CreatePayment :one
INSERT INTO payments (order_id, stripe_payment_indent, amount, currency, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdatePaymentStatus: one
UPDATE payments SET status = $1
WHERE order_id = $2
RETURNING *;
