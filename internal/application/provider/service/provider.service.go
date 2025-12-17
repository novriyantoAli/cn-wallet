package service

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProviderService interface {
	CreateProvider(ctx context.Context, req *dto.CreateProviderRequest) (*dto.ProviderResponse, error)
	GetProviderByID(ctx context.Context, id uint) (*dto.ProviderResponse, error)
	GetProviderByCode(ctx context.Context, code string) (*dto.ProviderResponse, error)
	GetAllProviders(ctx context.Context, filter *dto.ProviderFilter) (*dto.ProviderListResponse, error)
	UpdateProvider(ctx context.Context, id uint, req *dto.UpdateProviderRequest) (*dto.ProviderResponse, error)
	DeleteProvider(ctx context.Context, id uint) error
}

type providerService struct {
	repo   repository.ProviderRepository
	logger *zap.Logger
}

func NewProviderService(repo repository.ProviderRepository, logger *zap.Logger) ProviderService {
	return &providerService{
		repo:   repo,
		logger: logger,
	}
}

// CreateProvider creates a new provider
func (s *providerService) CreateProvider(ctx context.Context, req *dto.CreateProviderRequest) (*dto.ProviderResponse, error) {
	// Check if provider with same code already exists
	exists, err := s.repo.CodeExists(ctx, req.Code)
	if err != nil {
		s.logger.Error("Failed to check code existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("code already exists")
	}

	// Create provider
	provider := &entity.Provider{
		Name:      req.Name,
		Code:      req.Code,
		Logo:      req.Logo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, provider); err != nil {
		s.logger.Error("Failed to create provider", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(provider), nil
}

// GetProviderByID retrieves a provider by its ID
func (s *providerService) GetProviderByID(ctx context.Context, id uint) (*dto.ProviderResponse, error) {
	provider, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Provider not found", zap.Uint("id", id))
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	return s.entityToResponse(provider), nil
}

// GetProviderByCode retrieves a provider by its code
func (s *providerService) GetProviderByCode(ctx context.Context, code string) (*dto.ProviderResponse, error) {
	provider, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Provider not found", zap.String("code", code))
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	return s.entityToResponse(provider), nil
}

// GetAllProviders retrieves all providers with pagination
func (s *providerService) GetAllProviders(ctx context.Context, filter *dto.ProviderFilter) (*dto.ProviderListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	providers, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get all providers", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.ProviderResponse, 0, len(providers))
	for _, provider := range providers {
		responses = append(responses, *s.entityToResponse(&provider))
	}

	return &dto.ProviderListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

// UpdateProvider updates an existing provider
func (s *providerService) UpdateProvider(ctx context.Context, id uint, req *dto.UpdateProviderRequest) (*dto.ProviderResponse, error) {
	// Check if provider exists
	provider, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Provider not found", zap.Uint("id", id))
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	// Update provider fields
	if req.Name != "" {
		provider.Name = req.Name
	}
	if req.Code != "" {
		provider.Code = req.Code
	}
	if req.Logo != "" {
		provider.Logo = req.Logo
	}
	provider.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, provider); err != nil {
		s.logger.Error("Failed to update provider", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(provider), nil
}

// DeleteProvider deletes a provider by its ID
func (s *providerService) DeleteProvider(ctx context.Context, id uint) error {
	// Check if provider exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Provider not found", zap.Uint("id", id))
			return errors.New("provider not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// entityToResponse maps provider entity to response DTO
func (s *providerService) entityToResponse(provider *entity.Provider) *dto.ProviderResponse {
	response := &dto.ProviderResponse{
		ID:        provider.ID,
		Name:      provider.Name,
		Code:      provider.Code,
		Logo:      provider.Logo,
		CreatedAt: provider.CreatedAt,
		UpdatedAt: provider.UpdatedAt,
	}
	return response
}
