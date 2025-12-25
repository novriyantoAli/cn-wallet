package dto

import "time"

// CreateLedgerEntryRequest represents a request to create a ledger entry
type CreateLedgerEntryRequest struct {
	UserID        uint64 `json:"user_id" binding:"required"`
	ReferenceID   string `json:"reference_id" binding:"required"`
	ReferenceType string `json:"reference_type" binding:"required,oneof=transfer paylater payment repayment adjustment"`
	Debit         int64  `json:"debit" binding:"gte=0"`
	Credit        int64  `json:"credit" binding:"gte=0"`
	AccountType   string `json:"account_type" binding:"required,oneof=wallet paylater voucher"`
}

// GetLedgerEntryResponse represents a ledger entry response
type GetLedgerEntryResponse struct {
	ID            uint64    `json:"id"`
	UserID        uint64    `json:"user_id"`
	ReferenceID   string    `json:"reference_id"`
	ReferenceType string    `json:"reference_type"`
	Debit         int64     `json:"debit"`
	Credit        int64     `json:"credit"`
	AccountType   string    `json:"account_type"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListLedgerEntriesRequest represents a request to list ledger entries
type ListLedgerEntriesRequest struct {
	UserID        *uint64 `form:"user_id"`
	ReferenceType *string `form:"reference_type"`
	AccountType   *string `form:"account_type"`
	FromDate      *string `form:"from_date"`
	ToDate        *string `form:"to_date"`
	Page          int     `form:"page" binding:"min=1"`
	PageSize      int     `form:"page_size" binding:"min=1,max=100"`
}

// ListLedgerEntriesResponse represents paginated ledger entries response
type ListLedgerEntriesResponse struct {
	Data       []GetLedgerEntryResponse `json:"data"`
	TotalCount int64                    `json:"total_count"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// LedgerStatsResponse represents ledger statistics
type LedgerStatsResponse struct {
	UserID          uint64 `json:"user_id"`
	TotalDebit      int64  `json:"total_debit"`
	TotalCredit     int64  `json:"total_credit"`
	EntryCount      int64  `json:"entry_count"`
	WalletDebit     int64  `json:"wallet_debit"`
	WalletCredit    int64  `json:"wallet_credit"`
	PaylaterDebit   int64  `json:"paylater_debit"`
	PaylaterCredits int64  `json:"paylater_credit"`
}
