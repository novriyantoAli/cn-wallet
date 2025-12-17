package dto

import "time"

type CreateProviderRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
	Logo string `json:"logo" binding:"omitempty"`
}

type UpdateProviderRequest struct {
	Name string `json:"name" binding:"omitempty"`
	Code string `json:"code" binding:"omitempty"`
	Logo string `json:"logo" binding:"omitempty"`
}

type ProviderResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProviderListResponse struct {
	Data       []ProviderResponse `json:"data"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}

type ProviderFilter struct {
	Name     string `form:"name"`
	Code     string `form:"code"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
