-- name: GetPayoutByID :one
SELECT * FROM payouts
WHERE id = $1;

-- name: ListPayoutsBySupplier :many
SELECT * FROM payouts
WHERE supplier_id = $1
ORDER BY created_at DESC;

-- name: ListPendingPayouts :many
SELECT * FROM payouts
WHERE status = 'pending'
ORDER BY created_at DESC;

-- name: CreatePayout :one
INSERT INTO payouts (supplier_id, order_detail_id, amount, commission)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePayoutStatus :one
UPDATE payouts SET status = $1, paid_at = NOW()
WHERE id = $2
RETURNING *;
