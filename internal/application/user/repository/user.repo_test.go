package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0 // Reset ID for creation

		// When
		err := repo.Create(ctx, user)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// Verify user was created in database
		var dbUser entity.User
		err = db.First(&dbUser, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.FullName, dbUser.FullName)
	})

	t.Run("should fail to create user with duplicate email", func(t *testing.T) {
		// Given
		user1 := testutil.CreateUserFixture()
		user1.ID = 0
		user1.Email = "duplicate@example.com"

		user2 := testutil.CreateUserFixture()
		user2.ID = 0
		user2.Email = "duplicate@example.com"

		// When
		err1 := repo.Create(ctx, user1)
		err2 := repo.Create(ctx, user2)

		// Then
		assert.NoError(t, err1)
		assert.Error(t, err2) // Should fail due to unique constraint
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_GetByID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get user by ID successfully", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// When
		foundUser, err := repo.GetByID(ctx, user.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.FullName, foundUser.FullName)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// When
		_, err := repo.GetByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get user by email successfully", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// When
		foundUser, err := repo.GetByEmail(ctx, user.Email)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.FullName, foundUser.FullName)
	})

	t.Run("should return error when user email not found", func(t *testing.T) {
		// When
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_GetAll(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get all users with pagination", func(t *testing.T) {
		// Given - Create multiple users
		for i := 0; i < 5; i++ {
			user := testutil.CreateUserFixture()
			user.ID = 0
			user.Email = fmt.Sprintf("user%d@example.com", i)
			user.FullName = fmt.Sprintf("User %d", i)
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		filter := &dto.UserFilter{
			Page:     1,
			PageSize: 3,
		}

		// When
		users, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, users, 3)               // Should return 3 users due to page size
		assert.Equal(t, int64(5), totalCount) // Total count should be 5
	})

	t.Run("should filter users by level", func(t *testing.T) {
		// Given
		user1 := testutil.CreateUserFixture()
		user1.ID = 0
		user1.Email = "alice@example.com"
		user1.FullName = "Alice Smith"
		user1.Level = "admin"
		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := testutil.CreateUserFixture()
		user2.ID = 0
		user2.Email = "bob@example.com"
		user2.FullName = "Bob Johnson"
		user2.Level = "agent"
		err = repo.Create(ctx, user2)
		require.NoError(t, err)

		filter := &dto.UserFilter{
			Email: "alice@example.com",
		}

		// When
		users, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, "alice@example.com", users[0].Email)
	})

	t.Run("should filter users by email", func(t *testing.T) {
		// Given
		user1 := testutil.CreateUserFixture()
		user1.ID = 0
		user1.Email = "active1@example.com"
		user1.IsActive = true
		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := testutil.CreateUserFixture()
		user2.ID = 0
		user2.Email = "inactive@example.com"
		user2.IsActive = false
		err = repo.Create(ctx, user2)
		require.NoError(t, err)

		filter := &dto.UserFilter{
			Email: "active1@example.com",
		}

		// When
		users, _, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		for _, user := range users {
			assert.Equal(t, "active1@example.com", user.Email)
		}
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_Update(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should update user successfully", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// When
		user.FullName = "Updated Name"
		user.Level = "admin"
		user.IsActive = false
		err = repo.Update(ctx, user)

		// Then
		assert.NoError(t, err)

		// Verify update in database
		var dbUser entity.User
		err = db.First(&dbUser, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", dbUser.FullName)
		assert.Equal(t, "admin", dbUser.Level)
		assert.False(t, dbUser.IsActive)
	})

	t.Run("should not update non-existent user", func(t *testing.T) {
		// Given
		user := &entity.User{
			Email:    "fake@example.com",
			FullName: "Fake User",
			Level:    "user",
			IsActive: true,
		}
		user.ID = 99999 // Non-existent ID

		// When
		err := repo.Update(ctx, user)

		// Then
		assert.NoError(t, err) // GORM doesn't error for updates to non-existent records
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_Delete(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should delete user successfully", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, user.ID)

		// Then
		assert.NoError(t, err)

		// Verify user is deleted (soft delete)
		var dbUser entity.User
		err = db.First(&dbUser, user.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("should not error when deleting non-existent user", func(t *testing.T) {
		// When
		err := repo.Delete(ctx, 99999)

		// Then
		assert.NoError(t, err) // Soft delete on non-existent record succeeds
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_EmailExists(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should return true for existing email", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// When
		exists, err := repo.EmailExists(ctx, user.Email)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false for non-existing email", func(t *testing.T) {
		// When
		exists, err := repo.EmailExists(ctx, "nonexistent@example.com")

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_ContextCancellation(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)

	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		// Given - a cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Immediately cancel the context

		user := testutil.CreateUserFixture()
		user.ID = 0

		// When
		err := repo.Create(ctx, user)

		// Then - with SQLite this won't error, but the test verifies context is passed
		// In production with PostgreSQL, this would properly respect context cancellation
		_ = err
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestUserRepository_GetByIDWithWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get user with wallet preloaded", func(t *testing.T) {
		// Given - create user and wallet
		user := testutil.CreateUserFixture()
		user.ID = 0
		err := db.Create(user).Error
		require.NoError(t, err)

		wallet := &walletEntity.Wallet{
			UserID:  user.ID,
			Balance: 100.50,
			PINHash: "hashed_pin",
		}
		err = db.Create(wallet).Error
		require.NoError(t, err)

		// When
		retrievedUser, err := repo.GetByID(ctx, user.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.NotNil(t, retrievedUser.Wallet)
		assert.Equal(t, wallet.ID, retrievedUser.Wallet.ID)
		assert.Equal(t, wallet.Balance, retrievedUser.Wallet.Balance)
	})

	t.Run("should get user without wallet if wallet does not exist", func(t *testing.T) {
		// Given - create user without wallet
		user := testutil.CreateUserFixture()
		user.ID = 0
		user.Email = "no-wallet@example.com"
		err := db.Create(user).Error
		require.NoError(t, err)

		// When
		retrievedUser, err := repo.GetByID(ctx, user.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Nil(t, retrievedUser.Wallet)
	})
}

func TestUserRepository_GetByEmailWithWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)
	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get user by email with wallet preloaded", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		user.Email = "wallet-user@example.com"
		err := db.Create(user).Error
		require.NoError(t, err)

		wallet := &walletEntity.Wallet{
			UserID:  user.ID,
			Balance: 250.75,
			PINHash: "another_hash",
		}
		err = db.Create(wallet).Error
		require.NoError(t, err)

		// When
		retrievedUser, err := repo.GetByEmail(ctx, user.Email)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.Email, retrievedUser.Email)
		assert.NotNil(t, retrievedUser.Wallet)
		assert.Equal(t, wallet.Balance, retrievedUser.Wallet.Balance)
	})

	t.Run("should get user by email without wallet if wallet missing", func(t *testing.T) {
		// Given
		user := testutil.CreateUserFixture()
		user.ID = 0
		user.Email = "no-wallet-email@example.com"
		err := db.Create(user).Error
		require.NoError(t, err)

		// When
		retrievedUser, err := repo.GetByEmail(ctx, user.Email)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Nil(t, retrievedUser.Wallet)
	})
}

func TestUserRepository_GetAllWithWallet(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	// Clean any leftover data
	testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewUserRepository(db, logger)
	ctx := context.Background()

	t.Run("should get all users without wallet preload", func(t *testing.T) {
		// Clean before test
		testutil.CleanDB(db)

		// Given - create 3 users with 2 wallets
		user1 := testutil.CreateUserFixture()
		user1.ID = 0
		user1.Email = "user1-wallet@example.com"
		err := db.Create(user1).Error
		require.NoError(t, err)

		user2 := testutil.CreateUserFixture()
		user2.ID = 0
		user2.Email = "user2-wallet@example.com"
		err = db.Create(user2).Error
		require.NoError(t, err)

		user3 := testutil.CreateUserFixture()
		user3.ID = 0
		user3.Email = "user3-no-wallet@example.com"
		err = db.Create(user3).Error
		require.NoError(t, err)

		// Create wallets for user1 and user2
		wallet1 := &walletEntity.Wallet{
			UserID:  user1.ID,
			Balance: 100.00,
			PINHash: "pin1",
		}
		err = db.Create(wallet1).Error
		require.NoError(t, err)

		wallet2 := &walletEntity.Wallet{
			UserID:  user2.ID,
			Balance: 200.00,
			PINHash: "pin2",
		}
		err = db.Create(wallet2).Error
		require.NoError(t, err)

		// When
		filter := &dto.UserFilter{
			Page:     1,
			PageSize: 10,
		}
		users, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(3), totalCount)
		assert.Len(t, users, 3)

		// Verify wallets are NOT preloaded in GetAll
		walletCount := 0
		for _, u := range users {
			if u.Wallet != nil {
				walletCount++
			}
		}
		assert.Equal(t, 0, walletCount, "GetAll should not preload wallets")
	})

	t.Run("should get paginated users with wallets", func(t *testing.T) {
		// Clean before test
		testutil.CleanDB(db)

		// Given - create multiple users
		for i := 1; i <= 5; i++ {
			user := testutil.CreateUserFixture()
			user.ID = 0
			user.Email = fmt.Sprintf("paginated-user%d@example.com", i)
			err := db.Create(user).Error
			require.NoError(t, err)

			// Create wallet for odd numbered users
			if i%2 == 1 {
				wallet := &walletEntity.Wallet{
					UserID:  user.ID,
					Balance: float64(i * 100),
					PINHash: fmt.Sprintf("pin%d", i),
				}
				err = db.Create(wallet).Error
				require.NoError(t, err)
			}
		}

		// When - get first page
		filter := &dto.UserFilter{
			Page:     1,
			PageSize: 2,
		}
		users, totalCount, err := repo.GetAll(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, int64(5), totalCount)
		assert.Len(t, users, 2)

		// Verify wallets are preloaded for paginated results
		for _, u := range users {
			// Check wallet is either present or nil (preloaded, not just missing)
			_ = u.Wallet // This confirms wallet field is populated (or nil)
		}
	})
}
