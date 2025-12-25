package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPaylaterLoanService_CreateLoan(t *testing.T) {
	t.Run("Success: Create loan with sufficient credit", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			DueDate:  time.Now().Add(30 * 24 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    0,
			AvailableLimit: 1000000,
			Status:         entity.PaylaterStatusActive,
		}

		// Mock transaction execution
		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(nil)

		mockAccountRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)
		mockLoanRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan *entity.PaylaterLoan) bool {
			return loan.UserID == 1 && loan.Amount == 100000 && loan.Total == 105000
		})).Return(nil).Run(func(args mock.Arguments) {
			loan := args.Get(1).(*entity.PaylaterLoan)
			loan.ID = 1
		})
		mockAccountRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(acc *entity.PaylaterAccount) bool {
			return acc.Outstanding == 105000 && acc.AvailableLimit == 895000
		})).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.CreateLoan(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(100000), result.Amount)
		assert.Equal(t, int64(5000), result.Interest)
		assert.Equal(t, int64(105000), result.Total)
		mockLoanRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			DueDate:  time.Now().Add(30 * 24 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(errors.New("paylater account not found or inactive"))

		mockAccountRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(nil, errors.New("not found"))

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.CreateLoan(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
		mockAccountRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Account not active", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			DueDate:  time.Now().Add(30 * 24 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    0,
			AvailableLimit: 1000000,
			Status:         entity.PaylaterStatusSuspended,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(errors.New("paylater account is not active"))

		mockAccountRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.CreateLoan(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not active")
		mockAccountRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Insufficient credit limit", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   100000,
			Interest: 5000,
			DueDate:  time.Now().Add(30 * 24 * time.Hour),
			Source:   entity.PaylaterLoanSourceCheckout,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    900000,
			AvailableLimit: 100000,
			Status:         entity.PaylaterStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(errors.New("insufficient credit limit for this loan"))

		mockAccountRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.CreateLoan(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient")
		mockAccountRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_GetLoanByID(t *testing.T) {
	t.Run("Success: Get loan by ID", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusActive,
			CreatedAt: time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", ctx, uint(1)).Return(loan, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetLoanByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, int64(100000), result.Amount)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Error: Loan not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()

		mockLoanRepo.On("GetLoanByID", ctx, uint(999)).Return(nil, errors.New("loan not found"))

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetLoanByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_GetLoansByUserID(t *testing.T) {
	t.Run("Success: Get loans by user ID", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loans := []entity.PaylaterLoan{
			{
				ID:        1,
				UserID:    1,
				Amount:    100000,
				Interest:  5000,
				Total:     105000,
				DueDate:   time.Now().Add(30 * 24 * time.Hour),
				Source:    entity.PaylaterLoanSourceCheckout,
				Status:    entity.PaylaterLoanStatusActive,
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				UserID:    1,
				Amount:    50000,
				Interest:  2500,
				Total:     52500,
				DueDate:   time.Now().Add(30 * 24 * time.Hour),
				Source:    entity.PaylaterLoanSourceTransfer,
				Status:    entity.PaylaterLoanStatusActive,
				CreatedAt: time.Now(),
			},
		}

		mockLoanRepo.On("GetLoansByUserID", ctx, uint(1)).Return(loans, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetLoansByUserID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, uint(2), result[1].ID)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Success: Empty loans list", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loans := []entity.PaylaterLoan{}

		mockLoanRepo.On("GetLoansByUserID", ctx, uint(1)).Return(loans, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetLoansByUserID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_ListLoans(t *testing.T) {
	t.Run("Success: List loans with pagination", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.ListPaylaterLoansRequest{
			Page:     1,
			PageSize: 10,
		}

		loans := []entity.PaylaterLoan{
			{ID: 1, UserID: 1, Amount: 100000, Total: 105000, Status: entity.PaylaterLoanStatusActive},
			{ID: 2, UserID: 1, Amount: 50000, Total: 52500, Status: entity.PaylaterLoanStatusActive},
		}

		mockLoanRepo.On("ListLoans", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(loans, int64(2), nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.ListLoans(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int64(2), result.TotalCount)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 1, result.TotalPages)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Success: List loans with filters", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		userID := uint(1)
		status := entity.PaylaterLoanStatusActive
		req := &dto.ListPaylaterLoansRequest{
			UserID:   &userID,
			Status:   &status,
			Page:     1,
			PageSize: 10,
		}

		loans := []entity.PaylaterLoan{
			{ID: 1, UserID: 1, Amount: 100000, Total: 105000, Status: entity.PaylaterLoanStatusActive},
		}

		mockLoanRepo.On("ListLoans", ctx, mock.MatchedBy(func(filters map[string]interface{}) bool {
			return filters["user_id"] == uint(1) && filters["status"] == entity.PaylaterLoanStatusActive
		}), 1, 10).Return(loans, int64(1), nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.ListLoans(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Success: Default pagination values", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.ListPaylaterLoansRequest{
			Page:     0, // Should default to 1
			PageSize: 0, // Should default to 10
		}

		loans := []entity.PaylaterLoan{}

		mockLoanRepo.On("ListLoans", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(loans, int64(0), nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.ListLoans(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_UpdateLoanStatus(t *testing.T) {
	t.Run("Success: Update loan status", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusActive,
			CreatedAt: time.Now(),
		}

		req := &dto.UpdateLoanStatusRequest{
			Status: entity.PaylaterLoanStatusOverdue,
		}

		mockLoanRepo.On("GetLoanByID", ctx, uint(1)).Return(loan, nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(1), entity.PaylaterLoanStatusOverdue).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.UpdateLoanStatus(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, entity.PaylaterLoanStatusOverdue, result.Status)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Error: Loan not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		req := &dto.UpdateLoanStatusRequest{
			Status: entity.PaylaterLoanStatusPaid,
		}

		mockLoanRepo.On("GetLoanByID", ctx, uint(999)).Return(nil, errors.New("loan not found"))

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.UpdateLoanStatus(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_MarkLoanAsPaid(t *testing.T) {
	t.Run("Success: Mark loan as paid", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusActive,
			CreatedAt: time.Now(),
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    105000,
			AvailableLimit: 895000,
			Status:         entity.PaylaterStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(nil)

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockLoanRepo.On("UpdateLoan", mock.Anything, mock.MatchedBy(func(l *entity.PaylaterLoan) bool {
			return l.Status == entity.PaylaterLoanStatusPaid
		})).Return(nil)
		mockAccountRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)
		mockAccountRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(acc *entity.PaylaterAccount) bool {
			return acc.Outstanding == 0 && acc.AvailableLimit == 1000000
		})).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.MarkLoanAsPaid(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, entity.PaylaterLoanStatusPaid, result.Status)
		mockLoanRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Loan already paid", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusPaid,
			CreatedAt: time.Now(),
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			txFunc(ctx)
		}).Return(errors.New("loan is already marked as paid"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.MarkLoanAsPaid(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already")
		mockLoanRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_ProcessOverdueLoans(t *testing.T) {
	t.Run("Success: Process overdue loans", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		overdueLoans := []entity.PaylaterLoan{
			{ID: 1, Status: entity.PaylaterLoanStatusActive},
			{ID: 2, Status: entity.PaylaterLoanStatusActive},
			{ID: 3, Status: entity.PaylaterLoanStatusActive},
		}

		mockLoanRepo.On("GetOverdueLoans", ctx, mock.AnythingOfType("time.Time")).Return(overdueLoans, nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(1), entity.PaylaterLoanStatusOverdue).Return(nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(2), entity.PaylaterLoanStatusOverdue).Return(nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(3), entity.PaylaterLoanStatusOverdue).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		count, err := service.ProcessOverdueLoans(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 3, count)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Success: No overdue loans", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		overdueLoans := []entity.PaylaterLoan{}

		mockLoanRepo.On("GetOverdueLoans", ctx, mock.AnythingOfType("time.Time")).Return(overdueLoans, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		count, err := service.ProcessOverdueLoans(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Partial Success: Some loans fail to update", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		overdueLoans := []entity.PaylaterLoan{
			{ID: 1, Status: entity.PaylaterLoanStatusActive},
			{ID: 2, Status: entity.PaylaterLoanStatusActive},
			{ID: 3, Status: entity.PaylaterLoanStatusActive},
		}

		mockLoanRepo.On("GetOverdueLoans", ctx, mock.AnythingOfType("time.Time")).Return(overdueLoans, nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(1), entity.PaylaterLoanStatusOverdue).Return(nil)
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(2), entity.PaylaterLoanStatusOverdue).Return(errors.New("update failed"))
		mockLoanRepo.On("UpdateLoanStatus", ctx, uint(3), entity.PaylaterLoanStatusOverdue).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		count, err := service.ProcessOverdueLoans(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2, count) // Only 2 succeeded
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_GetUserLoanStats(t *testing.T) {
	t.Run("Success: Get loan statistics", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		stats := map[string]interface{}{
			"total_loans":       int64(10),
			"active_loans":      int64(5),
			"total_borrowed":    int64(1000000),
			"total_paid":        int64(500000),
			"total_outstanding": int64(500000),
			"overdue_loans":     int64(2),
		}

		mockLoanRepo.On("GetLoanStats", ctx, uint(1)).Return(stats, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetUserLoanStats(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(10), result.TotalLoans)
		assert.Equal(t, int64(5), result.ActiveLoans)
		assert.Equal(t, int64(1000000), result.TotalBorrowed)
		assert.Equal(t, int64(500000), result.TotalPaid)
		assert.Equal(t, int64(500000), result.TotalOutstanding)
		assert.Equal(t, int64(2), result.OverdueLoans)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Error: Failed to get stats", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()

		mockLoanRepo.On("GetLoanStats", ctx, uint(1)).Return(nil, errors.New("database error"))

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		result, err := service.GetUserLoanStats(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestPaylaterLoanService_DeleteLoan(t *testing.T) {
	t.Run("Success: Delete paid loan", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusPaid,
			CreatedAt: time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", ctx, uint(1)).Return(loan, nil)
		mockLoanRepo.On("DeleteLoan", ctx, uint(1)).Return(nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		err := service.DeleteLoan(ctx, 1)

		assert.NoError(t, err)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Error: Cannot delete active loan", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()
		loan := &entity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    100000,
			Interest:  5000,
			Total:     105000,
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Source:    entity.PaylaterLoanSourceCheckout,
			Status:    entity.PaylaterLoanStatusActive,
			CreatedAt: time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", ctx, uint(1)).Return(loan, nil)

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		err := service.DeleteLoan(ctx, 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only delete paid loans")
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("Error: Loan not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockAccountRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		ctx := context.Background()

		mockLoanRepo.On("GetLoanByID", ctx, uint(999)).Return(nil, errors.New("loan not found"))

		service := NewPaylaterLoanService(mockLoanRepo, mockAccountRepo, mockTxManager, logger)
		err := service.DeleteLoan(ctx, 999)

		assert.Error(t, err)
		mockLoanRepo.AssertExpectations(t)
	})
}
