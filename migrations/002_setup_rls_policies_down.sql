-- Migration Rollback: Remove Row Level Security (RLS) policies
-- Version: 002
-- Description: Drops all RLS policies and disables RLS on tables
-- Author: Digital Library Team
-- Date: 2025-01-17

-- ============================================
-- DROP RESERVATIONS POLICIES
-- ============================================
DROP POLICY IF EXISTS "reservations_delete_policy" ON reservations;
DROP POLICY IF EXISTS "reservations_update_staff_policy" ON reservations;
DROP POLICY IF EXISTS "reservations_update_own_policy" ON reservations;
DROP POLICY IF EXISTS "reservations_insert_policy" ON reservations;
DROP POLICY IF EXISTS "reservations_select_staff_policy" ON reservations;
DROP POLICY IF EXISTS "reservations_select_own_policy" ON reservations;

-- ============================================
-- DROP BORROW RECORDS POLICIES
-- ============================================
DROP POLICY IF EXISTS "borrow_records_delete_policy" ON borrow_records;
DROP POLICY IF EXISTS "borrow_records_update_policy" ON borrow_records;
DROP POLICY IF EXISTS "borrow_records_insert_policy" ON borrow_records;
DROP POLICY IF EXISTS "borrow_records_select_staff_policy" ON borrow_records;
DROP POLICY IF EXISTS "borrow_records_select_own_policy" ON borrow_records;

-- ============================================
-- DROP USERS POLICIES
-- ============================================
DROP POLICY IF EXISTS "users_delete_policy" ON users;
DROP POLICY IF EXISTS "users_update_admin_policy" ON users;
DROP POLICY IF EXISTS "users_update_own_policy" ON users;
DROP POLICY IF EXISTS "users_select_admin_policy" ON users;
DROP POLICY IF EXISTS "users_select_own_policy" ON users;

-- ============================================
-- DROP BOOK COPIES POLICIES
-- ============================================
DROP POLICY IF EXISTS "book_copies_delete_policy" ON book_copies;
DROP POLICY IF EXISTS "book_copies_update_policy" ON book_copies;
DROP POLICY IF EXISTS "book_copies_insert_policy" ON book_copies;
DROP POLICY IF EXISTS "book_copies_select_policy" ON book_copies;

-- ============================================
-- DROP BOOKS POLICIES
-- ============================================
DROP POLICY IF EXISTS "books_delete_policy" ON books;
DROP POLICY IF EXISTS "books_update_policy" ON books;
DROP POLICY IF EXISTS "books_insert_policy" ON books;
DROP POLICY IF EXISTS "books_select_policy" ON books;

-- ============================================
-- DISABLE ROW LEVEL SECURITY
-- ============================================
ALTER TABLE reservations DISABLE ROW LEVEL SECURITY;
ALTER TABLE borrow_records DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
ALTER TABLE book_copies DISABLE ROW LEVEL SECURITY;
ALTER TABLE books DISABLE ROW LEVEL SECURITY;