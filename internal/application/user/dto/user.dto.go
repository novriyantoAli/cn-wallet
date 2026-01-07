package dto

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type UpdateUserProviderRequest struct {
	ProviderID uint `json:"provider_id" binding:"required"`
}

type UpdateUserLevelRequest struct {
	Level string `json:"level" binding:"required,oneof=user provider reseller admin"`
}

type UpdateUserRequest struct {
	FullName     string `json:"full_name" binding:"omitempty"`
	Level        string `json:"level" binding:"omitempty,oneof=user provider reseller admin"`
	IsActive     bool   `json:"is_active"`
	ProviderName string `json:"provider_name" binding:"omitempty"`
	ProviderCode string `json:"provider_code" binding:"omitempty"`
	ProviderLogo string `json:"provider_logo" binding:"omitempty"`
}

type WalletInfo struct {
	ID        uint      `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID         uint        `json:"id"`
	ProviderID *uint       `json:"provider_id"`
	Email      string      `json:"email"`
	FullName   string      `json:"full_name"`
	Level      string      `json:"level"`
	IsActive   bool        `json:"is_active"`
	Wallet     *WalletInfo `json:"wallet,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type UserFilter struct {
	ProviderID uint   `form:"provider_id"`
	Name       string `form:"name"`
	Email      string `form:"email"`
	Level      string `form:"level"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}
