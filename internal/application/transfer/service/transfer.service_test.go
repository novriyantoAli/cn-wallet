package service

import (
	"context"
	"errors"
	"testing"
	"time"

	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Mock TransferRepository
type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) CreateTransfer(ctx context.Context, transfer *entity.Transfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) GetTransferByID(ctx context.Context, id uint) (*entity.Transfer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetTransfersByUserID(ctx context.Context, userID uint) ([]entity.Transfer, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]entity.Transfer, error) {
	args := m.Called(ctx, targetUserID)
	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListTransfers(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]entity.Transfer, int64, error) {
	args := m.Called(ctx, filters, page, pageSize)
	return args.Get(0).([]entity.Transfer), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransferRepository) UpdateTransferStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTransferRepository) GetUserTransferStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTransferRepository) CancelTransfer(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock WalletRepository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetWalletByUserID(ctx context.Context, userID uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, userID uint, balance float64) error {
	args := m.Called(ctx, userID, balance)
	return args.Error(0)
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, wallet *walletEntity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletByID(ctx context.Context, id uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetForUpdate(ctx context.Context, userID uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWallet(ctx context.Context, wallet *walletEntity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) DeleteWallet(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock PaylaterAccountRepository
type MockPaylaterAccountRepository struct {
	mock.Mock
}

func (m *MockPaylaterAccountRepository) GetForUpdate(ctx context.Context, userID uint) (*paylaterEntity.PaylaterAccount, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterAccount), args.Error(1)
}

func (m *MockPaylaterAccountRepository) UpdateAccount(ctx context.Context, account *paylaterEntity.PaylaterAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockPaylaterAccountRepository) GetByUserID(ctx context.Context, userID uint) (*paylaterEntity.PaylaterAccount, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterAccount), args.Error(1)
}

func (m *MockPaylaterAccountRepository) Create(ctx context.Context, account *paylaterEntity.PaylaterAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockPaylaterAccountRepository) CreateAccount(ctx context.Context, account *paylaterEntity.PaylaterAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockPaylaterAccountRepository) GetAccountByUserID(ctx context.Context, userID uint) (*paylaterEntity.PaylaterAccount, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterAccount), args.Error(1)
}

func (m *MockPaylaterAccountRepository) GetAccountByID(ctx context.Context, id uint) (*paylaterEntity.PaylaterAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterAccount), args.Error(1)
}

func (m *MockPaylaterAccountRepository) UpdateCreditLimit(ctx context.Context, userID uint, creditLimit int64) error {
	args := m.Called(ctx, userID, creditLimit)
	return args.Error(0)
}

func (m *MockPaylaterAccountRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	args := m.Called(ctx, userID, status)
	return args.Error(0)
}

func (m *MockPaylaterAccountRepository) DeleteAccount(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Error(0) == nil {
		// Execute the function if no error is expected
		return fn(ctx)
	}
	return args.Error(0)
}

func TestTransferService_CreateTransfer_WalletSource_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       string(entity.TransferSourceWallet),
	}

	senderWallet := &walletEntity.Wallet{
		ID:      1,
		UserID:  1,
		Balance: 500000,
	}

	receiverWallet := &walletEntity.Wallet{
		ID:      2,
		UserID:  2,
		Balance: 200000,
	}

	// Expectations
	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(1)).Return(senderWallet, nil)
	mockWalletRepo.On("UpdateBalance", ctx, uint(1), 400000.0).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
	mockWalletRepo.On("UpdateBalance", ctx, uint(2), 300000.0).Return(nil)
	mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, transfer)
	assert.Equal(t, uint(1), transfer.UserID)
	assert.Equal(t, uint(2), transfer.TargetUserID)
	assert.Equal(t, int64(100000), transfer.Amount)
	assert.Equal(t, entity.TransferSourceWallet, transfer.Source)
	assert.Equal(t, entity.TransferStatusCompleted, transfer.Status)

	mockTransferRepo.AssertExpectations(t)
	mockWalletRepo.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_PaylaterSource_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       string(entity.TransferSourcePaylater),
	}

	paylaterAccount := &paylaterEntity.PaylaterAccount{
		ID:             1,
		UserID:         1,
		CreditLimit:    1000000,
		Outstanding:    200000,
		AvailableLimit: 800000,
	}

	receiverWallet := &walletEntity.Wallet{
		ID:      2,
		UserID:  2,
		Balance: 200000,
	}

	// Expectations
	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPaylaterRepo.On("GetForUpdate", ctx, uint(1)).Return(paylaterAccount, nil)
	mockPaylaterRepo.On("UpdateAccount", ctx, mock.AnythingOfType("*entity.PaylaterAccount")).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
	mockWalletRepo.On("UpdateBalance", ctx, uint(2), 300000.0).Return(nil)
	mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, transfer)
	assert.Equal(t, uint(1), transfer.UserID)
	assert.Equal(t, uint(2), transfer.TargetUserID)
	assert.Equal(t, int64(100000), transfer.Amount)
	assert.Equal(t, entity.TransferSourcePaylater, transfer.Source)
	assert.Equal(t, entity.TransferStatusCompleted, transfer.Status)

	mockTransferRepo.AssertExpectations(t)
	mockWalletRepo.AssertExpectations(t)
	mockPaylaterRepo.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_SameUser_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 1,
		Amount:       100000,
		Source:       string(entity.TransferSourceWallet),
	}

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "cannot transfer to yourself", err.Error())
}

func TestTransferService_CreateTransfer_ZeroAmount_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       0,
		Source:       string(entity.TransferSourceWallet),
	}

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "transfer amount must be greater than 0", err.Error())
}

func TestTransferService_CreateTransfer_NegativeAmount_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       -100000,
		Source:       string(entity.TransferSourceWallet),
	}

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "transfer amount must be greater than 0", err.Error())
}

func TestTransferService_CreateTransfer_InsufficientWalletBalance_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       string(entity.TransferSourceWallet),
	}

	senderWallet := &walletEntity.Wallet{
		ID:      1,
		UserID:  1,
		Balance: 50000, // Less than requested amount
	}

	// Expectations
	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(1)).Return(senderWallet, nil)
	mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "insufficient wallet balance", err.Error())

	mockWalletRepo.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_InsufficientPaylaterCredit_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       string(entity.TransferSourcePaylater),
	}

	paylaterAccount := &paylaterEntity.PaylaterAccount{
		ID:             1,
		UserID:         1,
		CreditLimit:    1000000,
		Outstanding:    950000,
		AvailableLimit: 50000, // Less than requested amount
	}

	// Expectations
	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPaylaterRepo.On("GetForUpdate", ctx, uint(1)).Return(paylaterAccount, nil)
	mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "insufficient paylater credit", err.Error())

	mockPaylaterRepo.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_ReceiverWalletNotFound_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.CreateTransferRequest{
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       string(entity.TransferSourceWallet),
	}

	senderWallet := &walletEntity.Wallet{
		ID:      1,
		UserID:  1,
		Balance: 500000,
	}

	// Expectations
	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(1)).Return(senderWallet, nil)
	mockWalletRepo.On("UpdateBalance", ctx, uint(1), 400000.0).Return(nil)
	mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(nil, nil)
	mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

	// When
	transfer, err := service.CreateTransfer(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.Equal(t, "receiver wallet not found", err.Error())

	mockWalletRepo.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestTransferService_GetTransferByID_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(1)

	transfer := &entity.Transfer{
		ID:           1,
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       entity.TransferSourceWallet,
		Status:       entity.TransferStatusCompleted,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)

	// When
	response, err := service.GetTransferByID(ctx, transferID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, transfer.ID, response.ID)
	assert.Equal(t, transfer.UserID, response.UserID)
	assert.Equal(t, transfer.TargetUserID, response.TargetUserID)
	assert.Equal(t, transfer.Amount, response.Amount)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_GetTransferByID_NotFound_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(999)

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

	// When
	response, err := service.GetTransferByID(ctx, transferID)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_GetTransfersByUserID_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	userID := uint(1)

	transfers := []entity.Transfer{
		{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
		{ID: 2, UserID: 1, TargetUserID: 3, Amount: 200000, Source: entity.TransferSourcePaylater, Status: entity.TransferStatusCompleted},
	}

	mockTransferRepo.On("GetTransfersByUserID", ctx, userID).Return(transfers, nil)

	// When
	responses, err := service.GetTransfersByUserID(ctx, userID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, transfers[0].ID, responses[0].ID)
	assert.Equal(t, transfers[1].ID, responses[1].ID)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_GetTransfersByTargetUserID_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	targetUserID := uint(2)

	transfers := []entity.Transfer{
		{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
		{ID: 3, UserID: 3, TargetUserID: 2, Amount: 150000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
	}

	mockTransferRepo.On("GetTransfersByTargetUserID", ctx, targetUserID).Return(transfers, nil)

	// When
	responses, err := service.GetTransfersByTargetUserID(ctx, targetUserID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, targetUserID, responses[0].TargetUserID)
	assert.Equal(t, targetUserID, responses[1].TargetUserID)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_ListTransfers_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	userID := uint(1)
	req := &dto.ListTransfersRequest{
		UserID:   &userID,
		Page:     1,
		PageSize: 10,
	}

	transfers := []entity.Transfer{
		{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000},
		{ID: 2, UserID: 1, TargetUserID: 3, Amount: 200000},
	}

	mockTransferRepo.On("ListTransfers", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(transfers, int64(2), nil)

	// When
	response, err := service.ListTransfers(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, int64(2), response.TotalCount)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, 1, response.TotalPages)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_ListTransfers_DefaultPagination(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	req := &dto.ListTransfersRequest{
		Page:     0, // Should default to 1
		PageSize: 0, // Should default to 10
	}

	transfers := []entity.Transfer{}

	mockTransferRepo.On("ListTransfers", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(transfers, int64(0), nil)

	// When
	response, err := service.ListTransfers(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_UpdateTransferStatus_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(1)
	req := &dto.UpdateTransferStatusRequest{
		Status: string(entity.TransferStatusCompleted),
	}

	transfer := &entity.Transfer{
		ID:           1,
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       entity.TransferSourceWallet,
		Status:       entity.TransferStatusPending,
	}

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)
	mockTransferRepo.On("UpdateTransferStatus", ctx, transferID, req.Status).Return(nil)

	// When
	response, err := service.UpdateTransferStatus(ctx, transferID, req)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, transferID, response.ID)
	assert.Equal(t, req.Status, response.Status)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_UpdateTransferStatus_NotFound_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(999)
	req := &dto.UpdateTransferStatusRequest{
		Status: string(entity.TransferStatusCompleted),
	}

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

	// When
	response, err := service.UpdateTransferStatus(ctx, transferID, req)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_GetUserTransferStats_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	userID := uint(1)

	stats := map[string]interface{}{
		"total_sent":     int64(500000),
		"total_received": int64(300000),
		"count_sent":     int64(5),
		"count_received": int64(3),
		"wallet_sent":    int64(300000),
		"paylater_sent":  int64(200000),
		"completed_sent": int64(400000),
		"pending_sent":   int64(1),
		"failed_sent":    int64(0),
	}

	mockTransferRepo.On("GetUserTransferStats", ctx, userID).Return(stats, nil)

	// When
	response, err := service.GetUserTransferStats(ctx, userID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, int64(500000), response.TotalSent)
	assert.Equal(t, int64(300000), response.TotalReceived)
	assert.Equal(t, int64(5), response.CountSent)
	assert.Equal(t, int64(3), response.CountReceived)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_GetUserTransferStats_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	userID := uint(1)

	mockTransferRepo.On("GetUserTransferStats", ctx, userID).Return(map[string]interface{}{}, errors.New("database error"))

	// When
	response, err := service.GetUserTransferStats(ctx, userID)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_CancelTransfer_Success(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(1)

	transfer := &entity.Transfer{
		ID:           1,
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       entity.TransferSourceWallet,
		Status:       entity.TransferStatusPending,
	}

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)
	mockTransferRepo.On("CancelTransfer", ctx, transferID).Return(nil)

	// When
	err := service.CancelTransfer(ctx, transferID)

	// Then
	assert.NoError(t, err)

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_CancelTransfer_NotPending_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(1)

	transfer := &entity.Transfer{
		ID:           1,
		UserID:       1,
		TargetUserID: 2,
		Amount:       100000,
		Source:       entity.TransferSourceWallet,
		Status:       entity.TransferStatusCompleted, // Already completed
	}

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)

	// When
	err := service.CancelTransfer(ctx, transferID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "can only cancel pending transfers", err.Error())

	mockTransferRepo.AssertExpectations(t)
}

func TestTransferService_CancelTransfer_NotFound_Error(t *testing.T) {
	// Given
	mockTransferRepo := new(MockTransferRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockTxManager := new(MockTransactionManager)
	logger := zap.NewNop()

	service := NewTransferService(mockTransferRepo, mockWalletRepo, mockPaylaterRepo, mockTxManager, logger)

	ctx := context.Background()
	transferID := uint(999)

	mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

	// When
	err := service.CancelTransfer(ctx, transferID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockTransferRepo.AssertExpectations(t)
}
