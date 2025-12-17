package entity

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Code      string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Logo      string         `json:"logo" gorm:"type:varchar(500)"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p Provider) TableName() string {
	return "providers"
}
