package service

import (
	"context"
	"errors"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"

	"go.uber.org/zap"
)

// WalletService defines the interface for wallet business logic
type WalletService interface {
	CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*entity.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uint) (*dto.GetWalletResponse, error)
	GetWalletByID(ctx context.Context, id uint) (*dto.GetWalletResponse, error)
	AddBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error)
	DeductBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error)
	DeleteWallet(ctx context.Context, userID uint) error
}

type walletService struct {
	repo   repository.WalletRepository
	logger *zap.Logger
}

// NewWalletService creates a new wallet service
func NewWalletService(repo repository.WalletRepository, logger *zap.Logger) WalletService {
	return &walletService{
		repo:   repo,
		logger: logger,
	}
}

// CreateWallet creates a new wallet for a user
func (s *walletService) CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*entity.Wallet, error) {
	// Create wallet entity
	wallet := &entity.Wallet{
		UserID:  req.UserID,
		Balance: 0.00,
	}

	// Save to repository
	if err := s.repo.CreateWallet(ctx, wallet); err != nil {
		s.logger.Error("Failed to create wallet", zap.Error(err), zap.Uint("user_id", req.UserID))
		return nil, err
	}

	s.logger.Info("Wallet created successfully", zap.Uint("user_id", req.UserID), zap.Uint("wallet_id", wallet.ID))
	return wallet, nil
}

// GetWalletByUserID retrieves wallet information by user ID
func (s *walletService) GetWalletByUserID(ctx context.Context, userID uint) (*dto.GetWalletResponse, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	return s.entityToResponse(wallet), nil
}

// GetWalletByID retrieves wallet information by wallet ID
func (s *walletService) GetWalletByID(ctx context.Context, id uint) (*dto.GetWalletResponse, error) {
	wallet, err := s.repo.GetWalletByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	return s.entityToResponse(wallet), nil
}

// AddBalance adds funds to the wallet
func (s *walletService) AddBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error) {
	// Validate amount
	if amount <= 0 {
		s.logger.Warn("Invalid amount", zap.Uint("user_id", userID), zap.Float64("amount", amount))
		return nil, errors.New("amount must be positive")
	}

	// Get wallet
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	// Update balance
	wallet.Balance += amount
	if err := s.repo.UpdateWallet(ctx, wallet); err != nil {
		s.logger.Error("Failed to update balance", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Balance added to wallet",
		zap.Uint("user_id", userID),
		zap.Float64("amount", amount),
		zap.Float64("new_balance", wallet.Balance),
		zap.String("description", description))

	return s.entityToResponse(wallet), nil
}

// DeductBalance deducts funds from the wallet
func (s *walletService) DeductBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error) {
	// Validate amount
	if amount <= 0 {
		s.logger.Warn("Invalid amount", zap.Uint("user_id", userID), zap.Float64("amount", amount))
		return nil, errors.New("amount must be positive")
	}

	// Get wallet
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	// Check sufficient balance
	if wallet.Balance < amount {
		s.logger.Warn("Insufficient balance", zap.Uint("user_id", userID), zap.Float64("balance", wallet.Balance), zap.Float64("amount", amount))
		return nil, errors.New("insufficient balance")
	}

	// Update balance
	wallet.Balance -= amount
	if err := s.repo.UpdateWallet(ctx, wallet); err != nil {
		s.logger.Error("Failed to update balance", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	s.logger.Info("Balance deducted from wallet",
		zap.Uint("user_id", userID),
		zap.Float64("amount", amount),
		zap.Float64("new_balance", wallet.Balance),
		zap.String("description", description))

	return s.entityToResponse(wallet), nil
}

// DeleteWallet deletes a wallet
func (s *walletService) DeleteWallet(ctx context.Context, userID uint) error {
	if err := s.repo.DeleteWallet(ctx, userID); err != nil {
		s.logger.Error("Failed to delete wallet", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	s.logger.Info("Wallet deleted", zap.Uint("user_id", userID))
	return nil
}

// entityToResponse converts a wallet entity to a response DTO
func (s *walletService) entityToResponse(wallet *entity.Wallet) *dto.GetWalletResponse {
	return &dto.GetWalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
