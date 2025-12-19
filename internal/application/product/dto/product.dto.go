package dto

import "time"

type CreateProductRequest struct {
	ProviderID uint    `json:"provider_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Code       string  `json:"code" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	PriceBasic float64 `json:"price_basic" binding:"required,gt=0"`
	PriceSell  float64 `json:"price_sell" binding:"required,gt=0"`
	Type       string  `json:"type" binding:"required,oneof=PULSA DATA GAME PLN"`
	IsActive   bool    `json:"is_active" binding:"omitempty"`
	IconURL    string  `json:"icon_url" binding:"omitempty"`
}

type UpdateProductRequest struct {
	Name       string  `json:"name" binding:"omitempty"`
	Category   string  `json:"category" binding:"omitempty"`
	PriceBasic float64 `json:"price_basic" binding:"omitempty,gt=0"`
	PriceSell  float64 `json:"price_sell" binding:"omitempty,gt=0"`
	Type       string  `json:"type" binding:"omitempty,oneof=PULSA DATA GAME PLN"`
	IsActive   bool    `json:"is_active" binding:"omitempty"`
	IconURL    string  `json:"icon_url" binding:"omitempty"`
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
	Category   string        `json:"category"`
	PriceBasic float64       `json:"price_basic"`
	PriceSell  float64       `json:"price_sell"`
	Type       string        `json:"type"`
	IsActive   bool          `json:"is_active"`
	IconURL    string        `json:"icon_url"`
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
