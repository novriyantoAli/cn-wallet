package entity

import "time"

// PaylaterLoanSource constants
const (
	PaylaterLoanSourceTransfer PaylaterLoanSource = "transfer"
	PaylaterLoanSourceCheckout PaylaterLoanSource = "checkout"
)

type PaylaterLoanSource string

func (ps PaylaterLoanSource) String() string {
	return string(ps)
}

func (ps PaylaterLoanSource) IsValid() bool {
	switch ps {
	case PaylaterLoanSourceTransfer, PaylaterLoanSourceCheckout:
		return true
	default:
		return false
	}
}

// PaylaterLoanStatus constants
const (
	PaylaterLoanStatusPending   PaylaterLoanStatus = "pending"
	PaylaterLoanStatusActive    PaylaterLoanStatus = "active"
	PaylaterLoanStatusPaid      PaylaterLoanStatus = "paid"
	PaylaterLoanStatusOverdue   PaylaterLoanStatus = "overdue"
	PaylaterLoanStatusDefaulted PaylaterLoanStatus = "defaulted"
)

type PaylaterLoanStatus string

func (ps PaylaterLoanStatus) String() string {
	return string(ps)
}

func (ps PaylaterLoanStatus) IsValid() bool {
	switch ps {
	case PaylaterLoanStatusPending, PaylaterLoanStatusActive, PaylaterLoanStatusPaid, PaylaterLoanStatusOverdue, PaylaterLoanStatusDefaulted:
		return true
	default:
		return false
	}
}

// PaylaterLoan represents a loan taken by a user through paylater
type PaylaterLoan struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	UserID    uint               `gorm:"not null;index" json:"user_id"`
	Amount    int64              `gorm:"not null;check:amount > 0" json:"amount"`
	Interest  int64              `gorm:"not null;check:interest >= 0" json:"interest"`
	Total     int64              `gorm:"not null" json:"total"`
	DueDate   time.Time          `gorm:"type:date;not null" json:"due_date"`
	Source    PaylaterLoanSource `gorm:"type:varchar(20);not null;check:source IN ('transfer', 'checkout')" json:"source"`
	Status    PaylaterLoanStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for PaylaterLoan model
func (PaylaterLoan) TableName() string {
	return "paylater_loans"
}
