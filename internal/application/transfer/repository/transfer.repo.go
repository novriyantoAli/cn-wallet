package repository

import (
	"context"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, transfer *entity.Transfer) error
	GetTransferByID(ctx context.Context, id uint) (*entity.Transfer, error)
	GetTransfersByUserID(ctx context.Context, userID uint) ([]entity.Transfer, error)
	GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]entity.Transfer, error)
	ListTransfers(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]entity.Transfer, int64, error)
	UpdateTransferStatus(ctx context.Context, id uint, status string) error
	GetUserTransferStats(ctx context.Context, userID uint) (map[string]interface{}, error)
	CancelTransfer(ctx context.Context, id uint) error
}

type transferRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTransferRepository(db *gorm.DB, logger *zap.Logger) TransferRepository {
	return &transferRepository{
		db:     db,
		logger: logger,
	}
}

// CreateTransfer creates a new transfer record
func (r *transferRepository) CreateTransfer(ctx context.Context, transfer *entity.Transfer) error {
	r.logger.Info("Creating transfer", zap.Uint("user_id", transfer.UserID), zap.Uint("target_user_id", transfer.TargetUserID), zap.Int64("amount", transfer.Amount))
	db := database.GetDB(ctx, r.db)
	return db.Create(transfer).Error
}

// GetTransferByID retrieves a transfer by ID
func (r *transferRepository) GetTransferByID(ctx context.Context, id uint) (*entity.Transfer, error) {
	var transfer entity.Transfer
	db := database.GetDB(ctx, r.db)
	err := db.First(&transfer, id).Error
	if err != nil {
		r.logger.Error("Failed to get transfer by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &transfer, nil
}

// GetTransfersByUserID retrieves all transfers sent by a user
func (r *transferRepository) GetTransfersByUserID(ctx context.Context, userID uint) ([]entity.Transfer, error) {
	var transfers []entity.Transfer
	db := database.GetDB(ctx, r.db)
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transfers).Error
	if err != nil {
		r.logger.Error("Failed to get transfers by user ID", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}
	return transfers, nil
}

// GetTransfersByTargetUserID retrieves all transfers received by a user
func (r *transferRepository) GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]entity.Transfer, error) {
	var transfers []entity.Transfer
	db := database.GetDB(ctx, r.db)
	err := db.Where("target_user_id = ?", targetUserID).
		Order("created_at DESC").
		Find(&transfers).Error
	if err != nil {
		r.logger.Error("Failed to get transfers by target user ID", zap.Uint("target_user_id", targetUserID), zap.Error(err))
		return nil, err
	}
	return transfers, nil
}

// ListTransfers retrieves transfers with filters and pagination
func (r *transferRepository) ListTransfers(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]entity.Transfer, int64, error) {
	var transfers []entity.Transfer
	var totalCount int64

	db := database.GetDB(ctx, r.db)
	query := db.Model(&entity.Transfer{})

	// Apply filters
	if userID, ok := filters["user_id"].(uint); ok {
		query = query.Where("user_id = ?", userID)
	}
	if targetUserID, ok := filters["target_user_id"].(uint); ok {
		query = query.Where("target_user_id = ?", targetUserID)
	}
	if source, ok := filters["source"].(string); ok {
		query = query.Where("source = ?", source)
	}
	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"].(time.Time); ok {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"].(time.Time); ok {
		query = query.Where("created_at <= ?", toDate)
	}

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to count transfers", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transfers).Error
	if err != nil {
		r.logger.Error("Failed to list transfers", zap.Error(err))
		return nil, 0, err
	}

	return transfers, totalCount, nil
}

// UpdateTransferStatus updates the status of a transfer
func (r *transferRepository) UpdateTransferStatus(ctx context.Context, id uint, status string) error {
	r.logger.Info("Updating transfer status", zap.Uint("id", id), zap.String("status", status))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.Transfer{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// GetUserTransferStats retrieves transfer statistics for a user
func (r *transferRepository) GetUserTransferStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	db := database.GetDB(ctx, r.db)

	// Sent transfers stats
	var sentStats struct {
		TotalAmount     int64
		Count           int64
		WalletAmount    int64
		PaylaterAmount  int64
		CompletedAmount int64
		PendingCount    int64
		FailedCount     int64
	}

	err := db.
		Model(&entity.Transfer{}).
		Select(`
			COALESCE(SUM(amount), 0) as total_amount,
			COUNT(*) as count,
			COALESCE(SUM(CASE WHEN source = 'wallet' THEN amount ELSE 0 END), 0) as wallet_amount,
			COALESCE(SUM(CASE WHEN source = 'paylater' THEN amount ELSE 0 END), 0) as paylater_amount,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN amount ELSE 0 END), 0) as completed_amount,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
			COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) as failed_count
		`).
		Where("user_id = ?", userID).
		Scan(&sentStats).Error

	if err != nil {
		r.logger.Error("Failed to get sent transfer stats", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Received transfers stats
	var receivedStats struct {
		TotalAmount int64
		Count       int64
	}

	err = db.Model(&entity.Transfer{}).
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as count").
		Where("target_user_id = ? AND status = 'completed'", userID).
		Scan(&receivedStats).Error

	if err != nil {
		r.logger.Error("Failed to get received transfer stats", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	stats["total_sent"] = sentStats.TotalAmount
	stats["count_sent"] = sentStats.Count
	stats["total_received"] = receivedStats.TotalAmount
	stats["count_received"] = receivedStats.Count
	stats["wallet_sent"] = sentStats.WalletAmount
	stats["paylater_sent"] = sentStats.PaylaterAmount
	stats["completed_sent"] = sentStats.CompletedAmount
	stats["pending_sent"] = sentStats.PendingCount
	stats["failed_sent"] = sentStats.FailedCount

	return stats, nil
}

// CancelTransfer cancels a pending transfer
func (r *transferRepository) CancelTransfer(ctx context.Context, id uint) error {
	r.logger.Info("Cancelling transfer", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db)
	return db.Model(&entity.Transfer{}).
		Where("id = ? AND status = ?", id, entity.TransferStatusPending).
		Updates(map[string]interface{}{
			"status":     entity.TransferStatusCancelled,
			"updated_at": time.Now(),
		}).Error
}
