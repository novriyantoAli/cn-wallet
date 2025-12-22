package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productDto "github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerDto "github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	transactionDto "github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	securityEntity "github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	userDto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletDto "github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	wifiVoucherDto "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
)

// User fixtures
func CreateUserFixture() *userEntity.User {
	return &userEntity.User{
		ID:        1,
		Email:     "john@example.com",
		FullName:  "John Doe",
		Level:     "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateUserRequestFixture() *userDto.CreateUserRequest {
	return &userDto.CreateUserRequest{
		Email:    "john@example.com",
		FullName: "John Doe",
	}
}

func CreateUpdateUserRequestFixture() *userDto.UpdateUserRequest {
	return &userDto.UpdateUserRequest{
		FullName: "John Updated",
		Level:    "user",
		IsActive: true,
	}
}

// Payment fixtures
func CreatePaymentFixture() *entity.Payment {
	return &entity.Payment{
		ID:          1,
		Amount:      100.50,
		Currency:    "USD",
		Status:      entity.PaymentStatusPending,
		Description: "Test payment",
		UserID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func CreatePaymentRequestFixture() *dto.CreatePaymentRequest {
	return &dto.CreatePaymentRequest{
		Amount:      100.50,
		Currency:    "USD",
		Description: "Test payment",
		UserID:      1,
	}
}

func CreateUpdatePaymentRequestFixture() *dto.UpdatePaymentRequest {
	return &dto.UpdatePaymentRequest{
		Status:      entity.PaymentStatusCompleted.String(),
		Description: "Payment completed",
	}
}

func CreatePaymentFilterFixture() *dto.PaymentFilter {
	return &dto.PaymentFilter{
		Status:   "pending",
		Currency: "USD",
		UserID:   1,
		Page:     1,
		PageSize: 10,
	}
}

// Wallet fixtures
func CreateWalletFixture() *walletEntity.Wallet {
	return &walletEntity.Wallet{
		ID:        1,
		UserID:    1,
		Balance:   1000.00,
		PINHash:   "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92", // SHA256 of "123456"
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateWalletRequestFixture() *walletDto.CreateWalletRequest {
	return &walletDto.CreateWalletRequest{
		UserID: 1,
		PIN:    "123456",
	}
}

func CreateWalletResponseFixture() *walletDto.GetWalletResponse {
	return &walletDto.GetWalletResponse{
		ID:        1,
		UserID:    1,
		Balance:   1000.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Provider fixtures
func CreateProviderFixture() *providerEntity.Provider {
	return &providerEntity.Provider{
		ID:        1,
		Name:      "Telkomsel",
		Code:      "TSEL",
		Logo:      "https://example.com/logo.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateProviderRequestFixture() *providerDto.CreateProviderRequest {
	return &providerDto.CreateProviderRequest{
		Name: "Indosat",
		Code: "ISAT",
		Logo: "https://example.com/isat-logo.png",
	}
}

func CreateUpdateProviderRequestFixture() *providerDto.UpdateProviderRequest {
	return &providerDto.UpdateProviderRequest{
		Name: "Telkomsel Updated",
		Code: "TSEL",
		Logo: "https://example.com/updated-logo.png",
	}
}

func CreateProviderFilterFixture() *providerDto.ProviderFilter {
	return &providerDto.ProviderFilter{
		Page:     1,
		PageSize: 10,
	}
}

// Product fixtures
func CreateProductFixture() *productEntity.Product {
	return &productEntity.Product{
		ID:         1,
		ProviderID: 1,
		Name:       "Pulsa 10K",
		Code:       "PULSA10K",
		Category:   "pulsa",
		PriceBasic: 9500.00,
		PriceSell:  10500.00,
		IsActive:   true,
		IconURL:    "https://example.com/pulsa-icon.png",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func CreateProductRequestFixture() *productDto.CreateProductRequest {
	return &productDto.CreateProductRequest{
		ProviderID: 1,
		Name:       "Paket Data 5GB",
		Code:       "DATA5GB",
		Category:   "wifi",
		PriceBasic: 45000.00,
		PriceSell:  55000.00,
		IsActive:   true,
		IconURL:    "https://example.com/data-icon.png",
	}
}

func CreateUpdateProductRequestFixture() *productDto.UpdateProductRequest {
	return &productDto.UpdateProductRequest{
		Name:       "Paket Data 10GB",
		Category:   "wifi",
		PriceBasic: 70000.00,
		PriceSell:  80000.00,
		IsActive:   true,
		IconURL:    "https://example.com/data-10gb-icon.png",
	}
}

func CreateProductFilterFixture() *productDto.ProductFilter {
	return &productDto.ProductFilter{
		ProviderID: 1,
		Page:       1,
		PageSize:   10,
	}
}

// WifiVoucher fixtures
func CreateWifiVoucherFixture() *wifiVoucherEntity.WifiVoucher {
	userID := uint(1)
	soldAt := time.Now().Add(-24 * time.Hour)
	return &wifiVoucherEntity.WifiVoucher{
		ID:              1,
		Code:            "WIFI001",
		Password:        "pass123",
		DurationMinutes: 1440,
		BatchID:         "BATCH001",
		Status:          wifiVoucherEntity.StatusAvailable,
		SoldToUserID:    &userID,
		SoldAt:          &soldAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func CreateWifiVoucherRequestFixture() *wifiVoucherDto.CreateWifiVoucherRequest {
	return &wifiVoucherDto.CreateWifiVoucherRequest{
		Code:            "WIFI001",
		Password:        "pass123",
		DurationMinutes: 1440,
		BatchID:         "BATCH001",
	}
}

func CreateUpdateWifiVoucherRequestFixture() *wifiVoucherDto.UpdateWifiVoucherRequest {
	return &wifiVoucherDto.UpdateWifiVoucherRequest{
		Code:            "WIFI002",
		Password:        "pass456",
		DurationMinutes: 2880,
		BatchID:         "BATCH002",
		Status:          wifiVoucherEntity.StatusSold,
	}
}

func CreateWifiVoucherFilterFixture() *wifiVoucherDto.WifiVoucherFilter {
	return &wifiVoucherDto.WifiVoucherFilter{
		Status:   wifiVoucherEntity.StatusAvailable,
		Page:     1,
		PageSize: 10,
	}
}

// Transaction fixtures
func CreateTransactionFixture() *transactionEntity.Transaction {
	return &transactionEntity.Transaction{
		ID:                 uuid.New(),
		WalletID:           1,
		Type:               transactionEntity.TypeTopup,
		Amount:             100.00,
		Status:             transactionEntity.StatusPending,
		Description:        "Test transaction",
		PaymentMethod:      "credit_card",
		PaymentProviderRef: "ref123",
		ProductID:          nil,
		TargetNumber:       "",
		SerialNumber:       "",
		RelatedWalletID:    nil,
		CreatedAt:          time.Now(),
		Wallet:             nil,
		Product:            nil,
		RelatedWallet:      nil,
	}
}

func CreateTransactionRequestFixture() *transactionDto.CreateTransactionRequest {
	return &transactionDto.CreateTransactionRequest{
		WalletID:           1,
		Type:               transactionEntity.TypeTopup,
		Amount:             100.00,
		Description:        "Test transaction",
		PaymentMethod:      "credit_card",
		PaymentProviderRef: "ref123",
		ProductID:          nil,
		TargetNumber:       "",
		SerialNumber:       "",
		RelatedWalletID:    nil,
	}
}

func CreateUpdateTransactionRequestFixture() *transactionDto.UpdateTransactionRequest {
	return &transactionDto.UpdateTransactionRequest{
		Status:      transactionEntity.StatusSuccess,
		Description: "Updated transaction",
	}
}

func CreateTransactionFilterFixture() *transactionDto.TransactionFilter {
	return &transactionDto.TransactionFilter{
		WalletID: 1,
		Type:     transactionEntity.TypeTopup,
		Status:   transactionEntity.StatusPending,
		Page:     1,
		PageSize: 10,
	}
}

// Security fixtures
func CreateUserSecurityFixture() *securityEntity.UserSecurity {
	return &securityEntity.UserSecurity{
		UserID:        1,
		PinHash:       "hashed_pin_value",
		FailedAttempt: 0,
		LockedUntil:   nil,
	}
}
