package repository

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetByCode(ctx context.Context, code string) (*entity.Product, error)
	GetAll(ctx context.Context, filter *dto.ProductFilter) ([]entity.Product, int64, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint) error
	CodeExists(ctx context.Context, code string) (bool, error)
}

type productRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProductRepository(db *gorm.DB, logger *zap.Logger) ProductRepository {
	return &productRepository{
		db:     db,
		logger: logger,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	r.logger.Info("Creating product", zap.String("code", product.Code), zap.Uint("provider_id", product.ProviderID))
	db := database.GetDB(ctx, r.db)
	return db.Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	db := database.GetDB(ctx, r.db)
	err := db.Preload("Provider").First(&product, id).Error
	if err != nil {
		r.logger.Error("Failed to get product by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByCode(ctx context.Context, code string) (*entity.Product, error) {
	var product entity.Product
	db := database.GetDB(ctx, r.db)
	err := db.Preload("Provider").Where("code = ?", code).First(&product).Error
	if err != nil {
		r.logger.Error("Failed to get product by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, filter *dto.ProductFilter) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := database.GetDB(ctx, r.db)
	// query := db
	if filter.ProviderID > 0 {
		query = query.Where("provider_id = ?", filter.ProviderID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}

	if err := query.Model(&entity.Product{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products", zap.Error(err))
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Preload("Provider").Offset(offset).Limit(filter.PageSize).Find(&products).Error; err != nil {
		r.logger.Error("Failed to get all products", zap.Error(err))
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	r.logger.Info("Updating product", zap.Uint("id", product.ID))
	db := database.GetDB(ctx, r.db)
	return db.Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting product", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db)
	return db.Delete(&entity.Product{}, id).Error
}

func (r *productRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	db := database.GetDB(ctx, r.db)
	if err := db.Model(&entity.Product{}).Where("code = ?", code).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check if code exists", zap.String("code", code), zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
