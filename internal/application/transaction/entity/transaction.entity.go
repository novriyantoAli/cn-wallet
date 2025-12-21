package entity

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents the 'transactions' table.
type Transaction struct {
	ID          uuid.UUID `gorm:"type:TEXT;primaryKey" json:"id"`
	WalletID    uint      `gorm:"not null;index" json:"wallet_id"`
	Type        string    `gorm:"type:varchar(20);not null" json:"type"`
	Amount      float64   `gorm:"type:REAL;not null" json:"amount"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Description string    `gorm:"type:text" json:"description"`
	// Payment Details (TopUp)
	PaymentMethod      string `gorm:"type:varchar(50)" json:"payment_method"`
	PaymentProviderRef string `gorm:"type:varchar(100)" json:"payment_provider_ref"`
	// Purchase Details
	ProductID    *uint  `json:"product_id"` // Nullable
	TargetNumber string `gorm:"type:varchar(50)" json:"target_number"`
	SerialNumber string `gorm:"type:varchar(100)" json:"serial_number"`
	// Transfer Details
	RelatedWalletID *uint     `json:"related_wallet_id"` // Nullable
	CreatedAt       time.Time `json:"created_at"`
	// Relationships
	Wallet        interface{} `gorm:"-" json:"-"`
	Product       interface{} `gorm:"-" json:"product,omitempty"`
	RelatedWallet interface{} `gorm:"-" json:"related_wallet,omitempty"`
}

// Transaction Type Constants
const (
	TypeTopup       = "topup"
	TypePurchase    = "purchase"
	TypeTransferOut = "transfer_out"
	TypeTransferIn  = "transfer_in"
	TypeRefund      = "refund"
)

// Transaction Status Constants
const (
	StatusPending = "pending"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)
