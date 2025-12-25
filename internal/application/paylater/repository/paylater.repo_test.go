package repository

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPaylaterAccountRepository_CreateAccount(t *testing.T) {
	t.Run("should create paylater account successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}

		// When
		err = repo.CreateAccount(ctx, account)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, account.ID)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(10000), retrieved.CreditLimit)
	})

	t.Run("should return error when creating duplicate account", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}

		// Create first account
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// Try to create duplicate
		account2 := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 5000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account2)
		assert.Error(t, err)
	})
}

func TestPaylaterAccountRepository_GetAccountByUserID(t *testing.T) {
	t.Run("should get paylater account by user ID successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetAccountByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(10000), retrieved.CreditLimit)
		assert.Equal(t, entity.PaylaterStatusActive, retrieved.Status)
	})

	t.Run("should return error when account not found", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetAccountByUserID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, "paylater account not found", err.Error())
	})
}

func TestPaylaterAccountRepository_GetAccountByID(t *testing.T) {
	t.Run("should get paylater account by ID successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)
		accountID := account.ID

		// When
		retrieved, err := repo.GetAccountByID(ctx, accountID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, accountID, retrieved.ID)
		assert.Equal(t, uint(1), retrieved.UserID)
	})

	t.Run("should return error when account ID not found", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetAccountByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, "paylater account not found", err.Error())
	})
}

func TestPaylaterAccountRepository_GetForUpdate(t *testing.T) {
	t.Run("should get paylater account for update successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetForUpdate(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(10000), retrieved.CreditLimit)
	})

	t.Run("should return error when account not found for update", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetForUpdate(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, "paylater account not found", err.Error())
	})
}

func TestPaylaterAccountRepository_UpdateAccount(t *testing.T) {
	t.Run("should update paylater account successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		account.CreditLimit = 15000
		account.Status = entity.PaylaterStatusSuspended
		err = repo.UpdateAccount(ctx, account)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(15000), retrieved.CreditLimit)
		assert.Equal(t, entity.PaylaterStatusSuspended, retrieved.Status)
	})

	t.Run("should handle update on non-existent account", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			ID:          999,
			UserID:      999,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}

		// When
		err = repo.UpdateAccount(ctx, account)

		// Then
		assert.NoError(t, err) // GORM returns no error for update on non-existent records
	})
}

func TestPaylaterAccountRepository_UpdateCreditLimit(t *testing.T) {
	t.Run("should update credit limit successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		err = repo.UpdateCreditLimit(ctx, 1, 20000)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(20000), retrieved.CreditLimit)
	})

	t.Run("should handle update credit limit for non-existent user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		// When
		err = repo.UpdateCreditLimit(ctx, 999, 20000)

		// Then
		assert.NoError(t, err) // GORM returns no error for update with no matching records
	})
}

func TestPaylaterAccountRepository_UpdateStatus(t *testing.T) {
	t.Run("should update account status successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		err = repo.UpdateStatus(ctx, 1, string(entity.PaylaterStatusSuspended))

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, entity.PaylaterStatusSuspended, retrieved.Status)
	})

	t.Run("should transition from active to suspended", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		err = repo.UpdateStatus(ctx, 1, string(entity.PaylaterStatusSuspended))

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, entity.PaylaterStatusSuspended, retrieved.Status)
	})
}

func TestPaylaterAccountRepository_DeleteAccount(t *testing.T) {
	t.Run("should delete paylater account successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}
		err = repo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		// When
		err = repo.DeleteAccount(ctx, 1)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetAccountByUserID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("should handle delete for non-existent account", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		// When
		err = repo.DeleteAccount(ctx, 999)

		// Then
		assert.NoError(t, err) // GORM returns no error for delete with no matching records
	})
}

func TestPaylaterAccountRepository_ContextCancellation(t *testing.T) {
	t.Run("should handle multiple accounts independently", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterAccountRepository(db, logger)
		ctx := context.Background()

		account1 := &entity.PaylaterAccount{
			UserID:      1,
			CreditLimit: 10000,
			Status:      entity.PaylaterStatusActive,
		}

		account2 := &entity.PaylaterAccount{
			UserID:      2,
			CreditLimit: 5000,
			Status:      entity.PaylaterStatusActive,
		}

		// When
		err = repo.CreateAccount(ctx, account1)
		assert.NoError(t, err)

		err = repo.CreateAccount(ctx, account2)
		assert.NoError(t, err)

		// Then - retrieve both accounts separately
		retrieved1, err := repo.GetAccountByUserID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(10000), retrieved1.CreditLimit)

		retrieved2, err := repo.GetAccountByUserID(ctx, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(5000), retrieved2.CreditLimit)
	})
}
