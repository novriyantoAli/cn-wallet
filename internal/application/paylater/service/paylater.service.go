package service

import (
	"context"
	"errors"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
)

// PaylaterAccountService defines the interface for paylater account business logic
type PaylaterAccountService interface {
	CreateAccount(ctx context.Context, req *dto.CreatePaylaterAccountRequest) (*entity.PaylaterAccount, error)
	GetAccountByUserID(ctx context.Context, userID uint) (*dto.GetPaylaterAccountResponse, error)
	GetAccountByID(ctx context.Context, id uint) (*dto.GetPaylaterAccountResponse, error)
	UpdateCreditLimit(ctx context.Context, userID uint, req *dto.UpdateCreditLimitRequest) (*dto.GetPaylaterAccountResponse, error)
	UpdateStatus(ctx context.Context, userID uint, req *dto.UpdateStatusRequest) (*dto.GetPaylaterAccountResponse, error)
	UseCredit(ctx context.Context, userID uint, req *dto.UseCreditRequest) (*dto.UseCreditResponse, error)
	Repayment(ctx context.Context, userID uint, req *dto.RepaymentRequest) (*dto.RepaymentResponse, error)
	DeleteAccount(ctx context.Context, userID uint) error
}

type paylaterAccountService struct {
	repo      repository.PaylaterAccountRepository
	txManager database.TransactionManagerI
	logger    *zap.Logger
}

// NewPaylaterAccountService creates a new paylater account service
func NewPaylaterAccountService(
	repo repository.PaylaterAccountRepository,
	txManager database.TransactionManagerI,
	logger *zap.Logger,
) PaylaterAccountService {
	return &paylaterAccountService{
		repo:      repo,
		txManager: txManager,
		logger:    logger,
	}
}

// CreateAccount creates a new paylater account for a user
func (s *paylaterAccountService) CreateAccount(ctx context.Context, req *dto.CreatePaylaterAccountRequest) (*entity.PaylaterAccount, error) {
	// Check if account already exists
	existing, err := s.repo.GetAccountByUserID(ctx, req.UserID)
	if err == nil && existing != nil {
		s.logger.Warn("Paylater account already exists", zap.Uint("user_id", req.UserID))
		return nil, errors.New("paylater account already exists for this user")
	}

	// Create account entity
	account := &entity.PaylaterAccount{
		UserID:         req.UserID,
		CreditLimit:    req.CreditLimit,
		Outstanding:    0,
		AvailableLimit: req.CreditLimit,
		Status:         entity.PaylaterStatusActive,
	}

	// Save to repository
	if err := s.repo.CreateAccount(ctx, account); err != nil {
		s.logger.Error("Failed to create paylater account", zap.Error(err), zap.Uint("user_id", req.UserID))
		return nil, err
	}

	s.logger.Info("Paylater account created successfully", zap.Uint("user_id", req.UserID), zap.Uint("account_id", account.ID))
	return account, nil
}

// GetAccountByUserID retrieves paylater account information by user ID
func (s *paylaterAccountService) GetAccountByUserID(ctx context.Context, userID uint) (*dto.GetPaylaterAccountResponse, error) {
	account, err := s.repo.GetAccountByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get paylater account", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	return &dto.GetPaylaterAccountResponse{
		ID:             account.ID,
		UserID:         account.UserID,
		CreditLimit:    account.CreditLimit,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Status:         string(account.Status),
		CreatedAt:      account.CreatedAt,
	}, nil
}

// GetAccountByID retrieves paylater account information by account ID
func (s *paylaterAccountService) GetAccountByID(ctx context.Context, id uint) (*dto.GetPaylaterAccountResponse, error) {
	account, err := s.repo.GetAccountByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get paylater account", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	return &dto.GetPaylaterAccountResponse{
		ID:             account.ID,
		UserID:         account.UserID,
		CreditLimit:    account.CreditLimit,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Status:         string(account.Status),
		CreatedAt:      account.CreatedAt,
	}, nil
}

// UpdateCreditLimit updates the credit limit of a paylater account
func (s *paylaterAccountService) UpdateCreditLimit(ctx context.Context, userID uint, req *dto.UpdateCreditLimitRequest) (*dto.GetPaylaterAccountResponse, error) {
	var account *entity.PaylaterAccount
	var err error

	// Execute in transaction to ensure consistency
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get account with lock
		account, err = s.repo.GetForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		// Update credit limit
		account.CreditLimit = req.CreditLimit
		account.AvailableLimit = req.CreditLimit - account.Outstanding

		// Validate available limit is not negative
		if account.AvailableLimit < 0 {
			return errors.New("credit limit cannot be less than outstanding balance")
		}

		// Save changes
		if err := s.repo.UpdateAccount(txCtx, account); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to update credit limit", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Credit limit updated successfully", zap.Uint("user_id", userID), zap.Int64("new_limit", req.CreditLimit))

	return &dto.GetPaylaterAccountResponse{
		ID:             account.ID,
		UserID:         account.UserID,
		CreditLimit:    account.CreditLimit,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Status:         string(account.Status),
		CreatedAt:      account.CreatedAt,
	}, nil
}

// UpdateStatus updates the status of a paylater account
func (s *paylaterAccountService) UpdateStatus(ctx context.Context, userID uint, req *dto.UpdateStatusRequest) (*dto.GetPaylaterAccountResponse, error) {
	account, err := s.repo.GetAccountByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get paylater account", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	// Update status
	account.Status = entity.PaylaterStatus(req.Status)
	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		s.logger.Error("Failed to update status", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Paylater account status updated", zap.Uint("user_id", userID), zap.String("status", req.Status))

	return &dto.GetPaylaterAccountResponse{
		ID:             account.ID,
		UserID:         account.UserID,
		CreditLimit:    account.CreditLimit,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Status:         string(account.Status),
		CreatedAt:      account.CreatedAt,
	}, nil
}

// UseCredit deducts available credit and increases outstanding balance
func (s *paylaterAccountService) UseCredit(ctx context.Context, userID uint, req *dto.UseCreditRequest) (*dto.UseCreditResponse, error) {
	var account *entity.PaylaterAccount
	var err error

	// Execute in transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get account with lock
		account, err = s.repo.GetForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		// Validate account status
		if account.Status != entity.PaylaterStatusActive {
			return errors.New("account is not active")
		}

		// Check if sufficient credit is available
		if account.AvailableLimit < req.Amount {
			return errors.New("insufficient credit limit")
		}

		// Update outstanding and available limit
		account.Outstanding += req.Amount
		account.AvailableLimit = account.CreditLimit - account.Outstanding

		// Save changes
		if err := s.repo.UpdateAccount(txCtx, account); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to use credit", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Credit used successfully",
		zap.Uint("user_id", userID),
		zap.Int64("amount", req.Amount),
		zap.Int64("outstanding", account.Outstanding))

	return &dto.UseCreditResponse{
		UserID:         account.UserID,
		Amount:         req.Amount,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Message:        "Credit used successfully",
	}, nil
}

// Repayment processes repayment and reduces outstanding balance
func (s *paylaterAccountService) Repayment(ctx context.Context, userID uint, req *dto.RepaymentRequest) (*dto.RepaymentResponse, error) {
	var account *entity.PaylaterAccount
	var err error

	// Execute in transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get account with lock
		account, err = s.repo.GetForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		// Check if repayment amount exceeds outstanding
		if req.Amount > account.Outstanding {
			return errors.New("repayment amount cannot exceed outstanding balance")
		}

		// Update outstanding and available limit
		account.Outstanding -= req.Amount
		account.AvailableLimit = account.CreditLimit - account.Outstanding

		// Save changes
		if err := s.repo.UpdateAccount(txCtx, account); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to process repayment", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Repayment processed successfully",
		zap.Uint("user_id", userID),
		zap.Int64("amount", req.Amount),
		zap.Int64("outstanding", account.Outstanding))

	return &dto.RepaymentResponse{
		UserID:         account.UserID,
		Amount:         req.Amount,
		Outstanding:    account.Outstanding,
		AvailableLimit: account.AvailableLimit,
		Message:        "Repayment processed successfully",
	}, nil
}

// DeleteAccount deletes a paylater account
func (s *paylaterAccountService) DeleteAccount(ctx context.Context, userID uint) error {
	// Check if account has outstanding balance
	account, err := s.repo.GetAccountByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if account.Outstanding > 0 {
		return errors.New("cannot delete account with outstanding balance")
	}

	if err := s.repo.DeleteAccount(ctx, userID); err != nil {
		s.logger.Error("Failed to delete paylater account", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	s.logger.Info("Paylater account deleted successfully", zap.Uint("user_id", userID))
	return nil
}
