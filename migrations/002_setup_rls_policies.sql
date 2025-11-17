-- Migration: Set up Row Level Security (RLS) policies
-- Version: 002
-- Description: Implements RLS policies for secure data access control
-- Author: Digital Library Team
-- Date: 2025-01-17

-- Enable Row Level Security on all tables
ALTER TABLE books ENABLE ROW LEVEL SECURITY;
ALTER TABLE book_copies ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE borrow_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE reservations ENABLE ROW LEVEL SECURITY;

-- ============================================
-- BOOKS TABLE POLICIES
-- ============================================

-- Policy: Anyone can view non-deleted books
CREATE POLICY "books_select_policy" ON books
    FOR SELECT
    USING (deleted_at IS NULL);

-- Policy: Admin and Librarian can insert books
CREATE POLICY "books_insert_policy" ON books
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Admin and Librarian can update books
CREATE POLICY "books_update_policy" ON books
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Only Admin can delete books (soft delete)
CREATE POLICY "books_delete_policy" ON books
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role = 'admin'
            AND users.status = 'active'
        )
    )
    WITH CHECK (deleted_at IS NOT NULL);

-- ============================================
-- BOOK COPIES TABLE POLICIES
-- ============================================

-- Policy: Anyone can view book copies
CREATE POLICY "book_copies_select_policy" ON book_copies
    FOR SELECT
    USING (true);

-- Policy: Admin and Librarian can insert book copies
CREATE POLICY "book_copies_insert_policy" ON book_copies
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Admin and Librarian can update book copies
CREATE POLICY "book_copies_update_policy" ON book_copies
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Only Admin can delete book copies
CREATE POLICY "book_copies_delete_policy" ON book_copies
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role = 'admin'
            AND users.status = 'active'
        )
    );

-- ============================================
-- USERS TABLE POLICIES
-- ============================================

-- Policy: Users can view their own profile
CREATE POLICY "users_select_own_policy" ON users
    FOR SELECT
    USING (id = auth.uid());

-- Policy: Admin can view all users
CREATE POLICY "users_select_admin_policy" ON users
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM users u
            WHERE u.id = auth.uid()
            AND u.role = 'admin'
            AND u.status = 'active'
        )
    );

-- Policy: Users can update their own profile (limited fields)
CREATE POLICY "users_update_own_policy" ON users
    FOR UPDATE
    USING (id = auth.uid())
    WITH CHECK (
        id = auth.uid()
        AND role = (SELECT role FROM users WHERE id = auth.uid())  -- Cannot change role
        AND status = (SELECT status FROM users WHERE id = auth.uid())  -- Cannot change status
    );

-- Policy: Admin can update any user
CREATE POLICY "users_update_admin_policy" ON users
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users u
            WHERE u.id = auth.uid()
            AND u.role = 'admin'
            AND u.status = 'active'
        )
    );

-- Policy: Only Admin can delete users
CREATE POLICY "users_delete_policy" ON users
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM users u
            WHERE u.id = auth.uid()
            AND u.role = 'admin'
            AND u.status = 'active'
        )
    );

-- ============================================
-- BORROW RECORDS TABLE POLICIES
-- ============================================

-- Policy: Users can view their own borrow records
CREATE POLICY "borrow_records_select_own_policy" ON borrow_records
    FOR SELECT
    USING (user_id = auth.uid());

-- Policy: Admin and Librarian can view all borrow records
CREATE POLICY "borrow_records_select_staff_policy" ON borrow_records
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Admin and Librarian can create borrow records
CREATE POLICY "borrow_records_insert_policy" ON borrow_records
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Admin and Librarian can update borrow records
CREATE POLICY "borrow_records_update_policy" ON borrow_records
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Only Admin can delete borrow records
CREATE POLICY "borrow_records_delete_policy" ON borrow_records
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role = 'admin'
            AND users.status = 'active'
        )
    );

-- ============================================
-- RESERVATIONS TABLE POLICIES
-- ============================================

-- Policy: Users can view their own reservations
CREATE POLICY "reservations_select_own_policy" ON reservations
    FOR SELECT
    USING (user_id = auth.uid());

-- Policy: Admin and Librarian can view all reservations
CREATE POLICY "reservations_select_staff_policy" ON reservations
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Active users can create reservations for themselves
CREATE POLICY "reservations_insert_policy" ON reservations
    FOR INSERT
    WITH CHECK (
        user_id = auth.uid()
        AND EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.status = 'active'
        )
    );

-- Policy: Users can cancel their own reservations
CREATE POLICY "reservations_update_own_policy" ON reservations
    FOR UPDATE
    USING (user_id = auth.uid())
    WITH CHECK (status = 'cancelled');

-- Policy: Admin and Librarian can update any reservation
CREATE POLICY "reservations_update_staff_policy" ON reservations
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- Policy: Admin and Librarian can delete reservations
CREATE POLICY "reservations_delete_policy" ON reservations
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = auth.uid()
            AND users.role IN ('admin', 'librarian')
            AND users.status = 'active'
        )
    );

-- ============================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================

COMMENT ON POLICY "books_select_policy" ON books IS 'Allow anyone to view non-deleted books';
COMMENT ON POLICY "books_insert_policy" ON books IS 'Allow admin and librarian to add books';
COMMENT ON POLICY "books_update_policy" ON books IS 'Allow admin and librarian to update books';
COMMENT ON POLICY "books_delete_policy" ON books IS 'Allow only admin to soft delete books';

COMMENT ON POLICY "users_select_own_policy" ON users IS 'Users can view their own profile';
COMMENT ON POLICY "users_select_admin_policy" ON users IS 'Admin can view all user profiles';
COMMENT ON POLICY "users_update_own_policy" ON users IS 'Users can update their own profile (limited)';
COMMENT ON POLICY "users_update_admin_policy" ON users IS 'Admin can update any user profile';

COMMENT ON POLICY "borrow_records_select_own_policy" ON borrow_records IS 'Users can view their own borrowing history';
COMMENT ON POLICY "borrow_records_select_staff_policy" ON borrow_records IS 'Staff can view all borrowing records';

COMMENT ON POLICY "reservations_select_own_policy" ON reservations IS 'Users can view their own reservations';
COMMENT ON POLICY "reservations_insert_policy" ON reservations IS 'Active users can create reservations';