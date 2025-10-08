-- Add order_id column to payments table
ALTER TABLE payments ADD COLUMN IF NOT EXISTS order_id VARCHAR(255);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
