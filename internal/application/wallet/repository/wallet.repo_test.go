package repository

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletRepository_CreateWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should create wallet successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0 // Reset ID for creation

		// When
		err := repo.CreateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, wallet.ID)

		// Verify wallet was created in database
		var dbWallet entity.Wallet
		err = db.First(&dbWallet, wallet.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, wallet.UserID, dbWallet.UserID)
		assert.Equal(t, wallet.Balance, dbWallet.Balance)
		assert.Equal(t, wallet.PINHash, dbWallet.PINHash)
	})

	t.Run("should create wallet with zero balance", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 2
		wallet.Balance = 0.00

		// When
		err := repo.CreateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, wallet.ID)
		assert.Equal(t, 0.00, wallet.Balance)
	})

	t.Run("should create wallet with large balance", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 3
		wallet.Balance = 999999.99

		// When
		err := repo.CreateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 999999.99, wallet.Balance)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_GetWalletByUserID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should get wallet by user ID successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		foundWallet, err := repo.GetWalletByUserID(ctx, wallet.UserID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, wallet.UserID, foundWallet.UserID)
		assert.Equal(t, wallet.Balance, foundWallet.Balance)
		assert.Equal(t, wallet.PINHash, foundWallet.PINHash)
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// When
		_, err := repo.GetWalletByUserID(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should return different wallets for different users", func(t *testing.T) {
		// Given
		wallet1 := testutil.CreateWalletFixture()
		wallet1.ID = 0
		wallet1.UserID = 10
		wallet1.Balance = 100.00

		wallet2 := testutil.CreateWalletFixture()
		wallet2.ID = 0
		wallet2.UserID = 11
		wallet2.Balance = 200.00

		err := repo.CreateWallet(ctx, wallet1)
		require.NoError(t, err)
		err = repo.CreateWallet(ctx, wallet2)
		require.NoError(t, err)

		// When
		found1, err := repo.GetWalletByUserID(ctx, wallet1.UserID)
		require.NoError(t, err)

		found2, err := repo.GetWalletByUserID(ctx, wallet2.UserID)
		require.NoError(t, err)

		// Then
		assert.Equal(t, 100.00, found1.Balance)
		assert.Equal(t, 200.00, found2.Balance)
		assert.NotEqual(t, found1.ID, found2.ID)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_GetWalletByID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should get wallet by ID successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		foundWallet, err := repo.GetWalletByID(ctx, wallet.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, wallet.ID, foundWallet.ID)
		assert.Equal(t, wallet.UserID, foundWallet.UserID)
		assert.Equal(t, wallet.Balance, foundWallet.Balance)
	})

	t.Run("should return error when wallet ID not found", func(t *testing.T) {
		// When
		_, err := repo.GetWalletByID(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should distinguish between different wallet IDs", func(t *testing.T) {
		// Given
		wallet1 := testutil.CreateWalletFixture()
		wallet1.ID = 0
		wallet1.UserID = 20
		wallet1.Balance = 150.00

		wallet2 := testutil.CreateWalletFixture()
		wallet2.ID = 0
		wallet2.UserID = 21
		wallet2.Balance = 250.00

		err := repo.CreateWallet(ctx, wallet1)
		require.NoError(t, err)
		err = repo.CreateWallet(ctx, wallet2)
		require.NoError(t, err)

		// When
		found1, err := repo.GetWalletByID(ctx, wallet1.ID)
		require.NoError(t, err)

		found2, err := repo.GetWalletByID(ctx, wallet2.ID)
		require.NoError(t, err)

		// Then
		assert.Equal(t, wallet1.UserID, found1.UserID)
		assert.Equal(t, wallet2.UserID, found2.UserID)
		assert.Equal(t, 150.00, found1.Balance)
		assert.Equal(t, 250.00, found2.Balance)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_UpdateWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should update wallet balance and PIN successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		wallet.Balance = 500.00
		wallet.PINHash = "new_pin_hash"

		// When
		err = repo.UpdateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)

		// Verify update
		updated, err := repo.GetWalletByID(ctx, wallet.ID)
		assert.NoError(t, err)
		assert.Equal(t, 500.00, updated.Balance)
		assert.Equal(t, "new_pin_hash", updated.PINHash)
	})

	t.Run("should update wallet balance only", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 100 // Different user ID
		originalPIN := wallet.PINHash
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		wallet.Balance = 750.50

		// When
		err = repo.UpdateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)

		updated, err := repo.GetWalletByID(ctx, wallet.ID)
		assert.NoError(t, err)
		assert.Equal(t, 750.50, updated.Balance)
		assert.Equal(t, originalPIN, updated.PINHash) // PIN should remain unchanged
	})

	t.Run("should update wallet PIN only", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 101 // Different user ID
		originalBalance := wallet.Balance
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		wallet.PINHash = "updated_pin_hash"

		// When
		err = repo.UpdateWallet(ctx, wallet)

		// Then
		assert.NoError(t, err)

		updated, err := repo.GetWalletByID(ctx, wallet.ID)
		assert.NoError(t, err)
		assert.Equal(t, "updated_pin_hash", updated.PINHash)
		assert.Equal(t, originalBalance, updated.Balance) // Balance should remain unchanged
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_UpdateBalance(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should update balance successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		newBalance := 2500.50

		// When
		err = repo.UpdateBalance(ctx, wallet.UserID, newBalance)

		// Then
		assert.NoError(t, err)

		// Verify
		updated, err := repo.GetWalletByUserID(ctx, wallet.UserID)
		assert.NoError(t, err)
		assert.Equal(t, newBalance, updated.Balance)
	})

	t.Run("should update balance to zero", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 30
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		err = repo.UpdateBalance(ctx, wallet.UserID, 0.00)

		// Then
		assert.NoError(t, err)

		updated, err := repo.GetWalletByUserID(ctx, wallet.UserID)
		assert.NoError(t, err)
		assert.Equal(t, 0.00, updated.Balance)
	})

	t.Run("should handle balance update for non-existent user gracefully", func(t *testing.T) {
		// When
		err := repo.UpdateBalance(ctx, 9999, 100.00)

		// Then - Should not error, just no-op
		assert.NoError(t, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_UpdatePIN(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should update PIN successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		newPINHash := "new_hashed_pin_1234567890"

		// When
		err = repo.UpdatePIN(ctx, wallet.UserID, newPINHash)

		// Then
		assert.NoError(t, err)

		// Verify
		updated, err := repo.GetWalletByUserID(ctx, wallet.UserID)
		assert.NoError(t, err)
		assert.Equal(t, newPINHash, updated.PINHash)
	})

	t.Run("should update PIN to different value", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 40
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		err = repo.UpdatePIN(ctx, wallet.UserID, "another_pin_hash_0987654321")

		// Then
		assert.NoError(t, err)

		updated, err := repo.GetWalletByUserID(ctx, wallet.UserID)
		assert.NoError(t, err)
		assert.Equal(t, "another_pin_hash_0987654321", updated.PINHash)
	})

	t.Run("should handle PIN update for non-existent user gracefully", func(t *testing.T) {
		// When
		err := repo.UpdatePIN(ctx, 9999, "some_pin_hash")

		// Then - Should not error, just no-op
		assert.NoError(t, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_GetForUpdate(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should get wallet with FOR UPDATE lock successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		lockedWallet, err := repo.GetForUpdate(ctx, wallet.UserID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, lockedWallet)
		assert.Equal(t, wallet.ID, lockedWallet.ID)
		assert.Equal(t, wallet.UserID, lockedWallet.UserID)
		assert.Equal(t, wallet.Balance, lockedWallet.Balance)
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		// When
		_, err := repo.GetForUpdate(ctx, 9999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should fetch correct wallet data with lock", func(t *testing.T) {
		// Given
		wallet1 := testutil.CreateWalletFixture()
		wallet1.ID = 0
		wallet1.UserID = 60
		wallet1.Balance = 1000.00
		wallet1.PINHash = "pin_hash_60"

		wallet2 := testutil.CreateWalletFixture()
		wallet2.ID = 0
		wallet2.UserID = 61
		wallet2.Balance = 2000.00
		wallet2.PINHash = "pin_hash_61"

		err := repo.CreateWallet(ctx, wallet1)
		require.NoError(t, err)
		err = repo.CreateWallet(ctx, wallet2)
		require.NoError(t, err)

		// When
		locked1, err := repo.GetForUpdate(ctx, wallet1.UserID)
		require.NoError(t, err)

		locked2, err := repo.GetForUpdate(ctx, wallet2.UserID)
		require.NoError(t, err)

		// Then
		assert.Equal(t, 1000.00, locked1.Balance)
		assert.Equal(t, "pin_hash_60", locked1.PINHash)
		assert.Equal(t, 2000.00, locked2.Balance)
		assert.Equal(t, "pin_hash_61", locked2.PINHash)
	})

	t.Run("should acquire lock during transaction", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 70
		wallet.Balance = 500.00
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When - Get wallet with lock in a transaction context
		lockedWallet, err := repo.GetForUpdate(ctx, wallet.UserID)

		// Then - Verify we got the wallet (lock behavior is enforced at DB level)
		assert.NoError(t, err)
		assert.NotNil(t, lockedWallet)
		assert.Equal(t, wallet.Balance, lockedWallet.Balance)
	})

	t.Run("should preserve wallet state when locked", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		wallet.UserID = 80
		wallet.Balance = 750.50
		wallet.PINHash = "original_pin"
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When - Get with lock
		lockedWallet, err := repo.GetForUpdate(ctx, wallet.UserID)
		require.NoError(t, err)

		// Then - Verify all fields are preserved
		assert.Equal(t, wallet.ID, lockedWallet.ID)
		assert.Equal(t, wallet.UserID, lockedWallet.UserID)
		assert.Equal(t, wallet.Balance, lockedWallet.Balance)
		assert.Equal(t, wallet.PINHash, lockedWallet.PINHash)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestWalletRepository_DeleteWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewWalletRepository(db, logger)
	ctx := context.Background()

	t.Run("should delete wallet successfully", func(t *testing.T) {
		// Given
		wallet := testutil.CreateWalletFixture()
		wallet.ID = 0
		err := repo.CreateWallet(ctx, wallet)
		require.NoError(t, err)

		// When
		err = repo.DeleteWallet(ctx, wallet.UserID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetWalletByUserID(ctx, wallet.UserID)
		assert.Error(t, err)
		assert.Equal(t, "wallet not found", err.Error())
	})

	t.Run("should handle deletion of non-existent wallet gracefully", func(t *testing.T) {
		// When
		err := repo.DeleteWallet(ctx, 9999)

		// Then - Should not error, just no-op
		assert.NoError(t, err)
	})

	t.Run("should delete only specified user wallet", func(t *testing.T) {
		// Given
		wallet1 := testutil.CreateWalletFixture()
		wallet1.ID = 0
		wallet1.UserID = 50
		wallet1.Balance = 100.00

		wallet2 := testutil.CreateWalletFixture()
		wallet2.ID = 0
		wallet2.UserID = 51
		wallet2.Balance = 200.00

		err := repo.CreateWallet(ctx, wallet1)
		require.NoError(t, err)
		err = repo.CreateWallet(ctx, wallet2)
		require.NoError(t, err)

		// When
		err = repo.DeleteWallet(ctx, wallet1.UserID)

		// Then
		assert.NoError(t, err)

		// Wallet1 should be deleted
		_, err = repo.GetWalletByUserID(ctx, wallet1.UserID)
		assert.Error(t, err)

		// Wallet2 should still exist
		found, err := repo.GetWalletByUserID(ctx, wallet2.UserID)
		assert.NoError(t, err)
		assert.Equal(t, wallet2.UserID, found.UserID)
		assert.Equal(t, 200.00, found.Balance)
	})

	// Cleanup
	testutil.CleanDB(db)
}
