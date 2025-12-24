package entity

import "time"

// PaylaterAccount represents a user's paylater account with credit limit
type PaylaterAccount struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	CreditLimit    int64     `gorm:"not null;check:credit_limit > 0" json:"credit_limit"`
	Outstanding    int64     `gorm:"not null;default:0;check:outstanding >= 0" json:"outstanding"`
	AvailableLimit int64     `gorm:"not null" json:"available_limit"`
	Status         string    `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for PaylaterAccount model
func (PaylaterAccount) TableName() string {
	return "paylater_accounts"
}

// PaylaterStatus constants
const (
	PaylaterStatusActive    = "active"
	PaylaterStatusSuspended = "suspended"
	PaylaterStatusClosed    = "closed"
)
