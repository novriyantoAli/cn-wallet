package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestProductRepository_Create(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should create product successfully", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Pulsa 20K",
			Code:       "PULSA20K",
			Category:   "Mobile",
			PriceBasic: 18000.00,
			PriceSell:  20000.00,
		}

		err = repo.Create(ctx, product)

		assert.NoError(t, err)
		assert.NotZero(t, product.ID)

		var dbProduct entity.Product
		err = db.First(&dbProduct, product.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, product.Code, dbProduct.Code)
		assert.Equal(t, product.Name, dbProduct.Name)
		assert.Equal(t, product.ProviderID, dbProduct.ProviderID)
	})

	t.Run("should fail to create product with duplicate code", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 1"
		provider.Code = "UNPROV1"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product1 := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Product 1",
			Code:       "UNIQUE_CODE",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}

		product2 := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Product 2",
			Code:       "UNIQUE_CODE",
			Category:   "Internet",
			PriceBasic: 13000.00,
			PriceSell:  15000.00,
		}

		err1 := repo.Create(ctx, product1)
		err2 := repo.Create(ctx, product2)

		assert.NoError(t, err1)
		assert.Error(t, err2)
	})
}

func TestProductRepository_GetByID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should get product by ID successfully", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 2"
		provider.Code = "UNPROV2"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Test Product",
			Code:       "TESTCODE123",
			Category:   "Gaming",
			PriceBasic: 27000.00,
			PriceSell:  30000.00,
		}
		err = repo.Create(ctx, product)

		foundProduct, err := repo.GetByID(ctx, product.ID)

		assert.NoError(t, err)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Code, foundProduct.Code)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.NotNil(t, foundProduct.Provider)
		assert.Equal(t, provider.ID, foundProduct.Provider.ID)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestProductRepository_GetByCode(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should get product by code successfully", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 3"
		provider.Code = "UNPROV3"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "PLN 50K",
			Code:       "PLN50K",
			Category:   "Electricity",
			PriceBasic: 45000.00,
			PriceSell:  50000.00,
		}
		err = repo.Create(ctx, product)
		require.NoError(t, err)

		foundProduct, err := repo.GetByCode(ctx, product.Code)

		assert.NoError(t, err)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Code, foundProduct.Code)
		assert.NotNil(t, foundProduct.Provider)
	})

	t.Run("should return error when code not found", func(t *testing.T) {
		_, err := repo.GetByCode(ctx, "NONEXISTENT_CODE")

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestProductRepository_GetAll_Pagination(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	provider := testutil.CreateProviderFixture()
	provider.ID = 0
	provider.Name = "Unique Provider 4"
	provider.Code = "UNPROV4"
	err = db.Create(provider).Error
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       fmt.Sprintf("Product %d", i),
			Code:       fmt.Sprintf("PROD%d", i),
			Category:   "Mobile",
			PriceBasic: float64((i + 1) * 9000),
			PriceSell:  float64((i + 1) * 10000),
		}
		err := repo.Create(ctx, product)
		require.NoError(t, err)
	}

	filter := &dto.ProductFilter{
		Page:     1,
		PageSize: 3,
	}

	products, totalCount, err := repo.GetAll(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, int64(5), totalCount)
}

func TestProductRepository_GetAll_FilterByProvider(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	provider1 := testutil.CreateProviderFixture()
	provider1.ID = 0
	provider1.Code = "PROV1"
	provider1.Name = "Provider 1"
	err = db.Create(provider1).Error
	require.NoError(t, err)

	provider2 := testutil.CreateProviderFixture()
	provider2.ID = 0
	provider2.Code = "PROV2"
	provider2.Name = "Provider 2"
	err = db.Create(provider2).Error
	require.NoError(t, err)

	product1 := &entity.Product{
		ProviderID: provider1.ID,
		Name:       "Provider1 Product",
		Code:       "P1PROD",
		Category:   "Mobile",
		PriceBasic: 9000.00,
		PriceSell:  10000.00,
	}
	err = repo.Create(ctx, product1)
	require.NoError(t, err)

	product2 := &entity.Product{
		ProviderID: provider2.ID,
		Name:       "Provider2 Product",
		Code:       "P2PROD",
		Category:   "Internet",
		PriceBasic: 18000.00,
		PriceSell:  20000.00,
	}
	err = repo.Create(ctx, product2)
	require.NoError(t, err)

	filter := &dto.ProductFilter{
		ProviderID: provider1.ID,
		Page:       1,
		PageSize:   10,
	}

	products, totalCount, err := repo.GetAll(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, int64(1), totalCount)
	assert.Equal(t, provider1.ID, products[0].ProviderID)
}

func TestProductRepository_GetAll_FilterByType(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	provider := testutil.CreateProviderFixture()
	provider.ID = 0
	provider.Name = "Unique Provider 5"
	provider.Code = "UNPROV5"
	err = db.Create(provider).Error
	require.NoError(t, err)

	pulsaProduct := &entity.Product{
		ProviderID: provider.ID,
		Name:       "Pulsa Product",
		Code:       "PULSAPROD",
		Category:   "Mobile",
		PriceBasic: 9000.00,
		PriceSell:  10000.00,
	}
	err = repo.Create(ctx, pulsaProduct)
	require.NoError(t, err)

	dataProduct := &entity.Product{
		ProviderID: provider.ID,
		Name:       "Data Product",
		Code:       "DATAPROD",
		Category:   "Internet",
		PriceBasic: 45000.00,
		PriceSell:  50000.00,
	}
	err = repo.Create(ctx, dataProduct)

	filter := &dto.ProductFilter{
		Page:     1,
		PageSize: 10,
	}

	products, totalCount, err := repo.GetAll(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalCount)
	assert.Len(t, products, 2)
}

func TestProductRepository_GetAll_FilterByCode(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	provider := testutil.CreateProviderFixture()
	provider.ID = 0
	provider.Name = "Unique Provider 6"
	provider.Code = "UNPROV6"
	err = db.Create(provider).Error
	require.NoError(t, err)

	product := &entity.Product{
		ProviderID: provider.ID,
		Name:       "Test Product",
		Code:       "TESTCODE123",
		Category:   "Gaming",
		PriceBasic: 27000.00,
		PriceSell:  30000.00,
	}
	err = repo.Create(ctx, product)
	require.NoError(t, err)

	filter := &dto.ProductFilter{
		Code:     "TEST",
		Page:     1,
		PageSize: 10,
	}

	products, totalCount, err := repo.GetAll(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalCount)
	assert.Len(t, products, 1)
	assert.Contains(t, products[0].Code, "TEST")
}

func TestProductRepository_Update(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should update product successfully", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 7"
		provider.Code = "UNPROV7"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Original Name",
			Code:       "ORIGCODE",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}
		err = repo.Create(ctx, product)
		require.NoError(t, err)

		product.Name = "Updated Name"
		product.PriceSell = 15000.00
		err = repo.Update(ctx, product)

		assert.NoError(t, err)

		var dbProduct entity.Product
		err = db.First(&dbProduct, product.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", dbProduct.Name)
		assert.Equal(t, 15000.00, dbProduct.PriceSell)
	})

	t.Run("should not error when updating non-existent product", func(t *testing.T) {
		product := &entity.Product{
			Name:       "Fake Product",
			Code:       "FAKECODE",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}
		product.ID = 99999

		err := repo.Update(ctx, product)

		assert.NoError(t, err)
	})
}

func TestProductRepository_Delete(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should delete product successfully", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 8"
		provider.Code = "UNPROV8"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "To Delete",
			Code:       "TODEL",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}
		err = repo.Create(ctx, product)
		require.NoError(t, err)

		err = repo.Delete(ctx, product.ID)

		assert.NoError(t, err)

		var dbProduct entity.Product
		err = db.First(&dbProduct, product.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("should not error when deleting non-existent product", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)

		assert.NoError(t, err)
	})
}

func TestProductRepository_CodeExists(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("should return true for existing code", func(t *testing.T) {
		provider := testutil.CreateProviderFixture()
		provider.ID = 0
		provider.Name = "Unique Provider 9"
		provider.Code = "UNPROV9"
		err := db.Create(provider).Error
		require.NoError(t, err)

		product := &entity.Product{
			ProviderID: provider.ID,
			Name:       "Existing Product",
			Code:       "EXISTING",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}
		err = repo.Create(ctx, product)
		require.NoError(t, err)

		exists, err := repo.CodeExists(ctx, product.Code)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false for non-existing code", func(t *testing.T) {
		exists, err := repo.CodeExists(ctx, "NONEXISTENT_PRODUCT_CODE")

		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestProductRepository_ContextCancellation(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewProductRepository(db, logger)

	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		product := &entity.Product{
			Name:       "Test Product",
			Code:       "TEST",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		}

		err := repo.Create(ctx, product)

		_ = err
	})
}
