package testutil

import (
	"database/sql"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	userDto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/shopspring/decimal"
)

// User fixtures
func CreateUserFixture() *userEntity.User {
	return &userEntity.User{
		ID:           1,
		Email:        "john@example.com",
		PhoneNumber:  sql.NullString{String: "+1234567890", Valid: true},
		FullName:     sql.NullString{String: "John Doe", Valid: true},
		PasswordHash: sql.NullString{String: "$2a$10$example.hashed.password", Valid: true},
		PinHash:      "$2a$10$example.pin.hashed",
		Balance:      decimal.RequireFromString("1000.00"),
		Level:        "user",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func CreateUserRequestFixture() *userDto.CreateUserRequest {
	return &userDto.CreateUserRequest{
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		FullName:    "John Doe",
		Password:    "password123",
		PIN:         "123456",
	}
}

func CreateUpdateUserRequestFixture() *userDto.UpdateUserRequest {
	return &userDto.UpdateUserRequest{
		PhoneNumber: "+1987654321",
		FullName:    "John Updated",
		Level:       "user",
		IsActive:    true,
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
