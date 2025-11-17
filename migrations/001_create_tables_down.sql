-- Migration Rollback: Drop initial database schema for Digital Library System
-- Version: 001
-- Description: Drops all tables, indexes, triggers, and functions created in 001_create_tables.sql
-- Author: Digital Library Team
-- Date: 2025-01-17

-- Drop triggers first
DROP TRIGGER IF EXISTS update_reservations_updated_at ON reservations;
DROP TRIGGER IF EXISTS update_borrow_records_updated_at ON borrow_records;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_book_copies_updated_at ON book_copies;
DROP TRIGGER IF EXISTS update_books_updated_at ON books;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_reservations_status;
DROP INDEX IF EXISTS idx_reservations_book_id;
DROP INDEX IF EXISTS idx_reservations_user_id;
DROP INDEX IF EXISTS idx_borrow_records_checkout_date;
DROP INDEX IF EXISTS idx_borrow_records_due_date;
DROP INDEX IF EXISTS idx_borrow_records_status;
DROP INDEX IF EXISTS idx_borrow_records_book_copy_id;
DROP INDEX IF EXISTS idx_borrow_records_user_id;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_book_copies_status;
DROP INDEX IF EXISTS idx_book_copies_book_id;
DROP INDEX IF EXISTS idx_books_deleted_at;
DROP INDEX IF EXISTS idx_books_category;
DROP INDEX IF EXISTS idx_books_isbn;
DROP INDEX IF EXISTS idx_books_author;
DROP INDEX IF EXISTS idx_books_title;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS reservations CASCADE;
DROP TABLE IF EXISTS borrow_records CASCADE;
DROP TABLE IF EXISTS book_copies CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS books CASCADE;

-- Drop extension (optional - only if no other tables use it)
-- DROP EXTENSION IF EXISTS "pgcrypto";