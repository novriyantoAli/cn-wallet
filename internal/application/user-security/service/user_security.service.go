package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/repository"

	"go.uber.org/zap"
)

type UserSecurityService interface {
	SetPIN(ctx context.Context, req *dto.SetPINRequest) error
	VerifyPIN(ctx context.Context, req *dto.VerifyPINRequest) (bool, error)
	GetSecurity(ctx context.Context, userID uint) (*dto.UserSecurityResponse, error)
	IsAccountLocked(ctx context.Context, userID uint) (bool, error)
}

type userSecurityService struct {
	repo   repository.UserSecurityRepository
	logger *zap.Logger
}

func NewUserSecurityService(repo repository.UserSecurityRepository, logger *zap.Logger) UserSecurityService {
	return &userSecurityService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userSecurityService) SetPIN(ctx context.Context, req *dto.SetPINRequest) error {
	s.logger.Info("Setting PIN for user", zap.Uint("user_id", req.UserID))

	pinHash := s.hashPIN(req.PIN)

	security, err := s.repo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if security == nil {
		// Create new security record
		return s.repo.Create(ctx, &entity.UserSecurity{
			UserID:        req.UserID,
			PinHash:       pinHash,
			FailedAttempt: 0,
			LockedUntil:   nil,
		})
	}

	// Update existing PIN
	return s.repo.UpdatePIN(ctx, req.UserID, pinHash)
}

func (s *userSecurityService) VerifyPIN(ctx context.Context, req *dto.VerifyPINRequest) (bool, error) {
	s.logger.Info("Verifying PIN for user", zap.Uint("user_id", req.UserID))

	security, err := s.repo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return false, err
	}

	if security == nil {
		s.logger.Warn("Security record not found", zap.Uint("user_id", req.UserID))
		return false, nil
	}

	// Check if account is locked
	if security.LockedUntil != nil && security.LockedUntil.After(time.Now()) {
		s.logger.Warn("Account is locked", zap.Uint("user_id", req.UserID))
		return false, nil
	}

	pinHash := s.hashPIN(req.PIN)
	if security.PinHash == pinHash {
		// Reset failed attempts on successful verification
		if err := s.repo.ResetFailedAttempt(ctx, req.UserID); err != nil {
			s.logger.Error("Failed to reset failed attempts", zap.Error(err))
		}
		return true, nil
	}

	// Increment failed attempts
	if err := s.repo.IncrementFailedAttempt(ctx, req.UserID); err != nil {
		s.logger.Error("Failed to increment failed attempt", zap.Error(err))
	}

	// Lock account after 3 failed attempts for 15 minutes
	if security.FailedAttempt+1 >= 3 {
		if err := s.repo.LockAccount(ctx, req.UserID, 15*time.Minute); err != nil {
			s.logger.Error("Failed to lock account", zap.Error(err))
		}
	}

	return false, nil
}

func (s *userSecurityService) GetSecurity(ctx context.Context, userID uint) (*dto.UserSecurityResponse, error) {
	s.logger.Info("Getting security info for user", zap.Uint("user_id", userID))

	security, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if security == nil {
		return &dto.UserSecurityResponse{
			UserID:        userID,
			FailedAttempt: 0,
			IsLocked:      false,
		}, nil
	}

	isLocked := security.LockedUntil != nil && security.LockedUntil.After(time.Now())

	return &dto.UserSecurityResponse{
		UserID:        security.UserID,
		FailedAttempt: security.FailedAttempt,
		IsLocked:      isLocked,
	}, nil
}

func (s *userSecurityService) IsAccountLocked(ctx context.Context, userID uint) (bool, error) {
	security, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	if security == nil {
		return false, nil
	}

	return security.LockedUntil != nil && security.LockedUntil.After(time.Now()), nil
}

func (s *userSecurityService) hashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}
