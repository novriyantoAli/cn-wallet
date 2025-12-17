package testutil

import (
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
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
