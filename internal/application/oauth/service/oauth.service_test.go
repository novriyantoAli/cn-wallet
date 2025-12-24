package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/oauth/dto"
	userdto "github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/repository"
	"github.com/novriyantoAli/cn-wallet/internal/config"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, filter *userdto.UserFilter) ([]entity.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func newMockUserRepository() repository.UserRepository {
	return &MockUserRepository{}
}

func newMockJWTManager() *jwt.JWTManager {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret-key-for-jwt-generation-in-tests",
			Expiry:    time.Hour,
		},
	}
	return jwt.NewJWTManager(cfg)
}

func TestOAuthService_GetAuthorizationURL(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := newMockUserRepository()
	jwtManager := newMockJWTManager()
	config := OAuthConfig{
		Google: GoogleOAuthConfig{
			ClientID:    "test-client-id",
			RedirectURI: "http://localhost:8080/callback",
		},
		Github: GithubOAuthConfig{
			ClientID:    "test-github-id",
			RedirectURI: "http://localhost:8080/callback",
		},
	}

	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should generate Google authorization URL", func(t *testing.T) {
		resp, err := service.GetAuthorizationURL(context.Background(), dto.GoogleProvider, "http://localhost:3000/callback")
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AuthorizationURL)
		assert.NotEmpty(t, resp.State)
		assert.Contains(t, resp.AuthorizationURL, "accounts.google.com")
		assert.Contains(t, resp.AuthorizationURL, "test-client-id")

		// Verify state is a valid JWT token
		claims, err := jwtManager.VerifyToken(resp.State)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "oauth_state", claims.Level)
	})

	t.Run("should generate Github authorization URL", func(t *testing.T) {
		resp, err := service.GetAuthorizationURL(context.Background(), dto.GithubProvider, "http://localhost:3000/callback")
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AuthorizationURL)
		assert.NotEmpty(t, resp.State)
		assert.Contains(t, resp.AuthorizationURL, "github.com")
		assert.Contains(t, resp.AuthorizationURL, "test-github-id")
	})

	t.Run("should return error for unsupported provider", func(t *testing.T) {
		resp, err := service.GetAuthorizationURL(context.Background(), dto.OAuthProvider("invalid"), "http://localhost:3000/callback")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "unsupported OAuth provider")
	})
}

func TestOAuthService_ExchangeCodeForToken(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := newMockUserRepository()
	jwtManager := newMockJWTManager()
	config := OAuthConfig{
		Google: GoogleOAuthConfig{
			ClientID:     "test-google-client-id",
			ClientSecret: "test-google-client-secret",
			RedirectURI:  "http://localhost:8080/callback",
		},
	}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should return error for empty code", func(t *testing.T) {
		resp, err := service.ExchangeCodeForToken(context.Background(), dto.GoogleProvider, "")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "authorization code is required", err.Error())
	})

	t.Run("should handle exchange error gracefully", func(t *testing.T) {
		// Invalid code will fail at the OAuth provider
		resp, err := service.ExchangeCodeForToken(context.Background(), dto.GoogleProvider, "invalid-test-code")
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Error will come from OAuth provider, not from our validation
		assert.Contains(t, err.Error(), "failed to exchange code for token")
	})
}

func TestOAuthService_GetUserInfo(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := newMockUserRepository()
	jwtManager := newMockJWTManager()
	config := OAuthConfig{}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should return error for empty access token", func(t *testing.T) {
		resp, err := service.GetUserInfo(context.Background(), dto.GoogleProvider, "")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "access token is required", err.Error())
	})

	t.Run("should handle invalid access token", func(t *testing.T) {
		resp, err := service.GetUserInfo(context.Background(), dto.GoogleProvider, "invalid-token")
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Error will come from the real OAuth provider API
		assert.Contains(t, err.Error(), "failed to fetch user info")
	})
}

func TestOAuthService_Authenticate(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := newMockUserRepository()
	jwtManager := newMockJWTManager()
	config := OAuthConfig{
		Google: GoogleOAuthConfig{
			ClientID:     "test-google-client-id",
			ClientSecret: "test-google-client-secret",
			RedirectURI:  "http://localhost:8080/callback",
		},
	}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should return error for empty code", func(t *testing.T) {
		resp, err := service.Authenticate(context.Background(), dto.GoogleProvider, "", "state123")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "authorization code is required", err.Error())
	})

	t.Run("should return error for empty state", func(t *testing.T) {
		resp, err := service.Authenticate(context.Background(), dto.GoogleProvider, "code123", "")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "state parameter is required", err.Error())
	})

	t.Run("should return error for invalid state token", func(t *testing.T) {
		resp, err := service.Authenticate(context.Background(), dto.GoogleProvider, "code123", "invalid-state")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid or expired state token")
	})

	t.Run("should handle invalid auth code with valid state", func(t *testing.T) {
		// Generate a valid state token for oauth flow
		validState, err := jwtManager.GenerateToken(0, "", "oauth_state")
		assert.NoError(t, err)

		// Now test with invalid code but valid state
		resp, err := service.Authenticate(context.Background(), dto.GoogleProvider, "invalid-auth-code", validState)
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Error will come from token exchange
		assert.Contains(t, err.Error(), "failed to exchange code for token")
	})
}

func TestOAuthService_RefreshToken(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := newMockUserRepository()
	jwtManager := newMockJWTManager()
	config := OAuthConfig{
		Google: GoogleOAuthConfig{
			ClientID:     "test-google-client-id",
			ClientSecret: "test-google-client-secret",
			RedirectURI:  "http://localhost:8080/callback",
		},
	}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should return error for empty refresh token", func(t *testing.T) {
		resp, err := service.RefreshToken(context.Background(), dto.GoogleProvider, "")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "refresh token is required", err.Error())
	})

	t.Run("should handle invalid refresh token", func(t *testing.T) {
		resp, err := service.RefreshToken(context.Background(), dto.GoogleProvider, "invalid-refresh-token")
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Error will come from OAuth provider
		assert.Contains(t, err.Error(), "failed to refresh token")
	})

	t.Run("should return error for unsupported provider", func(t *testing.T) {
		resp, err := service.RefreshToken(context.Background(), dto.OAuthProvider("unsupported"), "some-token")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "unsupported OAuth provider")
	})
}

func TestOAuthService_GetCurrentUser(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	mockUserRepo := &MockUserRepository{}
	jwtManager := newMockJWTManager()
	config := OAuthConfig{}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

	t.Run("should return user for valid token", func(t *testing.T) {
		// Create a valid token
		userID := uint(123)
		token, err := jwtManager.GenerateToken(userID, "test@example.com", "user")
		assert.NoError(t, err)

		expectedUser := &entity.User{
			ID:       userID,
			Email:    "test@example.com",
			FullName: "Test User",
			Level:    "user",
			IsActive: true,
		}

		mockUserRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

		user, err := service.GetCurrentUser(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		user, err := service.GetCurrentUser(context.Background(), "invalid-token")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		userID := uint(999)
		token, err := jwtManager.GenerateToken(userID, "notfound@example.com", "user")
		assert.NoError(t, err)

		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

		mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, gorm.ErrRecordNotFound)

		user, err := service.GetCurrentUser(context.Background(), token)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error for database error", func(t *testing.T) {
		userID := uint(123)
		token, err := jwtManager.GenerateToken(userID, "test@example.com", "user")
		assert.NoError(t, err)

		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager)

		dbError := errors.New("database connection failed")
		mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, dbError)

		user, err := service.GetCurrentUser(context.Background(), token)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, dbError, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestOAuthService_Logout(t *testing.T) {
	// Create a real Redis client for testing
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret-key-for-jwt-generation-in-tests",
			Expiry:    time.Hour,
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       1, // Use test DB
		},
	}

	logger := testutil.NewTestLogger(t)
	mockUserRepo := &MockUserRepository{}

	// Use a simple JWT manager without Redis for unit tests
	jwtManagerSimple := jwt.NewJWTManager(cfg)
	config := OAuthConfig{}
	service := NewOAuthService(config, logger, mockUserRepo, jwtManagerSimple)

	t.Run("should return error for invalid token", func(t *testing.T) {
		err := service.Logout(context.Background(), "invalid-token")
		assert.Error(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("should handle valid token without Redis", func(t *testing.T) {
		// Generate a valid token
		userID := uint(123)
		token, err := jwtManagerSimple.GenerateToken(userID, "test@example.com", "user")
		assert.NoError(t, err)

		// Since we're using JWT manager without Redis, it will fail gracefully
		err = service.Logout(context.Background(), token)
		// The error depends on Redis availability
		// In unit test without Redis, it should fail
		if err != nil {
			assert.Contains(t, err.Error(), "failed to revoke token")
		}
	})
}

func TestOAuthService_CreateOrUpdateUser(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	jwtManager := newMockJWTManager()
	config := OAuthConfig{}

	userInfo := &dto.OAuthUserInfo{
		ID:        "oauth-123",
		Email:     "newuser@example.com",
		Name:      "New User",
		AvatarURL: "https://example.com/avatar.jpg",
		Provider:  "google",
	}

	t.Run("should create new user when user does not exist", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager).(*oauthService)

		mockUserRepo.On("GetByEmail", mock.Anything, userInfo.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Email == userInfo.Email && u.FullName == userInfo.Name && string(u.Level) == "user" && u.IsActive
		})).Return(nil)

		user, err := service.createOrUpdateUser(context.Background(), dto.GoogleProvider, userInfo)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userInfo.Email, user.Email)
		assert.Equal(t, userInfo.Name, user.FullName)
		assert.Equal(t, entity.UserLevel("user"), user.Level)
		assert.True(t, user.IsActive)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should update existing user when user exists", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager).(*oauthService)

		existingUser := &entity.User{
			ID:       uint(456),
			Email:    userInfo.Email,
			FullName: "Old Name",
			Level:    "user",
			IsActive: false,
		}

		mockUserRepo.On("GetByEmail", mock.Anything, userInfo.Email).Return(existingUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Email == userInfo.Email && u.FullName == userInfo.Name && u.IsActive
		})).Return(nil)

		user, err := service.createOrUpdateUser(context.Background(), dto.GoogleProvider, userInfo)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(456), user.ID)
		assert.Equal(t, userInfo.Name, user.FullName)
		assert.True(t, user.IsActive)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetByEmail fails with non-NotFound error", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager).(*oauthService)

		dbError := errors.New("database error")
		mockUserRepo.On("GetByEmail", mock.Anything, userInfo.Email).Return(nil, dbError)

		user, err := service.createOrUpdateUser(context.Background(), dto.GoogleProvider, userInfo)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to check existing user")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager).(*oauthService)

		existingUser := &entity.User{
			ID:       uint(456),
			Email:    userInfo.Email,
			FullName: "Old Name",
			Level:    "user",
			IsActive: false,
		}

		updateError := errors.New("update failed")
		mockUserRepo.On("GetByEmail", mock.Anything, userInfo.Email).Return(existingUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(updateError)

		user, err := service.createOrUpdateUser(context.Background(), dto.GoogleProvider, userInfo)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to update user")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		service := NewOAuthService(config, logger, mockUserRepo, jwtManager).(*oauthService)

		createError := errors.New("create failed")
		mockUserRepo.On("GetByEmail", mock.Anything, userInfo.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(createError)

		user, err := service.createOrUpdateUser(context.Background(), dto.GoogleProvider, userInfo)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to create user")
		mockUserRepo.AssertExpectations(t)
	})
}
