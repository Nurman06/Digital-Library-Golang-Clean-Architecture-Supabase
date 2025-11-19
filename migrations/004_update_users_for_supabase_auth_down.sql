-- Rollback migration: Restore password_hash column and remove auth.users foreign key

-- Remove foreign key constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_id_fkey;

-- Add back password_hash column
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';

-- Remove comments
COMMENT ON TABLE users IS NULL;
COMMENT ON COLUMN users.id IS NULL;