package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBorrowRecord_Validate(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(DefaultBorrowingPeriod)

	tests := []struct {
		name    string
		record  *BorrowRecord
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid borrow record",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
				RenewalCount: 0,
				LateFee:      0,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			record: &BorrowRecord{
				UserID:       "",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "missing book copy ID",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "book copy ID is required",
		},
		{
			name: "zero checkout date",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: time.Time{},
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "checkout date is required",
		},
		{
			name: "zero due date",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      time.Time{},
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "due date is required",
		},
		{
			name: "due date before checkout date",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      now.Add(-24 * time.Hour),
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "due date must be after checkout date",
		},
		{
			name: "missing status",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       "",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       "invalid",
			},
			wantErr: true,
			errMsg:  "invalid borrow status",
		},
		{
			name: "negative renewal count",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
				RenewalCount: -1,
			},
			wantErr: true,
			errMsg:  "renewal count cannot be negative",
		},
		{
			name: "negative late fee",
			record: &BorrowRecord{
				UserID:       "user-1",
				BookCopyID:   "copy-1",
				CheckoutDate: now,
				DueDate:      dueDate,
				Status:       BorrowStatusActive,
				LateFee:      -1.0,
			},
			wantErr: true,
			errMsg:  "late fee cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBorrowRecord_IsOverdue(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		record   *BorrowRecord
		expected bool
	}{
		{
			name: "not overdue - future due date",
			record: &BorrowRecord{
				DueDate: now.Add(24 * time.Hour),
				Status:  BorrowStatusActive,
			},
			expected: false,
		},
		{
			name: "overdue - past due date",
			record: &BorrowRecord{
				DueDate: now.Add(-24 * time.Hour),
				Status:  BorrowStatusActive,
			},
			expected: true,
		},
		{
			name: "not overdue - already returned",
			record: &BorrowRecord{
				DueDate: now.Add(-24 * time.Hour),
				Status:  BorrowStatusReturned,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.record.IsOverdue())
		})
	}
}

func TestBorrowRecord_CanRenew(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		record   *BorrowRecord
		expected bool
	}{
		{
			name: "can renew - active and under limit",
			record: &BorrowRecord{
				Status:       BorrowStatusActive,
				RenewalCount: 0,
				DueDate:      now.Add(24 * time.Hour),
			},
			expected: true,
		},
		{
			name: "cannot renew - at max renewals",
			record: &BorrowRecord{
				Status:       BorrowStatusActive,
				RenewalCount: MaxRenewalCount,
				DueDate:      now.Add(24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot renew - overdue",
			record: &BorrowRecord{
				Status:       BorrowStatusActive,
				RenewalCount: 0,
				DueDate:      now.Add(-24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot renew - already returned",
			record: &BorrowRecord{
				Status:       BorrowStatusReturned,
				RenewalCount: 0,
				DueDate:      now.Add(24 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.record.CanRenew())
		})
	}
}

func TestBorrowRecord_CalculateLateFee(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(-10 * 24 * time.Hour) // 10 days ago

	tests := []struct {
		name       string
		record     *BorrowRecord
		expected   float64
		description string
	}{
		{
			name: "no late fee - returned on time",
			record: &BorrowRecord{
				DueDate:    dueDate,
				ReturnDate: &dueDate,
			},
			expected:    0,
			description: "returned on due date",
		},
		{
			name: "no late fee - returned early",
			record: func() *BorrowRecord {
				early := dueDate.Add(-24 * time.Hour)
				return &BorrowRecord{
					DueDate:    dueDate,
					ReturnDate: &early,
				}
			}(),
			expected:    0,
			description: "returned before due date",
		},
		{
			name: "late fee - 5 days overdue",
			record: func() *BorrowRecord {
				late := dueDate.Add(5 * 24 * time.Hour)
				return &BorrowRecord{
					DueDate:    dueDate,
					ReturnDate: &late,
				}
			}(),
			expected:    2.50, // 5 days * $0.50
			description: "5 days late",
		},
		{
			name: "late fee capped at max",
			record: func() *BorrowRecord {
				late := dueDate.Add(30 * 24 * time.Hour)
				return &BorrowRecord{
					DueDate:    dueDate,
					ReturnDate: &late,
				}
			}(),
			expected:    MaxLateFee, // capped at $10.00
			description: "30 days late, capped at max",
		},
		{
			name: "no late fee - not yet returned",
			record: &BorrowRecord{
				DueDate:    dueDate,
				ReturnDate: nil,
			},
			expected:    0,
			description: "no return date set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee := tt.record.CalculateLateFee()
			assert.Equal(t, tt.expected, fee, tt.description)
		})
	}
}

func TestBorrowRecord_CalculateCurrentLateFee(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		record   *BorrowRecord
		minFee   float64
		maxFee   float64
	}{
		{
			name: "no fee - not overdue",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(24 * time.Hour),
			},
			minFee: 0,
			maxFee: 0,
		},
		{
			name: "no fee - already returned",
			record: &BorrowRecord{
				Status:  BorrowStatusReturned,
				DueDate: now.Add(-24 * time.Hour),
			},
			minFee: 0,
			maxFee: 0,
		},
		{
			name: "has fee - 5 days overdue",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(-5 * 24 * time.Hour),
			},
			minFee: 2.0,  // at least 4 days * $0.50
			maxFee: 3.0,  // at most 6 days * $0.50
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee := tt.record.CalculateCurrentLateFee()
			assert.GreaterOrEqual(t, fee, tt.minFee)
			assert.LessOrEqual(t, fee, tt.maxFee)
		})
	}
}

func TestBorrowRecord_GetOverdueDays(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		record   *BorrowRecord
		minDays  int
		maxDays  int
	}{
		{
			name: "not overdue",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(24 * time.Hour),
			},
			minDays: 0,
			maxDays: 0,
		},
		{
			name: "5 days overdue",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(-5 * 24 * time.Hour),
			},
			minDays: 4,
			maxDays: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := tt.record.GetOverdueDays()
			assert.GreaterOrEqual(t, days, tt.minDays)
			assert.LessOrEqual(t, days, tt.maxDays)
		})
	}
}

func TestBorrowRecord_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   BorrowRecordStatus
		expected bool
	}{
		{
			name:     "active",
			status:   BorrowStatusActive,
			expected: true,
		},
		{
			name:     "returned",
			status:   BorrowStatusReturned,
			expected: false,
		},
		{
			name:     "overdue",
			status:   BorrowStatusOverdue,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &BorrowRecord{Status: tt.status}
			assert.Equal(t, tt.expected, record.IsActive())
		})
	}
}

func TestBorrowRecord_IsReturned(t *testing.T) {
	tests := []struct {
		name     string
		status   BorrowRecordStatus
		expected bool
	}{
		{
			name:     "returned",
			status:   BorrowStatusReturned,
			expected: true,
		},
		{
			name:     "active",
			status:   BorrowStatusActive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &BorrowRecord{Status: tt.status}
			assert.Equal(t, tt.expected, record.IsReturned())
		})
	}
}

func TestBorrowRecord_Renew(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(7 * 24 * time.Hour)

	t.Run("success - renew active book", func(t *testing.T) {
		record := &BorrowRecord{
			Status:       BorrowStatusActive,
			RenewalCount: 0,
			DueDate:      dueDate,
		}

		err := record.Renew()
		
		require.NoError(t, err)
		assert.Equal(t, 1, record.RenewalCount)
		expectedDueDate := dueDate.Add(DefaultBorrowingPeriod)
		assert.Equal(t, expectedDueDate, record.DueDate)
	})

	t.Run("error - at max renewals", func(t *testing.T) {
		record := &BorrowRecord{
			Status:       BorrowStatusActive,
			RenewalCount: MaxRenewalCount,
			DueDate:      dueDate,
		}

		err := record.Renew()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be renewed")
	})

	t.Run("error - overdue", func(t *testing.T) {
		record := &BorrowRecord{
			Status:       BorrowStatusActive,
			RenewalCount: 0,
			DueDate:      now.Add(-24 * time.Hour),
		}

		err := record.Renew()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be renewed")
	})
}

func TestBorrowRecord_Return(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(-5 * 24 * time.Hour)

	t.Run("success - return book", func(t *testing.T) {
		record := &BorrowRecord{
			Status:  BorrowStatusActive,
			DueDate: dueDate,
		}

		err := record.Return()
		
		require.NoError(t, err)
		assert.Equal(t, BorrowStatusReturned, record.Status)
		assert.NotNil(t, record.ReturnDate)
		assert.Greater(t, record.LateFee, 0.0) // Should have late fee
	})

	t.Run("error - already returned", func(t *testing.T) {
		returnDate := now
		record := &BorrowRecord{
			Status:     BorrowStatusReturned,
			DueDate:    dueDate,
			ReturnDate: &returnDate,
		}

		err := record.Return()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already been returned")
	})
}

func TestBorrowRecord_UpdateStatus(t *testing.T) {
	now := time.Now()

	t.Run("update to overdue", func(t *testing.T) {
		record := &BorrowRecord{
			Status:  BorrowStatusActive,
			DueDate: now.Add(-24 * time.Hour),
		}

		record.UpdateStatus()
		
		assert.Equal(t, BorrowStatusOverdue, record.Status)
	})

	t.Run("no update - not overdue", func(t *testing.T) {
		record := &BorrowRecord{
			Status:  BorrowStatusActive,
			DueDate: now.Add(24 * time.Hour),
		}

		record.UpdateStatus()
		
		assert.Equal(t, BorrowStatusActive, record.Status)
	})
}

func TestBorrowRecord_GetDaysRemaining(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		record   *BorrowRecord
		minDays  int
		maxDays  int
	}{
		{
			name: "7 days remaining",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(7 * 24 * time.Hour),
			},
			minDays: 6,
			maxDays: 8,
		},
		{
			name: "overdue - negative days",
			record: &BorrowRecord{
				Status:  BorrowStatusActive,
				DueDate: now.Add(-5 * 24 * time.Hour),
			},
			minDays: -6,
			maxDays: -4,
		},
		{
			name: "returned - 0 days",
			record: &BorrowRecord{
				Status:  BorrowStatusReturned,
				DueDate: now.Add(7 * 24 * time.Hour),
			},
			minDays: 0,
			maxDays: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := tt.record.GetDaysRemaining()
			assert.GreaterOrEqual(t, days, tt.minDays)
			assert.LessOrEqual(t, days, tt.maxDays)
		})
	}
}

// Reservation tests

func TestReservation_Validate(t *testing.T) {
	now := time.Now()
	expiryDate := now.Add(DefaultReservationHoldPeriod)

	tests := []struct {
		name        string
		reservation *Reservation
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid reservation",
			reservation: &Reservation{
				UserID:          "user-1",
				BookID:          "book-1",
				ReservationDate: now,
				ExpiryDate:      expiryDate,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			reservation: &Reservation{
				UserID:          "",
				BookID:          "book-1",
				ReservationDate: now,
				ExpiryDate:      expiryDate,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "missing book ID",
			reservation: &Reservation{
				UserID:          "user-1",
				BookID:          "",
				ReservationDate: now,
				ExpiryDate:      expiryDate,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: true,
			errMsg:  "book ID is required",
		},
		{
			name: "invalid queue position",
			reservation: &Reservation{
				UserID:          "user-1",
				BookID:          "book-1",
				ReservationDate: now,
				ExpiryDate:      expiryDate,
				Status:          ReservationStatusPending,
				QueuePosition:   0,
			},
			wantErr: true,
			errMsg:  "queue position must be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.reservation.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReservation_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		res      *Reservation
		expected bool
	}{
		{
			name: "not expired - future expiry",
			res: &Reservation{
				Status:     ReservationStatusPending,
				ExpiryDate: now.Add(24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "expired - past expiry",
			res: &Reservation{
				Status:     ReservationStatusPending,
				ExpiryDate: now.Add(-24 * time.Hour),
			},
			expected: true,
		},
		{
			name: "not expired - fulfilled",
			res: &Reservation{
				Status:     ReservationStatusFulfilled,
				ExpiryDate: now.Add(-24 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.res.IsExpired())
		})
	}
}

func TestReservation_IsPending(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		res      *Reservation
		expected bool
	}{
		{
			name: "pending and not expired",
			res: &Reservation{
				Status:     ReservationStatusPending,
				ExpiryDate: now.Add(24 * time.Hour),
			},
			expected: true,
		},
		{
			name: "pending but expired",
			res: &Reservation{
				Status:     ReservationStatusPending,
				ExpiryDate: now.Add(-24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "fulfilled",
			res: &Reservation{
				Status:     ReservationStatusFulfilled,
				ExpiryDate: now.Add(24 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.res.IsPending())
		})
	}
}

func TestReservation_Cancel(t *testing.T) {
	t.Run("success - cancel pending", func(t *testing.T) {
		res := &Reservation{
			Status: ReservationStatusPending,
		}

		err := res.Cancel()
		
		require.NoError(t, err)
		assert.Equal(t, ReservationStatusCancelled, res.Status)
	})

	t.Run("error - not pending", func(t *testing.T) {
		res := &Reservation{
			Status: ReservationStatusFulfilled,
		}

		err := res.Cancel()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only pending reservations can be cancelled")
	})
}

func TestReservation_Fulfill(t *testing.T) {
	now := time.Now()

	t.Run("success - fulfill pending", func(t *testing.T) {
		res := &Reservation{
			Status:     ReservationStatusPending,
			ExpiryDate: now.Add(24 * time.Hour),
		}

		err := res.Fulfill()
		
		require.NoError(t, err)
		assert.Equal(t, ReservationStatusFulfilled, res.Status)
	})

	t.Run("error - expired", func(t *testing.T) {
		res := &Reservation{
			Status:     ReservationStatusPending,
			ExpiryDate: now.Add(-24 * time.Hour),
		}

		err := res.Fulfill()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be fulfilled")
	})
}

func TestReservation_Expire(t *testing.T) {
	t.Run("success - expire pending", func(t *testing.T) {
		res := &Reservation{
			Status: ReservationStatusPending,
		}

		err := res.Expire()
		
		require.NoError(t, err)
		assert.Equal(t, ReservationStatusExpired, res.Status)
	})

	t.Run("error - not pending", func(t *testing.T) {
		res := &Reservation{
			Status: ReservationStatusFulfilled,
		}

		err := res.Expire()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only pending reservations can be expired")
	})
}