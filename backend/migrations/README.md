# Database Migrations

## Migration Files

Migration files are numbered sequentially and should be executed in order:

1. `001_*.sql` - Initial schema
2. `002_*.sql` - Add users table
3. `003_*.sql` - Add products table
4. `004_*.sql` - Add transactions table
5. `005_*.sql` - Add payments & QRIS tables
6. `006_*.sql` - **Add order_id column to payments**
7. `007_*.sql` - Down migration for order_id
8. `008_*.sql` - **Add unique constraint for pending payments**
9. `009_*.sql` - Down migration for unique constraint
10. `010_*.sql` - **Cleanup duplicate payments (one-time)**
11. `011_*.sql` - **Replace qr_image with url field**

## Running Migrations

### Manual Migration (PostgreSQL)

```bash
# Connect to database
psql -U postgres -d qris_pos

# Run specific migration
\i migrations/006_add_order_id_to_payments.up.sql
\i migrations/008_add_unique_pending_payment_constraint.up.sql
\i migrations/010_cleanup_duplicate_payments.up.sql
```

### Using golang-migrate CLI

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all pending migrations
migrate -path migrations -database "postgresql://postgres:password@localhost:5432/qris_pos?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgresql://postgres:password@localhost:5432/qris_pos?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgresql://postgres:password@localhost:5432/qris_pos?sslmode=disable" version
```

## Important Migrations (Latest)

### Migration 006: Add order_id Column

**Purpose**: Store Midtrans order_id in payment table for status checking

**What it does**:
- Adds `order_id` VARCHAR(255) column to `payments` table
- Creates index `idx_payments_order_id` for faster lookups

**Required**: Yes - Payment status checking will fail without this

### Migration 008: Unique Constraint

**Purpose**: Prevent duplicate pending/success payments for same transaction

**What it does**:
- Creates partial unique index on `payments(transaction_id)`
- Only applies to non-deleted pending/success payments
- Prevents race conditions when generating QRIS

**Required**: Yes - Prevents duplicate payment creation

### Migration 010: Cleanup Duplicates

**Purpose**: One-time cleanup of existing duplicate payments

**What it does**:
- Deletes duplicate QRIS codes for duplicate pending payments
- Deletes duplicate pending payments (keeps oldest)
- Safe to run multiple times (idempotent)

**Required**: Recommended - Cleans up existing bad data

**Note**: This is a data migration, cannot be rolled back

## Migration Order for Fresh Database

```sql
-- 1. Create schema (run initial migrations 001-005)
\i migrations/001_*.up.sql
\i migrations/002_*.up.sql
\i migrations/003_*.up.sql
\i migrations/004_*.up.sql
\i migrations/005_*.up.sql

-- 2. Add order_id column
\i migrations/006_add_order_id_to_payments.up.sql

-- 3. Add unique constraint
\i migrations/008_add_unique_pending_payment_constraint.up.sql

-- 4. Cleanup duplicates (if any)
\i migrations/010_cleanup_duplicate_payments.up.sql
```

## For Existing Database (Already Running)

If you already have a running database with data:

```bash
# 1. Backup first!
pg_dump -U postgres qris_pos > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Connect to database
psql -U postgres -d qris_pos

# 3. Add order_id column (if not exists)
\i migrations/006_add_order_id_to_payments.up.sql

# 4. Cleanup duplicates BEFORE adding constraint
\i migrations/010_cleanup_duplicate_payments.up.sql

# 5. Add unique constraint
\i migrations/008_add_unique_pending_payment_constraint.up.sql

# 6. Verify
SELECT
  transaction_id,
  COUNT(*) as payment_count
FROM payments
WHERE deleted_at IS NULL
  AND status IN ('pending', 'success')
GROUP BY transaction_id
HAVING COUNT(*) > 1;
-- Should return 0 rows
```

## Troubleshooting

### Error: "duplicate key value violates unique constraint"

This means you still have duplicate payments. Run cleanup migration first:

```sql
\i migrations/010_cleanup_duplicate_payments.up.sql
```

Then retry the unique constraint migration.

### Check for Duplicates

```sql
-- Check duplicate payments
SELECT
  p1.transaction_id,
  p1.id as payment1_id,
  p1.status as payment1_status,
  p1.created_at as payment1_created,
  p2.id as payment2_id,
  p2.status as payment2_status,
  p2.created_at as payment2_created
FROM payments p1
INNER JOIN payments p2
  ON p1.transaction_id = p2.transaction_id
  AND p1.id < p2.id
WHERE p1.deleted_at IS NULL
  AND p2.deleted_at IS NULL
  AND (p1.status IN ('pending', 'success') OR p2.status IN ('pending', 'success'));
```

### Manual Cleanup (if migration fails)

```sql
BEGIN;

-- Delete QRIS codes for duplicates
DELETE FROM qris_codes
WHERE payment_id IN (
  SELECT id FROM payments
  WHERE transaction_id IN (
    SELECT transaction_id
    FROM payments
    WHERE deleted_at IS NULL
      AND status = 'pending'
    GROUP BY transaction_id
    HAVING COUNT(*) > 1
  )
  AND status = 'pending'
  AND id NOT IN (
    SELECT MIN(id)
    FROM payments
    WHERE deleted_at IS NULL
      AND status = 'pending'
    GROUP BY transaction_id
  )
);

-- Delete duplicate payments
DELETE FROM payments
WHERE id IN (
  SELECT id FROM payments
  WHERE transaction_id IN (
    SELECT transaction_id
    FROM payments
    WHERE deleted_at IS NULL
      AND status = 'pending'
    GROUP BY transaction_id
    HAVING COUNT(*) > 1
  )
  AND status = 'pending'
  AND id NOT IN (
    SELECT MIN(id)
    FROM payments
    WHERE deleted_at IS NULL
      AND status = 'pending'
    GROUP BY transaction_id
  )
);

COMMIT;
```

## Notes

- Always backup before running migrations on production
- Test migrations on development database first
- Data migrations (like 010) cannot be rolled back
- Unique constraint only applies to active (non-deleted) records
- Order_id is crucial for payment status checking with Midtrans
