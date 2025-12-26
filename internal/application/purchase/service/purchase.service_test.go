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
	"gorm.io/gorm"

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
			nil,
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
			nil,
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
			nil,
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
		jwtManager := testutil.CreateTestJWTManager()
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

		// Generate valid token
		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

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
		jwtManager := testutil.CreateTestJWTManager()
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
		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// Generate valid token
		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

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
		jwtManager := testutil.CreateTestJWTManager()
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
			Return(errors.New("wallet not found"))

		// Generate valid token
		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

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
		jwtManager := testutil.CreateTestJWTManager()
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
			Return(errors.New("insufficient wallet balance"))

		// Generate valid token
		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "insufficient wallet balance", err.Error())
	})

	t.Run("should return error when product is not wifi category", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			userRepo:    mockUserRepo,
			productRepo: mockProductRepo,
			jwtManager:  jwtManager,
			logger:      logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryPulsa // Not wifi

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product is not a wifi voucher", err.Error())
	})

	t.Run("should return error when product is inactive", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			userRepo:    mockUserRepo,
			productRepo: mockProductRepo,
			jwtManager:  jwtManager,
			logger:      logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryWifi
		product.IsActive = false

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product is inactive", err.Error())
	})

	t.Run("should return error when user is inactive", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()
		user := testutil.CreateUserFixture()
		user.IsActive = false

		service := &purchaseService{
			userRepo:   mockUserRepo,
			jwtManager: jwtManager,
			logger:     logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user is inactive", err.Error())
	})

	t.Run("should return error when product configuration is invalid", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()
		user := testutil.CreateUserFixture()

		service := &purchaseService{
			userRepo:    mockUserRepo,
			productRepo: mockProductRepo,
			jwtManager:  jwtManager,
			logger:      logger,
		}

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		product := testutil.CreateProductFixture()
		product.Category = productEntity.CategoryWifi
		product.ProviderID = 0 // Invalid

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchaseWifi(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid product configuration", err.Error())
	})
}

// Test ProcessPurchase success and failure cases
func TestPurchaseService_ProcessPurchase_Success(t *testing.T) {
	ctx := context.Background()

	t.Run("should process purchase successfully", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProviderClient := &MockProviderClient{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			mockTransactionRepo,
			mockUserRepo,
			nil,
			nil,
			mockProviderClient,
			mockTxManager,
			jwtManager,
			logger,
		)

		user := testutil.CreateUserFixture()
		product := testutil.CreateProductFixture()
		wallet := testutil.CreateWalletFixture()
		wallet.Balance = 100000.0

		req := &dto.PurchaseRequest{
			ProductID: product.ID,
			Phone:     "08123456789",
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, product.ID).Return(product, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, user.ID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
		mockWalletRepo.On("UpdateBalance", ctx, user.ID, mock.AnythingOfType("float64")).Return(nil)
		mockProviderClient.On("ProcessPurchase", ctx, product.Code, req.Phone, mock.AnythingOfType("string")).Return("SERIAL123", nil)
		mockTransactionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, transactionEntity.StatusSuccess, response.Status)
		assert.Equal(t, "SERIAL123", response.SerialNumber)
		assert.Equal(t, "Purchase successful", response.Message)
		mockProviderClient.AssertExpectations(t)
	})

	t.Run("should handle provider failure and refund", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockTransactionRepo := &testutil.MockTransactionRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockProviderClient := &MockProviderClient{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			mockTransactionRepo,
			mockUserRepo,
			nil,
			nil,
			mockProviderClient,
			mockTxManager,
			jwtManager,
			logger,
		)

		user := testutil.CreateUserFixture()
		product := testutil.CreateProductFixture()
		wallet := testutil.CreateWalletFixture()
		wallet.Balance = 100000.0

		req := &dto.PurchaseRequest{
			ProductID: product.ID,
			Phone:     "08123456789",
		}

		// Mock expectations
		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, product.ID).Return(product, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockWalletRepo.On("GetForUpdate", ctx, user.ID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
		mockWalletRepo.On("UpdateBalance", ctx, user.ID, mock.AnythingOfType("float64")).Return(nil)
		mockProviderClient.On("ProcessPurchase", ctx, product.Code, req.Phone, mock.AnythingOfType("string")).Return("", errors.New("provider error"))
		mockTransactionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, transactionEntity.StatusFailed, response.Status)
		assert.Equal(t, "Purchase failed, balance refunded", response.Message)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			mockProductRepo,
			nil,
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			nil,
			jwtManager,
			logger,
		)

		user := testutil.CreateUserFixture()

		req := &dto.PurchaseRequest{
			ProductID: 999,
			Phone:     "08123456789",
		}

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product not found", err.Error())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			nil,
			nil,
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			nil,
			jwtManager,
			logger,
		)

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		mockUserRepo.On("GetByID", ctx, uint(1)).Return(nil, gorm.ErrRecordNotFound)

		validToken := testutil.CreateValidJWTToken(jwtManager, 1)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			mockTxManager,
			jwtManager,
			logger,
		)

		user := testutil.CreateUserFixture()
		product := testutil.CreateProductFixture()

		req := &dto.PurchaseRequest{
			ProductID: product.ID,
			Phone:     "08123456789",
		}

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, product.ID).Return(product, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("wallet not found"))
		mockWalletRepo.On("GetForUpdate", ctx, user.ID).Return(nil, gorm.ErrRecordNotFound)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should return error when insufficient balance", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		mockWalletRepo := &testutil.MockWalletRepository{}
		mockUserRepo := &testutil.MockUserRepository{}
		mockTxManager := &testutil.MockTransactionManager{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := NewPurchaseService(
			mockProductRepo,
			mockWalletRepo,
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			mockTxManager,
			jwtManager,
			logger,
		)

		user := testutil.CreateUserFixture()
		product := testutil.CreateProductFixture()
		wallet := testutil.CreateWalletFixture()
		wallet.Balance = 100.0 // Less than product price

		req := &dto.PurchaseRequest{
			ProductID: product.ID,
			Phone:     "08123456789",
		}

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockProductRepo.On("GetByID", ctx, product.ID).Return(product, nil)
		mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("insufficient wallet balance"))
		mockWalletRepo.On("GetForUpdate", ctx, user.ID).Return(wallet, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		response, err := service.ProcessPurchase(ctx, validToken, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "insufficient wallet balance", err.Error())
	})
}

// Test Helper Functions
func TestPurchaseService_GetUserFromToken(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user successfully with valid token", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := &purchaseService{
			userRepo:   mockUserRepo,
			jwtManager: jwtManager,
			logger:     logger,
		}

		user := testutil.CreateUserFixture()

		mockUserRepo.On("GetByID", ctx, user.ID).Return(user, nil)

		validToken := testutil.CreateValidJWTToken(jwtManager, user.ID)

		// When
		result, err := service.getUserFromToken(ctx, validToken)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when jwt manager is nil", func(t *testing.T) {
		// Setup
		logger := testutil.NewSilentLogger()

		service := &purchaseService{
			jwtManager: nil,
			logger:     logger,
		}

		// When
		result, err := service.getUserFromToken(ctx, "token")

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		// Setup
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := &purchaseService{
			jwtManager: jwtManager,
			logger:     logger,
		}

		// When
		result, err := service.getUserFromToken(ctx, "invalid-token")

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := &purchaseService{
			userRepo:   mockUserRepo,
			jwtManager: jwtManager,
			logger:     logger,
		}

		mockUserRepo.On("GetByID", ctx, uint(1)).Return(nil, gorm.ErrRecordNotFound)

		validToken := testutil.CreateValidJWTToken(jwtManager, 1)

		// When
		result, err := service.getUserFromToken(ctx, validToken)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockUserRepo := &testutil.MockUserRepository{}
		logger := testutil.NewSilentLogger()
		jwtManager := testutil.CreateTestJWTManager()

		service := &purchaseService{
			userRepo:   mockUserRepo,
			jwtManager: jwtManager,
			logger:     logger,
		}

		mockUserRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("database error"))

		validToken := testutil.CreateValidJWTToken(jwtManager, 1)

		// When
		result, err := service.getUserFromToken(ctx, validToken)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestPurchaseService_GetProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should get product successfully", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()

		service := &purchaseService{
			productRepo: mockProductRepo,
			logger:      logger,
		}

		product := testutil.CreateProductFixture()

		mockProductRepo.On("GetByID", ctx, product.ID).Return(product, nil)

		// When
		result, err := service.getProduct(ctx, product.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, product.ID, result.ID)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()

		service := &purchaseService{
			productRepo: mockProductRepo,
			logger:      logger,
		}

		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		result, err := service.getProduct(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "product not found", err.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockProductRepo := &testutil.MockProductRepository{}
		logger := testutil.NewSilentLogger()

		service := &purchaseService{
			productRepo: mockProductRepo,
			logger:      logger,
		}

		mockProductRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("database error"))

		// When
		result, err := service.getProduct(ctx, 1)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
	})
}
