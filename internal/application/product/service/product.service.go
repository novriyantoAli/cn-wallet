package service

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/repository"
	providerRepo "github.com/novriyantoAli/cn-wallet/internal/application/provider/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
	GetProductByCode(ctx context.Context, code string) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.ProductListResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type productService struct {
	repo         repository.ProductRepository
	providerRepo providerRepo.ProviderRepository
	logger       *zap.Logger
}

func NewProductService(repo repository.ProductRepository, providerRepo providerRepo.ProviderRepository, logger *zap.Logger) ProductService {
	return &productService{
		repo:         repo,
		providerRepo: providerRepo,
		logger:       logger,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	provider, err := s.providerRepo.GetByID(ctx, req.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Provider not found", zap.Uint("provider_id", req.ProviderID))
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	exists, err := s.repo.CodeExists(ctx, req.Code)
	if err != nil {
		s.logger.Error("Failed to check code existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("code already exists")
	}

	product := &entity.Product{
		ProviderID:    req.ProviderID,
		Name:          req.Name,
		Code:          req.Code,
		Category:      req.Category,
		PriceBasic:    req.PriceBasic,
		PriceSell:     req.PriceSell,
		IsActive:      req.IsActive,
		IconURL:       req.IconURL,
		DurationHours: req.DurationHours,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	product.Provider = provider
	return s.entityToResponse(product), nil
}

func (s *productService) GetProductByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product not found", zap.Uint("id", id))
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return s.entityToResponse(product), nil
}

func (s *productService) GetProductByCode(ctx context.Context, code string) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product not found", zap.String("code", code))
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return s.entityToResponse(product), nil
}

func (s *productService) GetAllProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.ProductListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	if filter.ProviderID > 0 {
		_, err := s.providerRepo.GetByID(ctx, filter.ProviderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("Provider not found", zap.Uint("provider_id", filter.ProviderID))
				return nil, errors.New("provider not found")
			}
			return nil, err
		}
	}

	products, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get all products", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, *s.entityToResponse(&product))
	}

	return &dto.ProductListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product not found", zap.Uint("id", id))
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.PriceBasic > 0 {
		product.PriceBasic = req.PriceBasic
	}
	if req.PriceSell > 0 {
		product.PriceSell = req.PriceSell
	}
	if req.IconURL != "" {
		product.IconURL = req.IconURL
	}
	if req.DurationHours != nil && *req.DurationHours >= 0 {
		product.DurationHours = req.DurationHours
	}
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error("Failed to update product", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(product), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product not found", zap.Uint("id", id))
			return errors.New("product not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *productService) entityToResponse(product *entity.Product) *dto.ProductResponse {
	response := &dto.ProductResponse{
		ID:         product.ID,
		ProviderID: product.ProviderID,
		Name:       product.Name,
		Code:       product.Code,
		Category:   product.Category,
		PriceBasic: product.PriceBasic,
		PriceSell:  product.PriceSell,
		IsActive:   product.IsActive,
		IconURL:    product.IconURL,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}

	if product.Provider != nil {
		response.Provider = &dto.ProviderInfo{
			ID:   product.Provider.ID,
			Name: product.Provider.Name,
			Code: product.Provider.Code,
		}
	}

	return response
}
