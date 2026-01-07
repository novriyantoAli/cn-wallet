package repository

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPaylaterRepaymentRepository_Create(t *testing.T) {
	t.Run("should create paylater repayment successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}

		// When
		err = repo.Create(ctx, repayment)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, repayment.ID)

		// Verify
		retrieved, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint(1), retrieved.PaylaterLoanID)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(10000), retrieved.Amount)
	})

	t.Run("should create multiple repayments for same loan", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         5000,
		}

		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         5000,
		}

		// When
		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)

		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)

		// Then
		assert.NotZero(t, repayment1.ID)
		assert.NotZero(t, repayment2.ID)
		assert.NotEqual(t, repayment1.ID, repayment2.ID)
	})

	t.Run("should create repayment with large amount", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         1000000,
		}

		// When
		err = repo.Create(ctx, repayment)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, repayment.ID)

		// Verify
		retrieved, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1000000), retrieved.Amount)
	})
}

func TestPaylaterRepaymentRepository_GetByID(t *testing.T) {
	t.Run("should get repayment by ID successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         50000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)
		repaymentID := repayment.ID

		// When
		retrieved, err := repo.GetByID(ctx, repaymentID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, repaymentID, retrieved.ID)
		assert.Equal(t, uint(1), retrieved.PaylaterLoanID)
		assert.Equal(t, uint(1), retrieved.UserID)
		assert.Equal(t, int64(50000), retrieved.Amount)
	})

	t.Run("should return error when repayment ID not found", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, "repayment not found", err.Error())
	})

	t.Run("should get correct repayment from multiple records", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         20000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByID(ctx, repayment2.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, repayment2.ID, retrieved.ID)
		assert.Equal(t, int64(20000), retrieved.Amount)
	})
}

func TestPaylaterRepaymentRepository_Update(t *testing.T) {
	t.Run("should update repayment amount successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)

		// When
		repayment.Amount = 40000
		err = repo.Update(ctx, repayment)

		// Then
		assert.NoError(t, err)

		// Verify
		updated, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(40000), updated.Amount)
	})

	t.Run("should update all fields of repayment", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)

		// When
		repayment.PaylaterLoanID = 2
		repayment.UserID = 2
		repayment.Amount = 50000
		err = repo.Update(ctx, repayment)

		// Then
		assert.NoError(t, err)

		// Verify
		updated, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(2), updated.PaylaterLoanID)
		assert.Equal(t, uint(2), updated.UserID)
		assert.Equal(t, int64(50000), updated.Amount)
	})

	t.Run("should handle update with large amount", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)

		// When
		repayment.Amount = 5000000
		err = repo.Update(ctx, repayment)

		// Then
		assert.NoError(t, err)

		// Verify
		updated, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(5000000), updated.Amount)
	})

	t.Run("should handle update on non-existent repayment", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			ID:             999,
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}

		// When
		err = repo.Update(ctx, repayment)

		// Then
		assert.NoError(t, err) // GORM returns no error for update on non-existent records
	})
}

func TestPaylaterRepaymentRepository_Delete(t *testing.T) {
	t.Run("should delete repayment successfully", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         20000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)
		repaymentID := repayment.ID

		// When
		err = repo.Delete(ctx, repaymentID)

		// Then
		assert.NoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, repaymentID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("should delete correct repayment from multiple records", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         15000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         25000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)

		// When
		err = repo.Delete(ctx, repayment1.ID)

		// Then
		assert.NoError(t, err)

		// Verify - first repayment should be deleted
		_, err = repo.GetByID(ctx, repayment1.ID)
		assert.Error(t, err)

		// Verify - second repayment should still exist
		retrieved, err := repo.GetByID(ctx, repayment2.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, repayment2.ID, retrieved.ID)
	})

	t.Run("should handle delete for non-existent repayment", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// When
		err = repo.Delete(ctx, 999)

		// Then
		assert.NoError(t, err) // GORM returns no error for delete with no matching records
	})

	t.Run("should not affect other repayments when deleting", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         2,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         15000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		err = repo.Delete(ctx, repayment2.ID)

		// Then
		assert.NoError(t, err)

		// Verify - other repayments still exist
		retrieved1, err := repo.GetByID(ctx, repayment1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved1)

		retrieved3, err := repo.GetByID(ctx, repayment3.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved3)

		// Verify - deleted repayment is gone
		_, err = repo.GetByID(ctx, repayment2.ID)
		assert.Error(t, err)
	})
}

func TestPaylaterRepaymentRepository_CRUDOperations(t *testing.T) {
	t.Run("should perform complete CRUD cycle", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create
		repayment := &entity.PaylaterRepayment{
			PaylaterLoanID: 5,
			UserID:         5,
			Amount:         75000,
		}
		err = repo.Create(ctx, repayment)
		assert.NoError(t, err)
		assert.NotZero(t, repayment.ID)

		// Read
		retrieved, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, repayment.Amount, retrieved.Amount)

		// Update
		repayment.Amount = 85000
		err = repo.Update(ctx, repayment)
		assert.NoError(t, err)

		// Verify update
		updated, err := repo.GetByID(ctx, repayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(85000), updated.Amount)

		// Delete
		err = repo.Delete(ctx, repayment.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(ctx, repayment.ID)
		assert.Error(t, err)
	})

	t.Run("should handle multiple repayments independently", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}

		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         2,
			Amount:         20000,
		}

		// When - Create both
		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)

		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)

		// Then - retrieve both separately
		retrieved1, err := repo.GetByID(ctx, repayment1.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(10000), retrieved1.Amount)

		retrieved2, err := repo.GetByID(ctx, repayment2.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(20000), retrieved2.Amount)

		// Update one shouldn't affect the other
		repayment1.Amount = 15000
		err = repo.Update(ctx, repayment1)
		assert.NoError(t, err)

		updated2, err := repo.GetByID(ctx, repayment2.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(20000), updated2.Amount) // Should remain unchanged
	})
}

func TestPaylaterRepaymentRepository_GetByLoanID(t *testing.T) {
	t.Run("should get all repayments for a specific loan", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for loan 1
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         20000,
		}
		// Create repayment for loan 2
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         2,
			Amount:         15000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByLoanID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, 2, len(retrieved))

		// Verify both repayments for loan 1 are returned
		amounts := []int64{retrieved[0].Amount, retrieved[1].Amount}
		assert.Contains(t, amounts, int64(30000))
		assert.Contains(t, amounts, int64(20000))
	})

	t.Run("should return empty slice for loan with no repayments", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetByLoanID(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, 0, len(retrieved))
	})

	t.Run("should return repayments ordered by created_at DESC", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for same loan
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByLoanID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 3, len(retrieved))
		// Newest first (DESC order)
		assert.Equal(t, repayment3.ID, retrieved[0].ID)
		assert.Equal(t, repayment2.ID, retrieved[1].ID)
		assert.Equal(t, repayment1.ID, retrieved[2].ID)
	})

	t.Run("should only return repayments for specified loan", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for multiple loans
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         1,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         2,
			Amount:         15000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByLoanID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 2, len(retrieved))
		for _, r := range retrieved {
			assert.Equal(t, uint(1), r.PaylaterLoanID)
		}
	})
}

func TestPaylaterRepaymentRepository_GetByUserID(t *testing.T) {
	t.Run("should get all repayments for a specific user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for user 1
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         30000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         1,
			Amount:         20000,
		}
		// Create repayment for user 2
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 3,
			UserID:         2,
			Amount:         15000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, 2, len(retrieved))

		// Verify all repayments belong to user 1
		for _, r := range retrieved {
			assert.Equal(t, uint(1), r.UserID)
		}
	})

	t.Run("should return empty slice for user with no repayments", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// When
		retrieved, err := repo.GetByUserID(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, 0, len(retrieved))
	})

	t.Run("should return repayments ordered by created_at DESC", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for same user
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         1,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 3,
			UserID:         1,
			Amount:         30000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 3, len(retrieved))
		// Newest first (DESC order)
		assert.Equal(t, repayment3.ID, retrieved[0].ID)
		assert.Equal(t, repayment2.ID, retrieved[1].ID)
		assert.Equal(t, repayment1.ID, retrieved[2].ID)
	})

	t.Run("should only return repayments for specified user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for multiple users
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         2,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 3,
			UserID:         1,
			Amount:         15000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 2, len(retrieved))
		for _, r := range retrieved {
			assert.Equal(t, uint(1), r.UserID)
		}
	})

	t.Run("should return repayments from multiple loans for same user", func(t *testing.T) {
		// Setup
		db, err := testutil.SetupTestDB()
		assert.NoError(t, err)
		logger := testutil.NewTestLogger(t)
		repo := NewPaylaterRepaymentRepository(db, logger)
		ctx := context.Background()

		// Create repayments for same user from different loans
		repayment1 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         10000,
		}
		repayment2 := &entity.PaylaterRepayment{
			PaylaterLoanID: 2,
			UserID:         1,
			Amount:         20000,
		}
		repayment3 := &entity.PaylaterRepayment{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         5000,
		}

		err = repo.Create(ctx, repayment1)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment2)
		assert.NoError(t, err)
		err = repo.Create(ctx, repayment3)
		assert.NoError(t, err)

		// When
		retrieved, err := repo.GetByUserID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 3, len(retrieved))

		// Verify loans are from different sources
		loanIDs := []uint{retrieved[0].PaylaterLoanID, retrieved[1].PaylaterLoanID, retrieved[2].PaylaterLoanID}
		assert.Contains(t, loanIDs, uint(1))
		assert.Contains(t, loanIDs, uint(2))
	})
}
