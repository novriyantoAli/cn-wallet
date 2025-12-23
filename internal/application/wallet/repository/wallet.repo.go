package repository

import (
	"context"
	"errors"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WalletRepository defines the interface for wallet data access
type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) error
	GetWalletByUserID(ctx context.Context, userID uint) (*entity.Wallet, error)
	GetWalletByID(ctx context.Context, id uint) (*entity.Wallet, error)
	GetForUpdate(ctx context.Context, userID uint) (*entity.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *entity.Wallet) error
	UpdateBalance(ctx context.Context, userID uint, newBalance float64) error
	DeleteWallet(ctx context.Context, userID uint) error
}

type walletRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB, logger *zap.Logger) WalletRepository {
	return &walletRepository{
		db:     db,
		logger: logger,
	}
}

// CreateWallet creates a new wallet for a user
func (r *walletRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	r.logger.Info("Creating wallet", zap.Uint("user_id", wallet.UserID))
	db := database.GetDB(ctx, r.db)
	return db.Create(wallet).Error
}

// GetWalletByUserID retrieves wallet by user ID
func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	db := database.GetDB(ctx, r.db)
	result := db.Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Wallet not found", zap.Uint("user_id", userID))
			return nil, errors.New("wallet not found")
		}
		r.logger.Error("Failed to get wallet", zap.Error(result.Error), zap.Uint("user_id", userID))
		return nil, result.Error
	}
	return &wallet, nil
}

// GetWalletByID retrieves wallet by wallet ID
func (r *walletRepository) GetWalletByID(ctx context.Context, id uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	db := database.GetDB(ctx, r.db)
	result := db.First(&wallet, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Wallet not found", zap.Uint("id", id))
			return nil, errors.New("wallet not found")
		}
		r.logger.Error("Failed to get wallet", zap.Error(result.Error), zap.Uint("id", id))
		return nil, result.Error
	}
	return &wallet, nil
}

// GetForUpdate retrieves wallet by user ID with a FOR UPDATE lock for atomic updates
func (r *walletRepository) GetForUpdate(ctx context.Context, userID uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	db := database.GetDB(ctx, r.db)
	result := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Wallet not found", zap.Uint("user_id", userID))
			return nil, errors.New("wallet not found")
		}
		r.logger.Error("Failed to get wallet for update", zap.Error(result.Error), zap.Uint("user_id", userID))
		return nil, result.Error
	}
	return &wallet, nil
}

// UpdateWallet updates wallet information
func (r *walletRepository) UpdateWallet(ctx context.Context, wallet *entity.Wallet) error {
	db := database.GetDB(ctx, r.db)
	return db.Save(wallet).Error
}

// UpdateBalance updates the wallet balance
func (r *walletRepository) UpdateBalance(ctx context.Context, userID uint, newBalance float64) error {
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Update("balance", newBalance).Error
}

// DeleteWallet deletes a wallet
func (r *walletRepository) DeleteWallet(ctx context.Context, userID uint) error {
	db := database.GetDB(ctx, r.db)
	return db.Where("user_id = ?", userID).Delete(&entity.Wallet{}).Error
}
