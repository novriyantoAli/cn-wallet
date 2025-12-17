package entity

import (
	"time"

	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"gorm.io/gorm"
)

type Product struct {
	ID         uint                     `json:"id" gorm:"primaryKey"`
	ProviderID uint                     `json:"provider_id" gorm:"type:bigint;not null;index"`
	Provider   *providerEntity.Provider `json:"provider,omitempty" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name       string                   `json:"name" gorm:"type:varchar(255);not null"`
	Code       string                   `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Price      float64                  `json:"price" gorm:"type:decimal(10,2);not null"`
	Type       string                   `json:"type" gorm:"type:varchar(50);not null"`
	CreatedAt  time.Time                `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt  time.Time                `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt           `json:"deleted_at,omitempty" gorm:"index"`
}

func (p Product) TableName() string {
	return "products"
}
