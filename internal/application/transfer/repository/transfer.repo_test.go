package repository

import (
	"context"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTransferFixture(userID, targetUserID uint) *entity.Transfer {
	return &entity.Transfer{
		UserID:       userID,
		TargetUserID: targetUserID,
		Amount:       100000,
		Source:       entity.TransferSourceWallet,
		Status:       entity.TransferStatusCompleted,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestTransferRepository_CreateTransfer(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Create transfer", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0

		// When
		err := repo.CreateTransfer(ctx, transfer)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, transfer.ID)

		// Verify in database
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, transfer.UserID, dbTransfer.UserID)
		assert.Equal(t, transfer.TargetUserID, dbTransfer.TargetUserID)
		assert.Equal(t, transfer.Amount, dbTransfer.Amount)
		assert.Equal(t, transfer.Source, dbTransfer.Source)
		assert.Equal(t, transfer.Status, dbTransfer.Status)
	})

	t.Run("Success: Create transfer with paylater source", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Source = entity.TransferSourcePaylater
		transfer.Status = entity.TransferStatusPending

		// When
		err := repo.CreateTransfer(ctx, transfer)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, transfer.ID)
		assert.Equal(t, entity.TransferSourcePaylater, transfer.Source)
		assert.Equal(t, entity.TransferStatusPending, transfer.Status)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_GetTransferByID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Get transfer by ID", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		foundTransfer, err := repo.GetTransferByID(ctx, transfer.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, transfer.ID, foundTransfer.ID)
		assert.Equal(t, transfer.UserID, foundTransfer.UserID)
		assert.Equal(t, transfer.TargetUserID, foundTransfer.TargetUserID)
		assert.Equal(t, transfer.Amount, foundTransfer.Amount)
		assert.Equal(t, transfer.Source, foundTransfer.Source)
		assert.Equal(t, transfer.Status, foundTransfer.Status)
	})

	t.Run("Error: Transfer not found", func(t *testing.T) {
		// When
		_, err := repo.GetTransferByID(ctx, 999999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_GetTransfersByUserID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Get transfers by user ID", func(t *testing.T) {
		// Given - Create multiple transfers from same user
		userID := uint(1)
		for i := 0; i < 3; i++ {
			transfer := createTransferFixture(userID, uint(i+2))
			transfer.ID = 0
			transfer.Amount = int64((i + 1) * 100000)
			err := repo.CreateTransfer(ctx, transfer)
			require.NoError(t, err)
		}

		// Create transfers from different user
		transfer := createTransferFixture(99, 100)
		transfer.ID = 0
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		transfers, err := repo.GetTransfersByUserID(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, transfers, 3)
		for _, transfer := range transfers {
			assert.Equal(t, userID, transfer.UserID)
		}
	})

	t.Run("Success: Get empty list for user with no transfers", func(t *testing.T) {
		// When
		transfers, err := repo.GetTransfersByUserID(ctx, 888)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, transfers)
	})

	t.Run("Success: Transfers ordered by created_at DESC", func(t *testing.T) {
		// Given
		userID := uint(5)
		now := time.Now()

		transfer1 := createTransferFixture(userID, 6)
		transfer1.ID = 0
		transfer1.CreatedAt = now.Add(-2 * time.Hour)
		err := repo.CreateTransfer(ctx, transfer1)
		require.NoError(t, err)

		transfer2 := createTransferFixture(userID, 7)
		transfer2.ID = 0
		transfer2.CreatedAt = now
		err = repo.CreateTransfer(ctx, transfer2)
		require.NoError(t, err)

		// When
		transfers, err := repo.GetTransfersByUserID(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, transfers, 2)
		// Most recent should be first
		assert.True(t, transfers[0].CreatedAt.After(transfers[1].CreatedAt) || transfers[0].CreatedAt.Equal(transfers[1].CreatedAt))
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_GetTransfersByTargetUserID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Get transfers by target user ID", func(t *testing.T) {
		// Given - Create multiple transfers to same user
		targetUserID := uint(10)
		for i := 0; i < 3; i++ {
			transfer := createTransferFixture(uint(i+1), targetUserID)
			transfer.ID = 0
			transfer.Amount = int64((i + 1) * 50000)
			err := repo.CreateTransfer(ctx, transfer)
			require.NoError(t, err)
		}

		// Create transfers to different user
		transfer := createTransferFixture(99, 100)
		transfer.ID = 0
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		transfers, err := repo.GetTransfersByTargetUserID(ctx, targetUserID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, transfers, 3)
		for _, transfer := range transfers {
			assert.Equal(t, targetUserID, transfer.TargetUserID)
		}
	})

	t.Run("Success: Get empty list for user with no received transfers", func(t *testing.T) {
		// When
		transfers, err := repo.GetTransfersByTargetUserID(ctx, 777)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, transfers)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_ListTransfers(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: List all transfers with pagination", func(t *testing.T) {
		// Given - Create 5 transfers
		for i := 0; i < 5; i++ {
			transfer := createTransferFixture(uint(i+1), uint(i+10))
			transfer.ID = 0
			err := repo.CreateTransfer(ctx, transfer)
			require.NoError(t, err)
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, map[string]interface{}{}, 1, 3)

		// Then
		assert.NoError(t, err)
		assert.Len(t, transfers, 3)
		assert.Equal(t, int64(5), totalCount)
	})

	t.Run("Success: List transfers with pagination page 2", func(t *testing.T) {
		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, map[string]interface{}{}, 2, 3)

		// Then
		assert.NoError(t, err)
		assert.Len(t, transfers, 2)
		assert.Equal(t, int64(5), totalCount)
	})

	t.Run("Success: Filter by user_id", func(t *testing.T) {
		// Given
		userID := uint(20)
		for i := 0; i < 3; i++ {
			transfer := createTransferFixture(userID, uint(i+100))
			transfer.ID = 0
			err := repo.CreateTransfer(ctx, transfer)
			require.NoError(t, err)
		}

		filters := map[string]interface{}{
			"user_id": userID,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(3), totalCount)
		for _, transfer := range transfers {
			assert.Equal(t, userID, transfer.UserID)
		}
	})

	t.Run("Success: Filter by target_user_id", func(t *testing.T) {
		// Given
		targetUserID := uint(30)
		for i := 0; i < 2; i++ {
			transfer := createTransferFixture(uint(i+200), targetUserID)
			transfer.ID = 0
			err := repo.CreateTransfer(ctx, transfer)
			require.NoError(t, err)
		}

		filters := map[string]interface{}{
			"target_user_id": targetUserID,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(2), totalCount)
		for _, transfer := range transfers {
			assert.Equal(t, targetUserID, transfer.TargetUserID)
		}
	})

	t.Run("Success: Filter by source", func(t *testing.T) {
		// Given
		transfer1 := createTransferFixture(40, 41)
		transfer1.ID = 0
		transfer1.Source = entity.TransferSourceWallet
		err := repo.CreateTransfer(ctx, transfer1)
		require.NoError(t, err)

		transfer2 := createTransferFixture(42, 43)
		transfer2.ID = 0
		transfer2.Source = entity.TransferSourcePaylater
		err = repo.CreateTransfer(ctx, transfer2)
		require.NoError(t, err)

		filters := map[string]interface{}{
			"user_id": uint(42),
			"source":  entity.TransferSourcePaylater,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(1), totalCount)
		for _, transfer := range transfers {
			assert.Equal(t, entity.TransferSourcePaylater, transfer.Source)
		}
	})

	t.Run("Success: Filter by status", func(t *testing.T) {
		// Given
		transfer1 := createTransferFixture(50, 51)
		transfer1.ID = 0
		transfer1.Status = entity.TransferStatusCompleted
		err := repo.CreateTransfer(ctx, transfer1)
		require.NoError(t, err)

		transfer2 := createTransferFixture(52, 53)
		transfer2.ID = 0
		transfer2.Status = entity.TransferStatusPending
		err = repo.CreateTransfer(ctx, transfer2)
		require.NoError(t, err)

		filters := map[string]interface{}{
			"user_id": uint(52),
			"status":  entity.TransferStatusPending,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(1), totalCount)
		for _, transfer := range transfers {
			assert.Equal(t, entity.TransferStatusPending, transfer.Status)
		}
	})

	t.Run("Success: Filter by date range", func(t *testing.T) {
		// Given
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		transfer := createTransferFixture(60, 61)
		transfer.ID = 0
		transfer.CreatedAt = now
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		filters := map[string]interface{}{
			"from_date": yesterday,
			"to_date":   tomorrow,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 100)

		// Then
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, int64(1))
		for _, transfer := range transfers {
			assert.True(t, transfer.CreatedAt.After(yesterday) || transfer.CreatedAt.Equal(yesterday))
			assert.True(t, transfer.CreatedAt.Before(tomorrow) || transfer.CreatedAt.Equal(tomorrow))
		}
	})

	t.Run("Success: Multiple filters combined", func(t *testing.T) {
		// Given
		userID := uint(70)
		transfer := createTransferFixture(userID, 71)
		transfer.ID = 0
		transfer.Source = entity.TransferSourceWallet
		transfer.Status = entity.TransferStatusCompleted
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		filters := map[string]interface{}{
			"user_id": userID,
			"source":  entity.TransferSourceWallet,
			"status":  entity.TransferStatusCompleted,
		}

		// When
		transfers, totalCount, err := repo.ListTransfers(ctx, filters, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, int64(1))
		for _, transfer := range transfers {
			assert.Equal(t, userID, transfer.UserID)
			assert.Equal(t, entity.TransferSourceWallet, transfer.Source)
			assert.Equal(t, entity.TransferStatusCompleted, transfer.Status)
		}
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_UpdateTransferStatus(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Update transfer status", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Status = entity.TransferStatusPending
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		err = repo.UpdateTransferStatus(ctx, transfer.ID, string(entity.TransferStatusCompleted))

		// Then
		assert.NoError(t, err)

		// Verify in database
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.TransferStatusCompleted, dbTransfer.Status)
		assert.True(t, dbTransfer.UpdatedAt.After(transfer.UpdatedAt) || dbTransfer.UpdatedAt.Equal(transfer.UpdatedAt))
	})

	t.Run("Success: Update to failed status", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Status = entity.TransferStatusPending
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		err = repo.UpdateTransferStatus(ctx, transfer.ID, string(entity.TransferStatusFailed))

		// Then
		assert.NoError(t, err)

		// Verify in database
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.TransferStatusFailed, dbTransfer.Status)
	})

	t.Run("Success: Update non-existent transfer (no error)", func(t *testing.T) {
		// When
		err := repo.UpdateTransferStatus(ctx, 999999, string(entity.TransferStatusCompleted))

		// Then
		assert.NoError(t, err) // GORM doesn't error for updates to non-existent records
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_GetUserTransferStats(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Get user transfer statistics", func(t *testing.T) {
		// Given
		userID := uint(100)

		// Create sent transfers
		transfer1 := createTransferFixture(userID, 101)
		transfer1.ID = 0
		transfer1.Amount = 100000
		transfer1.Source = entity.TransferSourceWallet
		transfer1.Status = entity.TransferStatusCompleted
		err := repo.CreateTransfer(ctx, transfer1)
		require.NoError(t, err)

		transfer2 := createTransferFixture(userID, 102)
		transfer2.ID = 0
		transfer2.Amount = 200000
		transfer2.Source = entity.TransferSourcePaylater
		transfer2.Status = entity.TransferStatusCompleted
		err = repo.CreateTransfer(ctx, transfer2)
		require.NoError(t, err)

		transfer3 := createTransferFixture(userID, 103)
		transfer3.ID = 0
		transfer3.Amount = 50000
		transfer3.Source = entity.TransferSourceWallet
		transfer3.Status = entity.TransferStatusPending
		err = repo.CreateTransfer(ctx, transfer3)
		require.NoError(t, err)

		transfer4 := createTransferFixture(userID, 104)
		transfer4.ID = 0
		transfer4.Amount = 30000
		transfer4.Source = entity.TransferSourceWallet
		transfer4.Status = entity.TransferStatusFailed
		err = repo.CreateTransfer(ctx, transfer4)
		require.NoError(t, err)

		// Create received transfers
		receivedTransfer1 := createTransferFixture(105, userID)
		receivedTransfer1.ID = 0
		receivedTransfer1.Amount = 150000
		receivedTransfer1.Status = entity.TransferStatusCompleted
		err = repo.CreateTransfer(ctx, receivedTransfer1)
		require.NoError(t, err)

		receivedTransfer2 := createTransferFixture(106, userID)
		receivedTransfer2.ID = 0
		receivedTransfer2.Amount = 80000
		receivedTransfer2.Status = entity.TransferStatusCompleted
		err = repo.CreateTransfer(ctx, receivedTransfer2)
		require.NoError(t, err)

		// When
		stats, err := repo.GetUserTransferStats(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(380000), stats["total_sent"])     // 100000 + 200000 + 50000 + 30000
		assert.Equal(t, int64(4), stats["count_sent"])          // 4 sent transfers
		assert.Equal(t, int64(230000), stats["total_received"]) // 150000 + 80000
		assert.Equal(t, int64(2), stats["count_received"])      // 2 received transfers
		assert.Equal(t, int64(180000), stats["wallet_sent"])    // 100000 + 50000 + 30000
		assert.Equal(t, int64(200000), stats["paylater_sent"])  // 200000
		assert.Equal(t, int64(300000), stats["completed_sent"]) // 100000 + 200000
		assert.Equal(t, int64(1), stats["pending_sent"])        // 1 pending
		assert.Equal(t, int64(1), stats["failed_sent"])         // 1 failed
	})

	t.Run("Success: Get stats for user with no transfers", func(t *testing.T) {
		// When
		stats, err := repo.GetUserTransferStats(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), stats["total_sent"])
		assert.Equal(t, int64(0), stats["count_sent"])
		assert.Equal(t, int64(0), stats["total_received"])
		assert.Equal(t, int64(0), stats["count_received"])
	})

	t.Run("Success: Get stats with only sent transfers", func(t *testing.T) {
		// Given
		userID := uint(200)
		transfer := createTransferFixture(userID, 201)
		transfer.ID = 0
		transfer.Amount = 100000
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		stats, err := repo.GetUserTransferStats(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(100000), stats["total_sent"])
		assert.Equal(t, int64(1), stats["count_sent"])
		assert.Equal(t, int64(0), stats["total_received"])
		assert.Equal(t, int64(0), stats["count_received"])
	})

	t.Run("Success: Get stats with only received transfers", func(t *testing.T) {
		// Given
		userID := uint(300)
		transfer := createTransferFixture(301, userID)
		transfer.ID = 0
		transfer.Amount = 50000
		transfer.Status = entity.TransferStatusCompleted
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		stats, err := repo.GetUserTransferStats(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), stats["total_sent"])
		assert.Equal(t, int64(0), stats["count_sent"])
		assert.Equal(t, int64(50000), stats["total_received"])
		assert.Equal(t, int64(1), stats["count_received"])
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransferRepository_CancelTransfer(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransferRepository(db, logger)
	ctx := context.Background()

	t.Run("Success: Cancel pending transfer", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Status = entity.TransferStatusPending
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		err = repo.CancelTransfer(ctx, transfer.ID)

		// Then
		assert.NoError(t, err)

		// Verify in database
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.TransferStatusCancelled, dbTransfer.Status)
		assert.True(t, dbTransfer.UpdatedAt.After(transfer.UpdatedAt) || dbTransfer.UpdatedAt.Equal(transfer.UpdatedAt))
	})

	t.Run("Success: Cannot cancel completed transfer", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Status = entity.TransferStatusCompleted
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		err = repo.CancelTransfer(ctx, transfer.ID)

		// Then
		assert.NoError(t, err) // No error, but status shouldn't change

		// Verify status didn't change
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.TransferStatusCompleted, dbTransfer.Status) // Still completed
	})

	t.Run("Success: Cannot cancel failed transfer", func(t *testing.T) {
		// Given
		transfer := createTransferFixture(1, 2)
		transfer.ID = 0
		transfer.Status = entity.TransferStatusFailed
		err := repo.CreateTransfer(ctx, transfer)
		require.NoError(t, err)

		// When
		err = repo.CancelTransfer(ctx, transfer.ID)

		// Then
		assert.NoError(t, err) // No error, but status shouldn't change

		// Verify status didn't change
		var dbTransfer entity.Transfer
		err = db.First(&dbTransfer, transfer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.TransferStatusFailed, dbTransfer.Status) // Still failed
	})

	t.Run("Success: Cancel non-existent transfer (no error)", func(t *testing.T) {
		// When
		err := repo.CancelTransfer(ctx, 999999)

		// Then
		assert.NoError(t, err) // GORM doesn't error for updates to non-existent records
	})

	// Cleanup
	testutil.CleanDB(db)
}
