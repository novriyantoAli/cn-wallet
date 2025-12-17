package dto

import "time"

type CreateProductRequest struct {
	ProviderID uint    `json:"provider_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Code       string  `json:"code" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Type       string  `json:"type" binding:"required,oneof=PULSA DATA GAME PLN"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name" binding:"omitempty"`
	Price float64 `json:"price" binding:"omitempty,gt=0"`
	Type  string  `json:"type" binding:"omitempty,oneof=PULSA DATA GAME PLN"`
}

type ProviderInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ProductResponse struct {
	ID         uint          `json:"id"`
	ProviderID uint          `json:"provider_id"`
	Provider   *ProviderInfo `json:"provider,omitempty"`
	Name       string        `json:"name"`
	Code       string        `json:"code"`
	Price      float64       `json:"price"`
	Type       string        `json:"type"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type ProductListResponse struct {
	Data       []ProductResponse `json:"data"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

type ProductFilter struct {
	ProviderID uint   `form:"provider_id"`
	Type       string `form:"type"`
	Code       string `form:"code"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}
