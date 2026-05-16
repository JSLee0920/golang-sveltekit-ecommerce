-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('customer', 'supplier', 'admin'));

ALTER TABLE suppliers
ADD CONSTRAINT suppliers_role_check
CHECK (status IN ('pending', 'active', 'suspended'));

ALTER TABLE orders
ADD CONSTRAINT order_status_check
CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled'));

ALTER TABLE payments
ADD CONSTRAINT payment_status_check
CHECK (status IN ('pending', 'succeeded', 'failed'));

ALTER TABLE payouts
    ADD CONSTRAINT payouts_status_check
    CHECK (status IN ('pending', 'paid'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS suppliers_role_check;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS order_status_check;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payment_status_check;
ALTER TABLE payouts DROP CONSTRAINT IF EXISTS payouts_status_check;

-- +goose StatementEnd
