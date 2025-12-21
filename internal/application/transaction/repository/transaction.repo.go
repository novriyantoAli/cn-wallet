package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction data access.
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetAll(ctx context.Context, filter *dto.TransactionFilter) ([]entity.Transaction, int64, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWalletID(ctx context.Context, walletID uint, page, pageSize int) ([]entity.Transaction, int64, error)
}

type transactionRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewTransactionRepository creates a new instance of TransactionRepository.
func NewTransactionRepository(db *gorm.DB, logger *zap.Logger) TransactionRepository {
	return &transactionRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new transaction in the database.
func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	r.logger.Info("Creating transaction", zap.String("id", transaction.ID.String()), zap.String("type", transaction.Type))
	db := database.GetDB(ctx, r.db)
	return db.Create(transaction).Error
}

// GetByID retrieves a transaction by ID.
func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction
	db := database.GetDB(ctx, r.db)
	err := db.First(&transaction, "id = ?", id).Error
	if err != nil {
		r.logger.Error("Failed to get transaction by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}
	return &transaction, nil
}

// GetAll retrieves all transactions with filters and pagination.
func (r *transactionRepository) GetAll(ctx context.Context, filter *dto.TransactionFilter) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var totalCount int64

	db := database.GetDB(ctx, r.db)
	query := db.Model(&entity.Transaction{})

	// Apply filters
	if filter.WalletID != 0 {
		query = query.Where("wallet_id = ?", filter.WalletID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to count transactions", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		r.logger.Error("Failed to get transactions", zap.Error(err))
		return nil, 0, err
	}

	return transactions, totalCount, nil
}

// Update updates an existing transaction.
func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	r.logger.Info("Updating transaction", zap.String("id", transaction.ID.String()))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.Transaction{}).Where("id = ?", transaction.ID).Updates(transaction).Error
}

// Delete deletes a transaction by ID.
func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Info("Deleting transaction", zap.String("id", id.String()))
	db := database.GetDB(ctx, r.db)
	return db.Delete(&entity.Transaction{}, "id = ?", id).Error
}

// GetByWalletID retrieves all transactions for a specific wallet.
func (r *transactionRepository) GetByWalletID(ctx context.Context, walletID uint, page, pageSize int) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var totalCount int64

	db := database.GetDB(ctx, r.db)
	query := db.Where("wallet_id = ?", walletID)

	// Get total count
	if err := query.Model(&entity.Transaction{}).Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to count transactions for wallet", zap.Uint("wallet_id", walletID), zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		r.logger.Error("Failed to get transactions for wallet", zap.Uint("wallet_id", walletID), zap.Error(err))
		return nil, 0, err
	}

	return transactions, totalCount, nil
}
