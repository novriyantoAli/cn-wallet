package entity

import (
	"time"

	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"gorm.io/gorm"
)

// WifiVoucher represents the 'wifi_vouchers' table.
type WifiVoucher struct {
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string `gorm:"type:varchar(50);unique;not null" json:"code"`
	Password        string `gorm:"type:varchar(50)" json:"password"`
	DurationMinutes int    `gorm:"not null" json:"duration_minutes"`
	BatchID         string `gorm:"type:varchar(50)" json:"batch_id"`
	Status          string `gorm:"type:varchar(20);default:'available';check:status IN ('available', 'sold', 'used')" json:"status"`

	// Usage Logs
	SoldToUserID *uint      `gorm:"index" json:"sold_to_user_id"`
	SoldAt       *time.Time `json:"sold_at"`
	UsedAt       *time.Time `json:"used_at"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	SoldToUser *userEntity.User `gorm:"foreignKey:SoldToUserID" json:"sold_to_user,omitempty"`
}

// TableName specifies the table name for WifiVoucher
func (WifiVoucher) TableName() string {
	return "wifi_vouchers"
}

// Status constants
const (
	StatusAvailable = "available"
	StatusSold      = "sold"
	StatusUsed      = "used"
)
