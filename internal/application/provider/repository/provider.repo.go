package repository

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProviderRepository interface {
	Create(ctx context.Context, provider *entity.Provider) error
	GetByID(ctx context.Context, id uint) (*entity.Provider, error)
	GetByCode(ctx context.Context, code string) (*entity.Provider, error)
	GetAll(ctx context.Context, filter *dto.ProviderFilter) ([]entity.Provider, int64, error)
	Update(ctx context.Context, provider *entity.Provider) error
	Delete(ctx context.Context, id uint) error
	CodeExists(ctx context.Context, code string) (bool, error)
}

type providerRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProviderRepository(db *gorm.DB, logger *zap.Logger) ProviderRepository {
	return &providerRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new provider
func (r *providerRepository) Create(ctx context.Context, provider *entity.Provider) error {
	r.logger.Info("Creating provider", zap.String("code", provider.Code))
	db := database.GetDB(ctx, r.db)
	return db.Create(provider).Error
}

// GetByID retrieves a provider by its ID
func (r *providerRepository) GetByID(ctx context.Context, id uint) (*entity.Provider, error) {
	var provider entity.Provider
	db := database.GetDB(ctx, r.db)
	err := db.First(&provider, id).Error
	if err != nil {
		r.logger.Error("Failed to get provider by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &provider, nil
}

// GetByCode retrieves a provider by its code
func (r *providerRepository) GetByCode(ctx context.Context, code string) (*entity.Provider, error) {
	var provider entity.Provider
	db := database.GetDB(ctx, r.db)
	err := db.Where("code = ?", code).First(&provider).Error
	if err != nil {
		r.logger.Error("Failed to get provider by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}
	return &provider, nil
}

// GetAll retrieves all providers with pagination (without products for performance)
func (r *providerRepository) GetAll(ctx context.Context, filter *dto.ProviderFilter) ([]entity.Provider, int64, error) {
	var providers []entity.Provider
	var total int64

	db := database.GetDB(ctx, r.db)

	// Get total count
	if err := db.Model(&entity.Provider{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count providers", zap.Error(err))
		return nil, 0, err
	}

	// Get paginated results
	offset := (filter.Page - 1) * filter.PageSize
	if err := db.Offset(offset).Limit(filter.PageSize).Find(&providers).Error; err != nil {
		r.logger.Error("Failed to get all providers", zap.Error(err))
		return nil, 0, err
	}

	return providers, total, nil
}

// Update updates an existing provider
func (r *providerRepository) Update(ctx context.Context, provider *entity.Provider) error {
	r.logger.Info("Updating provider", zap.Uint("id", provider.ID))
	db := database.GetDB(ctx, r.db)
	return db.Save(provider).Error
}

// Delete deletes a provider by its ID
func (r *providerRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting provider", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db)
	return db.Delete(&entity.Provider{}, id).Error
}

// CodeExists checks if a provider with the given code already exists
func (r *providerRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	db := database.GetDB(ctx, r.db)
	if err := db.Model(&entity.Provider{}).Where("code = ?", code).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check if code exists", zap.String("code", code), zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
