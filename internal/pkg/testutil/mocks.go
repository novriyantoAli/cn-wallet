package testutil

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	userDto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletDto "github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *userEntity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*userEntity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*userEntity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(filter *userDto.UserFilter) ([]userEntity.User, int64, error) {
	args := m.Called(filter)
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

func (m *MockUserRepository) Update(user *userEntity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(payment *entity.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(id uint) (*entity.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetAll(filter *dto.PaymentFilter) ([]entity.Payment, int64, error) {
	args := m.Called(filter)
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

func (m *MockPaymentRepository) Update(payment *entity.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByUserID(userID uint) ([]entity.Payment, error) {
	args := m.Called(userID)
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

func (m *MockUserService) CreateUser(req *userDto.CreateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*userDto.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*userDto.UserResponse, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUsers(filter *userDto.UserFilter) (*userDto.UserListResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req *userDto.UpdateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUserPassword(id uint, req *userDto.UpdateUserPasswordRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockUserService) UpdateUserPIN(id uint, req *userDto.UpdateUserPINRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

// MockWalletRepository is a mock implementation of WalletRepository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(wallet *walletEntity.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByID(id uint) (*walletEntity.Wallet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetByUserID(userID uint) (*walletEntity.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletEntity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetAll(filter *walletDto.WalletFilter) ([]walletEntity.Wallet, int64, error) {
	args := m.Called(filter)
	var wallets []walletEntity.Wallet
	if args.Get(0) != nil {
		wallets = args.Get(0).([]walletEntity.Wallet)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return wallets, count, args.Error(2)
}

func (m *MockWalletRepository) Update(wallet *walletEntity.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWalletRepository) AddBalance(id uint, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) SubtractBalance(id uint, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

// MockWalletService is a mock implementation of WalletService
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(req *walletDto.CreateWalletRequest) (*walletDto.WalletResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.WalletResponse), args.Error(1)
}

func (m *MockWalletService) GetWalletByID(id uint) (*walletDto.WalletResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.WalletResponse), args.Error(1)
}

func (m *MockWalletService) GetWalletByUserID(userID uint) (*walletDto.WalletResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.WalletResponse), args.Error(1)
}

func (m *MockWalletService) GetWallets(filter *walletDto.WalletFilter) (*walletDto.WalletListResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletDto.WalletListResponse), args.Error(1)
}

func (m *MockWalletService) UpdateWalletPIN(id uint, req *walletDto.UpdateWalletPINRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockWalletService) AddBalance(id uint, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func (m *MockWalletService) WithdrawBalance(id uint, amount decimal.Decimal, pin string) error {
	args := m.Called(id, amount, pin)
	return args.Error(0)
}

func (m *MockWalletService) TransferBalance(fromID, toUserID uint, amount decimal.Decimal, pin string) error {
	args := m.Called(fromID, toUserID, amount, pin)
	return args.Error(0)
}

func (m *MockWalletService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
