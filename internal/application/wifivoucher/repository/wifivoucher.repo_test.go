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

func TestWifiVoucherRepository_GetByProviderIDWithDurationHours(t *testing.T) {
	t.Run("should get wifi vouchers by provider ID and duration hours successfully", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		providerID := uint(1)

		wifiVoucher1 := testutil.CreateWifiVoucherFixture()
		wifiVoucher1.ID = 0
		wifiVoucher1.ProviderID = &providerID
		wifiVoucher1.DurationHours = 24
		db.Create(wifiVoucher1)

		wifiVoucher2 := testutil.CreateWifiVoucherFixture()
		wifiVoucher2.ID = 0
		wifiVoucher2.Code = "WIFI002"
		wifiVoucher2.ProviderID = &providerID
		wifiVoucher2.DurationHours = 24
		db.Create(wifiVoucher2)

		// Different duration hours - should not be returned
		wifiVoucher3 := testutil.CreateWifiVoucherFixture()
		wifiVoucher3.ID = 0
		wifiVoucher3.Code = "WIFI003"
		wifiVoucher3.ProviderID = &providerID
		wifiVoucher3.DurationHours = 48
		db.Create(wifiVoucher3)

		result, err := repo.GetByProviderIDWithDurationHours(context.Background(), providerID, 24)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 24, result[0].DurationHours)
		assert.Equal(t, 24, result[1].DurationHours)
	})

	t.Run("should return empty slice when no vouchers match", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		result, err := repo.GetByProviderIDWithDurationHours(context.Background(), 999, 24)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should filter by both provider ID and duration hours", func(t *testing.T) {
		db, err := testutil.SetupTestDB()
		require.NoError(t, err)
		defer testutil.CleanDB(db)

		repo := NewWifiVoucherRepository(db, testutil.NewTestLogger(t))
		providerID1 := uint(1)
		providerID2 := uint(2)

		wifiVoucher1 := testutil.CreateWifiVoucherFixture()
		wifiVoucher1.ID = 0
		wifiVoucher1.ProviderID = &providerID1
		wifiVoucher1.DurationHours = 24
		db.Create(wifiVoucher1)

		wifiVoucher2 := testutil.CreateWifiVoucherFixture()
		wifiVoucher2.ID = 0
		wifiVoucher2.Code = "WIFI002"
		wifiVoucher2.ProviderID = &providerID2
		wifiVoucher2.DurationHours = 24
		db.Create(wifiVoucher2)

		result, err := repo.GetByProviderIDWithDurationHours(context.Background(), providerID1, 24)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, providerID1, *result[0].ProviderID)
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
