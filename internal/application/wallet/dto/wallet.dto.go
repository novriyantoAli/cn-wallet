package dto

import "time"

// CreateWalletRequest represents a request to create a wallet
type CreateWalletRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	PIN    string `json:"pin" binding:"required,len=6,numeric"`
}

// SetPINRequest represents a request to set or update wallet PIN
type SetPINRequest struct {
	CurrentPIN string `json:"current_pin" binding:"required,len=6,numeric"`
	NewPIN     string `json:"new_pin" binding:"required,len=6,numeric"`
}

// VerifyPINRequest represents a request to verify wallet PIN
type VerifyPINRequest struct {
	PIN string `json:"pin" binding:"required,len=6,numeric"`
}

// UpdateBalanceRequest represents a request to update wallet balance
type UpdateBalanceRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

// TransferRequest represents a request to transfer balance between wallets
type TransferRequest struct {
	ToUserID    uint    `json:"to_user_id" binding:"required,gt=0"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// TransferResponse represents the response after a successful transfer
type TransferResponse struct {
	FromUserID  uint    `json:"from_user_id"`
	ToUserID    uint    `json:"to_user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
	Message     string  `json:"message"`
}

// GetWalletResponse represents wallet information response
type GetWalletResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WalletResponse wraps wallet data in a response
type WalletResponse struct {
	Data    *GetWalletResponse `json:"data,omitempty"`
	Message string             `json:"message,omitempty"`
	Error   string             `json:"error,omitempty"`
}
