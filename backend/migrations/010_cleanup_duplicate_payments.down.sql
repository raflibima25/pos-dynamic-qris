-- This is a data cleanup migration, cannot be reversed
-- Down migration does nothing as we cannot restore deleted duplicates
-- This is intentional - duplicate data should not be restored
SELECT 1; -- No-op
