package dto

import (
	"time"
)

// CreateWifiVoucherRequest represents the request body for creating a wifi voucher.
type CreateWifiVoucherRequest struct {
	Code          string `json:"code" binding:"required,max=50"`
	Password      string `json:"password" binding:"max=50"`
	DurationHours int    `json:"duration_hours" binding:"required,gt=0"`
	BatchID       string `json:"batch_id" binding:"max=50"`
}

// UpdateWifiVoucherRequest represents the request body for updating a wifi voucher.
type UpdateWifiVoucherRequest struct {
	Code          string     `json:"code" binding:"max=50"`
	Password      string     `json:"password" binding:"max=50"`
	DurationHours int        `json:"duration_hours" binding:"gt=0"`
	BatchID       string     `json:"batch_id" binding:"max=50"`
	Status        string     `json:"status" binding:"omitempty,oneof=available sold used"`
	SoldToUserID  *uint      `json:"sold_to_user_id"`
	SoldAt        *time.Time `json:"sold_at"`
	UsedAt        *time.Time `json:"used_at"`
}

// WifiVoucherResponse represents the response body for wifi voucher.
type WifiVoucherResponse struct {
	ID            uint       `json:"id"`
	Code          string     `json:"code"`
	Password      string     `json:"password"`
	DurationHours int        `json:"duration_hours"`
	BatchID       string     `json:"batch_id"`
	Status        string     `json:"status"`
	SoldToUserID  *uint      `json:"sold_to_user_id"`
	SoldAt        *time.Time `json:"sold_at"`
	UsedAt        *time.Time `json:"used_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// WifiVoucherListResponse represents the list response for wifi vouchers.
type WifiVoucherListResponse struct {
	Data       []WifiVoucherResponse `json:"data"`
	TotalCount int64                 `json:"total_count"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
}

// WifiVoucherFilter represents filter options for listing wifi vouchers.
type WifiVoucherFilter struct {
	Status   string
	Code     string
	BatchID  string
	Page     int
	PageSize int
}
