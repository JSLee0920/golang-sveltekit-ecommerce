-- name: GetPayoutByID :one
SELECT * FROM payouts
WHERE id = $1;

-- name: ListPayoutsBySupplier :many
SELECT * FROM payouts
WHERE supplier_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPayoutsBySupplier :one
SELECT count(*) FROM payouts
WHERE supplier_id = $1;

-- name: ListPendingPayouts :many
SELECT * FROM payouts
WHERE status = 'pending'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreatePayout :one
INSERT INTO payouts (supplier_id, order_detail_id, amount, commission)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePayoutStatus :one
UPDATE payouts SET status = $1, paid_at = now()
WHERE id = $2
RETURNING *;
