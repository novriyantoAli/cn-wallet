package repository

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WifiVoucherRepository defines the interface for wifi voucher data access.
type WifiVoucherRepository interface {
	Create(ctx context.Context, wifiVoucher *entity.WifiVoucher) error
	GetByID(ctx context.Context, id uint) (*entity.WifiVoucher, error)
	GetByCode(ctx context.Context, code string) (*entity.WifiVoucher, error)
	GetAll(ctx context.Context, filter *dto.WifiVoucherFilter) ([]entity.WifiVoucher, int64, error)
	Update(ctx context.Context, wifiVoucher *entity.WifiVoucher) error
	Delete(ctx context.Context, id uint) error
	CodeExists(ctx context.Context, code string) (bool, error)
}

type wifiVoucherRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewWifiVoucherRepository creates a new instance of WifiVoucherRepository.
func NewWifiVoucherRepository(db *gorm.DB, logger *zap.Logger) WifiVoucherRepository {
	return &wifiVoucherRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new wifi voucher in the database.
func (r *wifiVoucherRepository) Create(ctx context.Context, wifiVoucher *entity.WifiVoucher) error {
	r.logger.Info("Creating wifi voucher", zap.String("code", wifiVoucher.Code))
	db := database.GetDB(ctx, r.db)
	return db.Create(wifiVoucher).Error
}

// GetByID retrieves a wifi voucher by ID.
func (r *wifiVoucherRepository) GetByID(ctx context.Context, id uint) (*entity.WifiVoucher, error) {
	var wifiVoucher entity.WifiVoucher
	db := database.GetDB(ctx, r.db)
	err := db.Preload("SoldToUser").First(&wifiVoucher, id).Error
	if err != nil {
		r.logger.Error("Failed to get wifi voucher by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &wifiVoucher, nil
}

// GetByCode retrieves a wifi voucher by code.
func (r *wifiVoucherRepository) GetByCode(ctx context.Context, code string) (*entity.WifiVoucher, error) {
	var wifiVoucher entity.WifiVoucher
	db := database.GetDB(ctx, r.db)
	err := db.Preload("SoldToUser").Where("code = ?", code).First(&wifiVoucher).Error
	if err != nil {
		r.logger.Error("Failed to get wifi voucher by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}
	return &wifiVoucher, nil
}

// GetAll retrieves all wifi vouchers with filters and pagination.
func (r *wifiVoucherRepository) GetAll(ctx context.Context, filter *dto.WifiVoucherFilter) ([]entity.WifiVoucher, int64, error) {
	var wifiVouchers []entity.WifiVoucher
	var totalCount int64

	db := database.GetDB(ctx, r.db)
	query := db.Model(&entity.WifiVoucher{})

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.BatchID != "" {
		query = query.Where("batch_id = ?", filter.BatchID)
	}

	query.Count(&totalCount)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Preload("SoldToUser").Find(&wifiVouchers).Error
	if err != nil {
		r.logger.Error("Failed to get wifi vouchers", zap.Error(err))
		return nil, 0, err
	}

	return wifiVouchers, totalCount, nil
}

// Update updates an existing wifi voucher.
func (r *wifiVoucherRepository) Update(ctx context.Context, wifiVoucher *entity.WifiVoucher) error {
	r.logger.Info("Updating wifi voucher", zap.Uint("id", wifiVoucher.ID))
	db := database.GetDB(ctx, r.db)
	return db.Save(wifiVoucher).Error
}

// Delete soft deletes a wifi voucher.
func (r *wifiVoucherRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting wifi voucher", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db)
	return db.Delete(&entity.WifiVoucher{}, id).Error
}

// CodeExists checks if a wifi voucher code already exists.
func (r *wifiVoucherRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	db := database.GetDB(ctx, r.db)
	err := db.Model(&entity.WifiVoucher{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
