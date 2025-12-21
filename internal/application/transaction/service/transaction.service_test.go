package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupTransactionServiceWithMocks() (TransactionService, *testutil.MockTransactionRepository) {
	mockRepo := &testutil.MockTransactionRepository{}
	logger := testutil.NewSilentLogger()
	service := NewTransactionService(mockRepo, logger)
	return service, mockRepo
}

func setupTransactionService(mockRepo *testutil.MockTransactionRepository) TransactionService {
	logger := testutil.NewSilentLogger()
	service := NewTransactionService(mockRepo, logger)
	return service
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should create transaction successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		req := testutil.CreateTransactionRequestFixture()

		// Mock expectations
		mockRepo.On("Create", ctx, mock.MatchedBy(func(txn *entity.Transaction) bool {
			return txn.WalletID == req.WalletID && txn.Amount == req.Amount && txn.Type == req.Type
		})).Return(nil).Run(func(args mock.Arguments) {
			txn := args.Get(1).(*entity.Transaction)
			txn.ID = uuid.New()
		})

		// When
		result, err := service.CreateTransaction(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.WalletID, result.WalletID)
		assert.Equal(t, req.Type, result.Type)
		assert.Equal(t, req.Amount, result.Amount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when amount is zero or negative", func(t *testing.T) {
		// Setup
		service, _ := setupTransactionServiceWithMocks()
		req := testutil.CreateTransactionRequestFixture()
		req.Amount = 0

		// When
		result, err := service.CreateTransaction(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should fail when repository create fails", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		req := testutil.CreateTransactionRequestFixture()

		mockRepo.On("Create", ctx, mock.MatchedBy(func(txn *entity.Transaction) bool {
			return true
		})).Return(errors.New("database error"))

		// When
		result, err := service.CreateTransaction(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetTransactionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get transaction by ID successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()

		mockRepo.On("GetByID", ctx, transaction.ID).Return(transaction, nil)

		// When
		result, err := service.GetTransactionByID(ctx, transaction.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, transaction.ID, result.ID)
		assert.Equal(t, transaction.WalletID, result.WalletID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when transaction not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transactionID := uuid.New()

		mockRepo.On("GetByID", ctx, transactionID).Return(nil, gorm.ErrRecordNotFound)

		// When
		result, err := service.GetTransactionByID(ctx, transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetAllTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all transactions with pagination", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transactions := []entity.Transaction{
			*testutil.CreateTransactionFixture(),
			*testutil.CreateTransactionFixture(),
		}
		transactions[0].ID = uuid.New()
		transactions[1].ID = uuid.New()

		filter := &dto.TransactionFilter{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", ctx, filter).Return(transactions, int64(2), nil)

		// When
		result, err := service.GetAllTransactions(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.TotalCount)
		assert.Len(t, result.Data, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should filter transactions by wallet ID", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()
		transaction.WalletID = 1

		filter := &dto.TransactionFilter{
			WalletID: 1,
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", ctx, filter).Return([]entity.Transaction{*transaction}, int64(1), nil)

		// When
		result, err := service.GetAllTransactions(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.TotalCount)
		assert.Len(t, result.Data, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()

		filter := &dto.TransactionFilter{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", ctx, filter).Return(nil, int64(0), errors.New("database error"))

		// When
		result, err := service.GetAllTransactions(ctx, filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should update transaction successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()

		updateReq := testutil.CreateUpdateTransactionRequestFixture()

		mockRepo.On("GetByID", ctx, transaction.ID).Return(transaction, nil)

		mockRepo.On("Update", ctx, mock.MatchedBy(func(txn *entity.Transaction) bool {
			return txn.ID == transaction.ID && txn.Status == updateReq.Status
		})).Return(nil)

		// When
		result, err := service.UpdateTransaction(ctx, transaction.ID, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, transaction.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when transaction not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transactionID := uuid.New()
		updateReq := testutil.CreateUpdateTransactionRequestFixture()

		mockRepo.On("GetByID", ctx, transactionID).Return(nil, gorm.ErrRecordNotFound)

		// When
		result, err := service.UpdateTransaction(ctx, transactionID, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when update fails", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()
		updateReq := testutil.CreateUpdateTransactionRequestFixture()

		mockRepo.On("GetByID", ctx, transaction.ID).Return(transaction, nil)

		mockRepo.On("Update", ctx, mock.MatchedBy(func(txn *entity.Transaction) bool {
			return true
		})).Return(errors.New("database error"))

		// When
		result, err := service.UpdateTransaction(ctx, transaction.ID, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete transaction successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()

		mockRepo.On("GetByID", ctx, transaction.ID).Return(transaction, nil)

		mockRepo.On("Delete", ctx, transaction.ID).Return(nil)

		// When
		err := service.DeleteTransaction(ctx, transaction.ID)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when transaction not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transactionID := uuid.New()

		mockRepo.On("GetByID", ctx, transactionID).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.DeleteTransaction(ctx, transactionID)

		// Then
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when delete fails", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transaction := testutil.CreateTransactionFixture()
		transaction.ID = uuid.New()

		mockRepo.On("GetByID", ctx, transaction.ID).Return(transaction, nil)

		mockRepo.On("Delete", ctx, transaction.ID).Return(errors.New("database error"))

		// When
		err := service.DeleteTransaction(ctx, transaction.ID)

		// Then
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetWalletTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("should get wallet transactions with pagination", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()
		transactions := []entity.Transaction{
			*testutil.CreateTransactionFixture(),
			*testutil.CreateTransactionFixture(),
		}
		transactions[0].ID = uuid.New()
		transactions[0].WalletID = 1
		transactions[1].ID = uuid.New()
		transactions[1].WalletID = 1

		mockRepo.On("GetByWalletID", ctx, uint(1), 1, 10).Return(transactions, int64(2), nil)

		// When
		result, err := service.GetWalletTransactions(ctx, 1, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.TotalCount)
		assert.Len(t, result.Data, 2)
		for _, txn := range result.Data {
			assert.Equal(t, uint(1), txn.WalletID)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return empty list when wallet has no transactions", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()

		mockRepo.On("GetByWalletID", ctx, uint(999), 1, 10).Return([]entity.Transaction{}, int64(0), nil)

		// When
		result, err := service.GetWalletTransactions(ctx, 999, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.TotalCount)
		assert.Len(t, result.Data, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// Setup
		service, mockRepo := setupTransactionServiceWithMocks()

		mockRepo.On("GetByWalletID", ctx, uint(1), 1, 10).Return(nil, int64(0), errors.New("database error"))

		// When
		result, err := service.GetWalletTransactions(ctx, 1, 1, 10)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
