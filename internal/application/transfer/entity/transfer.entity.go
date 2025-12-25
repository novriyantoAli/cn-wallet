package entity

import (
	"time"
)

// Transfer status constants
const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

type TransferStatus string

func (ts TransferStatus) String() string {
	return string(ts)
}

func (ts TransferStatus) IsValid() bool {
	switch ts {
	case TransferStatusPending, TransferStatusCompleted, TransferStatusFailed, TransferStatusCancelled:
		return true
	default:
		return false
	}
}

// Transfer source constants
const (
	TransferSourceWallet   TransferSource = "wallet"
	TransferSourcePaylater TransferSource = "paylater"
)

type TransferSource string

func (ts TransferSource) String() string {
	return string(ts)
}

func (ts TransferSource) IsValid() bool {
	switch ts {
	case TransferSourceWallet, TransferSourcePaylater:
		return true
	default:
		return false
	}
}

// Transfer represents a transfer record between users
type Transfer struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	TargetUserID uint           `gorm:"not null;index" json:"target_user_id"`
	Amount       int64          `gorm:"not null;check:amount > 0" json:"amount"`
	Source       TransferSource `gorm:"type:varchar(20);not null;check:source IN ('wallet', 'paylater')" json:"source"`
	Status       TransferStatus `gorm:"type:varchar(20);not null;index;check:status IN ('pending', 'completed', 'failed', 'cancelled')" json:"status"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for Transfer
func (Transfer) TableName() string {
	return "transfers"
}
