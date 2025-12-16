package entity

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	Email        string          `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	PhoneNumber  sql.NullString  `json:"phone_number" gorm:"uniqueIndex;type:varchar(20)"`
	FullName     sql.NullString  `json:"full_name" gorm:"type:varchar(100)"`
	PasswordHash sql.NullString  `json:"-" gorm:"type:varchar(255)"`
	PinHash      string          `json:"-" gorm:"type:varchar(255);not null"` // 6-digit Security PIN (Bcrypt)
	Balance      decimal.Decimal `json:"balance" gorm:"type:numeric(15,2);default:0.00;check:balance >= 0"`
	Level        string          `json:"level" gorm:"type:varchar(20);default:'user';check:level IN ('user','agent','admin')"`
	IsActive     bool            `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}

func (u User) TableName() string {
	return "users"
}
