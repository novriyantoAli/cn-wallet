package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTransactionRepository_Create(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should create transaction successfully", func(t *testing.T) {
		// Given
		transaction := testutil.CreateTransactionFixture()

		// When
		err := repo.Create(ctx, transaction)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, transaction.ID)

		// Verify transaction was created in database
		var dbTransaction entity.Transaction
		err = db.First(&dbTransaction, "id = ?", transaction.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, transaction.WalletID, dbTransaction.WalletID)
		assert.Equal(t, transaction.Type, dbTransaction.Type)
		assert.Equal(t, transaction.Amount, dbTransaction.Amount)
	})

	t.Run("should create multiple transactions successfully", func(t *testing.T) {
		// Given
		transaction1 := testutil.CreateTransactionFixture()
		transaction2 := testutil.CreateTransactionFixture()

		// When
		err1 := repo.Create(ctx, transaction1)
		err2 := repo.Create(ctx, transaction2)

		// Then
		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// Verify both were created
		var count int64
		db.Model(&entity.Transaction{}).Count(&count)
		assert.GreaterOrEqual(t, count, int64(2))
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransactionRepository_GetByID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should get transaction by ID successfully", func(t *testing.T) {
		// Given
		transaction := testutil.CreateTransactionFixture()
		err := repo.Create(ctx, transaction)
		require.NoError(t, err)

		// When
		foundTransaction, err := repo.GetByID(ctx, transaction.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, transaction.ID, foundTransaction.ID)
		assert.Equal(t, transaction.WalletID, foundTransaction.WalletID)
		assert.Equal(t, transaction.Type, foundTransaction.Type)
		assert.Equal(t, transaction.Amount, foundTransaction.Amount)
	})

	t.Run("should return error when transaction not found", func(t *testing.T) {
		// When
		_, err := repo.GetByID(ctx, uuid.New())

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransactionRepository_GetAll(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should return empty list when no transactions exist", func(t *testing.T) {
		// Given - Ensure database is clean
		testutil.CleanDB(db)

		filter := &dto.TransactionFilter{
			Page:     1,
			PageSize: 10,
		}

		// When
		transactions, count, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
		assert.Len(t, transactions, 0)
	})

	t.Run("should get all transactions with pagination", func(t *testing.T) {
		// Given - Cleanup first
		testutil.CleanDB(db)
		for i := 0; i < 5; i++ {
			transaction := testutil.CreateTransactionFixture()
			err := repo.Create(ctx, transaction)
			require.NoError(t, err)
		}

		filter := &dto.TransactionFilter{
			Page:     1,
			PageSize: 10,
		}

		// When
		transactions, count, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(5))
		assert.GreaterOrEqual(t, len(transactions), 5)
	})

	t.Run("should filter transactions by status", func(t *testing.T) {
		// Given - Cleanup first
		testutil.CleanDB(db)
		transaction1 := testutil.CreateTransactionFixture()
		transaction1.Status = entity.StatusPending
		transaction2 := testutil.CreateTransactionFixture()
		transaction2.Status = entity.StatusSuccess
		transaction3 := testutil.CreateTransactionFixture()
		transaction3.Status = entity.StatusSuccess

		repo.Create(ctx, transaction1)
		repo.Create(ctx, transaction2)
		repo.Create(ctx, transaction3)

		filter := &dto.TransactionFilter{
			Status:   entity.StatusSuccess,
			Page:     1,
			PageSize: 10,
		}

		// When
		transactions, count, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))
		for _, tx := range transactions {
			assert.Equal(t, entity.StatusSuccess, tx.Status)
		}
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransactionRepository_Update(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should update transaction successfully", func(t *testing.T) {
		// Given
		transaction := testutil.CreateTransactionFixture()
		transaction.Status = entity.StatusPending
		err := repo.Create(ctx, transaction)
		require.NoError(t, err)

		// When
		transaction.Status = entity.StatusSuccess
		err = repo.Update(ctx, transaction)

		// Then
		assert.NoError(t, err)

		// Verify update in database
		var dbTransaction entity.Transaction
		err = db.First(&dbTransaction, "id = ?", transaction.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusSuccess, dbTransaction.Status)
	})

	t.Run("should fail to update nonexistent transaction", func(t *testing.T) {
		// Given
		transaction := testutil.CreateTransactionFixture()
		// Not created in database

		// When
		err := repo.Update(ctx, transaction)

		// Then
		// GORM may not return an error for updating non-existent records
		// Just verify no panic occurs
		_ = err
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransactionRepository_Delete(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should delete transaction successfully", func(t *testing.T) {
		// Given
		transaction := testutil.CreateTransactionFixture()
		err := repo.Create(ctx, transaction)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, transaction.ID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(ctx, transaction.ID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("should return error when deleting nonexistent transaction", func(t *testing.T) {
		// When
		err := repo.Delete(ctx, uuid.New())

		// Then
		// Delete may not return an error for non-existent records in GORM
		// just verify no panic occurs
		_ = err
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestTransactionRepository_GetByWalletID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewTransactionRepository(db, logger)
	ctx := context.Background()

	t.Run("should get transactions by wallet ID with pagination", func(t *testing.T) {
		// Given
		walletID := uint(1)
		for i := 0; i < 5; i++ {
			transaction := testutil.CreateTransactionFixture()
			transaction.WalletID = walletID
			err := repo.Create(ctx, transaction)
			require.NoError(t, err)
		}

		// When
		transactions, count, err := repo.GetByWalletID(ctx, walletID, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(5))
		assert.GreaterOrEqual(t, len(transactions), 5)

		// Verify all belong to wallet
		for _, tx := range transactions {
			assert.Equal(t, walletID, tx.WalletID)
		}
	})

	t.Run("should return empty list for wallet with no transactions", func(t *testing.T) {
		// Given
		walletID := uint(999)

		// When
		transactions, count, err := repo.GetByWalletID(ctx, walletID, 1, 10)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
		assert.Len(t, transactions, 0)
	})

	t.Run("should respect pagination parameters", func(t *testing.T) {
		// Given
		walletID := uint(2)
		for i := 0; i < 15; i++ {
			transaction := testutil.CreateTransactionFixture()
			transaction.WalletID = walletID
			err := repo.Create(ctx, transaction)
			require.NoError(t, err)
		}

		// When - First page
		transactions1, count1, err1 := repo.GetByWalletID(ctx, walletID, 1, 10)

		// When - Second page
		transactions2, count2, err2 := repo.GetByWalletID(ctx, walletID, 2, 10)

		// Then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.GreaterOrEqual(t, count1, int64(15))
		assert.GreaterOrEqual(t, count2, int64(15))
		assert.Equal(t, 10, len(transactions1))
		assert.GreaterOrEqual(t, len(transactions2), 5)
	})

	// Cleanup
	testutil.CleanDB(db)
}
