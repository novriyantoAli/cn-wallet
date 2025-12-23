package service

import (
	"context"
	"errors"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletService_CreateWallet(t *testing.T) {
	ctx := context.Background()

	t.Run("should create wallet successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		req := &dto.CreateWalletRequest{
			UserID: 1,
			PIN:    "123456",
		}

		// Mock expectations
		mockRepo.On("CreateWallet", ctx, mock.AnythingOfType("*entity.Wallet")).Return(nil).Run(func(args mock.Arguments) {
			wallet := args.Get(1).(*entity.Wallet)
			wallet.ID = 1
		})

		// When
		wallet, err := service.CreateWallet(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, uint(1), wallet.ID)
		assert.Equal(t, uint(1), wallet.UserID)
		assert.Equal(t, 0.00, wallet.Balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		req := &dto.CreateWalletRequest{
			UserID: 1,
			PIN:    "123456",
		}

		// Mock expectations
		mockRepo.On("CreateWallet", ctx, mock.AnythingOfType("*entity.Wallet")).Return(errors.New("database error"))

		// When
		wallet, err := service.CreateWallet(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, wallet)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_GetWalletByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get wallet by user ID successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		walletFixture := testutil.CreateWalletFixture()

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, walletFixture.UserID).Return(walletFixture, nil)

		// When
		response, err := service.GetWalletByUserID(ctx, walletFixture.UserID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, walletFixture.UserID, response.UserID)
		assert.Equal(t, walletFixture.Balance, response.Balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, uint(9999)).Return(nil, errors.New("wallet not found"))

		// When
		response, err := service.GetWalletByUserID(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_GetWalletByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get wallet by ID successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		walletFixture := testutil.CreateWalletFixture()

		// Mock expectations
		mockRepo.On("GetWalletByID", ctx, walletFixture.ID).Return(walletFixture, nil)

		// When
		response, err := service.GetWalletByID(ctx, walletFixture.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, walletFixture.ID, response.ID)
		assert.Equal(t, walletFixture.UserID, response.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when wallet ID not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("GetWalletByID", ctx, uint(9999)).Return(nil, errors.New("wallet not found"))

		// When
		response, err := service.GetWalletByID(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_AddBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("should add balance successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		walletFixture := testutil.CreateWalletFixture()
		amount := 100.00

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, walletFixture.UserID).Return(walletFixture, nil)
		mockRepo.On("UpdateWallet", ctx, mock.AnythingOfType("*entity.Wallet")).Return(nil)

		// When
		response, err := service.AddBalance(ctx, walletFixture.UserID, amount, "test deposit")

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		// The response wallet should have the updated balance (1000 + 100 = 1100)
		assert.Equal(t, 1100.00, response.Balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for negative amount", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// When
		response, err := service.AddBalance(ctx, 1, -100.00, "test")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "amount must be positive", err.Error())
	})

	t.Run("should return error for zero amount", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// When
		response, err := service.AddBalance(ctx, 1, 0.00, "test")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "amount must be positive", err.Error())
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, uint(9999)).Return(nil, errors.New("wallet not found"))

		// When
		response, err := service.AddBalance(ctx, 9999, 100.00, "test")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_DeductBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("should deduct balance successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		walletFixture := testutil.CreateWalletFixture()
		amount := 100.00

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, walletFixture.UserID).Return(walletFixture, nil)
		mockRepo.On("UpdateWallet", ctx, mock.AnythingOfType("*entity.Wallet")).Return(nil)

		// When
		response, err := service.DeductBalance(ctx, walletFixture.UserID, amount, "test withdrawal")

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		// The response wallet should have the updated balance (1000 - 100 = 900)
		assert.Equal(t, 900.00, response.Balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for insufficient balance", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		walletFixture := testutil.CreateWalletFixture()

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, walletFixture.UserID).Return(walletFixture, nil)

		// When
		response, err := service.DeductBalance(ctx, walletFixture.UserID, 2000.00, "test withdrawal")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "insufficient balance", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for negative amount", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// When
		response, err := service.DeductBalance(ctx, 1, -100.00, "test")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "amount must be positive", err.Error())
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("GetWalletByUserID", ctx, uint(9999)).Return(nil, errors.New("wallet not found"))

		// When
		response, err := service.DeductBalance(ctx, 9999, 100.00, "test")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_DeleteWallet(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete wallet successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("DeleteWallet", ctx, uint(1)).Return(nil)

		// When
		err := service.DeleteWallet(ctx, 1)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockWalletRepository{}
		logger := testutil.NewSilentLogger()
		service := NewWalletService(mockRepo, logger)

		// Mock expectations
		mockRepo.On("DeleteWallet", ctx, uint(1)).Return(errors.New("database error"))

		// When
		err := service.DeleteWallet(ctx, 1)

		// Then
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
