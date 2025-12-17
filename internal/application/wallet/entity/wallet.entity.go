package entity

import "time"

// Wallet represents a user's wallet with balance and security PIN
type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance   float64   `gorm:"type:decimal(15,2);default:0.00;check:balance >= 0" json:"balance"`
	PINHash   string    `gorm:"type:varchar(255);not null" json:"pin_hash"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Wallet model
func (Wallet) TableName() string {
	return "wallets"
}
