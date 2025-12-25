package entity

import "time"

// PaylaterLoan represents a loan taken by a user through paylater
type PaylaterLoan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Amount    int64     `gorm:"not null;check:amount > 0" json:"amount"`
	Interest  int64     `gorm:"not null;check:interest >= 0" json:"interest"`
	Total     int64     `gorm:"not null" json:"total"`
	DueDate   time.Time `gorm:"type:date;not null" json:"due_date"`
	Source    string    `gorm:"type:varchar(20);not null;check:source IN ('transfer', 'checkout')" json:"source"`
	Status    string    `gorm:"type:varchar(20);not null;index" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for PaylaterLoan model
func (PaylaterLoan) TableName() string {
	return "paylater_loans"
}

// PaylaterLoanSource constants
const (
	PaylaterLoanSourceTransfer = "transfer"
	PaylaterLoanSourceCheckout = "checkout"
)

// PaylaterLoanStatus constants
const (
	PaylaterLoanStatusPending   = "pending"
	PaylaterLoanStatusActive    = "active"
	PaylaterLoanStatusPaid      = "paid"
	PaylaterLoanStatusOverdue   = "overdue"
	PaylaterLoanStatusDefaulted = "defaulted"
)
