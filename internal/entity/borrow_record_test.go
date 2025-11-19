package entity

import (
	"testing"
	"time"
)

func TestBorrowRecord_Validate(t *testing.T) {
	now := time.Now()
	future := now.Add(14 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid borrow record",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
				RenewalCount: 0,
				LateFee:      0,
				Status:       BorrowStatusActive,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			record: &BorrowRecord{
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "missing book copy ID",
			record: &BorrowRecord{
				UserID:       "user-id",
				CheckoutDate: now,
				DueDate:      future,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "book copy ID is required",
		},
		{
			name: "missing checkout date",
			record: &BorrowRecord{
				UserID:     "user-id",
				BookCopyID: "copy-id",
				DueDate:    future,
				Status:     BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "checkout date is required",
		},
		{
			name: "missing due date",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "due date is required",
		},
		{
			name: "due date before checkout date",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: future,
				DueDate:      now,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "due date must be after checkout date",
		},
		{
			name: "missing status",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
				Status:       "invalid-status",
			},
			wantErr: true,
			errMsg:  "invalid borrow status",
		},
		{
			name: "negative renewal count",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
				RenewalCount: -1,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "renewal count cannot be negative",
		},
		{
			name: "negative late fee",
			record: &BorrowRecord{
				UserID:       "user-id",
				BookCopyID:   "copy-id",
				CheckoutDate: now,
				DueDate:      future,
				LateFee:      -1.0,
				Status:       BorrowStatusActive,
			},
			wantErr: true,
			errMsg:  "late fee cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BorrowRecord.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("BorrowRecord.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBorrowRecord_IsOverdue(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * 24 * time.Hour)
	future := now.Add(1 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		want    bool
	}{
		{
			name: "active record past due date",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusActive,
			},
			want: true,
		},
		{
			name: "active record not yet due",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusActive,
			},
			want: false,
		},
		{
			name: "returned record past due date",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusReturned,
			},
			want: false,
		},
		{
			name: "overdue status",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusOverdue,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.IsOverdue(); got != tt.want {
				t.Errorf("BorrowRecord.IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorrowRecord_CanRenew(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * 24 * time.Hour)
	future := now.Add(1 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		want    bool
	}{
		{
			name: "active record, not overdue, under renewal limit",
			record: &BorrowRecord{
				DueDate:      future,
				RenewalCount: 0,
				Status:       BorrowStatusActive,
			},
			want: true,
		},
		{
			name: "active record, not overdue, at renewal limit",
			record: &BorrowRecord{
				DueDate:      future,
				RenewalCount: MaxRenewalCount,
				Status:       BorrowStatusActive,
			},
			want: false,
		},
		{
			name: "active record, overdue",
			record: &BorrowRecord{
				DueDate:      past,
				RenewalCount: 0,
				Status:       BorrowStatusActive,
			},
			want: false,
		},
		{
			name: "returned record",
			record: &BorrowRecord{
				DueDate:      future,
				RenewalCount: 0,
				Status:       BorrowStatusReturned,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.CanRenew(); got != tt.want {
				t.Errorf("BorrowRecord.CanRenew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorrowRecord_CalculateLateFee(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(-5 * 24 * time.Hour) // 5 days overdue

	tests := []struct {
		name       string
		record     *BorrowRecord
		wantLateFee float64
	}{
		{
			name: "no return date",
			record: &BorrowRecord{
				DueDate:    dueDate,
				ReturnDate: nil,
			},
			wantLateFee: 0,
		},
		{
			name: "returned on time",
			record: &BorrowRecord{
				DueDate:    now,
				ReturnDate: &now,
			},
			wantLateFee: 0,
		},
		{
			name: "returned 5 days late",
			record: &BorrowRecord{
				DueDate:    dueDate,
				ReturnDate: &now,
			},
			wantLateFee: 2.50, // 5 days * $0.50
		},
		{
			name: "returned 25 days late (exceeds max fee)",
			record: &BorrowRecord{
				DueDate:    now.Add(-25 * 24 * time.Hour),
				ReturnDate: &now,
			},
			wantLateFee: MaxLateFee, // Capped at $10.00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.record.CalculateLateFee()
			if got != tt.wantLateFee {
				t.Errorf("BorrowRecord.CalculateLateFee() = %v, want %v", got, tt.wantLateFee)
			}
		})
	}
}

func TestBorrowRecord_GetOverdueDays(t *testing.T) {
	now := time.Now()
	past := now.Add(-5 * 24 * time.Hour)
	future := now.Add(5 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		want    int
	}{
		{
			name: "5 days overdue",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusActive,
			},
			want: 5,
		},
		{
			name: "not overdue",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusActive,
			},
			want: 0,
		},
		{
			name: "returned record",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusReturned,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.GetOverdueDays(); got != tt.want {
				t.Errorf("BorrowRecord.GetOverdueDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorrowRecord_StatusChecks(t *testing.T) {
	tests := []struct {
		name   string
		status BorrowRecordStatus
		checks map[string]bool
	}{
		{
			name:   "active record",
			status: BorrowStatusActive,
			checks: map[string]bool{
				"IsActive":   true,
				"IsReturned": false,
			},
		},
		{
			name:   "returned record",
			status: BorrowStatusReturned,
			checks: map[string]bool{
				"IsActive":   false,
				"IsReturned": true,
			},
		},
		{
			name:   "overdue record",
			status: BorrowStatusOverdue,
			checks: map[string]bool{
				"IsActive":   false,
				"IsReturned": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &BorrowRecord{Status: tt.status}

			if got := record.IsActive(); got != tt.checks["IsActive"] {
				t.Errorf("IsActive() = %v, want %v", got, tt.checks["IsActive"])
			}
			if got := record.IsReturned(); got != tt.checks["IsReturned"] {
				t.Errorf("IsReturned() = %v, want %v", got, tt.checks["IsReturned"])
			}
		})
	}
}

func TestBorrowRecord_Renew(t *testing.T) {
	now := time.Now()
	future := now.Add(14 * 24 * time.Hour)

	tests := []struct {
		name         string
		record       *BorrowRecord
		wantErr      bool
		wantRenewalCount int
	}{
		{
			name: "successful renewal",
			record: &BorrowRecord{
				DueDate:      future,
				RenewalCount: 0,
				Status:       BorrowStatusActive,
			},
			wantErr:      false,
			wantRenewalCount: 1,
		},
		{
			name: "renewal at limit",
			record: &BorrowRecord{
				DueDate:      future,
				RenewalCount: MaxRenewalCount,
				Status:       BorrowStatusActive,
			},
			wantErr:      true,
			wantRenewalCount: MaxRenewalCount,
		},
		{
			name: "renewal when overdue",
			record: &BorrowRecord{
				DueDate:      now.Add(-1 * 24 * time.Hour),
				RenewalCount: 0,
				Status:       BorrowStatusActive,
			},
			wantErr:      true,
			wantRenewalCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldDueDate := tt.record.DueDate
			err := tt.record.Renew()

			if (err != nil) != tt.wantErr {
				t.Errorf("BorrowRecord.Renew() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.record.RenewalCount != tt.wantRenewalCount {
				t.Errorf("RenewalCount = %v, want %v", tt.record.RenewalCount, tt.wantRenewalCount)
			}

			if !tt.wantErr {
				expectedDueDate := oldDueDate.Add(DefaultBorrowingPeriod)
				if !tt.record.DueDate.Equal(expectedDueDate) {
					t.Errorf("DueDate = %v, want %v", tt.record.DueDate, expectedDueDate)
				}
			}
		})
	}
}

func TestBorrowRecord_Return(t *testing.T) {
	now := time.Now()
	future := now.Add(14 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		wantErr bool
	}{
		{
			name: "successful return",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusActive,
			},
			wantErr: false,
		},
		{
			name: "already returned",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusReturned,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Return()

			if (err != nil) != tt.wantErr {
				t.Errorf("BorrowRecord.Return() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if tt.record.Status != BorrowStatusReturned {
					t.Errorf("Status = %v, want %v", tt.record.Status, BorrowStatusReturned)
				}
				if tt.record.ReturnDate == nil {
					t.Error("ReturnDate should not be nil after return")
				}
			}
		})
	}
}

func TestBorrowRecord_UpdateStatus(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * 24 * time.Hour)
	future := now.Add(1 * 24 * time.Hour)

	tests := []struct {
		name       string
		record     *BorrowRecord
		wantStatus BorrowRecordStatus
	}{
		{
			name: "active and overdue becomes overdue",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusActive,
			},
			wantStatus: BorrowStatusOverdue,
		},
		{
			name: "active and not overdue stays active",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusActive,
			},
			wantStatus: BorrowStatusActive,
		},
		{
			name: "returned stays returned",
			record: &BorrowRecord{
				DueDate: past,
				Status:  BorrowStatusReturned,
			},
			wantStatus: BorrowStatusReturned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.record.UpdateStatus()
			if tt.record.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", tt.record.Status, tt.wantStatus)
			}
		})
	}
}

func TestBorrowRecord_GetDaysRemaining(t *testing.T) {
	now := time.Now()
	future := now.Add(5 * 24 * time.Hour)

	tests := []struct {
		name    string
		record  *BorrowRecord
		wantMin int
		wantMax int
	}{
		{
			name: "5 days remaining",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusActive,
			},
			wantMin: 4,
			wantMax: 5,
		},
		{
			name: "overdue returns negative",
			record: &BorrowRecord{
				DueDate: now.Add(-5 * 24 * time.Hour),
				Status:  BorrowStatusActive,
			},
			wantMin: -6,
			wantMax: -4,
		},
		{
			name: "returned record returns 0",
			record: &BorrowRecord{
				DueDate: future,
				Status:  BorrowStatusReturned,
			},
			wantMin: 0,
			wantMax: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.record.GetDaysRemaining()
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("BorrowRecord.GetDaysRemaining() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestReservation_Validate(t *testing.T) {
	now := time.Now()
	future := now.Add(3 * 24 * time.Hour)

	tests := []struct {
		name    string
		res     *Reservation
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid reservation",
			res: &Reservation{
				UserID:          "user-id",
				BookID:          "book-id",
				ReservationDate: now,
				ExpiryDate:      future,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			res: &Reservation{
				BookID:          "book-id",
				ReservationDate: now,
				ExpiryDate:      future,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "missing book ID",
			res: &Reservation{
				UserID:          "user-id",
				ReservationDate: now,
				ExpiryDate:      future,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: true,
			errMsg:  "book ID is required",
		},
		{
			name: "expiry before reservation",
			res: &Reservation{
				UserID:          "user-id",
				BookID:          "book-id",
				ReservationDate: future,
				ExpiryDate:      now,
				Status:          ReservationStatusPending,
				QueuePosition:   1,
			},
			wantErr: true,
			errMsg:  "expiry date must be after reservation date",
		},
		{
			name: "invalid queue position",
			res: &Reservation{
				UserID:          "user-id",
				BookID:          "book-id",
				ReservationDate: now,
				ExpiryDate:      future,
				Status:          ReservationStatusPending,
				QueuePosition:   0,
			},
			wantErr: true,
			errMsg:  "queue position must be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.res.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reservation.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Reservation.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestReservation_IsExpired(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * 24 * time.Hour)
	future := now.Add(1 * 24 * time.Hour)

	tests := []struct {
		name string
		res  *Reservation
		want bool
	}{
		{
			name: "pending and expired",
			res: &Reservation{
				ExpiryDate: past,
				Status:     ReservationStatusPending,
			},
			want: true,
		},
		{
			name: "pending and not expired",
			res: &Reservation{
				ExpiryDate: future,
				Status:     ReservationStatusPending,
			},
			want: false,
		},
		{
			name: "fulfilled",
			res: &Reservation{
				ExpiryDate: past,
				Status:     ReservationStatusFulfilled,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.IsExpired(); got != tt.want {
				t.Errorf("Reservation.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReservation_Cancel(t *testing.T) {
	tests := []struct {
		name    string
		status  ReservationStatus
		wantErr bool
	}{
		{
			name:    "pending can be cancelled",
			status:  ReservationStatusPending,
			wantErr: false,
		},
		{
			name:    "fulfilled cannot be cancelled",
			status:  ReservationStatusFulfilled,
			wantErr: true,
		},
		{
			name:    "expired cannot be cancelled",
			status:  ReservationStatusExpired,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Reservation{Status: tt.status}
			err := res.Cancel()

			if (err != nil) != tt.wantErr {
				t.Errorf("Reservation.Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && res.Status != ReservationStatusCancelled {
				t.Errorf("Status = %v, want %v", res.Status, ReservationStatusCancelled)
			}
		})
	}
}

func TestReservation_Fulfill(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * 24 * time.Hour)
	past := now.Add(-1 * 24 * time.Hour)

	tests := []struct {
		name    string
		res     *Reservation
		wantErr bool
	}{
		{
			name: "pending and not expired can be fulfilled",
			res: &Reservation{
				ExpiryDate: future,
				Status:     ReservationStatusPending,
			},
			wantErr: false,
		},
		{
			name: "expired cannot be fulfilled",
			res: &Reservation{
				ExpiryDate: past,
				Status:     ReservationStatusPending,
			},
			wantErr: true,
		},
		{
			name: "cancelled cannot be fulfilled",
			res: &Reservation{
				ExpiryDate: future,
				Status:     ReservationStatusCancelled,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.res.Fulfill()

			if (err != nil) != tt.wantErr {
				t.Errorf("Reservation.Fulfill() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.res.Status != ReservationStatusFulfilled {
				t.Errorf("Status = %v, want %v", tt.res.Status, ReservationStatusFulfilled)
			}
		})
	}
}