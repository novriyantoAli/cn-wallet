package testutil

import (
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productDto "github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerDto "github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	userDto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletDto "github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
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
		Price:      10000.00,
		Type:       "PULSA",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func CreateProductRequestFixture() *productDto.CreateProductRequest {
	return &productDto.CreateProductRequest{
		ProviderID: 1,
		Name:       "Paket Data 5GB",
		Code:       "DATA5GB",
		Price:      50000.00,
		Type:       "DATA",
	}
}

func CreateUpdateProductRequestFixture() *productDto.UpdateProductRequest {
	return &productDto.UpdateProductRequest{
		Name:  "Paket Data 10GB",
		Price: 75000.00,
		Type:  "DATA",
	}
}

func CreateProductFilterFixture() *productDto.ProductFilter {
	return &productDto.ProductFilter{
		ProviderID: 1,
		Page:       1,
		PageSize:   10,
	}
}
