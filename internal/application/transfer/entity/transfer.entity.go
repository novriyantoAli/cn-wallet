package entity

import "time"

// TransferStatus represents the status of a transfer
type TransferStatus string

// TransferSource represents the source of funds for a transfer
type TransferSource string

// Transfer status constants
const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

// Transfer source constants
const (
	TransferSourceWallet   TransferSource = "wallet"
	TransferSourcePaylater TransferSource = "paylater"
)

// Transfer represents a transfer record between users
type Transfer struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID         string         `json:"uuid" gorm:"type:char(36);not null;uniqueIndex"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	TargetUserID uint           `json:"target_user_id" gorm:"not null;index"`
	Amount       int64          `json:"amount" gorm:"not null;check:amount > 0"`
	Source       TransferSource `json:"source" gorm:"type:varchar(20);not null;check:source IN ('wallet', 'paylater')"`
	Status       TransferStatus `json:"status" gorm:"type:varchar(20);not null;index;check:status IN ('pending', 'completed', 'failed', 'cancelled')"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// TableName specifies the table name for Transfer
func (Transfer) TableName() string {
	return "transfers"
}

// String returns the string representation of TransferStatus
func (ts TransferStatus) String() string {
	return string(ts)
}

// IsValid validates if TransferStatus is a valid status
func (ts TransferStatus) IsValid() bool {
	switch ts {
	case TransferStatusPending, TransferStatusCompleted, TransferStatusFailed, TransferStatusCancelled:
		return true
	default:
		return false
	}
}

// String returns the string representation of TransferSource
func (ts TransferSource) String() string {
	return string(ts)
}

// IsValid validates if TransferSource is a valid source
func (ts TransferSource) IsValid() bool {
	switch ts {
	case TransferSourceWallet, TransferSourcePaylater:
		return true
	default:
		return false
	}
}
