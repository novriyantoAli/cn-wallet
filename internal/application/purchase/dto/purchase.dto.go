package dto

import (
	"github.com/google/uuid"
)

// PurchaseRequest represents a request to purchase a product
type PurchaseRequest struct {
	ProductID uint   `json:"product_id" binding:"required,gt=0"`
	Phone     string `json:"phone" binding:"required"`
}

// PurchaseResponse represents the response after a purchase attempt
type PurchaseResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"` // pending, success, failed
	SerialNumber  string    `json:"serial_number,omitempty"`
	Message       string    `json:"message"`
}

// PurchaseWifiRequest represents a request to purchase wifi voucher product
type PurchaseWifiRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
}

// PurchaseWifiResponse represents the response after a wifi purchase attempt
type PurchaseWifiResponse struct {
	TransactionID   uuid.UUID `json:"transaction_id"`
	VoucherID       uint      `json:"voucher_id"`
	VoucherCode     string    `json:"voucher_code"`
	VoucherPassword string    `json:"voucher_password"`
	Status          string    `json:"status"` // pending, success, failed
	Message         string    `json:"message"`
}

// PurchaseHistory represents a purchase transaction history
type PurchaseHistory struct {
	ID           uuid.UUID `json:"id"`
	WalletID     uint      `json:"wallet_id"`
	ProductID    uint      `json:"product_id"`
	Amount       float64   `json:"amount"`
	Phone        string    `json:"phone"`
	Status       string    `json:"status"`
	SerialNumber string    `json:"serial_number,omitempty"`
	ProviderRef  string    `json:"provider_ref,omitempty"`
	CreatedAt    string    `json:"created_at"`
}

// PurchaseHistoryList represents a list of purchase histories
type PurchaseHistoryList struct {
	Data      []PurchaseHistory `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}

// PurchaseHistoryFilter represents filters for purchase history
type PurchaseHistoryFilter struct {
	WalletID  uint
	ProductID uint
	Status    string
	Page      int
	PageSize  int
}
