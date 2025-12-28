package entity

import (
	"time"

	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"gorm.io/gorm"
)

type User struct {
	ID         uint                     `json:"id" gorm:"primaryKey"`
	Email      string                   `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	FullName   string                   `json:"full_name" gorm:"type:varchar(255);not null"`
	Level      UserLevel                `json:"level" gorm:"type:varchar(20);default:'user';check:level IN ('user','provider','reseller','admin')"`
	IsActive   bool                     `json:"is_active" gorm:"default:true"`
	Wallet     *walletEntity.Wallet     `json:"wallet,omitempty" gorm:"foreignKey:UserID;references:ID"`
	ProviderID *uint                    `json:"provider_id" gorm:"uniqueIndex"`
	Provider   *providerEntity.Provider `json:"provider,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time                `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt  time.Time                `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt           `json:"deleted_at,omitempty" gorm:"index"`
}

type UserLevel string

const (
	UserLevelUser     UserLevel = "user"
	UserLevelProvider UserLevel = "provider"
	UserLevelReseller UserLevel = "reseller"
	UserLevelAdmin    UserLevel = "admin"
)

func (u User) TableName() string {
	return "users"
}

func (ul UserLevel) String() string {
	return string(ul)
}

func (ul UserLevel) IsValid() bool {
	switch ul {
	case UserLevelUser, UserLevelProvider, UserLevelReseller, UserLevelAdmin:
		return true
	default:
		return false
	}
}
