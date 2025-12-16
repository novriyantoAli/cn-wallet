package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/user/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/repository"
	"github.com/shopspring/decimal"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetUsers(filter *dto.UserFilter) (*dto.UserListResponse, error)
	UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateUserPassword(id uint, req *dto.UpdateUserPasswordRequest) error
	UpdateUserPIN(id uint, req *dto.UpdateUserPINRequest) error
	DeleteUser(id uint) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		s.logger.Error("Failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	var passwordHash sql.NullString
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("Failed to hash password", zap.Error(err))
			return nil, err
		}
		passwordHash = sql.NullString{String: string(hashedPassword), Valid: true}
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash PIN", zap.Error(err))
		return nil, err
	}

	user := &entity.User{
		Email:        req.Email,
		PhoneNumber:  sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		FullName:     sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		PasswordHash: passwordHash,
		PinHash:      string(hashedPIN),
		Balance:      decimal.RequireFromString("0.00"),
		Level:        "user",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repo.Create(user)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) GetUsers(filter *dto.UserFilter) (*dto.UserListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	users, totalCount, err := s.repo.GetAll(filter)
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

func (s *userService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.PhoneNumber != "" {
		user.PhoneNumber = sql.NullString{String: req.PhoneNumber, Valid: true}
	}
	if req.FullName != "" {
		user.FullName = sql.NullString{String: req.FullName, Valid: true}
	}
	if req.Level != "" {
		user.Level = req.Level
	}
	user.IsActive = req.IsActive
	user.UpdatedAt = time.Now()

	err = s.repo.Update(user)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(user), nil
}

func (s *userService) UpdateUserPassword(id uint, req *dto.UpdateUserPasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if !user.PasswordHash.Valid {
		return errors.New("user does not have a password set")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash new password", zap.Error(err))
		return err
	}

	user.PasswordHash = sql.NullString{String: string(hashedPassword), Valid: true}
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

func (s *userService) UpdateUserPIN(id uint, req *dto.UpdateUserPINRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PinHash), []byte(req.PIN))
	if err != nil {
		return errors.New("current PIN is incorrect")
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.NewPIN), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash new PIN", zap.Error(err))
		return err
	}

	user.PinHash = string(hashedPIN)
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *userService) entityToResponse(user *entity.User) *dto.UserResponse {
	phoneNumber := ""
	if user.PhoneNumber.Valid {
		phoneNumber = user.PhoneNumber.String
	}

	fullName := ""
	if user.FullName.Valid {
		fullName = user.FullName.String
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: phoneNumber,
		FullName:    fullName,
		Balance:     user.Balance.String(),
		Level:       user.Level,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
