package entity

import (
	"time"

	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"gorm.io/gorm"
)

type ProductCategory string

type Product struct {
	ID            uint                     `json:"id" gorm:"primaryKey"`
	ProviderID    uint                     `json:"provider_id" gorm:"type:bigint;not null;index"`
	Provider      *providerEntity.Provider `json:"provider,omitempty" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name          string                   `json:"name" gorm:"type:varchar(255);not null"`
	Code          string                   `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Category      ProductCategory          `json:"category" gorm:"type:varchar(50);not null"`
	PriceBasic    float64                  `json:"price_basic" gorm:"type:decimal(12,2);not null"`
	PriceSell     float64                  `json:"price_sell" gorm:"type:decimal(12,2);not null"`
	IsActive      bool                     `json:"is_active" gorm:"default:true"`
	IconURL       string                   `json:"icon_url" gorm:"type:text"`
	DurationHours int                      `json:"duration_hours"`
	CreatedAt     time.Time                `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt     time.Time                `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt     gorm.DeletedAt           `json:"deleted_at,omitempty" gorm:"index"`
}

// Product Category Constants
const (
	CategoryPulsa ProductCategory = "pulsa"
	CategoryWifi  ProductCategory = "wifi"
)

func (p Product) TableName() string {
	return "products"
}

func (pc ProductCategory) String() string {
	return string(pc)
}

func (pc ProductCategory) IsValid() bool {
	switch pc {
	case CategoryPulsa, CategoryWifi:
		return true
	default:
		return false
	}
}
