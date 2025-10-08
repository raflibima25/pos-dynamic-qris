-- Rollback: Restore qr_image and remove url

BEGIN;

-- Add back qr_image column
ALTER TABLE qris_codes ADD COLUMN IF NOT EXISTS qr_image TEXT;

-- Remove url column
ALTER TABLE qris_codes DROP COLUMN IF EXISTS url;

COMMIT;
