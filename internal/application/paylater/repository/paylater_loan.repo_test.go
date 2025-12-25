package repository

import (
	"context"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaylaterLoanRepository_CreateLoan(t *testing.T) {
	t.Run("should create paylater loan successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}

		// When
		err = repo.CreateLoan(ctx, loan)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, loan.ID)
		assert.Equal(t, uint(1), loan.UserID)
		assert.Equal(t, int64(100000), loan.Amount)
		assert.Equal(t, int64(5000), loan.Interest)
		assert.Equal(t, int64(105000), loan.Total)
	})

	t.Run("should create multiple loans for same user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		loan1 := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}

		loan2 := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}

		// When
		err = repo.CreateLoan(ctx, loan1)
		assert.NoError(t, err)

		err = repo.CreateLoan(ctx, loan2)
		assert.NoError(t, err)

		// Then
		loans, err := repo.GetLoansByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, loans, 2)
	})
}

func TestPaylaterLoanRepository_GetLoanByID(t *testing.T) {
	t.Run("should get paylater loan by ID successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		retrieved, err := repo.GetLoanByID(ctx, loan.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, loan.ID, retrieved.ID)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(100000), retrieved.Amount)
		assert.Equal(t, int64(5000), retrieved.Interest)
		assert.Equal(t, int64(105000), retrieved.Total)
		assert.Equal(t, entity.PaylaterLoanSourceCheckout, retrieved.Source)
		assert.Equal(t, entity.PaylaterLoanStatusActive, retrieved.Status)
	})

	t.Run("should return error when loan not found", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetLoanByID(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestPaylaterLoanRepository_GetLoansByUserID(t *testing.T) {
	t.Run("should get all loans for a user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create loans for user 1
		for i := 0; i < 3; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create loan for user 2
		loan := &entity.PaylaterLoan{
			UserID:   2,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		loans, err := repo.GetLoansByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 3)
		for _, l := range loans {
			assert.Equal(t, uint(1), l.UserID)
		}
	})

	t.Run("should return empty array when user has no loans", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		// When
		loans, err := repo.GetLoansByUserID(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, loans)
	})
}

func TestPaylaterLoanRepository_ListLoans(t *testing.T) {
	t.Run("should list loans with pagination", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create 15 loans
		for i := 0; i < 15; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// When - Page 1
		req := &dto.ListPaylaterLoansRequest{
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 10)
		assert.Equal(t, int64(15), totalCount)

		// When - Page 2
		req.Page = 2
		loans, totalCount, err = repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 5)
		assert.Equal(t, int64(15), totalCount)
	})

	t.Run("should filter loans by user_id", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create loans for user 1
		for i := 0; i < 3; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create loans for user 2
		for i := 0; i < 2; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   2,
				Amount:   50000,
				Interest: 2500,
				Total:    52500,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceTransfer,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// When
		userID := uint(1)
		req := &dto.ListPaylaterLoansRequest{
			UserID:   &userID,
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 3)
		assert.Equal(t, int64(3), totalCount)
		for _, l := range loans {
			assert.Equal(t, uint(1), l.UserID)
		}
	})

	t.Run("should filter loans by status", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create active loans
		for i := 0; i < 3; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create paid loans
		for i := 0; i < 2; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   50000,
				Interest: 2500,
				Total:    52500,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceTransfer,
				Status:   entity.PaylaterLoanStatusPaid,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// When
		status := string(entity.PaylaterLoanStatusActive)
		req := &dto.ListPaylaterLoansRequest{
			Status:   &status,
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 3)
		assert.Equal(t, int64(3), totalCount)
		for _, l := range loans {
			assert.Equal(t, entity.PaylaterLoanStatusActive, l.Status)
		}
	})

	t.Run("should filter loans by source", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create checkout loans
		for i := 0; i < 3; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create transfer loans
		for i := 0; i < 2; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   50000,
				Interest: 2500,
				Total:    52500,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceTransfer,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// When
		source := string(entity.PaylaterLoanSourceCheckout)
		req := &dto.ListPaylaterLoansRequest{
			Source:   &source,
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 3)
		assert.Equal(t, int64(3), totalCount)
		for _, l := range loans {
			assert.Equal(t, entity.PaylaterLoanSourceCheckout, l.Source)
		}
	})

	t.Run("should filter loans by date range", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		now := time.Now()
		dueDate := now.Add(30 * 24 * time.Hour)

		// Create old loan (2 days ago)
		oldLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, oldLoan)
		require.NoError(t, err)
		// Update created_at to 2 days ago (simulate old loan)
		db.Model(oldLoan).Update("created_at", now.Add(-48*time.Hour))

		// Create recent loan
		recentLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, recentLoan)
		require.NoError(t, err)

		// When - filter from yesterday
		fromDate := now.Add(-24 * time.Hour).Format("2006-01-02")
		req := &dto.ListPaylaterLoansRequest{
			FromDate: &fromDate,
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 1)
		assert.Equal(t, int64(1), totalCount)
	})

	t.Run("should apply multiple filters", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create loans with different combinations
		loan1 := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan1)
		require.NoError(t, err)

		loan2 := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan2)
		require.NoError(t, err)

		loan3 := &entity.PaylaterLoan{
			UserID:   2,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan3)
		require.NoError(t, err)

		// When - filter by user_id and source
		userID := uint(1)
		source := string(entity.PaylaterLoanSourceCheckout)
		req := &dto.ListPaylaterLoansRequest{
			UserID:   &userID,
			Source:   &source,
			Page:     1,
			PageSize: 10,
		}
		loans, totalCount, err := repo.ListLoans(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.Len(t, loans, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, uint(1), loans[0].UserID)
		assert.Equal(t, entity.PaylaterLoanSourceCheckout, loans[0].Source)
	})
}

func TestPaylaterLoanRepository_UpdateLoanStatus(t *testing.T) {
	t.Run("should update loan status successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		err = repo.UpdateLoanStatus(ctx, loan.ID, string(entity.PaylaterLoanStatusPaid))

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetLoanByID(ctx, loan.ID)
		assert.NoError(t, err)
		assert.Equal(t, entity.PaylaterLoanStatusPaid, retrieved.Status)
	})

	t.Run("should return error when loan not found", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		// When
		err = repo.UpdateLoanStatus(ctx, 9999, string(entity.PaylaterLoanStatusPaid))

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestPaylaterLoanRepository_UpdateLoan(t *testing.T) {
	t.Run("should update loan successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		loan.Status = entity.PaylaterLoanStatusPaid
		loan.Amount = 120000
		err = repo.UpdateLoan(ctx, loan)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetLoanByID(ctx, loan.ID)
		assert.NoError(t, err)
		assert.Equal(t, entity.PaylaterLoanStatusPaid, retrieved.Status)
		assert.Equal(t, int64(120000), retrieved.Amount)
	})
}

func TestPaylaterLoanRepository_GetOverdueLoans(t *testing.T) {
	t.Run("should get overdue loans", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		now := time.Now()

		// Create overdue loan (due yesterday)
		overdueLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  now.Add(-24 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, overdueLoan)
		require.NoError(t, err)

		// Create current loan (due tomorrow)
		currentLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  now.Add(24 * time.Hour),
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, currentLoan)
		require.NoError(t, err)

		// Create paid overdue loan (should not be included)
		paidOverdueLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   30000,
			Interest: 1500,
			Total:    31500,
			DueDate:  now.Add(-48 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusPaid,
		}
		err = repo.CreateLoan(ctx, paidOverdueLoan)
		require.NoError(t, err)

		// When
		overdueLoans, err := repo.GetOverdueLoans(ctx, now)

		// Then
		assert.NoError(t, err)
		assert.Len(t, overdueLoans, 1)
		assert.Equal(t, overdueLoan.ID, overdueLoans[0].ID)
	})

	t.Run("should return empty when no overdue loans", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		now := time.Now()

		// Create current loan (due tomorrow)
		currentLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   50000,
			Interest: 2500,
			Total:    52500,
			DueDate:  now.Add(24 * time.Hour),
			Source:   entity.PaylaterLoanSourceTransfer,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, currentLoan)
		require.NoError(t, err)

		// When
		overdueLoans, err := repo.GetOverdueLoans(ctx, now)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, overdueLoans)
	})
}

func TestPaylaterLoanRepository_GetLoanStats(t *testing.T) {
	t.Run("should get loan statistics for user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)

		// Create 3 active loans
		for i := 0; i < 3; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   100000,
				Interest: 5000,
				Total:    105000,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceCheckout,
				Status:   entity.PaylaterLoanStatusActive,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create 2 paid loans
		for i := 0; i < 2; i++ {
			loan := &entity.PaylaterLoan{
				UserID:   1,
				Amount:   50000,
				Interest: 2500,
				Total:    52500,
				DueDate:  dueDate,
				Source:   entity.PaylaterLoanSourceTransfer,
				Status:   entity.PaylaterLoanStatusPaid,
			}
			err = repo.CreateLoan(ctx, loan)
			require.NoError(t, err)
		}

		// Create 1 overdue loan
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   30000,
			Interest: 1500,
			Total:    31500,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusOverdue,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		stats, err := repo.GetLoanStats(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(6), stats["total_loans"])
		assert.Equal(t, int64(3), stats["active_loans"])
		assert.Equal(t, int64(1), stats["overdue_loans"])
		assert.Equal(t, int64(3*105000+2*52500+31500), stats["total_borrowed"])
		assert.Equal(t, int64(2*52500), stats["total_paid"])
		assert.Equal(t, int64(3*105000+31500), stats["total_outstanding"])
	})

	t.Run("should return zero stats for user with no loans", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		// When
		stats, err := repo.GetLoanStats(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), stats["total_loans"])
		assert.Equal(t, int64(0), stats["active_loans"])
		assert.Equal(t, int64(0), stats["overdue_loans"])
		assert.Equal(t, int64(0), stats["total_borrowed"])
		assert.Equal(t, int64(0), stats["total_paid"])
		assert.Equal(t, int64(0), stats["total_outstanding"])
	})
}

func TestPaylaterLoanRepository_DeleteLoan(t *testing.T) {
	t.Run("should delete loan successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		dueDate := time.Now().Add(30 * 24 * time.Hour)
		loan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			Total:    105000,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceCheckout,
			Status:   entity.PaylaterLoanStatusActive,
		}
		err = repo.CreateLoan(ctx, loan)
		require.NoError(t, err)

		// When
		err = repo.DeleteLoan(ctx, loan.ID)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetLoanByID(ctx, loan.ID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("should return error when deleting non-existent loan", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterLoanRepository(db, logger)
		ctx := context.Background()

		// When
		err = repo.DeleteLoan(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
