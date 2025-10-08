-- Cleanup duplicate payments (one-time data migration)
-- This migration removes duplicate pending payments for the same transaction
-- Keeps the first payment (by ID) and removes subsequent duplicates

BEGIN;

-- Step 1: Delete QRIS codes associated with duplicate pending payments
DELETE FROM qris_codes 
WHERE payment_id IN (
  SELECT p2.id 
  FROM payments p1
  INNER JOIN payments p2 
    ON p1.transaction_id = p2.transaction_id 
    AND p1.id < p2.id  -- Keep first payment, delete later ones
  WHERE p2.status = 'pending'
    AND p1.status IN ('pending', 'success')
    AND p1.deleted_at IS NULL 
    AND p2.deleted_at IS NULL
);

-- Step 2: Delete duplicate pending payments
-- Keep the oldest payment record (lowest ID)
DELETE FROM payments p2
USING payments p1
WHERE p1.transaction_id = p2.transaction_id 
  AND p1.id < p2.id  -- p1 is older (lower ID)
  AND p2.status = 'pending'
  AND p1.deleted_at IS NULL 
  AND p2.deleted_at IS NULL;

COMMIT;
