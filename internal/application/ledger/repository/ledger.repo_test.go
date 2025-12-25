package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
)

func setupRepositoryTest(t *testing.T) (LedgerRepository, *gorm.DB) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	repo := NewLedgerRepository(db)
	return repo, db
}

func TestCreateEntry(t *testing.T) {
	t.Run("Success with Debit", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		// Create a user first
		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		entry := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
		}

		err = repo.CreateEntry(ctx, entry)

		assert.NoError(t, err)
		assert.NotZero(t, entry.ID)
		assert.False(t, entry.CreatedAt.IsZero())

		// Verify entry was saved
		var savedEntry entity.LedgerEntry
		err = db.WithContext(ctx).First(&savedEntry, "id = ?", entry.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entry.UserID, savedEntry.UserID)
		assert.Equal(t, int64(50000), savedEntry.Debit)
	})

	t.Run("Success with Credit", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 2, Email: "test2@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		entry := &entity.LedgerEntry{
			UserID:        2,
			ReferenceID:   101,
			ReferenceType: entity.ReferenceTypePayment,
			Debit:         0,
			Credit:        30000,
			AccountType:   entity.AccountTypeWallet,
		}

		err = repo.CreateEntry(ctx, entry)

		assert.NoError(t, err)
		assert.NotZero(t, entry.ID)

		var savedEntry entity.LedgerEntry
		err = db.WithContext(ctx).First(&savedEntry, "id = ?", entry.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(30000), savedEntry.Credit)
	})
}

func TestGetEntryByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create entry
		entry := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
		}
		err = repo.CreateEntry(ctx, entry)
		require.NoError(t, err)

		// Get entry
		retrieved, err := repo.GetEntryByID(ctx, entry.ID)

		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, entry.ID, retrieved.ID)
		assert.Equal(t, uint64(1), retrieved.UserID)
		assert.Equal(t, int64(50000), retrieved.Debit)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo, _ := setupRepositoryTest(t)
		ctx := context.Background()

		retrieved, err := repo.GetEntryByID(ctx, 99999)

		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestGetEntriesByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create multiple entries for same user
		entries := []entity.LedgerEntry{
			{
				UserID:        1,
				ReferenceID:   100,
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         50000,
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
			},
			{
				UserID:        1,
				ReferenceID:   101,
				ReferenceType: entity.ReferenceTypePayment,
				Debit:         0,
				Credit:        30000,
				AccountType:   entity.AccountTypeWallet,
			},
		}

		for i := range entries {
			err := repo.CreateEntry(ctx, &entries[i])
			require.NoError(t, err)
		}

		// Retrieve all entries for user
		retrieved, err := repo.GetEntriesByUserID(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, retrieved, 2)
		assert.Equal(t, uint64(1), retrieved[0].UserID)
	})

	t.Run("Empty Results", func(t *testing.T) {
		repo, _ := setupRepositoryTest(t)
		ctx := context.Background()

		retrieved, err := repo.GetEntriesByUserID(ctx, 99999)

		assert.NoError(t, err)
		assert.Len(t, retrieved, 0)
	})
}

func TestGetEntriesByReference(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create entries with same reference
		entry1 := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
		}
		entry2 := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         0,
			Credit:        25000,
			AccountType:   entity.AccountTypeWallet,
		}

		err = repo.CreateEntry(ctx, entry1)
		require.NoError(t, err)
		err = repo.CreateEntry(ctx, entry2)
		require.NoError(t, err)

		// Get entries by reference
		retrieved, err := repo.GetEntriesByReference(ctx, "transfer", 100)

		assert.NoError(t, err)
		assert.Len(t, retrieved, 2)
		assert.Equal(t, uint64(100), retrieved[0].ReferenceID)
	})
}

func TestListEntries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create multiple entries
		for i := 0; i < 5; i++ {
			entry := &entity.LedgerEntry{
				UserID:        1,
				ReferenceID:   uint64(100 + i),
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         int64(10000 * (i + 1)),
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
			}
			err := repo.CreateEntry(ctx, entry)
			require.NoError(t, err)
		}

		// List with pagination
		req := &dto.ListLedgerEntriesRequest{
			Page:     1,
			PageSize: 3,
		}

		entries, totalCount, err := repo.ListEntries(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), totalCount)
		assert.Len(t, entries, 3)
	})

	t.Run("With Filter", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create entries with different account types
		entry1 := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
		}
		entry2 := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   101,
			ReferenceType: entity.ReferenceTypePaylater,
			Debit:         30000,
			Credit:        0,
			AccountType:   entity.AccountTypePaylater,
		}

		err = repo.CreateEntry(ctx, entry1)
		require.NoError(t, err)
		err = repo.CreateEntry(ctx, entry2)
		require.NoError(t, err)

		// List with filter
		accountType := "wallet"
		req := &dto.ListLedgerEntriesRequest{
			UserID:      pointer(uint64(1)),
			AccountType: &accountType,
			Page:        1,
			PageSize:    10,
		}

		entries, totalCount, err := repo.ListEntries(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), totalCount)
		assert.Len(t, entries, 1)
		assert.Equal(t, "wallet", string(entries[0].AccountType))
	})

	t.Run("With Date Range", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create entry
		entry := &entity.LedgerEntry{
			UserID:        1,
			ReferenceID:   100,
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
		}
		err = repo.CreateEntry(ctx, entry)
		require.NoError(t, err)

		// List with date range
		fromDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		toDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
		req := &dto.ListLedgerEntriesRequest{
			FromDate: &fromDate,
			ToDate:   &toDate,
			Page:     1,
			PageSize: 10,
		}

		entries, totalCount, err := repo.ListEntries(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), totalCount)
		assert.Len(t, entries, 1)
	})
}

func TestGetUserStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 1, Email: "test@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Create entries with debits and credits
		entries := []entity.LedgerEntry{
			{
				UserID:        1,
				ReferenceID:   100,
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         50000,
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
			},
			{
				UserID:        1,
				ReferenceID:   101,
				ReferenceType: entity.ReferenceTypePayment,
				Debit:         0,
				Credit:        30000,
				AccountType:   entity.AccountTypeWallet,
			},
			{
				UserID:        1,
				ReferenceID:   102,
				ReferenceType: entity.ReferenceTypePaylater,
				Debit:         20000,
				Credit:        0,
				AccountType:   entity.AccountTypePaylater,
			},
		}

		for i := range entries {
			err := repo.CreateEntry(ctx, &entries[i])
			require.NoError(t, err)
		}

		// Get stats
		stats, err := repo.GetUserStats(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, uint64(1), stats.UserID)
		assert.Equal(t, int64(70000), stats.TotalDebit) // 50000 + 20000
		assert.Equal(t, int64(30000), stats.TotalCredit)
		assert.Equal(t, int64(3), stats.EntryCount)
		assert.Equal(t, int64(50000), stats.WalletDebit)
		assert.Equal(t, int64(30000), stats.WalletCredit)
		assert.Equal(t, int64(20000), stats.PaylaterDebit)
	})

	t.Run("No Entries", func(t *testing.T) {
		repo, db := setupRepositoryTest(t)
		ctx := context.Background()

		user := &userEntity.User{ID: 2, Email: "test2@example.com"}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Get stats for user with no entries
		stats, err := repo.GetUserStats(ctx, 2)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, uint64(2), stats.UserID)
		assert.Equal(t, int64(0), stats.TotalDebit)
		assert.Equal(t, int64(0), stats.TotalCredit)
		assert.Equal(t, int64(0), stats.EntryCount)
	})
}

// Helper functions
func pointer[T any](v T) *T {
	return &v
}
