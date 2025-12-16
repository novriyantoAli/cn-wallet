package entity

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Email       string          `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber sql.NullString  `json:"phone_number" gorm:"uniqueIndex"`
	FullName    sql.NullString  `json:"full_name"`
	PasswordHash sql.NullString `json:"-"`
	PinHash     string          `json:"-" gorm:"not null"` // 6-digit Security PIN (Bcrypt)
	Balance     datatypes.Decimal `json:"balance" gorm:"type:decimal(15,2);default:0.00;check:balance>=0"`
	Level       string          `json:"level" gorm:"type:varchar(20);default:'user';check:level IN ('user','agent','admin')"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}

func (u User) TableName() string {
	return "users"
}
