package entity

import "time"

type PaylaterRepayment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PaylaterLoanID uint      `gorm:"not null;index" json:"ploan_id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	Amount         int64     `gorm:"not null;check:amount > 0" json:"amount"`
	PaymentSource  string    `gorm:"not null;type:varchar(50)" json:"payment_source"`
	Status         string    `gorm:"not null;type:varchar(50);default:'pending'" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
