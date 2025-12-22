package repository

import (
	"context"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSecurityRepository_Create(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should create user security successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1

		// When
		err := repo.Create(ctx, security)

		// Then
		assert.NoError(t, err)

		// Verify security was created in database
		var dbSecurity entity.UserSecurity
		err = db.First(&dbSecurity, "user_id = ?", security.UserID).Error
		assert.NoError(t, err)
		assert.Equal(t, security.UserID, dbSecurity.UserID)
		assert.Equal(t, security.PinHash, dbSecurity.PinHash)
		assert.Equal(t, security.FailedAttempt, dbSecurity.FailedAttempt)
	})

	t.Run("should fail to create security with duplicate user_id", func(t *testing.T) {
		// Given
		security1 := testutil.CreateUserSecurityFixture()
		security1.UserID = 2

		security2 := testutil.CreateUserSecurityFixture()
		security2.UserID = 2

		// When
		err1 := repo.Create(ctx, security1)
		err2 := repo.Create(ctx, security2)

		// Then
		assert.NoError(t, err1)
		assert.Error(t, err2) // Should fail due to primary key constraint
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_GetByUserID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should get user security by user_id successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		foundSecurity, err := repo.GetByUserID(ctx, security.UserID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, foundSecurity)
		assert.Equal(t, security.UserID, foundSecurity.UserID)
		assert.Equal(t, security.PinHash, foundSecurity.PinHash)
		assert.Equal(t, security.FailedAttempt, foundSecurity.FailedAttempt)
	})

	t.Run("should return nil when user security not found", func(t *testing.T) {
		// When
		foundSecurity, err := repo.GetByUserID(ctx, 999)

		// Then
		assert.NoError(t, err)
		assert.Nil(t, foundSecurity)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_UpdatePIN(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should update PIN successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		newPinHash := "new_hashed_pin_value"

		// When
		err = repo.UpdatePIN(ctx, security.UserID, newPinHash)

		// Then
		assert.NoError(t, err)

		// Verify PIN was updated
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Equal(t, newPinHash, updatedSecurity.PinHash)
	})

	t.Run("should not fail when updating PIN for non-existent user", func(t *testing.T) {
		// When
		err := repo.UpdatePIN(ctx, 999, "new_pin_hash")

		// Then
		assert.NoError(t, err) // GORM doesn't error on update with no rows affected
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_IncrementFailedAttempt(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should increment failed attempt successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		security.FailedAttempt = 0
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		err = repo.IncrementFailedAttempt(ctx, security.UserID)

		// Then
		assert.NoError(t, err)

		// Verify failed attempt was incremented
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Equal(t, 1, updatedSecurity.FailedAttempt)
	})

	t.Run("should increment failed attempt multiple times", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 2
		security.FailedAttempt = 0
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		for i := 0; i < 3; i++ {
			err = repo.IncrementFailedAttempt(ctx, security.UserID)
			require.NoError(t, err)
		}

		// Then
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Equal(t, 3, updatedSecurity.FailedAttempt)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_ResetFailedAttempt(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should reset failed attempt to zero", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		security.FailedAttempt = 3
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		err = repo.ResetFailedAttempt(ctx, security.UserID)

		// Then
		assert.NoError(t, err)

		// Verify failed attempt was reset and locked_until cleared
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Equal(t, 0, updatedSecurity.FailedAttempt)
		assert.Nil(t, updatedSecurity.LockedUntil)
	})

	t.Run("should clear locked_until when resetting failed attempt", func(t *testing.T) {
		// Given
		lockedUntil := time.Now().Add(15 * time.Minute)
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 2
		security.FailedAttempt = 3
		security.LockedUntil = &lockedUntil
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		err = repo.ResetFailedAttempt(ctx, security.UserID)

		// Then
		assert.NoError(t, err)

		// Verify locked_until was cleared
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Nil(t, updatedSecurity.LockedUntil)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_LockAccount(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should lock account successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		duration := 15 * time.Minute

		// When
		err = repo.LockAccount(ctx, security.UserID, duration)

		// Then
		assert.NoError(t, err)

		// Verify account is locked
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, updatedSecurity.LockedUntil)
		assert.True(t, updatedSecurity.LockedUntil.After(time.Now()))
		assert.True(t, updatedSecurity.LockedUntil.Before(time.Now().Add(duration+1*time.Second)))
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_Unlock(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should unlock account successfully", func(t *testing.T) {
		// Given
		lockedUntil := time.Now().Add(15 * time.Minute)
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		security.LockedUntil = &lockedUntil
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		err = repo.Unlock(ctx, security.UserID)

		// Then
		assert.NoError(t, err)

		// Verify account is unlocked
		updatedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Nil(t, updatedSecurity.LockedUntil)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserSecurityRepository_Delete(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserSecurityRepository(db, logger)
	ctx := context.Background()

	t.Run("should delete user security successfully", func(t *testing.T) {
		// Given
		security := testutil.CreateUserSecurityFixture()
		security.UserID = 1
		err := repo.Create(ctx, security)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, security.UserID)

		// Then
		assert.NoError(t, err)

		// Verify security was deleted
		deletedSecurity, err := repo.GetByUserID(ctx, security.UserID)
		assert.NoError(t, err)
		assert.Nil(t, deletedSecurity)
	})

	t.Run("should not fail when deleting non-existent security record", func(t *testing.T) {
		// When
		err := repo.Delete(ctx, 999)

		// Then
		assert.NoError(t, err) // GORM doesn't error on delete with no rows affected
	})

	// Cleanup
	testutil.CleanDB(db)
}
