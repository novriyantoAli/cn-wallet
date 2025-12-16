package entity

import (
"time"

"github.com/shopspring/decimal"
)

type Wallet struct {
ID        uint            `json:"id" gorm:"primaryKey"`
UserID    uint            `json:"user_id" gorm:"uniqueIndex;not null;foreignKey:UserID;references:ID;onDelete:CASCADE"`
Balance   decimal.Decimal `json:"balance" gorm:"type:numeric(15,2);default:0.00;check:balance >= 0"`
PinHash   string          `json:"-" gorm:"type:varchar(255);not null"` // 6-digit Security PIN
CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime:milli"`
UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (w Wallet) TableName() string {
return "wallets"
}
