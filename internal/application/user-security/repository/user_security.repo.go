package repository

import (
	"context"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserSecurityRepository interface {
	GetByUserID(ctx context.Context, userID uint) (*entity.UserSecurity, error)
	Create(ctx context.Context, security *entity.UserSecurity) error
	UpdatePIN(ctx context.Context, userID uint, pinHash string) error
	IncrementFailedAttempt(ctx context.Context, userID uint) error
	ResetFailedAttempt(ctx context.Context, userID uint) error
	LockAccount(ctx context.Context, userID uint, duration time.Duration) error
	Unlock(ctx context.Context, userID uint) error
	Delete(ctx context.Context, userID uint) error
}

type userSecurityRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserSecurityRepository(db *gorm.DB, logger *zap.Logger) UserSecurityRepository {
	return &userSecurityRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userSecurityRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserSecurity, error) {
	r.logger.Info("Getting user security", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	var security entity.UserSecurity
	if err := db.Where("user_id = ?", userID).First(&security).Error; err != nil {
		return nil, err
	}
	return &security, nil
}

func (r *userSecurityRepository) Create(ctx context.Context, security *entity.UserSecurity) error {
	r.logger.Info("Creating user security", zap.Uint("user_id", security.UserID))
	db := database.GetDB(ctx, r.db)
	return db.Create(security).Error
}

func (r *userSecurityRepository) UpdatePIN(ctx context.Context, userID uint, pinHash string) error {
	r.logger.Info("Updating user PIN", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.UserSecurity{}).Where("user_id = ?", userID).Update("pin_hash", pinHash).Error
}

func (r *userSecurityRepository) IncrementFailedAttempt(ctx context.Context, userID uint) error {
	r.logger.Info("Incrementing failed attempt", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.UserSecurity{}).Where("user_id = ?", userID).Update("failed_attempt", gorm.Expr("failed_attempt + ?", 1)).Error
}

func (r *userSecurityRepository) ResetFailedAttempt(ctx context.Context, userID uint) error {
	r.logger.Info("Resetting failed attempt", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.UserSecurity{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"failed_attempt": 0,
		"locked_until":   nil,
	}).Error
}

func (r *userSecurityRepository) LockAccount(ctx context.Context, userID uint, duration time.Duration) error {
	r.logger.Info("Locking user account", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	lockedUntil := time.Now().Add(duration)
	return db.Model(&entity.UserSecurity{}).Where("user_id = ?", userID).Update("locked_until", lockedUntil).Error
}

func (r *userSecurityRepository) Unlock(ctx context.Context, userID uint) error {
	r.logger.Info("Unlocking user account", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.UserSecurity{}).Where("user_id = ?", userID).Update("locked_until", nil).Error
}

func (r *userSecurityRepository) Delete(ctx context.Context, userID uint) error {
	r.logger.Info("Deleting user security", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	return db.Where("user_id = ?", userID).Delete(&entity.UserSecurity{}).Error
}
