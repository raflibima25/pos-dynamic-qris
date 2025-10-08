-- Add unique constraint to prevent multiple pending payments for same transaction
-- This ensures only one active payment can exist per transaction
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_pending_payment_per_transaction 
ON payments (transaction_id) 
WHERE status IN ('pending', 'success') AND deleted_at IS NULL;
