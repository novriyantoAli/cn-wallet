package dto

import "time"

// CreateTransferRequest represents a request to create a transfer
type CreateTransferRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	TargetUserID uint   `json:"target_user_id" binding:"required"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Source       string `json:"source" binding:"required,oneof=wallet paylater"`
}

// GetTransferResponse represents transfer information response
type GetTransferResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	TargetUserID uint      `json:"target_user_id"`
	Amount       int64     `json:"amount"`
	Source       string    `json:"source"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListTransfersRequest represents a request to list transfers with filters
type ListTransfersRequest struct {
	UserID       *uint   `form:"user_id"`
	TargetUserID *uint   `form:"target_user_id"`
	Source       *string `form:"source"`
	Status       *string `form:"status"`
	FromDate     *string `form:"from_date"`
	ToDate       *string `form:"to_date"`
	Page         int     `form:"page" binding:"min=1"`
	PageSize     int     `form:"page_size" binding:"min=1,max=100"`
}

// ListTransfersResponse represents paginated transfer list response
type ListTransfersResponse struct {
	Data       []GetTransferResponse `json:"data"`
	TotalCount int64                 `json:"total_count"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// UpdateTransferStatusRequest represents a request to update transfer status
type UpdateTransferStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed failed cancelled"`
}

// TransferStatsResponse represents transfer statistics
type TransferStatsResponse struct {
	UserID        uint  `json:"user_id"`
	TotalSent     int64 `json:"total_sent"`
	TotalReceived int64 `json:"total_received"`
	CountSent     int64 `json:"count_sent"`
	CountReceived int64 `json:"count_received"`
	WalletSent    int64 `json:"wallet_sent"`
	PaylaterSent  int64 `json:"paylater_sent"`
	CompletedSent int64 `json:"completed_sent"`
	PendingSent   int64 `json:"pending_sent"`
	FailedSent    int64 `json:"failed_sent"`
}
