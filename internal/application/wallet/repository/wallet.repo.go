package repository

import (
"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
"github.com/shopspring/decimal"
"go.uber.org/zap"
"gorm.io/gorm"
)

type WalletRepository interface {
Create(wallet *entity.Wallet) error
GetByID(id uint) (*entity.Wallet, error)
GetByUserID(userID uint) (*entity.Wallet, error)
GetAll(filter *dto.WalletFilter) ([]entity.Wallet, int64, error)
Update(wallet *entity.Wallet) error
Delete(id uint) error
AddBalance(id uint, amount decimal.Decimal) error
SubtractBalance(id uint, amount decimal.Decimal) error
}

type walletRepository struct {
db     *gorm.DB
logger *zap.Logger
}

func NewWalletRepository(db *gorm.DB, logger *zap.Logger) WalletRepository {
return &walletRepository{
db:     db,
logger: logger,
}
}

func (r *walletRepository) Create(wallet *entity.Wallet) error {
r.logger.Info("Creating wallet", zap.Uint("user_id", wallet.UserID))
return r.db.Create(wallet).Error
}

func (r *walletRepository) GetByID(id uint) (*entity.Wallet, error) {
var wallet entity.Wallet
err := r.db.First(&wallet, id).Error
if err != nil {
r.logger.Error("Failed to get wallet by ID", zap.Uint("id", id), zap.Error(err))
return nil, err
}
return &wallet, nil
}

func (r *walletRepository) GetByUserID(userID uint) (*entity.Wallet, error) {
var wallet entity.Wallet
err := r.db.Where("user_id = ?", userID).First(&wallet).Error
if err != nil {
r.logger.Error("Failed to get wallet by user ID", zap.Uint("user_id", userID), zap.Error(err))
return nil, err
}
return &wallet, nil
}

func (r *walletRepository) GetAll(filter *dto.WalletFilter) ([]entity.Wallet, int64, error) {
var wallets []entity.Wallet
var totalCount int64

query := r.db

// Get total count
if err := query.Model(&entity.Wallet{}).Count(&totalCount).Error; err != nil {
r.logger.Error("Failed to count wallets", zap.Error(err))
return nil, 0, err
}

// Apply pagination
offset := (filter.Page - 1) * filter.PageSize
if err := query.Offset(offset).Limit(filter.PageSize).Find(&wallets).Error; err != nil {
r.logger.Error("Failed to get wallets", zap.Error(err))
return nil, 0, err
}

return wallets, totalCount, nil
}

func (r *walletRepository) Update(wallet *entity.Wallet) error {
r.logger.Info("Updating wallet", zap.Uint("id", wallet.ID))
return r.db.Save(wallet).Error
}

func (r *walletRepository) Delete(id uint) error {
r.logger.Info("Deleting wallet", zap.Uint("id", id))
return r.db.Delete(&entity.Wallet{}, id).Error
}

func (r *walletRepository) AddBalance(id uint, amount decimal.Decimal) error {
r.logger.Info("Adding balance to wallet", zap.Uint("id", id), zap.String("amount", amount.String()))
return r.db.Model(&entity.Wallet{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *walletRepository) SubtractBalance(id uint, amount decimal.Decimal) error {
r.logger.Info("Subtracting balance from wallet", zap.Uint("id", id), zap.String("amount", amount.String()))
return r.db.Model(&entity.Wallet{}).Where("id = ?", id).Update("balance", gorm.Expr("balance - ?", amount)).Error
}
