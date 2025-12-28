package service

import (
	"context"
	"errors"
	"testing"

	providerDto "github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockProviderService struct {
	mock.Mock
}

func (m *MockProviderService) CreateProvider(ctx context.Context, req *providerDto.CreateProviderRequest) (*providerDto.ProviderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerDto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetProviderByID(ctx context.Context, id uint) (*providerDto.ProviderResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerDto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetProviderByCode(ctx context.Context, code string) (*providerDto.ProviderResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerDto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetAllProviders(ctx context.Context, filter *providerDto.ProviderFilter) (*providerDto.ProviderListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerDto.ProviderListResponse), args.Error(1)
}

func (m *MockProviderService) UpdateProvider(ctx context.Context, id uint, req *providerDto.UpdateProviderRequest) (*providerDto.ProviderResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providerDto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) DeleteProvider(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupUserServiceWithMocks() (UserService, *testutil.MockUserRepository, *testutil.MockProviderRepository) {
	mockRepo := &testutil.MockUserRepository{}
	mockProviderRepo := &testutil.MockProviderRepository{}
	logger := testutil.NewSilentLogger()
	service := NewUserService(mockRepo, mockProviderRepo, logger)
	return service, mockRepo, mockProviderRepo
}

func setupUserService(mockRepo *testutil.MockUserRepository, mockProviderRepo *testutil.MockProviderRepository) UserService {
	logger := testutil.NewSilentLogger()
	service := NewUserService(mockRepo, mockProviderRepo, logger)
	return service
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		// Setup
		service, mockRepo, _ := setupUserServiceWithMocks()

		req := &dto.CreateUserRequest{
			Email:    "test@example.com",
			FullName: "Test User",
		}

		// Mock expectations
		mockRepo.On("EmailExists", ctx, req.Email).Return(false, nil)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
			return u.Email == req.Email && u.Wallet != nil && u.Wallet.Balance == 0
		})).Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(1).(*entity.User)
			user.ID = 1
			if user.Wallet != nil {
				user.Wallet.UserID = user.ID
				user.Wallet.ID = 1
			}
		})

		// When
		response, err := service.CreateUser(ctx, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.FullName, response.FullName)
		assert.Equal(t, req.Email, response.Email)
		assert.NotNil(t, response.Wallet)
		assert.Equal(t, float64(0), response.Wallet.Balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		req := &dto.CreateUserRequest{
			Email:    "existing@example.com",
			FullName: "Existing User",
		}

		// Mock expectations
		mockRepo.On("EmailExists", ctx, req.Email).Return(true, nil)

		// When
		response, err := service.CreateUser(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		req := &dto.CreateUserRequest{
			Email:    "test@example.com",
			FullName: "Test User",
		}

		// Mock expectations
		mockRepo.On("EmailExists", ctx, req.Email).Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

		// When
		response, err := service.CreateUser(ctx, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by ID successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(user, nil)

		// When
		response, err := service.GetUserByID(ctx, 1)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.FullName, response.FullName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetUserByID(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by email successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		// Mock expectations
		mockRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

		// When
		response, err := service.GetUserByEmail(ctx, "test@example.com")

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, user.Email, response.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when email not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		// Mock expectations
		mockRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetUserByEmail(ctx, "nonexistent@example.com")

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("should get users with pagination", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		users := []entity.User{
			{
				ID:       1,
				Email:    "user1@example.com",
				FullName: "User 1",
				Level:    "user",
				IsActive: true,
			},
			{
				ID:       2,
				Email:    "user2@example.com",
				FullName: "User 2",
				Level:    "user",
				IsActive: true,
			},
		}

		filter := &dto.UserFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", ctx, filter).Return(users, int64(2), nil)

		// When
		response, err := service.GetUsers(ctx, filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		filter := &dto.UserFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", ctx, filter).Return(nil, int64(0), errors.New("database error"))

		// When
		response, err := service.GetUsers(ctx, filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should update user successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserRequest{
			FullName: "Updated User",
			Level:    "admin",
			IsActive: false,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		// When
		response, err := service.UpdateUser(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Updated User", response.FullName)
		assert.Equal(t, "admin", response.Level)
		assert.False(t, response.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should update user with provider info (graceful failure)", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserRequest{
			FullName:     "Provider User",
			Level:        "provider",
			IsActive:     true,
			ProviderName: "CN Hotspot",
			ProviderCode: "CNHOTSPOT",
			ProviderLogo: "https://example.com/logo.png",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		// When
		response, err := service.UpdateUser(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "provider", response.Level)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		updateReq := &dto.UpdateUserRequest{
			FullName: "Updated User",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateUser(ctx, 999, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete user successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		// When
		err := service.DeleteUser(ctx, 1)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.DeleteUser(ctx, 999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user response when user already has provider", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		providerID := uint(5)
		existingUser := &entity.User{
			ID:         1,
			Email:      "test@example.com",
			FullName:   "Test User",
			Level:      "user",
			IsActive:   true,
			ProviderID: &providerID,
		}

		updateReq := &dto.UpdateUserProviderRequest{
			ProviderID: 10,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)

		// When
		response, err := service.UpdateUserProvider(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should set provider ID when user has no provider", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		provider := &providerEntity.Provider{ID: 10}

		updateReq := &dto.UpdateUserProviderRequest{
			ProviderID: 10,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockProviderRepo.On("GetByID", ctx, uint(10)).Return(provider, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
			return u.ID == 1 && u.ProviderID != nil && *u.ProviderID == 10
		})).Return(nil)

		// When
		response, err := service.UpdateUserProvider(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockRepo.AssertExpectations(t)
		mockProviderRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		updateReq := &dto.UpdateUserProviderRequest{
			ProviderID: 10,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateUserProvider(ctx, 999, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when provider not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserProviderRequest{
			ProviderID: 999,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockProviderRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		_, err := service.UpdateUserProvider(ctx, 1, updateReq)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "provider not found", err.Error())
		mockRepo.AssertExpectations(t)
		mockProviderRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		provider := &providerEntity.Provider{ID: 10}

		updateReq := &dto.UpdateUserProviderRequest{
			ProviderID: 10,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockProviderRepo.On("GetByID", ctx, uint(10)).Return(provider, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("update failed"))

		// When
		response, err := service.UpdateUserProvider(ctx, 1, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
		mockProviderRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("should update user level successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserLevelRequest{
			Level: "provider",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
			return u.ID == 1 && u.Level == "provider"
		})).Return(nil)

		// When
		response, err := service.UpdateLevel(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "provider", response.Level)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should update level to admin", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserLevelRequest{
			Level: "admin",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
			return u.ID == 1 && u.Level == "admin"
		})).Return(nil)

		// When
		response, err := service.UpdateLevel(ctx, 1, updateReq)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "admin", response.Level)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		updateReq := &dto.UpdateUserLevelRequest{
			Level: "provider",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateLevel(ctx, 999, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockUserRepository{}
		mockProviderRepo := &testutil.MockProviderRepository{}
		service := setupUserService(mockRepo, mockProviderRepo)

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		updateReq := &dto.UpdateUserLevelRequest{
			Level: "provider",
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("update failed"))

		// When
		response, err := service.UpdateLevel(ctx, 1, updateReq)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}
