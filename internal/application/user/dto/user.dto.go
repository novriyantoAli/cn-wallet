package dto

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"omitempty"`
	Level    string `json:"level" binding:"omitempty,oneof=user agent admin"`
	IsActive bool   `json:"is_active"`
}

type WalletInfo struct {
	ID        uint      `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        uint        `json:"id"`
	Email     string      `json:"email"`
	FullName  string      `json:"full_name"`
	Level     string      `json:"level"`
	IsActive  bool        `json:"is_active"`
	Wallet    *WalletInfo `json:"wallet,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type UserFilter struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
