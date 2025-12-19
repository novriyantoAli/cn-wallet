package service

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupProductServiceWithMocks() (ProductService, *testutil.MockProductRepository, *testutil.MockProviderRepository) {
	mockProductRepo := &testutil.MockProductRepository{}
	mockProviderRepo := &testutil.MockProviderRepository{}
	logger := testutil.NewSilentLogger()
	service := NewProductService(mockProductRepo, mockProviderRepo, logger)
	return service, mockProductRepo, mockProviderRepo
}

func setupProductService(mockProductRepo *testutil.MockProductRepository, mockProviderRepo *testutil.MockProviderRepository) ProductService {
	logger := testutil.NewSilentLogger()
	service := NewProductService(mockProductRepo, mockProviderRepo, logger)
	return service
}

func TestProductService_CreateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should create product successfully", func(t *testing.T) {
		service, mockProductRepo, mockProviderRepo := setupProductServiceWithMocks()

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		req := &dto.CreateProductRequest{
			ProviderID: 1,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProviderRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockProductRepo.On("CodeExists", ctx, "PULSA10K").Return(false, nil)
		mockProductRepo.On("Create", ctx, mock.MatchedBy(func(p *entity.Product) bool {
			return p.ProviderID == req.ProviderID && p.Code == req.Code && p.PriceSell == req.PriceSell
		})).Return(nil).Run(func(args mock.Arguments) {
			product := args.Get(1).(*entity.Product)
			product.ID = 1
		})

		response, err := service.CreateProduct(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.Code, response.Code)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.PriceSell, response.PriceSell)
		mockProviderRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		req := &dto.CreateProductRequest{
			ProviderID: 999,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProviderRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		response, err := service.CreateProduct(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "provider not found", err.Error())
		mockProviderRepo.AssertExpectations(t)
	})

	t.Run("should return error when code already exists", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		req := &dto.CreateProductRequest{
			ProviderID: 1,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProviderRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockProductRepo.On("CodeExists", ctx, "PULSA10K").Return(true, nil)

		response, err := service.CreateProduct(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "code already exists", err.Error())
		mockProviderRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get product by ID successfully", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		product := &entity.Product{
			ID:         1,
			ProviderID: 1,
			Provider:   provider,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)

		response, err := service.GetProductByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, product.ID, response.ID)
		assert.Equal(t, product.Code, response.Code)
		assert.NotNil(t, response.Provider)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		response, err := service.GetProductByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product not found", err.Error())
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("should get product by code successfully", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		product := &entity.Product{
			ID:         1,
			ProviderID: 1,
			Provider:   provider,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProductRepo.On("GetByCode", ctx, "PULSA10K").Return(product, nil)

		response, err := service.GetProductByCode(ctx, "PULSA10K")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, product.Code, response.Code)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when code not found", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		mockProductRepo.On("GetByCode", ctx, "NONEXISTENT").Return(nil, gorm.ErrRecordNotFound)

		response, err := service.GetProductByCode(ctx, "NONEXISTENT")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product not found", err.Error())
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all products with pagination", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		products := []entity.Product{
			{
				ID:         1,
				ProviderID: 1,
				Provider:   provider,
				Name:       "Pulsa 10K",
				Code:       "PULSA10K",
				Category:   "Mobile",
				PriceBasic: 9000.00,
				PriceSell:  10000.00,
				Type:       "PULSA",
			},
			{
				ID:         2,
				ProviderID: 1,
				Provider:   provider,
				Name:       "Data 5GB",
				Code:       "DATA5GB",
				Category:   "Internet",
				PriceBasic: 45000.00,
				PriceSell:  50000.00,
				Type:       "DATA",
			},
		}

		filter := &dto.ProductFilter{
			Page:     1,
			PageSize: 10,
		}

		mockProductRepo.On("GetAll", ctx, filter).Return(products, int64(2), nil)

		response, err := service.GetAllProducts(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should filter products by provider ID", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		products := []entity.Product{
			{
				ID:         1,
				ProviderID: 1,
				Provider:   provider,
				Name:       "Pulsa 10K",
				Code:       "PULSA10K",
				Category:   "Mobile",
				PriceBasic: 9000.00,
				PriceSell:  10000.00,
				Type:       "PULSA",
			},
		}

		filter := &dto.ProductFilter{
			ProviderID: 1,
			Page:       1,
			PageSize:   10,
		}

		mockProviderRepo.On("GetByID", ctx, uint(1)).Return(provider, nil)
		mockProductRepo.On("GetAll", ctx, filter).Return(products, int64(1), nil)

		response, err := service.GetAllProducts(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 1)
		assert.Equal(t, int64(1), response.TotalCount)
		mockProviderRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found for filter", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		filter := &dto.ProductFilter{
			ProviderID: 999,
			Page:       1,
			PageSize:   10,
		}

		mockProviderRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		response, err := service.GetAllProducts(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "provider not found", err.Error())
		mockProviderRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should update product successfully", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		existingProduct := &entity.Product{
			ID:         1,
			ProviderID: 1,
			Provider:   provider,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		updateReq := &dto.UpdateProductRequest{
			Name:      "Pulsa 20K",
			PriceSell: 20000.00,
			Type:      "DATA",
		}

		mockProductRepo.On("GetByID", ctx, uint(1)).Return(existingProduct, nil)
		mockProductRepo.On("Update", ctx, mock.AnythingOfType("*entity.Product")).Return(nil)

		response, err := service.UpdateProduct(ctx, 1, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Pulsa 20K", response.Name)
		assert.Equal(t, 20000.00, response.PriceSell)
		assert.Equal(t, "DATA", response.Type)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		updateReq := &dto.UpdateProductRequest{
			Name: "Updated Product",
		}

		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		response, err := service.UpdateProduct(ctx, 999, updateReq)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "product not found", err.Error())
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete product successfully", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		provider := &providerEntity.Provider{
			ID:   1,
			Name: "Telkomsel",
			Code: "TSEL",
		}

		product := &entity.Product{
			ID:         1,
			ProviderID: 1,
			Provider:   provider,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Type:       "PULSA",
		}

		mockProductRepo.On("GetByID", ctx, uint(1)).Return(product, nil)
		mockProductRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := service.DeleteProduct(ctx, 1)

		assert.NoError(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		mockProductRepo := &testutil.MockProductRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupProductService(mockProductRepo, mockProviderRepo)

		mockProductRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.DeleteProduct(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
		mockProductRepo.AssertExpectations(t)
	})
}
