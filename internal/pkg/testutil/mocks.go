package testutil

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	ledgerDto "github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	paylaterDto "github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productDto "github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerDto "github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	purchaseDto "github.com/novriyantoAli/cn-wallet/internal/application/purchase/dto"
	transactionDto "github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	transferDto "github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	transferEntity "github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	userSecurityDto "github.com/novriyantoAli/cn-wallet/internal/application/user-security/dto"
	userSecurityEntity "github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	userDto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletDto "github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	wifiVoucherDto "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"

	"github.com/stretchr/testify/mock"
)

// Common test errors
var ErrRecordNotFound = errors.New("record not found")

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *userEntity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*userEntity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, filter *userDto.UserFilter) ([]userEntity.User, int64, error) {
	args := m.Called(ctx, filter)
	var users []userEntity.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userEntity.User)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return users, count, args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *userEntity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(ctx context.Context, id uint) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetAll(ctx context.Context, filter *dto.PaymentFilter) ([]entity.Payment, int64, error) {
	args := m.Called(ctx, filter)
	var payments []entity.Payment
	if args.Get(0) != nil {
		payments = args.Get(0).([]entity.Payment)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return payments, count, args.Error(2)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.Payment, error) {
	args := m.Called(ctx, userID)
	var payments []entity.Payment
	if args.Get(0) != nil {
		payments = args.Get(0).([]entity.Payment)
	}
	return payments, args.Error(1)
}

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *userDto.CreateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*userDto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*userDto.UserResponse, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUsers(ctx context.Context, filter *userDto.UserFilter) (*userDto.UserListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uint, req *userDto.UpdateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockWalletService is a mock implementation of WalletService
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(ctx context.Context, req *walletDto.CreateWalletRequest) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletService) GetWalletByUserID(ctx context.Context, userID uint) (*walletDto.GetWalletResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.GetWalletResponse), args.Error(1)
}

func (m *MockWalletService) GetWalletByID(ctx context.Context, id uint) (*walletDto.GetWalletResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.GetWalletResponse), args.Error(1)
}

func (m *MockWalletService) SetPIN(ctx context.Context, userID uint, req *walletDto.SetPINRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockWalletService) VerifyPIN(ctx context.Context, userID uint, pin string) error {
	args := m.Called(ctx, userID, pin)
	return args.Error(0)
}

func (m *MockWalletService) AddBalance(ctx context.Context, userID uint, amount float64, description string) (*walletDto.GetWalletResponse, error) {
	args := m.Called(ctx, userID, amount, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.GetWalletResponse), args.Error(1)
}

func (m *MockWalletService) DeductBalance(ctx context.Context, userID uint, amount float64, description string) (*walletDto.GetWalletResponse, error) {
	args := m.Called(ctx, userID, amount, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.GetWalletResponse), args.Error(1)
}

func (m *MockWalletService) Transfer(ctx context.Context, token string, req *walletDto.TransferRequest) (*walletDto.TransferResponse, error) {
	args := m.Called(ctx, token, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.TransferResponse), args.Error(1)
}

func (m *MockWalletService) DeleteWallet(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockWalletRepository is a mock implementation of WalletRepository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, wallet *walletEntity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletByUserID(ctx context.Context, userID uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetWalletByID(ctx context.Context, id uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWallet(ctx context.Context, wallet *walletEntity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, userID uint, newBalance float64) error {
	args := m.Called(ctx, userID, newBalance)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdatePIN(ctx context.Context, userID uint, pinHash string) error {
	args := m.Called(ctx, userID, pinHash)
	return args.Error(0)
}

func (m *MockWalletRepository) DeleteWallet(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockWalletRepository) GetForUpdate(ctx context.Context, userID uint) (*walletEntity.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *productEntity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uint) (*productEntity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productEntity.Product), args.Error(1)
}

func (m *MockProductRepository) GetByCode(ctx context.Context, code string) (*productEntity.Product, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productEntity.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter *productDto.ProductFilter) ([]productEntity.Product, int64, error) {
	args := m.Called(ctx, filter)
	var products []productEntity.Product
	if args.Get(0) != nil {
		products = args.Get(0).([]productEntity.Product)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return products, count, args.Error(2)
}

func (m *MockProductRepository) Update(ctx context.Context, product *productEntity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

// MockProviderRepository is a mock implementation of ProviderRepository
type MockProviderRepository struct {
	mock.Mock
}

func (m *MockProviderRepository) Create(ctx context.Context, provider *providerEntity.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepository) GetByID(ctx context.Context, id uint) (*providerEntity.Provider, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerEntity.Provider), args.Error(1)
}

func (m *MockProviderRepository) GetByCode(ctx context.Context, code string) (*providerEntity.Provider, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerEntity.Provider), args.Error(1)
}

func (m *MockProviderRepository) GetAll(ctx context.Context, filter *providerDto.ProviderFilter) ([]providerEntity.Provider, int64, error) {
	args := m.Called(ctx, filter)
	var providers []providerEntity.Provider
	if args.Get(0) != nil {
		providers = args.Get(0).([]providerEntity.Provider)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return providers, count, args.Error(2)
}

func (m *MockProviderRepository) Update(ctx context.Context, provider *providerEntity.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProviderRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *productDto.CreateProductRequest) (*productDto.ProductResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProductByID(ctx context.Context, id uint) (*productDto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProductByCode(ctx context.Context, code string) (*productDto.ProductResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetAllProducts(ctx context.Context, filter *productDto.ProductFilter) (*productDto.ProductListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductListResponse), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, id uint, req *productDto.UpdateProductRequest) (*productDto.ProductResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductResponse), args.Error(1)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *transactionEntity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*transactionEntity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionEntity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(ctx context.Context, filter *transactionDto.TransactionFilter) ([]transactionEntity.Transaction, int64, error) {
	args := m.Called(ctx, filter)
	var transactions []transactionEntity.Transaction
	if args.Get(0) != nil {
		transactions = args.Get(0).([]transactionEntity.Transaction)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return transactions, count, args.Error(2)
}

// MockTransactionService is a mock implementation of TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, req *transactionDto.CreateTransactionRequest) (*transactionDto.TransactionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionDto.TransactionResponse), args.Error(1)
}

func (m *MockTransactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*transactionDto.TransactionResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionDto.TransactionResponse), args.Error(1)
}

func (m *MockTransactionService) GetAllTransactions(ctx context.Context, filter *transactionDto.TransactionFilter) (*transactionDto.TransactionListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionDto.TransactionListResponse), args.Error(1)
}

func (m *MockTransactionService) UpdateTransaction(ctx context.Context, id uuid.UUID, req *transactionDto.UpdateTransactionRequest) (*transactionDto.TransactionResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionDto.TransactionResponse), args.Error(1)
}

func (m *MockTransactionService) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionService) GetWalletTransactions(ctx context.Context, walletID uint, page, pageSize int) (*transactionDto.TransactionListResponse, error) {
	args := m.Called(ctx, walletID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionDto.TransactionListResponse), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *transactionEntity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByWalletID(ctx context.Context, walletID uint, page, pageSize int) ([]transactionEntity.Transaction, int64, error) {
	args := m.Called(ctx, walletID, page, pageSize)
	var transactions []transactionEntity.Transaction
	if args.Get(0) != nil {
		transactions = args.Get(0).([]transactionEntity.Transaction)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return transactions, count, args.Error(2)
}

// MockTransactionManager is a mock implementation of TransactionManagerI
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// MockPurchaseService is a mock implementation of PurchaseService
type MockPurchaseService struct {
	mock.Mock
}

func (m *MockPurchaseService) ProcessPurchase(ctx context.Context, token string, req *purchaseDto.PurchaseRequest) (*purchaseDto.PurchaseResponse, error) {
	args := m.Called(ctx, token, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*purchaseDto.PurchaseResponse), args.Error(1)
}

func (m *MockPurchaseService) ProcessPurchaseWifi(ctx context.Context, token string, req *purchaseDto.PurchaseWifiRequest) (*purchaseDto.PurchaseResponse, error) {
	args := m.Called(ctx, token, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*purchaseDto.PurchaseResponse), args.Error(1)
}

func (m *MockPurchaseService) GetPurchaseHistory(ctx context.Context, filter *purchaseDto.PurchaseHistoryFilter) (*purchaseDto.PurchaseHistoryList, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*purchaseDto.PurchaseHistoryList), args.Error(1)
}

// MockUserSecurityRepository is a mock implementation of UserSecurityRepository
type MockUserSecurityRepository struct {
	mock.Mock
}

func (m *MockUserSecurityRepository) GetByUserID(ctx context.Context, userID uint) (*userSecurityEntity.UserSecurity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userSecurityEntity.UserSecurity), args.Error(1)
}

func (m *MockUserSecurityRepository) Create(ctx context.Context, security *userSecurityEntity.UserSecurity) error {
	args := m.Called(ctx, security)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) UpdatePIN(ctx context.Context, userID uint, pinHash string) error {
	args := m.Called(ctx, userID, pinHash)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) IncrementFailedAttempt(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) ResetFailedAttempt(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) LockAccount(ctx context.Context, userID uint, duration time.Duration) error {
	args := m.Called(ctx, userID, duration)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) Unlock(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserSecurityRepository) Delete(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockUserSecurityService is a mock implementation of UserSecurityService
type MockUserSecurityService struct {
	mock.Mock
}

func (m *MockUserSecurityService) SetPIN(ctx context.Context, req *userSecurityDto.SetPINRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserSecurityService) VerifyPIN(ctx context.Context, req *userSecurityDto.VerifyPINRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserSecurityService) GetSecurity(ctx context.Context, userID uint) (*userSecurityDto.UserSecurityResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userSecurityDto.UserSecurityResponse), args.Error(1)
}

func (m *MockUserSecurityService) IsAccountLocked(ctx context.Context, userID uint) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// MockWifiVoucherRepository is a mock implementation of WifiVoucherRepository
type MockWifiVoucherRepository struct {
	mock.Mock
}

func (m *MockWifiVoucherRepository) Create(ctx context.Context, wifiVoucher *wifiVoucherEntity.WifiVoucher) error {
	args := m.Called(ctx, wifiVoucher)
	return args.Error(0)
}

func (m *MockWifiVoucherRepository) GetByID(ctx context.Context, id uint) (*wifiVoucherEntity.WifiVoucher, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wifiVoucherEntity.WifiVoucher), args.Error(1)
}

func (m *MockWifiVoucherRepository) GetByCode(ctx context.Context, code string) (*wifiVoucherEntity.WifiVoucher, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wifiVoucherEntity.WifiVoucher), args.Error(1)
}

func (m *MockWifiVoucherRepository) GetAll(ctx context.Context, filter *wifiVoucherDto.WifiVoucherFilter) ([]wifiVoucherEntity.WifiVoucher, int64, error) {
	args := m.Called(ctx, filter)
	var vouchers []wifiVoucherEntity.WifiVoucher
	if args.Get(0) != nil {
		vouchers = args.Get(0).([]wifiVoucherEntity.WifiVoucher)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return vouchers, count, args.Error(2)
}

func (m *MockWifiVoucherRepository) Update(ctx context.Context, wifiVoucher *wifiVoucherEntity.WifiVoucher) error {
	args := m.Called(ctx, wifiVoucher)
	return args.Error(0)
}

func (m *MockWifiVoucherRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWifiVoucherRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockWifiVoucherRepository) GetByProviderAndDurationHours(ctx context.Context, providerID uint, durationHours int, filter *wifiVoucherDto.WifiVoucherFilter) ([]wifiVoucherEntity.WifiVoucher, int64, error) {
	args := m.Called(ctx, providerID, durationHours, filter)
	var vouchers []wifiVoucherEntity.WifiVoucher
	if args.Get(0) != nil {
		vouchers = args.Get(0).([]wifiVoucherEntity.WifiVoucher)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return vouchers, count, args.Error(2)
}

func (m *MockWifiVoucherRepository) GetByProviderIDWithDurationHours(ctx context.Context, providerID uint, durationHours int) ([]wifiVoucherEntity.WifiVoucher, error) {
	args := m.Called(ctx, providerID, durationHours)
	var vouchers []wifiVoucherEntity.WifiVoucher
	if args.Get(0) != nil {
		vouchers = args.Get(0).([]wifiVoucherEntity.WifiVoucher)
	}
	return vouchers, args.Error(1)
}

func (m *MockWifiVoucherRepository) GetForUpdate(ctx context.Context, id uint) (*wifiVoucherEntity.WifiVoucher, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wifiVoucherEntity.WifiVoucher), args.Error(1)
}

// MockPaylaterAccountRepository is a mock implementation of PaylaterAccountRepository
type MockPaylaterAccountRepository struct {
	mock.Mock
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

// MockPaylaterAccountService is a mock implementation of PaylaterAccountService
type MockPaylaterAccountService struct {
	mock.Mock
}

func (m *MockPaylaterAccountService) CreateAccount(ctx context.Context, req *paylaterDto.CreatePaylaterAccountRequest) (*paylaterEntity.PaylaterAccount, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterAccount), args.Error(1)
}

func (m *MockPaylaterAccountService) GetAccountByUserID(ctx context.Context, userID uint) (*paylaterDto.GetPaylaterAccountResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterAccountResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) GetAccountByID(ctx context.Context, id uint) (*paylaterDto.GetPaylaterAccountResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterAccountResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) UpdateCreditLimit(ctx context.Context, userID uint, req *paylaterDto.UpdateCreditLimitRequest) (*paylaterDto.GetPaylaterAccountResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterAccountResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) UpdateStatus(ctx context.Context, userID uint, req *paylaterDto.UpdateStatusRequest) (*paylaterDto.GetPaylaterAccountResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterAccountResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) UseCredit(ctx context.Context, userID uint, req *paylaterDto.UseCreditRequest) (*paylaterDto.UseCreditResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.UseCreditResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) Repayment(ctx context.Context, userID uint, req *paylaterDto.RepaymentRequest) (*paylaterDto.RepaymentResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.RepaymentResponse), args.Error(1)
}

func (m *MockPaylaterAccountService) DeleteAccount(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockPaylaterLoanRepository is a mock implementation of PaylaterLoanRepository
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

func (m *MockPaylaterLoanRepository) ListLoans(ctx context.Context, req *paylaterDto.ListPaylaterLoansRequest) ([]paylaterEntity.PaylaterLoan, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
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

// MockPaylaterLoanService is a mock implementation of PaylaterLoanService
type MockPaylaterLoanService struct {
	mock.Mock
}

func (m *MockPaylaterLoanService) CreateLoan(ctx context.Context, req *paylaterDto.CreatePaylaterLoanRequest) (*paylaterEntity.PaylaterLoan, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterEntity.PaylaterLoan), args.Error(1)
}

func (m *MockPaylaterLoanService) GetLoanByID(ctx context.Context, id uint) (*paylaterDto.GetPaylaterLoanResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterLoanResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) GetLoansByUserID(ctx context.Context, userID uint) ([]paylaterDto.GetPaylaterLoanResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]paylaterDto.GetPaylaterLoanResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) ListLoans(ctx context.Context, req *paylaterDto.ListPaylaterLoansRequest) (*paylaterDto.ListPaylaterLoansResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.ListPaylaterLoansResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) UpdateLoanStatus(ctx context.Context, id uint, req *paylaterDto.UpdateLoanStatusRequest) (*paylaterDto.GetPaylaterLoanResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterLoanResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) MarkLoanAsPaid(ctx context.Context, id uint) (*paylaterDto.GetPaylaterLoanResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.GetPaylaterLoanResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) ProcessOverdueLoans(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockPaylaterLoanService) GetUserLoanStats(ctx context.Context, userID uint) (*paylaterDto.PaylaterLoanStatsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paylaterDto.PaylaterLoanStatsResponse), args.Error(1)
}

func (m *MockPaylaterLoanService) DeleteLoan(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTransferService is a mock implementation of TransferService
type MockTransferService struct {
	mock.Mock
}

func (m *MockTransferService) CreateTransfer(ctx context.Context, req *transferDto.CreateTransferRequest) (*transferEntity.Transfer, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transferEntity.Transfer), args.Error(1)
}

func (m *MockTransferService) GetTransferByID(ctx context.Context, id uint) (*transferDto.GetTransferResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transferDto.GetTransferResponse), args.Error(1)
}

func (m *MockTransferService) GetTransfersByUserID(ctx context.Context, userID uint) ([]transferDto.GetTransferResponse, error) {
	args := m.Called(ctx, userID)
	var transfers []transferDto.GetTransferResponse
	if args.Get(0) != nil {
		transfers = args.Get(0).([]transferDto.GetTransferResponse)
	}
	return transfers, args.Error(1)
}

func (m *MockTransferService) GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]transferDto.GetTransferResponse, error) {
	args := m.Called(ctx, targetUserID)
	var transfers []transferDto.GetTransferResponse
	if args.Get(0) != nil {
		transfers = args.Get(0).([]transferDto.GetTransferResponse)
	}
	return transfers, args.Error(1)
}

func (m *MockTransferService) ListTransfers(ctx context.Context, req *transferDto.ListTransfersRequest) (*transferDto.ListTransfersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transferDto.ListTransfersResponse), args.Error(1)
}

func (m *MockTransferService) UpdateTransferStatus(ctx context.Context, id uint, req *transferDto.UpdateTransferStatusRequest) (*transferDto.GetTransferResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transferDto.GetTransferResponse), args.Error(1)
}

func (m *MockTransferService) GetUserTransferStats(ctx context.Context, userID uint) (*transferDto.TransferStatsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transferDto.TransferStatsResponse), args.Error(1)
}

func (m *MockTransferService) CancelTransfer(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLedgerService is a mock implementation of LedgerService
type MockLedgerService struct {
	mock.Mock
}

func (m *MockLedgerService) CreateEntry(ctx context.Context, req *ledgerDto.CreateLedgerEntryRequest) (*ledgerDto.GetLedgerEntryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerDto.GetLedgerEntryResponse), args.Error(1)
}

func (m *MockLedgerService) GetEntryByID(ctx context.Context, id uint64) (*ledgerDto.GetLedgerEntryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerDto.GetLedgerEntryResponse), args.Error(1)
}

func (m *MockLedgerService) GetEntriesByUserID(ctx context.Context, userID uint64) ([]ledgerDto.GetLedgerEntryResponse, error) {
	args := m.Called(ctx, userID)
	var entries []ledgerDto.GetLedgerEntryResponse
	if args.Get(0) != nil {
		entries = args.Get(0).([]ledgerDto.GetLedgerEntryResponse)
	}
	return entries, args.Error(1)
}

func (m *MockLedgerService) GetEntriesByReference(ctx context.Context, referenceType string, referenceID string) ([]ledgerDto.GetLedgerEntryResponse, error) {
	args := m.Called(ctx, referenceType, referenceID)
	var entries []ledgerDto.GetLedgerEntryResponse
	if args.Get(0) != nil {
		entries = args.Get(0).([]ledgerDto.GetLedgerEntryResponse)
	}
	return entries, args.Error(1)
}

func (m *MockLedgerService) ListEntries(ctx context.Context, req *ledgerDto.ListLedgerEntriesRequest) (*ledgerDto.ListLedgerEntriesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerDto.ListLedgerEntriesResponse), args.Error(1)
}

func (m *MockLedgerService) GetUserStats(ctx context.Context, userID uint64) (*ledgerDto.LedgerStatsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledgerDto.LedgerStatsResponse), args.Error(1)
}
