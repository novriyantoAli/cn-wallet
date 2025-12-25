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

func TestPaylaterAccountService_CreateAccount(t *testing.T) {
	t.Run("Success: Create new paylater account", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		req := &dto.CreatePaylaterAccountRequest{
			UserID:      1,
			CreditLimit: 1000000,
		}

		expectedAccount := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    0,
			AvailableLimit: 1000000,
			Status:         entity.PaylaterStatusActive,
			CreatedAt:      time.Now(),
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
		mockRepo.On("CreateAccount", mock.Anything, mock.MatchedBy(func(a *entity.PaylaterAccount) bool {
			return a.UserID == 1 && a.CreditLimit == 1000000
		})).Return(nil).Run(func(args mock.Arguments) {
			account := args.Get(1).(*entity.PaylaterAccount)
			account.ID = 1
			account.CreatedAt = expectedAccount.CreatedAt
		})

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.CreateAccount(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(1000000), result.CreditLimit)
		assert.Equal(t, int64(0), result.Outstanding)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Account already exists", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		req := &dto.CreatePaylaterAccountRequest{
			UserID:      1,
			CreditLimit: 1000000,
		}

		existingAccount := &entity.PaylaterAccount{
			ID:     1,
			UserID: 1,
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(existingAccount, nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.CreateAccount(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "paylater account already exists for this user", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_GetAccountByUserID(t *testing.T) {
	t.Run("Success: Get account by user ID", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
			CreatedAt:      time.Now(),
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(account, nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.GetAccountByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.UserID)
		assert.Equal(t, int64(1000000), result.CreditLimit)
		assert.Equal(t, int64(500000), result.Outstanding)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.GetAccountByUserID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_GetAccountByID(t *testing.T) {
	t.Run("Success: Get account by ID", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Status:         entity.PaylaterStatusActive,
			CreatedAt:      time.Now(),
		}

		mockRepo.On("GetAccountByID", mock.Anything, uint(1)).Return(account, nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.GetAccountByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		mockRepo.On("GetAccountByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.GetAccountByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_UpdateCreditLimit(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Success: Update credit limit", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.UpdateCreditLimitRequest{
			CreditLimit: 2000000,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)
		mockRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(a *entity.PaylaterAccount) bool {
			return a.CreditLimit == 2000000
		})).Run(func(args mock.Arguments) {
			acc := args.Get(1).(*entity.PaylaterAccount)
			acc.CreditLimit = 2000000
			acc.AvailableLimit = 1500000
		}).Return(nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.UpdateCreditLimit(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Credit limit less than outstanding", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.UpdateCreditLimitRequest{
			CreditLimit: 300000, // Less than outstanding 500000
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("credit limit cannot be less than outstanding balance"))

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.UpdateCreditLimit(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_UpdateStatus(t *testing.T) {
	t.Run("Success: Update status to suspended", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		req := &dto.UpdateStatusRequest{
			Status: string(entity.PaylaterStatusSuspended),
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(account, nil)
		mockRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(a *entity.PaylaterAccount) bool {
			return a.Status == entity.PaylaterStatusSuspended
		})).Run(func(args mock.Arguments) {
			acc := args.Get(1).(*entity.PaylaterAccount)
			acc.Status = entity.PaylaterStatusSuspended
		}).Return(nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.UpdateStatus(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, string(entity.PaylaterStatusSuspended), result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		req := &dto.UpdateStatusRequest{
			Status: string(entity.PaylaterStatusSuspended),
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		result, err := service.UpdateStatus(context.Background(), 999, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_UseCredit(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Success: Use credit", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.UseCreditRequest{
			Amount: 200000,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)
		mockRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(a *entity.PaylaterAccount) bool {
			return a.Outstanding == 500000
		})).Run(func(args mock.Arguments) {
			acc := args.Get(1).(*entity.PaylaterAccount)
			acc.Outstanding = 500000
			acc.AvailableLimit = 500000
		}).Return(nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.UseCredit(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(200000), result.Amount)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Account not active", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.UseCreditRequest{
			Amount: 200000,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Status:         entity.PaylaterStatusSuspended,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("account is not active"))

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.UseCredit(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Insufficient credit limit", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.UseCreditRequest{
			Amount: 800000, // More than available 700000
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("insufficient credit limit"))

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.UseCredit(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_Repayment(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Success: Process repayment", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.RepaymentRequest{
			Amount: 200000,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)
		mockRepo.On("UpdateAccount", mock.Anything, mock.MatchedBy(func(a *entity.PaylaterAccount) bool {
			return a.Outstanding == 300000
		})).Run(func(args mock.Arguments) {
			acc := args.Get(1).(*entity.PaylaterAccount)
			acc.Outstanding = 300000
			acc.AvailableLimit = 700000
		}).Return(nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.Repayment(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(200000), result.Amount)
		assert.Equal(t, int64(300000), result.Outstanding)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Error: Repayment exceeds outstanding", func(t *testing.T) {
		mockRepo := new(testutil.MockPaylaterAccountRepository)
		mockTxManager := new(testutil.MockTransactionManager)

		req := &dto.RepaymentRequest{
			Amount: 600000, // More than outstanding 500000
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetForUpdate", mock.Anything, uint(1)).Return(account, nil)

		mockTxManager.On("WithinTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(errors.New("repayment amount cannot exceed outstanding balance"))

		service := NewPaylaterAccountService(mockRepo, mockTxManager, logger)
		result, err := service.Repayment(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestPaylaterAccountService_DeleteAccount(t *testing.T) {
	t.Run("Success: Delete account with zero outstanding", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    0,
			AvailableLimit: 1000000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(account, nil)
		mockRepo.On("DeleteAccount", mock.Anything, uint(1)).Return(nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		err := service.DeleteAccount(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Delete account with outstanding balance", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         entity.PaylaterStatusActive,
		}

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(1)).Return(account, nil)

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		err := service.DeleteAccount(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, "cannot delete account with outstanding balance", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		logger := zap.NewNop()
		txManager := new(testutil.MockTransactionManager)
		mockRepo := new(testutil.MockPaylaterAccountRepository)

		mockRepo.On("GetAccountByUserID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		service := NewPaylaterAccountService(mockRepo, txManager, logger)
		err := service.DeleteAccount(context.Background(), 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
