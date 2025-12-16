package dto

import "time"

type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,e164"`
	FullName    string `json:"full_name" binding:"required"`
	Password    string `json:"password" binding:"omitempty,min=8"`
	PIN         string `json:"pin" binding:"required,len=6,numeric"`
}

type UpdateUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"omitempty,e164"`
	FullName    string `json:"full_name" binding:"omitempty"`
	Level       string `json:"level" binding:"omitempty,oneof=user agent admin"`
	IsActive    bool   `json:"is_active"`
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UpdateUserPINRequest struct {
	PIN    string `json:"pin" binding:"required,len=6,numeric"`
	NewPIN string `json:"new_pin" binding:"required,len=6,numeric"`
}

type UserResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	FullName    string    `json:"full_name"`
	Balance     string    `json:"balance"`
	Level       string    `json:"level"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
