package repository

import (
	"context"
	"errors"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PaylaterAccountRepository defines the interface for paylater account data access
type PaylaterAccountRepository interface {
	CreateAccount(ctx context.Context, account *entity.PaylaterAccount) error
	GetAccountByUserID(ctx context.Context, userID uint) (*entity.PaylaterAccount, error)
	GetAccountByID(ctx context.Context, id uint) (*entity.PaylaterAccount, error)
	GetForUpdate(ctx context.Context, userID uint) (*entity.PaylaterAccount, error)
	UpdateAccount(ctx context.Context, account *entity.PaylaterAccount) error
	UpdateCreditLimit(ctx context.Context, userID uint, creditLimit int64) error
	UpdateStatus(ctx context.Context, userID uint, status string) error
	DeleteAccount(ctx context.Context, userID uint) error
}

type paylaterAccountRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewPaylaterAccountRepository creates a new paylater account repository
func NewPaylaterAccountRepository(db *gorm.DB, logger *zap.Logger) PaylaterAccountRepository {
	return &paylaterAccountRepository{
		db:     db,
		logger: logger,
	}
}

// CreateAccount creates a new paylater account for a user
func (r *paylaterAccountRepository) CreateAccount(ctx context.Context, account *entity.PaylaterAccount) error {
	r.logger.Info("Creating paylater account", zap.Uint("user_id", account.UserID))
	db := database.GetDB(ctx, r.db)
	return db.Create(account).Error
}

// GetAccountByUserID retrieves paylater account by user ID
func (r *paylaterAccountRepository) GetAccountByUserID(ctx context.Context, userID uint) (*entity.PaylaterAccount, error) {
	var account entity.PaylaterAccount
	db := database.GetDB(ctx, r.db)
	result := db.Where("user_id = ?", userID).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Paylater account not found", zap.Uint("user_id", userID))
			return nil, errors.New("paylater account not found")
		}
		r.logger.Error("Failed to get paylater account", zap.Error(result.Error), zap.Uint("user_id", userID))
		return nil, result.Error
	}
	return &account, nil
}

// GetAccountByID retrieves paylater account by account ID
func (r *paylaterAccountRepository) GetAccountByID(ctx context.Context, id uint) (*entity.PaylaterAccount, error) {
	var account entity.PaylaterAccount
	db := database.GetDB(ctx, r.db)
	result := db.First(&account, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Paylater account not found", zap.Uint("id", id))
			return nil, errors.New("paylater account not found")
		}
		r.logger.Error("Failed to get paylater account", zap.Error(result.Error), zap.Uint("id", id))
		return nil, result.Error
	}
	return &account, nil
}

// GetForUpdate retrieves paylater account by user ID with a FOR UPDATE lock for atomic updates
func (r *paylaterAccountRepository) GetForUpdate(ctx context.Context, userID uint) (*entity.PaylaterAccount, error) {
	var account entity.PaylaterAccount
	db := database.GetDB(ctx, r.db)
	result := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Paylater account not found", zap.Uint("user_id", userID))
			return nil, errors.New("paylater account not found")
		}
		r.logger.Error("Failed to get paylater account for update", zap.Error(result.Error), zap.Uint("user_id", userID))
		return nil, result.Error
	}
	return &account, nil
}

// UpdateAccount updates paylater account information
func (r *paylaterAccountRepository) UpdateAccount(ctx context.Context, account *entity.PaylaterAccount) error {
	db := database.GetDB(ctx, r.db)
	return db.Save(account).Error
}

// UpdateCreditLimit updates the credit limit
func (r *paylaterAccountRepository) UpdateCreditLimit(ctx context.Context, userID uint, creditLimit int64) error {
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.PaylaterAccount{}).Where("user_id = ?", userID).Update("credit_limit", creditLimit).Error
}

// UpdateStatus updates the account status
func (r *paylaterAccountRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.PaylaterAccount{}).Where("user_id = ?", userID).Update("status", status).Error
}

// DeleteAccount deletes a paylater account
func (r *paylaterAccountRepository) DeleteAccount(ctx context.Context, userID uint) error {
	db := database.GetDB(ctx, r.db)
	return db.Where("user_id = ?", userID).Delete(&entity.PaylaterAccount{}).Error
}
