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
