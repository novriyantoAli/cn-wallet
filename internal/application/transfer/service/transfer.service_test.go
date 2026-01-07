package service

import (
	"context"
	"errors"
	"testing"
	"time"

	ledgerDTO "github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	ledgerEntity "github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	paylaterDTO "github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	userDTO "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ============================================================================
// Mock Repositories and Helpers
// ============================================================================

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]entity.Transfer, error) {
	args := m.Called(ctx, targetUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListTransfers(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]entity.Transfer, int64, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.Transfer), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransferRepository) UpdateTransferStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTransferRepository) GetUserTransferStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTransferRepository) CancelTransfer(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*userEntity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *userEntity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *userEntity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]userEntity.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, filter *userDTO.UserFilter) ([]userEntity.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]userEntity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
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

// Mock PaylaterLoanRepository
type MockPaylaterLoanRepository struct {
	mock.Mock
}

func (m *MockPaylaterLoanRepository) CreateLoan(ctx context.Context, loan *paylaterEntity.PaylaterLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockPaylaterLoanRepository) GetLoanByID(ctx context.Context, id uint) (*paylaterEntity.PaylaterLoan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterLoan), args.Error(1)
}

func (m *MockPaylaterLoanRepository) GetLoansByUserID(ctx context.Context, userID uint) ([]paylaterEntity.PaylaterLoan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]paylaterEntity.PaylaterLoan), args.Error(1)
}

func (m *MockPaylaterLoanRepository) ListLoans(ctx context.Context, filters *paylaterDTO.ListPaylaterLoansRequest) ([]paylaterEntity.PaylaterLoan, int64, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]paylaterEntity.PaylaterLoan), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaylaterLoanRepository) UpdateLoanStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaylaterLoanRepository) UpdateLoan(ctx context.Context, loan *paylaterEntity.PaylaterLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockPaylaterLoanRepository) GetOverdueLoans(ctx context.Context, asOf time.Time) ([]paylaterEntity.PaylaterLoan, error) {
	args := m.Called(ctx, asOf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]paylaterEntity.PaylaterLoan), args.Error(1)
}

func (m *MockPaylaterLoanRepository) GetLoanStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockPaylaterLoanRepository) DeleteLoan(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock LedgerRepository
type MockLedgerRepository struct {
	mock.Mock
}

func (m *MockLedgerRepository) CreateEntry(ctx context.Context, entry *ledgerEntity.LedgerEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockLedgerRepository) GetEntryByID(ctx context.Context, id uint64) (*ledgerEntity.LedgerEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerEntity.LedgerEntry), args.Error(1)
}

func (m *MockLedgerRepository) GetEntriesByUserID(ctx context.Context, userID uint64) ([]ledgerEntity.LedgerEntry, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ledgerEntity.LedgerEntry), args.Error(1)
}

func (m *MockLedgerRepository) GetEntriesByReference(ctx context.Context, referenceType string, referenceID string) ([]ledgerEntity.LedgerEntry, error) {
	args := m.Called(ctx, referenceType, referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ledgerEntity.LedgerEntry), args.Error(1)
}

func (m *MockLedgerRepository) ListEntries(ctx context.Context, req *ledgerDTO.ListLedgerEntriesRequest) ([]ledgerEntity.LedgerEntry, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]ledgerEntity.LedgerEntry), args.Get(1).(int64), args.Error(2)
}

func (m *MockLedgerRepository) GetUserStats(ctx context.Context, userID uint64) (*ledgerDTO.LedgerStatsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerDTO.LedgerStatsResponse), args.Error(1)
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

// ============================================================================
// Setup Helpers
// ============================================================================

// setupTransferServiceWithMocks creates a TransferService with mocked dependencies
func setupTransferServiceWithMocks() (TransferService, *MockTransferRepository, *MockUserRepository, *MockWalletRepository, *MockPaylaterAccountRepository, *MockPaylaterLoanRepository, *MockLedgerRepository, *MockTransactionManager) {
	mockTransferRepo := new(MockTransferRepository)
	mockUserRepo := new(MockUserRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockPaylaterRepo := new(MockPaylaterAccountRepository)
	mockPaylaterLoanRepo := new(MockPaylaterLoanRepository)
	mockLedgerRepo := new(MockLedgerRepository)
	mockTxManager := new(MockTransactionManager)
	logger := testutil.NewSilentLogger()

	service := NewTransferService(mockTransferRepo, mockUserRepo, mockWalletRepo, mockPaylaterRepo, mockPaylaterLoanRepo, mockLedgerRepo, mockTxManager, logger)

	return service, mockTransferRepo, mockUserRepo, mockWalletRepo, mockPaylaterRepo, mockPaylaterLoanRepo, mockLedgerRepo, mockTxManager
}

// ============================================================================
// Test Cases
// ============================================================================

func TestTransferService_CreateTransfer(t *testing.T) {
	ctx := context.Background()

	t.Run("should create transfer from wallet successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, mockUserRepo, mockWalletRepo, _, _, _, mockTxManager := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		receiverUser := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: true,
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

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(receiverUser, nil)
		mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, uint(1)).Return(senderWallet, nil)
		mockWalletRepo.On("UpdateBalance", ctx, uint(1), 400000.0).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, uint(2)).Return(receiverWallet, nil)
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
		mockUserRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should create transfer from paylater successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, mockUserRepo, mockWalletRepo, mockPaylaterRepo, mockPaylaterLoanRepo, mockLedgerRepo, mockTxManager := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourcePaylater),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		receiverUser := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: true,
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

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(receiverUser, nil)
		mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockPaylaterRepo.On("GetForUpdate", ctx, uint(1)).Return(paylaterAccount, nil)
		mockPaylaterLoanRepo.On("CreateLoan", ctx, mock.MatchedBy(func(loan *paylaterEntity.PaylaterLoan) bool {
			return loan != nil && loan.UserID == 1 && loan.Amount == 100000
		})).Return(nil)
		// Match any PaylaterAccount for update
		mockPaylaterRepo.On("UpdateAccount", ctx, mock.MatchedBy(func(acc *paylaterEntity.PaylaterAccount) bool {
			return acc != nil && acc.UserID == 1
		})).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, uint(2)).Return(receiverWallet, nil)
		mockWalletRepo.On("UpdateBalance", ctx, uint(2), 300000.0).Return(nil)
		mockLedgerRepo.On("CreateEntry", ctx, mock.MatchedBy(func(entry *ledgerEntity.LedgerEntry) bool {
			return entry != nil
		})).Return(nil)
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
		mockUserRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
		mockPaylaterRepo.AssertExpectations(t)
		mockPaylaterLoanRepo.AssertExpectations(t)
		mockLedgerRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should reject transfer to same user", func(t *testing.T) {
		// Setup
		service, _, _, _, _, _, _, _ := setupTransferServiceWithMocks()

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
	})

	t.Run("should reject zero amount", func(t *testing.T) {
		// Setup
		service, _, _, _, _, _, _, _ := setupTransferServiceWithMocks()

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
	})

	t.Run("should reject negative amount", func(t *testing.T) {
		// Setup
		service, _, _, _, _, _, _, _ := setupTransferServiceWithMocks()

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
	})

	t.Run("should reject transfer when sender user not found", func(t *testing.T) {
		// Setup
		service, _, mockUserRepo, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       999,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Contains(t, err.Error(), "failed to get user")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should reject transfer when sender user is inactive", func(t *testing.T) {
		// Setup
		service, _, mockUserRepo, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		inactiveUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: false,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(inactiveUser, nil)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Equal(t, "user is not active", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should reject transfer when receiver user not found", func(t *testing.T) {
		// Setup
		service, _, mockUserRepo, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 999,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Contains(t, err.Error(), "failed to get target user")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should reject transfer when receiver user is inactive", func(t *testing.T) {
		// Setup
		service, _, mockUserRepo, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		inactiveReceiver := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: false,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(inactiveReceiver, nil)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Equal(t, "target user is not active", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should reject transfer when receiver wallet not found", func(t *testing.T) {
		// Setup
		service, _, mockUserRepo, mockWalletRepo, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		receiverUser := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: true,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(receiverUser, nil)
		mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(nil, gorm.ErrRecordNotFound)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Contains(t, err.Error(), "failed to get target user wallet")

		mockUserRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("should reject transfer when wallet balance insufficient", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, mockUserRepo, mockWalletRepo, _, _, _, mockTxManager := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourceWallet),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		receiverUser := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: true,
		}

		senderWallet := &walletEntity.Wallet{
			ID:      1,
			UserID:  1,
			Balance: 50000,
		}

		receiverWallet := &walletEntity.Wallet{
			ID:      2,
			UserID:  2,
			Balance: 200000,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(receiverUser, nil)
		mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, uint(1)).Return(senderWallet, nil)
		mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Equal(t, "insufficient wallet balance", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should reject transfer when paylater credit insufficient", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, mockUserRepo, mockWalletRepo, mockPaylaterRepo, _, _, mockTxManager := setupTransferServiceWithMocks()

		req := &dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       string(entity.TransferSourcePaylater),
		}

		senderUser := &userEntity.User{
			ID:       1,
			Email:    "sender@test.com",
			FullName: "Sender",
			IsActive: true,
		}

		receiverUser := &userEntity.User{
			ID:       2,
			Email:    "receiver@test.com",
			FullName: "Receiver",
			IsActive: true,
		}

		paylaterAccount := &paylaterEntity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    950000,
			AvailableLimit: 50000,
		}

		receiverWallet := &walletEntity.Wallet{
			ID:      2,
			UserID:  2,
			Balance: 200000,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, uint(1)).Return(senderUser, nil)
		mockUserRepo.On("GetByID", ctx, uint(2)).Return(receiverUser, nil)
		mockWalletRepo.On("GetWalletByUserID", ctx, uint(2)).Return(receiverWallet, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockPaylaterRepo.On("GetForUpdate", ctx, uint(1)).Return(paylaterAccount, nil)
		mockTransferRepo.On("CreateTransfer", ctx, mock.AnythingOfType("*entity.Transfer")).Return(nil)

		// When
		transfer, err := service.CreateTransfer(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transfer)
		assert.Equal(t, "insufficient paylater credit", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
		mockPaylaterRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestTransferService_GetTransferByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get transfer by id successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

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

		// Mock expectations
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
	})

	t.Run("should return error when transfer not found", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(999)

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetTransferByID(ctx, transferID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_GetTransfersByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get transfers by user id successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(1)
		transfers := []entity.Transfer{
			{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
			{ID: 2, UserID: 1, TargetUserID: 3, Amount: 200000, Source: entity.TransferSourcePaylater, Status: entity.TransferStatusCompleted},
		}

		// Mock expectations
		mockTransferRepo.On("GetTransfersByUserID", ctx, userID).Return(transfers, nil)

		// When
		responses, err := service.GetTransfersByUserID(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, responses, 2)
		assert.Equal(t, transfers[0].ID, responses[0].ID)
		assert.Equal(t, transfers[1].ID, responses[1].ID)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should return empty list when no transfers found", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(999)

		// Mock expectations
		mockTransferRepo.On("GetTransfersByUserID", ctx, userID).Return([]entity.Transfer{}, nil)

		// When
		responses, err := service.GetTransfersByUserID(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, responses, 0)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(1)

		// Mock expectations
		mockTransferRepo.On("GetTransfersByUserID", ctx, userID).Return(nil, errors.New("database error"))

		// When
		responses, err := service.GetTransfersByUserID(ctx, userID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, responses)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_GetTransfersByTargetUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get transfers by target user id successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		targetUserID := uint(2)
		transfers := []entity.Transfer{
			{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
			{ID: 3, UserID: 3, TargetUserID: 2, Amount: 150000, Source: entity.TransferSourceWallet, Status: entity.TransferStatusCompleted},
		}

		// Mock expectations
		mockTransferRepo.On("GetTransfersByTargetUserID", ctx, targetUserID).Return(transfers, nil)

		// When
		responses, err := service.GetTransfersByTargetUserID(ctx, targetUserID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, responses, 2)
		assert.Equal(t, targetUserID, responses[0].TargetUserID)
		assert.Equal(t, targetUserID, responses[1].TargetUserID)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should return empty list when no transfers received", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		targetUserID := uint(999)

		// Mock expectations
		mockTransferRepo.On("GetTransfersByTargetUserID", ctx, targetUserID).Return([]entity.Transfer{}, nil)

		// When
		responses, err := service.GetTransfersByTargetUserID(ctx, targetUserID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, responses, 0)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_ListTransfers(t *testing.T) {
	ctx := context.Background()

	t.Run("should list transfers with pagination successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(1)
		req := &dto.TransferFilter{
			UserID:   &userID,
			Page:     1,
			PageSize: 10,
		}

		transfers := []entity.Transfer{
			{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000},
			{ID: 2, UserID: 1, TargetUserID: 3, Amount: 200000},
		}

		// Mock expectations
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
	})

	t.Run("should use default pagination when not provided", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.TransferFilter{
			Page:     0,
			PageSize: 0,
		}

		transfers := []entity.Transfer{}

		// Mock expectations
		mockTransferRepo.On("ListTransfers", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(transfers, int64(0), nil)

		// When
		response, err := service.ListTransfers(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should calculate total pages correctly", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		req := &dto.TransferFilter{
			Page:     1,
			PageSize: 10,
		}

		transfers := make([]entity.Transfer, 0)

		// Mock expectations
		mockTransferRepo.On("ListTransfers", ctx, mock.AnythingOfType("map[string]interface {}"), 1, 10).Return(transfers, int64(25), nil)

		// When
		response, err := service.ListTransfers(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 3, response.TotalPages)
		assert.Equal(t, int64(25), response.TotalCount)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_UpdateTransferStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("should update transfer status successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

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

		// Mock expectations
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
	})

	t.Run("should return error when transfer not found", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(999)
		req := &dto.UpdateTransferStatusRequest{
			Status: string(entity.TransferStatusCompleted),
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateTransferStatus(ctx, transferID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should handle update status error", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(1)
		req := &dto.UpdateTransferStatusRequest{
			Status: string(entity.TransferStatusFailed),
		}

		transfer := &entity.Transfer{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       entity.TransferSourceWallet,
			Status:       entity.TransferStatusPending,
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)
		mockTransferRepo.On("UpdateTransferStatus", ctx, transferID, req.Status).Return(errors.New("database error"))

		// When
		response, err := service.UpdateTransferStatus(ctx, transferID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_GetUserTransferStats(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user transfer stats successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(1)
		stats := map[string]interface{}{
			"total_sent":     int64(500000),
			"total_received": int64(300000),
			"count_sent":     int64(5),
			"count_received": int64(3),
			"wallet_sent":    int64(300000),
			"paylater_sent":  int64(200000),
			"completed_sent": int64(4),
			"pending_sent":   int64(1),
			"failed_sent":    int64(0),
		}

		// Mock expectations
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
		assert.Equal(t, int64(300000), response.WalletSent)
		assert.Equal(t, int64(200000), response.PaylaterSent)
		assert.Equal(t, int64(4), response.CompletedSent)
		assert.Equal(t, int64(1), response.PendingSent)
		assert.Equal(t, int64(0), response.FailedSent)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should return error when getting stats fails", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		userID := uint(1)

		// Mock expectations
		mockTransferRepo.On("GetUserTransferStats", ctx, userID).Return(nil, errors.New("database error"))

		// When
		response, err := service.GetUserTransferStats(ctx, userID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)

		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_CancelTransfer(t *testing.T) {
	ctx := context.Background()

	t.Run("should cancel pending transfer successfully", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(1)
		transfer := &entity.Transfer{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       entity.TransferSourceWallet,
			Status:       entity.TransferStatusPending,
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)
		mockTransferRepo.On("CancelTransfer", ctx, transferID).Return(nil)

		// When
		err := service.CancelTransfer(ctx, transferID)

		// Then
		assert.NoError(t, err)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should reject cancelling completed transfer", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(1)
		transfer := &entity.Transfer{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       entity.TransferSourceWallet,
			Status:       entity.TransferStatusCompleted,
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)

		// When
		err := service.CancelTransfer(ctx, transferID)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "can only cancel pending transfers", err.Error())

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should reject cancelling failed transfer", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(1)
		transfer := &entity.Transfer{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       entity.TransferSourceWallet,
			Status:       entity.TransferStatusFailed,
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)

		// When
		err := service.CancelTransfer(ctx, transferID)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "can only cancel pending transfers", err.Error())

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should return error when transfer not found", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(999)

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.CancelTransfer(ctx, transferID)

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("should handle cancel error from repository", func(t *testing.T) {
		// Setup
		service, mockTransferRepo, _, _, _, _, _, _ := setupTransferServiceWithMocks()

		transferID := uint(1)
		transfer := &entity.Transfer{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       entity.TransferSourceWallet,
			Status:       entity.TransferStatusPending,
		}

		// Mock expectations
		mockTransferRepo.On("GetTransferByID", ctx, transferID).Return(transfer, nil)
		mockTransferRepo.On("CancelTransfer", ctx, transferID).Return(errors.New("database error"))

		// When
		err := service.CancelTransfer(ctx, transferID)

		// Then
		assert.Error(t, err)

		mockTransferRepo.AssertExpectations(t)
	})
}
