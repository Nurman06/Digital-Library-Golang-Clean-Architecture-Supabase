-- Update users table to integrate with Supabase Auth
-- Remove password_hash column (auth handled by Supabase Auth)
-- Add foreign key reference to auth.users

-- Remove password_hash column if it exists
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- Add foreign key constraint to auth.users
-- Note: This assumes auth.users table exists from Supabase Auth
ALTER TABLE users 
  DROP CONSTRAINT IF EXISTS users_id_fkey,
  ADD CONSTRAINT users_id_fkey 
    FOREIGN KEY (id) 
    REFERENCES auth.users(id) 
    ON DELETE CASCADE;

-- Add comment to document the relationship
COMMENT ON TABLE users IS 'Extended user information for library system. References auth.users for authentication.';
COMMENT ON COLUMN users.id IS 'Foreign key to auth.users(id). Authentication handled by Supabase Auth.';