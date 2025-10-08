-- Replace qr_image with url field in qris_codes table
-- qr_image is no longer needed as frontend generates QR code client-side
-- url stores Midtrans simulator URL for testing payment

BEGIN;

-- Add url column for Midtrans simulator URL
ALTER TABLE qris_codes ADD COLUMN IF NOT EXISTS url TEXT;

-- Drop qr_image column as it's deprecated (no longer used)
ALTER TABLE qris_codes DROP COLUMN IF EXISTS qr_image;

COMMIT;
