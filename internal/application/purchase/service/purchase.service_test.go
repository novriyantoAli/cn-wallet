package service

import (
	"context"
	"errors"
	"testing"

	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/dto"
	transactionDto "github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProviderClient is a mock implementation of ProviderClient
type MockProviderClient struct {
	mock.Mock
}

func (m *MockProviderClient) ProcessPurchase(ctx context.Context, productCode string, phone string, reference string) (serialNumber string, err error) {
	args := m.Called(ctx, productCode, phone, reference)
	return args.String(0), args.Error(1)
}

func TestPurchaseService_ProcessPurchase_RequestValidation(t *testing.T) {
	ctx := context.Background()

	t.Run("should return error when product_id is missing", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockProviderClient := &MockProviderClient{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			mockTransactionRepo,
			mockUserRepo,
			mockWifiVoucherRepo,
			mockProviderClient,
			mockTxManager,
			nil,
			logger,
		)

		req := &dto.PurchaseRequest{
			ProductID: 0,
			Phone:     "08123456789",
		}

		// When
		response, err := service.ProcessPurchase(ctx, "token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product_id is required", err.Error())
	})

	t.Run("should return error when phone is missing", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockProviderClient := &MockProviderClient{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			mockTransactionRepo,
			mockUserRepo,
			mockWifiVoucherRepo,
			mockProviderClient,
			mockTxManager,
			nil,
			logger,
		)

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "",
		}

		// When
		response, err := service.ProcessPurchase(ctx, "token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "phone is required", err.Error())
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockProviderClient := &MockProviderClient{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			mockTransactionRepo,
			mockUserRepo,
			mockWifiVoucherRepo,
			mockProviderClient,
			mockTxManager,
			nil,
			logger,
		)

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		// When - JWT manager is nil, so it will fail with "invalid or expired token"
		response, err := service.ProcessPurchase(ctx, "invalid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		// The error comes from the nil JWT manager check
		assert.Equal(t, "invalid or expired token", err.Error())
	})
}

func TestPurchaseService_GetPurchaseHistory(t *testing.T) {
	ctx := context.Background()

	t.Run("should get purchase history with pagination", func(t *testing.T) {
		// Setup
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			transactionRepo: mockTransactionRepo,
			logger:          logger,
		}

		walletID := uint(1)

		transactions := []transactionEntity.Transaction{
			{
				WalletID:     walletID,
				Type:         transactionEntity.TypePurchase,
				Amount:       100.00,
				Status:       transactionEntity.StatusSuccess,
				TargetNumber: "08123456789",
				SerialNumber: "SN123456",
			},
		}

		filter := &dto.PurchaseHistoryFilter{
			WalletID: walletID,
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockTransactionRepo.On("GetAll", ctx, mock.MatchedBy(func(f *transactionDto.TransactionFilter) bool {
			return f.Page == 1 && f.PageSize == 10
		})).Return(transactions, int64(1), nil)

		// When
		response, err := service.GetPurchaseHistory(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, 100.00, response.Data[0].Amount)
		assert.Equal(t, "08123456789", response.Data[0].Phone)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("should set default page and page size", func(t *testing.T) {
		// Setup
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			transactionRepo: mockTransactionRepo,
			logger:          logger,
		}

		walletID := uint(1)

		filter := &dto.PurchaseHistoryFilter{
			WalletID: walletID,
			Page:     0,
			PageSize: 0,
		}

		// Mock expectations
		mockTransactionRepo.On("GetAll", ctx, mock.MatchedBy(func(f *transactionDto.TransactionFilter) bool {
			return f.Page == 1 && f.PageSize == 10
		})).Return([]transactionEntity.Transaction{}, int64(0), nil)

		// When
		response, err := service.GetPurchaseHistory(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("should calculate total pages correctly", func(t *testing.T) {
		// Setup
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			transactionRepo: mockTransactionRepo,
			logger:          logger,
		}

		walletID := uint(1)

		transactions := make([]transactionEntity.Transaction, 25)

		filter := &dto.PurchaseHistoryFilter{
			WalletID: walletID,
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockTransactionRepo.On("GetAll", ctx, mock.AnythingOfType("*dto.TransactionFilter")).Return(transactions, int64(25), nil)

		// When
		response, err := service.GetPurchaseHistory(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(25), response.Total)
		assert.Equal(t, 3, response.TotalPage) // 25/10 = 2 full pages + 1 partial page
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			transactionRepo: mockTransactionRepo,
			logger:          logger,
		}

		walletID := uint(1)

		filter := &dto.PurchaseHistoryFilter{
			WalletID: walletID,
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockTransactionRepo.On("GetAll", ctx, mock.AnythingOfType("*dto.TransactionFilter")).Return(nil, int64(0), errors.New("database error"))

		// When
		response, err := service.GetPurchaseHistory(ctx, filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("should convert transactions to purchase history correctly", func(t *testing.T) {
		// Setup
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			transactionRepo: mockTransactionRepo,
			logger:          logger,
		}

		walletID := uint(1)
		productID := uint(1)

		transactions := []transactionEntity.Transaction{
			{
				WalletID:     walletID,
				Type:         transactionEntity.TypePurchase,
				Amount:       250.00,
				Status:       transactionEntity.StatusSuccess,
				TargetNumber: "08987654321",
				SerialNumber: "SN789012",
				ProductID:    &productID,
			},
		}

		filter := &dto.PurchaseHistoryFilter{
			WalletID: walletID,
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockTransactionRepo.On("GetAll", ctx, mock.AnythingOfType("*dto.TransactionFilter")).Return(transactions, int64(1), nil)

		// When
		response, err := service.GetPurchaseHistory(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Data))
		hist := response.Data[0]
		assert.Equal(t, 250.00, hist.Amount)
		assert.Equal(t, "08987654321", hist.Phone)
		assert.Equal(t, "SN789012", hist.SerialNumber)
		assert.Equal(t, uint(1), hist.ProductID)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestPurchaseService_ProcessPurchaseWifi(t *testing.T) {
	ctx := context.Background()

	t.Run("should return error when product_id is missing", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 0,
		}

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product_id is required", err.Error())
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		logger := testutil.NewSilentLogger()
		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			jwtManager:      nil, // JWT manager is nil
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "invalid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("should return error when no available wifi vouchers", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.NewMockJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			userRepo:        mockUserRepo,
			productRepo:     mockProductRepo,
			jwtManager:      jwtManager,
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryWifi
		providerID := uint(1)
		durationHours := 24
		product.ProviderID = providerID
		product.DurationHours = &durationHours

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)
		mockWifiVoucherRepo.On("GetByProviderIDWithDurationHours", ctx, uint(1), 24).Return([]wifiVoucherEntity.WifiVoucher{}, nil)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "valid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "no available wifi vouchers", err.Error())
		mockWifiVoucherRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.NewMockJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			userRepo:        mockUserRepo,
			productRepo:     mockProductRepo,
			jwtManager:      jwtManager,
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 999,
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, errors.New("record not found"))

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "valid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "product not found")
		mockUserRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.NewMockJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			userRepo:        mockUserRepo,
			productRepo:     mockProductRepo,
			walletRepo:      mockWalletRepo,
			txManager:       mockTxManager,
			jwtManager:      jwtManager,
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		// Create voucher and product
		voucher := testutil.CreateWifiVoucherFixture()
		voucher.Status = wifiVoucherEntity.StatusAvailable
		vouchers := []wifiVoucherEntity.WifiVoucher{*voucher}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryWifi
		providerID := uint(1)
		durationHours := 24
		product.ProviderID = providerID
		product.DurationHours = &durationHours

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)
		mockWifiVoucherRepo.On("GetByProviderIDWithDurationHours", ctx, uint(1), 24).Return(vouchers, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)
				fn(ctx)
			}).
			Return(errors.New("wallet not found"))

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "valid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should return error when insufficient balance", func(t *testing.T) {
		// Setup
		mockWifiVoucherRepo := &testutil.MockWifiVoucherRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.NewMockJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			wifiVoucherRepo: mockWifiVoucherRepo,
			userRepo:        mockUserRepo,
			productRepo:     mockProductRepo,
			walletRepo:      mockWalletRepo,
			txManager:       mockTxManager,
			jwtManager:      jwtManager,
			logger:          logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		// Create voucher and product
		voucher := testutil.CreateWifiVoucherFixture()
		voucher.Status = wifiVoucherEntity.StatusAvailable
		vouchers := []wifiVoucherEntity.WifiVoucher{*voucher}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryWifi
		providerID := uint(1)
		durationHours := 24
		product.ProviderID = providerID
		product.DurationHours = &durationHours
		product.PriceSell = 100.0

		// Wallet with insufficient balance
		wallet := testutil.CreateWalletFixture()
		wallet.Balance = 50.0

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)
		mockWifiVoucherRepo.On("GetByProviderIDWithDurationHours", ctx, uint(1), 24).Return(vouchers, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)
				mockWalletRepo.On("GetForUpdate", ctx, user.ID).Return(wallet, nil).Once()
				fn(ctx)
			}).
			Return(errors.New("insufficient wallet balance"))

		// When
		response, err := service.ProcessPurchaseWifi(ctx, "valid-token", req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "insufficient wallet balance", err.Error())
	})
}
