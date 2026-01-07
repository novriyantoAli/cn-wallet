package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPaylaterRepaymentService_ProcessRepayment(t *testing.T) {
	t.Run("Success: Process repayment with valid loan", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:        1,
			UserID:    1,
			Amount:    1000000,
			Interest:  100000,
			Total:     1100000,
			DueDate:   time.Now().AddDate(0, 0, 30),
			Source:    paylaterEntity.PaylaterLoanSourceCheckout,
			Status:    paylaterEntity.PaylaterLoanStatusActive,
			CreatedAt: time.Now(),
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockRepaymentRepo.On("Create", mock.Anything, mock.MatchedBy(func(r *paylaterEntity.PaylaterRepayment) bool {
			return r.UserID == 1 && r.Amount == 500000 && r.PaymentSource == "wallet"
		})).Run(func(args mock.Arguments) {
			repayment := args.Get(1).(*paylaterEntity.PaylaterRepayment)
			repayment.ID = 1
		}).Return(nil)

		mockLedgerRepo.On("CreateEntry", mock.Anything, mock.MatchedBy(func(e *entity.LedgerEntry) bool {
			return e.UserID == 1 && e.Credit == 500000
		})).Return(nil)

		mockLedgerRepo.On("CreateEntry", mock.Anything, mock.MatchedBy(func(e *entity.LedgerEntry) bool {
			return e.UserID == 1 && e.Debit == 500000
		})).Return(nil)

		mockLoanRepo.On("UpdateLoanStatus", mock.Anything, uint(1), "paid").Return(nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.PaylaterLoanID)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(500000), result.Amount)
		assert.Equal(t, "wallet", result.PaymentSource)
		assert.Equal(t, "success", result.Status)
		mockLoanRepo.AssertExpectations(t)
		mockRepaymentRepo.AssertExpectations(t)
		mockLedgerRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Loan not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 999,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("loan not found"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "loan not found", err.Error())
		mockLoanRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Loan does not belong to user", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   2, // Different user
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("unauthorized repayment"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "unauthorized repayment", err.Error())
		mockLoanRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Loan is not active or overdue", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusPaid,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("cannot repay loan with status: paid"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Repayment amount is zero or negative", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         -100000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("repayment amount must be greater than 0"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Failed to create repayment record", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("failed to create repayment: database error"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockRepaymentRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
		mockRepaymentRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Failed to create cash inflow ledger entry", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("failed to create cash inflow ledger entry: database error"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockRepaymentRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			repayment := args.Get(1).(*paylaterEntity.PaylaterRepayment)
			repayment.ID = 1
		}).Return(nil)
		mockLedgerRepo.On("CreateEntry", mock.Anything, mock.MatchedBy(func(e *entity.LedgerEntry) bool {
			return e.Credit == 500000
		})).Return(errors.New("database error"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
		mockRepaymentRepo.AssertExpectations(t)
		mockLedgerRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Failed to update loan status", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusActive,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("failed to update loan status: database error"))

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockRepaymentRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			repayment := args.Get(1).(*paylaterEntity.PaylaterRepayment)
			repayment.ID = 1
		}).Return(nil)
		mockLedgerRepo.On("CreateEntry", mock.Anything, mock.Anything).Return(nil)
		mockLoanRepo.On("UpdateLoanStatus", mock.Anything, uint(1), "paid").Return(errors.New("database error"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockLoanRepo.AssertExpectations(t)
		mockRepaymentRepo.AssertExpectations(t)
		mockLedgerRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Success: Process repayment with external payment source", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         250000,
			PaymentSource:  "external_payment",
		}

		loan := &paylaterEntity.PaylaterLoan{
			ID:       1,
			UserID:   1,
			Amount:   1000000,
			Interest: 100000,
			Total:    1100000,
			DueDate:  time.Now().AddDate(0, 0, 30),
			Source:   paylaterEntity.PaylaterLoanSourceCheckout,
			Status:   paylaterEntity.PaylaterLoanStatusOverdue,
		}

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)

		mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).Return(loan, nil)
		mockRepaymentRepo.On("Create", mock.Anything, mock.MatchedBy(func(r *paylaterEntity.PaylaterRepayment) bool {
			return r.PaymentSource == "external_payment" && r.Amount == 250000
		})).Run(func(args mock.Arguments) {
			repayment := args.Get(1).(*paylaterEntity.PaylaterRepayment)
			repayment.ID = 2
		}).Return(nil)

		mockLedgerRepo.On("CreateEntry", mock.Anything, mock.Anything).Return(nil)
		mockLoanRepo.On("UpdateLoanStatus", mock.Anything, uint(1), "paid").Return(nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.ProcessRepayment(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "external_payment", result.PaymentSource)
		assert.Equal(t, int64(250000), result.Amount)
		mockLoanRepo.AssertExpectations(t)
		mockRepaymentRepo.AssertExpectations(t)
		mockLedgerRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentService_GetRepaymentByID(t *testing.T) {
	t.Run("Success: Get repayment by ID", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayment := &paylaterEntity.PaylaterRepayment{
			ID:             1,
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockRepaymentRepo.On("GetByID", mock.Anything, uint(1)).Return(repayment, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.GetRepaymentByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.PaylaterLoanID)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(500000), result.Amount)
		assert.Equal(t, "wallet", result.PaymentSource)
		assert.Equal(t, "success", result.Status)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Error: Repayment not found", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		mockRepaymentRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.GetRepaymentByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: Get different repayment by ID", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayment := &paylaterEntity.PaylaterRepayment{
			ID:             5,
			PaylaterLoanID: 2,
			UserID:         3,
			Amount:         1000000,
			PaymentSource:  "external_payment",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockRepaymentRepo.On("GetByID", mock.Anything, uint(5)).Return(repayment, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		result, err := service.GetRepaymentByID(context.Background(), 5)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(5), result.ID)
		assert.Equal(t, uint(2), result.PaylaterLoanID)
		assert.Equal(t, uint(3), result.UserID)
		mockRepaymentRepo.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentService_ListRepaymentsByLoan(t *testing.T) {
	t.Run("Success: List all repayments for a loan", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayments := []paylaterEntity.PaylaterRepayment{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         300000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-1 * time.Hour),
			},
		}

		mockRepaymentRepo.On("GetByLoanID", mock.Anything, uint(1)).Return(repayments, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.ListRepaymentsByLoan(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.Equal(t, uint(1), results[0].ID)
		assert.Equal(t, uint(2), results[1].ID)
		assert.Equal(t, int64(300000), results[0].Amount)
		assert.Equal(t, int64(200000), results[1].Amount)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: List returns empty slice for loan with no repayments", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		mockRepaymentRepo.On("GetByLoanID", mock.Anything, uint(999)).Return([]paylaterEntity.PaylaterRepayment{}, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.ListRepaymentsByLoan(context.Background(), 999)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 0)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Error: Database error when listing repayments", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		mockRepaymentRepo.On("GetByLoanID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.ListRepaymentsByLoan(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, results)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: List multiple repayments for the same loan with different users", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayments := []paylaterEntity.PaylaterRepayment{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-2 * time.Hour),
			},
			{
				ID:             3,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         1000000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-4 * time.Hour),
			},
		}

		mockRepaymentRepo.On("GetByLoanID", mock.Anything, uint(1)).Return(repayments, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.ListRepaymentsByLoan(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 3)
		totalAmount := results[0].Amount + results[1].Amount + results[2].Amount
		assert.Equal(t, int64(2000000), totalAmount)
		mockRepaymentRepo.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentService_GetRepaymentsByUser(t *testing.T) {
	t.Run("Success: Get all repayments for a user", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayments := []paylaterEntity.PaylaterRepayment{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         300000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 2,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-1 * time.Hour),
			},
		}

		mockRepaymentRepo.On("GetByUserID", mock.Anything, uint(1)).Return(repayments, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.GetRepaymentsByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.Equal(t, uint(1), results[0].ID)
		assert.Equal(t, uint(2), results[1].ID)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: Empty slice for user with no repayments", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		mockRepaymentRepo.On("GetByUserID", mock.Anything, uint(999)).Return([]paylaterEntity.PaylaterRepayment{}, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.GetRepaymentsByUser(context.Background(), 999)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 0)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Error: Database error when getting repayments", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		mockRepaymentRepo.On("GetByUserID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.GetRepaymentsByUser(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, results)
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: User with repayments from multiple loans", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayments := []paylaterEntity.PaylaterRepayment{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-2 * time.Hour),
			},
			{
				ID:             3,
				PaylaterLoanID: 2,
				UserID:         1,
				Amount:         1000000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-4 * time.Hour),
			},
			{
				ID:             4,
				PaylaterLoanID: 3,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-6 * time.Hour),
			},
		}

		mockRepaymentRepo.On("GetByUserID", mock.Anything, uint(1)).Return(repayments, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.GetRepaymentsByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 4)
		// Verify loans from different loans
		loanIDs := make(map[uint]int)
		for _, r := range results {
			loanIDs[r.PaylaterLoanID]++
		}
		assert.Equal(t, 2, loanIDs[1]) // 2 repayments for loan 1
		assert.Equal(t, 1, loanIDs[2]) // 1 repayment for loan 2
		assert.Equal(t, 1, loanIDs[3]) // 1 repayment for loan 3
		mockRepaymentRepo.AssertExpectations(t)
	})

	t.Run("Success: Get repayments for different user", func(t *testing.T) {
		logger := zap.NewNop()
		mockRepaymentRepo := new(testutil.MockPaylaterRepaymentRepository)
		mockLoanRepo := new(testutil.MockPaylaterLoanRepository)
		mockLedgerRepo := new(testutil.MockLedgerRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		repayments := []paylaterEntity.PaylaterRepayment{
			{
				ID:             10,
				PaylaterLoanID: 5,
				UserID:         2,
				Amount:         750000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
		}

		mockRepaymentRepo.On("GetByUserID", mock.Anything, uint(2)).Return(repayments, nil)

		service := NewPaylaterRepaymentService(mockRepaymentRepo, mockLoanRepo, mockLedgerRepo, mockTxManager, logger)
		results, err := service.GetRepaymentsByUser(context.Background(), 2)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, uint(2), results[0].UserID)
		assert.Equal(t, uint(5), results[0].PaylaterLoanID)
		mockRepaymentRepo.AssertExpectations(t)
	})
}
