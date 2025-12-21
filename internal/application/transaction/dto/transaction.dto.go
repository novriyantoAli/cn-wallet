package dto

import (
	"github.com/google/uuid"
)

// CreateTransactionRequest represents the request to create a transaction.
type CreateTransactionRequest struct {
	WalletID           uint    `json:"wallet_id" binding:"required"`
	Type               string  `json:"type" binding:"required,oneof=topup purchase transfer_out transfer_in refund"`
	Amount             float64 `json:"amount" binding:"required,gt=0"`
	Description        string  `json:"description"`
	PaymentMethod      string  `json:"payment_method"`
	PaymentProviderRef string  `json:"payment_provider_ref"`
	ProductID          *uint   `json:"product_id"`
	TargetNumber       string  `json:"target_number"`
	SerialNumber       string  `json:"serial_number"`
	RelatedWalletID    *uint   `json:"related_wallet_id"`
}

// UpdateTransactionRequest represents the request to update a transaction.
type UpdateTransactionRequest struct {
	Status      string `json:"status" binding:"oneof=pending success failed"`
	Description string `json:"description"`
}

// TransactionResponse represents the response for a transaction.
type TransactionResponse struct {
	ID                 uuid.UUID `json:"id"`
	WalletID           uint      `json:"wallet_id"`
	Type               string    `json:"type"`
	Amount             float64   `json:"amount"`
	Status             string    `json:"status"`
	Description        string    `json:"description"`
	PaymentMethod      string    `json:"payment_method"`
	PaymentProviderRef string    `json:"payment_provider_ref"`
	ProductID          *uint     `json:"product_id"`
	TargetNumber       string    `json:"target_number"`
	SerialNumber       string    `json:"serial_number"`
	RelatedWalletID    *uint     `json:"related_wallet_id"`
	CreatedAt          string    `json:"created_at"`
}

// TransactionListResponse represents the response for paginated transactions.
type TransactionListResponse struct {
	Data       []TransactionResponse `json:"data"`
	TotalCount int64                 `json:"total_count"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
}

// TransactionFilter represents filters for querying transactions.
type TransactionFilter struct {
	WalletID uint   `json:"wallet_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Page     int    `json:"page" binding:"min=1"`
	PageSize int    `json:"page_size" binding:"min=1,max=100"`
}
