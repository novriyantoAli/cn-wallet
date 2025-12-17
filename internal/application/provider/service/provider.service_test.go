package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockProviderRepo struct {
	mock.Mock
}

func (m *MockProviderRepo) Create(ctx context.Context, provider *entity.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepo) GetByID(ctx context.Context, id uint) (*entity.Provider, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Provider), args.Error(1)
}

func (m *MockProviderRepo) GetByCode(ctx context.Context, code string) (*entity.Provider, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Provider), args.Error(1)
}

func (m *MockProviderRepo) GetAll(ctx context.Context, filter *dto.ProviderFilter) ([]entity.Provider, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]entity.Provider), args.Get(1).(int64), args.Error(2)
}

func (m *MockProviderRepo) Update(ctx context.Context, provider *entity.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProviderRepo) CodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

var _ repository.ProviderRepository = (*MockProviderRepo)(nil)

func setupProviderService() (ProviderService, *MockProviderRepo) {
	mockRepo := new(MockProviderRepo)
	logger := testutil.NewSilentLogger()
	return NewProviderService(mockRepo, logger), mockRepo
}

func TestProviderService_CreateProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("should create provider successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		req := &dto.CreateProviderRequest{
			Name: "Telkomsel",
			Code: "TSEL",
			Logo: "https://example.com/logo.png",
		}

		// Mock expectations
		mockRepo.On("CodeExists", ctx, req.Code).Return(false, nil)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(p *entity.Provider) bool {
			return p.Name == req.Name && p.Code == req.Code && p.Logo == req.Logo
		})).Return(nil).Run(func(args mock.Arguments) {
			provider := args.Get(1).(*entity.Provider)
			provider.ID = 1
		})

		// When
		response, err := service.CreateProvider(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Code, response.Code)
		assert.Equal(t, req.Logo, response.Logo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when code already exists", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		req := &dto.CreateProviderRequest{
			Name: "Telkomsel",
			Code: "TSEL",
			Logo: "https://example.com/logo.png",
		}

		// Mock expectations
		mockRepo.On("CodeExists", ctx, req.Code).Return(true, nil)

		// When
		response, err := service.CreateProvider(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "code already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails during create", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		req := &dto.CreateProviderRequest{
			Name: "Telkomsel",
			Code: "TSEL",
			Logo: "https://example.com/logo.png",
		}

		// Mock expectations
		mockRepo.On("CodeExists", ctx, req.Code).Return(false, nil)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("database error"))

		// When
		response, err := service.CreateProvider(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProviderService_GetProviderByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get provider by ID successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)

		// When
		response, err := service.GetProviderByID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, provider.ID, response.ID)
		assert.Equal(t, provider.Name, response.Name)
		assert.Equal(t, provider.Code, response.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetProviderByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "provider not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProviderService_GetProviderByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("should get provider by code successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock expectations
		mockRepo.On("GetByCode", ctx, "TSEL").Return(provider, nil)

		// When
		response, err := service.GetProviderByCode(ctx, "TSEL")

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, provider.ID, response.ID)
		assert.Equal(t, provider.Code, response.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider code not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		// Mock expectations
		mockRepo.On("GetByCode", ctx, "NOTFOUND").Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetProviderByCode(ctx, "NOTFOUND")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "provider not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProviderService_GetAllProviders(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all providers with pagination", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		providers := []entity.Provider{
			{ID: 1, Name: "Provider 1", Code: "PRV1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Name: "Provider 2", Code: "PRV2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		filter := &dto.ProviderFilter{Page: 1, PageSize: 10}

		// Mock expectations
		mockRepo.On("GetAll", ctx, mock.MatchedBy(func(f *dto.ProviderFilter) bool {
			return f.Page == 1 && f.PageSize == 10
		})).Return(providers, int64(2), nil)

		// When
		response, err := service.GetAllProviders(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should use default page values", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		providers := []entity.Provider{
			{ID: 1, Name: "Provider 1", Code: "PRV1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		filter := &dto.ProviderFilter{Page: 0, PageSize: 0}

		// Mock expectations
		mockRepo.On("GetAll", ctx, mock.MatchedBy(func(f *dto.ProviderFilter) bool {
			return f.Page == 1 && f.PageSize == 10
		})).Return(providers, int64(1), nil)

		// When
		response, err := service.GetAllProviders(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		filter := &dto.ProviderFilter{Page: 1, PageSize: 10}

		// Mock expectations
		mockRepo.On("GetAll", ctx, mock.Anything).Return(nil, int64(0), errors.New("database error"))

		// When
		response, err := service.GetAllProviders(ctx, filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestProviderService_UpdateProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("should update provider successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := &dto.UpdateProviderRequest{
			Name: "Telkomsel Updated",
			Code: "TSEL",
			Logo: "https://example.com/updated-logo.png",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(p *entity.Provider) bool {
			return p.ID == 1 && p.Name == req.Name && p.Logo == req.Logo
		})).Return(nil)

		// When
		response, err := service.UpdateProvider(ctx, 1, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.Name, response.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		req := &dto.UpdateProviderRequest{
			Name: "Updated",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateProvider(ctx, 999, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "provider not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should update only provided fields", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := &dto.UpdateProviderRequest{
			Name: "Telkomsel Updated",
			Code: "",
			Logo: "",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(p *entity.Provider) bool {
			return p.Name == "Telkomsel Updated" && p.Code == "TSEL" && p.Logo == "https://example.com/logo.png"
		})).Return(nil)

		// When
		response, err := service.UpdateProvider(ctx, 1, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestProviderService_DeleteProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete provider successfully", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		// When
		err := service.DeleteProvider(ctx, 1)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.DeleteProvider(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "provider not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails during delete", func(t *testing.T) {
		// Setup
		service, mockRepo := setupProviderService()

		provider := &entity.Provider{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("database error"))

		// When
		err := service.DeleteProvider(ctx, 1)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
