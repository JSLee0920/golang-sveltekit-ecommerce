-- +goose Up
ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('customer', 'supplier', 'admin'));

ALTER TABLE suppliers
ADD CONSTRAINT suppliers_role_check
CHECK (status IN ('pending', 'active', 'suspended'));

-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS suppliers_role_check;
