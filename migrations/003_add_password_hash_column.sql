-- Add password_hash column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';

-- Add comment to document the column
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password for user authentication';