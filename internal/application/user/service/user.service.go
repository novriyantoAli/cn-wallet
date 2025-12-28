package service

import (
	"context"
	"errors"
	"time"

	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	providerRepo "github.com/novriyantoAli/cn-wallet/internal/application/provider/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/repository"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.UserListResponse, error)
	UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateUserProvider(ctx context.Context, id uint, req *dto.UpdateUserProviderRequest) (*dto.UserResponse, error)
	UpdateLevel(ctx context.Context, id uint, req *dto.UpdateUserLevelRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	repo         repository.UserRepository
	providerRepo providerRepo.ProviderRepository
	logger       *zap.Logger
}

func NewUserService(repo repository.UserRepository, providerRepo providerRepo.ProviderRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:         repo,
		providerRepo: providerRepo,
		logger:       logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create wallet entity for the new user
	wallet := &walletEntity.Wallet{
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user := &entity.User{
		Email:     req.Email,
		FullName:  req.FullName,
		Level:     "user",
		IsActive:  true,
		Wallet:    wallet, // Attach wallet to user entity
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.UserListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	users, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, *s.entityToResponse(&user))
	}

	return &dto.UserListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (s *userService) UpdateUserProvider(ctx context.Context, id uint, req *dto.UpdateUserProviderRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if user.ProviderID != nil {
		return s.entityToResponse(user), nil
	}

	// get provider from providerRepo
	provider, err := s.providerRepo.GetByID(ctx, req.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, errors.New("internal server error")
	}

	user.ProviderID = &provider.ID

	err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user provider info", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) UpdateLevel(ctx context.Context, id uint, req *dto.UpdateUserLevelRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Level = entity.UserLevel(req.Level)
	user.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user level", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Level != "" {
		user.Level = entity.UserLevel(req.Level)
	}
	user.IsActive = req.IsActive
	user.UpdatedAt = time.Now()

	// Handle provider information update if any provider fields are provided
	if req.ProviderName != "" && req.ProviderCode != "" {
		user.Provider = &providerEntity.Provider{
			Name:      req.ProviderName,
			Code:      req.ProviderCode,
			Logo:      req.ProviderLogo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *userService) entityToResponse(user *entity.User) *dto.UserResponse {
	var walletInfo *dto.WalletInfo
	if user.Wallet != nil {
		walletInfo = &dto.WalletInfo{
			ID:        user.Wallet.ID,
			Balance:   user.Wallet.Balance,
			CreatedAt: user.Wallet.CreatedAt,
			UpdatedAt: user.Wallet.UpdatedAt,
		}
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Level:     user.Level.String(),
		IsActive:  user.IsActive,
		Wallet:    walletInfo,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
