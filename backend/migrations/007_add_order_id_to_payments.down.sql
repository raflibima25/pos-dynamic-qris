-- Remove order_id column from payments table
DROP INDEX IF EXISTS idx_payments_order_id;
ALTER TABLE payments DROP COLUMN IF EXISTS order_id;
