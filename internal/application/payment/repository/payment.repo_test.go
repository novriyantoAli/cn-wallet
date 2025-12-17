package repository

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPaymentRepository_Create(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should create payment successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0 // Reset ID for creation

		// When
		err := repo.Create(ctx, payment)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, payment.ID)

		// Verify payment was created in database
		var dbPayment entity.Payment
		err = db.First(&dbPayment, payment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, payment.Amount, dbPayment.Amount)
		assert.Equal(t, payment.Currency, dbPayment.Currency)
		assert.Equal(t, payment.UserID, dbPayment.UserID)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_GetByID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should get payment by ID successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0
		err := repo.Create(ctx, payment)
		require.NoError(t, err)

		// When
		foundPayment, err := repo.GetByID(ctx, payment.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, payment.ID, foundPayment.ID)
		assert.Equal(t, payment.Amount, foundPayment.Amount)
		assert.Equal(t, payment.Currency, foundPayment.Currency)
		assert.Equal(t, payment.UserID, foundPayment.UserID)
	})

	t.Run("should return error when payment not found", func(t *testing.T) {
		// When
		ctx := context.Background()
		_, err := repo.GetByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_GetAll(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	// Clean up function
	cleanup := func() {
		db.Exec("DELETE FROM payments")
	}

	t.Run("should get all payments with pagination", func(t *testing.T) {
		cleanup() // Clean before test
		ctx := context.Background()
		// Given - Create multiple payments
		for i := 0; i < 5; i++ {
			payment := testutil.CreatePaymentFixture()
			payment.ID = 0
			payment.Amount = float64(100 + i)
			payment.UserID = uint(i + 1)
			err := repo.Create(ctx, payment)
			require.NoError(t, err)
		}

		filter := &dto.PaymentFilter{
			Page:     1,
			PageSize: 3,
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 3)            // Should return 3 payments due to page size
		assert.Equal(t, int64(5), totalCount) // Total count should be 5
	})

	t.Run("should filter payments by status", func(t *testing.T) {
		cleanup() // Clean before test
		ctx := context.Background()
		// Given
		payment1 := testutil.CreatePaymentFixture()
		payment1.ID = 0
		payment1.Status = entity.PaymentStatusPending
		payment1.UserID = 1
		err := repo.Create(ctx, payment1)
		require.NoError(t, err)

		payment2 := testutil.CreatePaymentFixture()
		payment2.ID = 0
		payment2.Status = entity.PaymentStatusCompleted
		payment2.UserID = 2
		err = repo.Create(ctx, payment2)
		require.NoError(t, err)

		filter := &dto.PaymentFilter{
			Status: entity.PaymentStatusPending.String(),
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, entity.PaymentStatusPending, payments[0].Status)
	})

	t.Run("should filter payments by currency", func(t *testing.T) {
		cleanup() // Clean before test
		ctx := context.Background()
		// Given
		payment1 := testutil.CreatePaymentFixture()
		payment1.ID = 0
		payment1.Currency = "USD"
		payment1.UserID = 1
		err := repo.Create(ctx, payment1)
		require.NoError(t, err)

		payment2 := testutil.CreatePaymentFixture()
		payment2.ID = 0
		payment2.Currency = "EUR"
		payment2.UserID = 2
		err = repo.Create(ctx, payment2)
		require.NoError(t, err)

		filter := &dto.PaymentFilter{
			Currency: "USD",
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, "USD", payments[0].Currency)
	})

	t.Run("should filter payments by user ID", func(t *testing.T) {
		cleanup() // Clean before test
		ctx := context.Background()
		// Given
		payment1 := testutil.CreatePaymentFixture()
		payment1.ID = 0
		payment1.UserID = 1
		err := repo.Create(ctx, payment1)
		require.NoError(t, err)

		payment2 := testutil.CreatePaymentFixture()
		payment2.ID = 0
		payment2.UserID = 2
		err = repo.Create(ctx, payment2)
		require.NoError(t, err)

		filter := &dto.PaymentFilter{
			UserID: 1,
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, uint(1), payments[0].UserID)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_Update(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should update payment successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0
		err := repo.Create(ctx, payment)
		require.NoError(t, err)

		// When
		payment.Status = entity.PaymentStatusCompleted
		payment.Description = "Updated description"
		err = repo.Update(ctx, payment)

		// Then
		assert.NoError(t, err)

		// Verify update in database
		var dbPayment entity.Payment
		err = db.First(&dbPayment, payment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, entity.PaymentStatusCompleted, dbPayment.Status)
		assert.Equal(t, "Updated description", dbPayment.Description)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_Delete(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should delete payment successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0
		err := repo.Create(ctx, payment)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, payment.ID)

		// Then
		assert.NoError(t, err)

		// Verify payment is deleted (soft delete with GORM)
		var dbPayment entity.Payment
		err = db.First(&dbPayment, payment.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_GetByUserID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should get payments by user ID successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		userID := uint(1)
		for i := 0; i < 3; i++ {
			payment := testutil.CreatePaymentFixture()
			payment.ID = 0
			payment.UserID = userID
			payment.Amount = float64(100 + i)
			err := repo.Create(ctx, payment)
			require.NoError(t, err)
		}

		// Create payment for different user
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0
		payment.UserID = 2
		err = repo.Create(ctx, payment)
		require.NoError(t, err)

		// When
		payments, err := repo.GetByUserID(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 3) // Should return only payments for user 1
		for _, p := range payments {
			assert.Equal(t, userID, p.UserID)
		}
	})

	t.Run("should return empty slice for user with no payments", func(t *testing.T) {
		// When
		ctx := context.Background()
		payments, err := repo.GetByUserID(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, payments)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_GetAllWithMultipleFilters(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	cleanup := func() {
		db.Exec("DELETE FROM payments")
	}

	t.Run("should filter by status and user ID together", func(t *testing.T) {
		cleanup()
		ctx := context.Background()
		// Given
		payment1 := testutil.CreatePaymentFixture()
		payment1.ID = 0
		payment1.Status = entity.PaymentStatusPending
		payment1.UserID = 1
		err := repo.Create(ctx, payment1)
		require.NoError(t, err)

		payment2 := testutil.CreatePaymentFixture()
		payment2.ID = 0
		payment2.Status = entity.PaymentStatusCompleted
		payment2.UserID = 1
		err = repo.Create(ctx, payment2)
		require.NoError(t, err)

		payment3 := testutil.CreatePaymentFixture()
		payment3.ID = 0
		payment3.Status = entity.PaymentStatusPending
		payment3.UserID = 2
		err = repo.Create(ctx, payment3)
		require.NoError(t, err)

		filter := &dto.PaymentFilter{
			Status: entity.PaymentStatusPending.String(),
			UserID: 1,
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, entity.PaymentStatusPending, payments[0].Status)
		assert.Equal(t, uint(1), payments[0].UserID)
	})

	t.Run("should apply pagination correctly", func(t *testing.T) {
		cleanup()
		ctx := context.Background()
		// Given - Create 10 payments
		for i := 0; i < 10; i++ {
			payment := testutil.CreatePaymentFixture()
			payment.ID = 0
			payment.UserID = uint(i + 1)
			err := repo.Create(ctx, payment)
			require.NoError(t, err)
		}

		// Test page 2
		filter := &dto.PaymentFilter{
			Page:     2,
			PageSize: 3,
		}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 3)
		assert.Equal(t, int64(10), totalCount)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_CreateAndRetrieve(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should create payment with all fields and retrieve them", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := testutil.CreatePaymentFixture()
		payment.ID = 0
		payment.Amount = 1234.56
		payment.Currency = "GBP"
		payment.Status = entity.PaymentStatusPending
		payment.Description = "Test payment with full details"
		payment.UserID = 42

		// When - Create
		err := repo.Create(ctx, payment)
		require.NoError(t, err)
		assert.NotZero(t, payment.ID)

		// When - Retrieve
		retrieved, err := repo.GetByID(ctx, payment.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, payment.ID, retrieved.ID)
		assert.Equal(t, payment.Amount, retrieved.Amount)
		assert.Equal(t, payment.Currency, retrieved.Currency)
		assert.Equal(t, payment.Status, retrieved.Status)
		assert.Equal(t, payment.Description, retrieved.Description)
		assert.Equal(t, payment.UserID, retrieved.UserID)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_GetAllWithoutFilters(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should get all payments without any filters", func(t *testing.T) {
		// Given
		ctx := context.Background()
		db.Exec("DELETE FROM payments")
		for i := 0; i < 3; i++ {
			payment := testutil.CreatePaymentFixture()
			payment.ID = 0
			payment.UserID = uint(i + 1)
			err := repo.Create(ctx, payment)
			require.NoError(t, err)
		}

		filter := &dto.PaymentFilter{}

		// When
		payments, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, payments, 3)
		assert.Equal(t, int64(3), totalCount)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_UpdateNonExistentPayment(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should update non-existent payment without error (GORM behavior)", func(t *testing.T) {
		// Given
		ctx := context.Background()
		payment := &entity.Payment{
			ID:          999,
			Amount:      100.0,
			Currency:    "USD",
			Status:      entity.PaymentStatusPending,
			Description: "Non-existent payment",
			UserID:      1,
		}

		// When
		err := repo.Update(ctx, payment)

		// Then - GORM Save doesn't error on non-existent records
		assert.NoError(t, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestPaymentRepository_DeleteNonExistentPayment(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewPaymentRepository(db, logger)

	t.Run("should delete non-existent payment without error (GORM behavior)", func(t *testing.T) {
		// Given
		ctx := context.Background()

		// When
		err := repo.Delete(ctx, 999)

		// Then - GORM Delete doesn't error on non-existent records
		assert.NoError(t, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}
