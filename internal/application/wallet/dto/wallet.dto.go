package dto

import "time"

type CreateWalletRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	PIN    string `json:"pin" binding:"required,len=6,numeric"`
}

type UpdateWalletPINRequest struct {
	CurrentPIN string `json:"current_pin" binding:"required,len=6,numeric"`
	NewPIN     string `json:"new_pin" binding:"required,len=6,numeric"`
}

type AddBalanceRequest struct {
	Amount string `json:"amount" binding:"required"` // Decimal amount as string for precision
}

type WithdrawBalanceRequest struct {
	Amount string `json:"amount" binding:"required"` // Decimal amount as string for precision
	PIN    string `json:"pin" binding:"required,len=6,numeric"`
}

type TransferBalanceRequest struct {
	ToUserID uint   `json:"to_user_id" binding:"required"`
	Amount   string `json:"amount" binding:"required"` // Decimal amount as string for precision
	PIN      string `json:"pin" binding:"required,len=6,numeric"`
}

type WalletResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletListResponse struct {
	Data       []WalletResponse `json:"data"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

type WalletFilter struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
