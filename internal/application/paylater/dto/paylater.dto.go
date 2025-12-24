package dto

import "time"

// CreatePaylaterAccountRequest represents a request to create a paylater account
type CreatePaylaterAccountRequest struct {
	UserID      uint  `json:"user_id" binding:"required"`
	CreditLimit int64 `json:"credit_limit" binding:"required,gt=0"`
}

// UpdateCreditLimitRequest represents a request to update credit limit
type UpdateCreditLimitRequest struct {
	CreditLimit int64 `json:"credit_limit" binding:"required,gt=0"`
}

// UpdateStatusRequest represents a request to update account status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active suspended closed"`
}

// UseCredit represents a request to use credit from paylater account
type UseCreditRequest struct {
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
}

// RepaymentRequest represents a request to repay outstanding balance
type RepaymentRequest struct {
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
}

// GetPaylaterAccountResponse represents paylater account information response
type GetPaylaterAccountResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	CreditLimit    int64     `json:"credit_limit"`
	Outstanding    int64     `json:"outstanding"`
	AvailableLimit int64     `json:"available_limit"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// UseCreditResponse represents the response after using credit
type UseCreditResponse struct {
	UserID         uint   `json:"user_id"`
	Amount         int64  `json:"amount"`
	Outstanding    int64  `json:"outstanding"`
	AvailableLimit int64  `json:"available_limit"`
	Message        string `json:"message"`
}

// RepaymentResponse represents the response after repayment
type RepaymentResponse struct {
	UserID         uint   `json:"user_id"`
	Amount         int64  `json:"amount"`
	Outstanding    int64  `json:"outstanding"`
	AvailableLimit int64  `json:"available_limit"`
	Message        string `json:"message"`
}
