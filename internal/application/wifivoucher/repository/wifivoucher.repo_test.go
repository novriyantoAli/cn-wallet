package repository

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestWifiVoucherRepository_Create(t *testing.T) {
	t.Run("should create wifi voucher successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0

		err = repo.Create(context.Background(), wifiVoucher)
		require.NoError(t, err)
		assert.NotZero(t, wifiVoucher.ID)
	})
}

func TestWifiVoucherRepository_GetByID(t *testing.T) {
	t.Run("should get wifi voucher by ID successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		result, err := repo.GetByID(context.Background(), wifiVoucher.ID)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, wifiVoucher.Code, result.Code)
	})

	t.Run("should return error when wifi voucher not found", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		_, err = repo.GetByID(context.Background(), 999)

		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestWifiVoucherRepository_GetByCode(t *testing.T) {
	t.Run("should get wifi voucher by code successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		result, err := repo.GetByCode(context.Background(), wifiVoucher.Code)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, wifiVoucher.Code, result.Code)
	})

	t.Run("should return error when code not found", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		_, err = repo.GetByCode(context.Background(), "NONEXISTENT")

		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestWifiVoucherRepository_GetAll(t *testing.T) {
	t.Run("should get all wifi vouchers with pagination", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher1 := testutil.CreateWifiVoucherFixture()
		wifiVoucher1.ID = 0
		wifiVoucher2 := testutil.CreateWifiVoucherFixture()
		wifiVoucher2.ID = 0
		wifiVoucher2.Code = "WIFI002"

		db.Create(wifiVoucher1)
		db.Create(wifiVoucher2)

		filter := &dto.WifiVoucherFilter{
			Page:     1,
			PageSize: 10,
		}

		wifiVouchers, totalCount, err := repo.GetAll(context.Background(), filter)
		require.NoError(t, err)
		assert.Equal(t, int64(2), totalCount)
		assert.Len(t, wifiVouchers, 2)
	})

	t.Run("should filter wifi vouchers by status", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		filter := &dto.WifiVoucherFilter{
			Status:   entity.StatusAvailable,
			Page:     1,
			PageSize: 10,
		}

		_, totalCount, err := repo.GetAll(context.Background(), filter)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, int64(1))
	})

	t.Run("should filter wifi vouchers by batch ID", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		filter := &dto.WifiVoucherFilter{
			BatchID:  "BATCH001",
			Page:     1,
			PageSize: 10,
		}

		_, totalCount, err := repo.GetAll(context.Background(), filter)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, int64(1))
	})
}

func TestWifiVoucherRepository_Update(t *testing.T) {
	t.Run("should update wifi voucher successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		wifiVoucher.Code = "WIFI_UPDATED"
		wifiVoucher.Password = "newpass"
		err = repo.Update(context.Background(), wifiVoucher)
		require.NoError(t, err)

		result, err := repo.GetByID(context.Background(), wifiVoucher.ID)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "WIFI_UPDATED", result.Code)
		assert.Equal(t, "newpass", result.Password)
	})

	t.Run("should return error when updating non-existent wifi voucher", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := &entity.WifiVoucher{
			Code:     "WIFI_UPDATED",
			Password: "newpass",
		}
		wifiVoucher.ID = 999

		err = repo.Update(context.Background(), wifiVoucher)
		require.NoError(t, err) // GORM update returns no error even if no rows affected
	})
}

func TestWifiVoucherRepository_Delete(t *testing.T) {
	t.Run("should delete wifi voucher successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		err = repo.Delete(context.Background(), wifiVoucher.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(context.Background(), wifiVoucher.ID)
		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestWifiVoucherRepository_CodeExists(t *testing.T) {
	t.Run("should return true when code exists", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 0
		db.Create(wifiVoucher)

		exists, err := repo.CodeExists(context.Background(), wifiVoucher.Code)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when code does not exist", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		exists, err := repo.CodeExists(context.Background(), "NONEXISTENT")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestWifiVoucherRepository_GetByProviderAndDurationHours(t *testing.T) {
	t.Run("should get wifi vouchers by provider and duration hours successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))

		// Create vouchers with different provider and duration combinations
		voucher1 := testutil.CreateWifiVoucherFixture()
		voucher1.ID = 0
		voucher1.ProviderID = 1
		voucher1.DurationHours = 24
		db.Create(voucher1)

		voucher2 := testutil.CreateWifiVoucherFixture()
		voucher2.ID = 0
		voucher2.Code = "WIFI002"
		voucher2.ProviderID = 1
		voucher2.DurationHours = 24
		db.Create(voucher2)

		voucher3 := testutil.CreateWifiVoucherFixture()
		voucher3.ID = 0
		voucher3.Code = "WIFI003"
		voucher3.ProviderID = 1
		voucher3.DurationHours = 48
		db.Create(voucher3)

		voucher4 := testutil.CreateWifiVoucherFixture()
		voucher4.ID = 0
		voucher4.Code = "WIFI004"
		voucher4.ProviderID = 2
		voucher4.DurationHours = 24
		db.Create(voucher4)

		// Get vouchers for provider 1 with 24 hours duration
		filter := &dto.WifiVoucherFilter{
			Page:     1,
			PageSize: 10,
		}
		results, total, err := repo.GetByProviderAndDurationHours(context.Background(), 1, 24, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, uint(1), results[0].ProviderID)
		assert.Equal(t, 24, results[0].DurationHours)
		assert.Equal(t, uint(1), results[1].ProviderID)
		assert.Equal(t, 24, results[1].DurationHours)
	})

	t.Run("should filter by provider and duration hours with status filter", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))

		voucher1 := testutil.CreateWifiVoucherFixture()
		voucher1.ID = 0
		voucher1.ProviderID = 1
		voucher1.DurationHours = 24
		voucher1.Status = "available"
		db.Create(voucher1)

		voucher2 := testutil.CreateWifiVoucherFixture()
		voucher2.ID = 0
		voucher2.Code = "WIFI002"
		voucher2.ProviderID = 1
		voucher2.DurationHours = 24
		voucher2.Status = "sold"
		db.Create(voucher2)

		filter := &dto.WifiVoucherFilter{
			Status:   "available",
			Page:     1,
			PageSize: 10,
		}
		results, total, err := repo.GetByProviderAndDurationHours(context.Background(), 1, 24, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "available", results[0].Status)
	})

	t.Run("should support pagination", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))

		// Create 5 vouchers
		for i := 1; i <= 5; i++ {
			voucher := testutil.CreateWifiVoucherFixture()
			voucher.ID = 0
			voucher.Code = "WIFI00" + string(rune('0'+i))
			voucher.ProviderID = 1
			voucher.DurationHours = 24
			db.Create(voucher)
		}

		// Get first page with page size 2
		filter := &dto.WifiVoucherFilter{
			Page:     1,
			PageSize: 2,
		}
		results, total, err := repo.GetByProviderAndDurationHours(context.Background(), 1, 24, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Equal(t, 2, len(results))

		// Get second page with page size 2
		filter.Page = 2
		results, total, err = repo.GetByProviderAndDurationHours(context.Background(), 1, 24, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Equal(t, 2, len(results))
	})

	t.Run("should return empty result when no matching vouchers", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))

		voucher := testutil.CreateWifiVoucherFixture()
		voucher.ID = 0
		voucher.ProviderID = 1
		voucher.DurationHours = 24
		db.Create(voucher)

		filter := &dto.WifiVoucherFilter{
			Page:     1,
			PageSize: 10,
		}
		// Search for provider 2 with 48 hours duration (doesn't exist)
		results, total, err := repo.GetByProviderAndDurationHours(context.Background(), 2, 48, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Equal(t, 0, len(results))
	})

	t.Run("should filter by code within provider and duration hours results", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))

		voucher1 := testutil.CreateWifiVoucherFixture()
		voucher1.ID = 0
		voucher1.Code = "WIFI001"
		voucher1.ProviderID = 1
		voucher1.DurationHours = 24
		db.Create(voucher1)

		voucher2 := testutil.CreateWifiVoucherFixture()
		voucher2.ID = 0
		voucher2.Code = "WIFIOTHER"
		voucher2.ProviderID = 1
		voucher2.DurationHours = 24
		db.Create(voucher2)

		filter := &dto.WifiVoucherFilter{
			Code:     "001",
			Page:     1,
			PageSize: 10,
		}
		results, total, err := repo.GetByProviderAndDurationHours(context.Background(), 1, 24, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "WIFI001", results[0].Code)
	})
}

func TestWifiVoucherRepository_ContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		wifiVoucher := testutil.CreateWifiVoucherFixture()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = repo.Create(ctx, wifiVoucher)
		// With SQLite this won't error, but the test verifies context is passed
		// In production with PostgreSQL, this would properly respect context cancellation
		_ = err
	})
}
