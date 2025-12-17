package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"

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
	SetPIN(ctx context.Context, userID uint, req *dto.SetPINRequest) error
	VerifyPIN(ctx context.Context, userID uint, pin string) error
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
	// Validate PIN format
	if err := validatePIN(req.PIN); err != nil {
		s.logger.Error("Invalid PIN format", zap.Error(err), zap.Uint("user_id", req.UserID))
		return nil, err
	}

	// Hash the PIN
	pinHash := hashPIN(req.PIN)

	// Create wallet entity
	wallet := &entity.Wallet{
		UserID:  req.UserID,
		Balance: 0.00,
		PINHash: pinHash,
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

// SetPIN sets or updates the wallet PIN
func (s *walletService) SetPIN(ctx context.Context, userID uint, req *dto.SetPINRequest) error {
	// Get wallet to verify current PIN
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	// Verify current PIN
	if !verifyPIN(req.CurrentPIN, wallet.PINHash) {
		s.logger.Warn("Invalid current PIN", zap.Uint("user_id", userID))
		return errors.New("invalid current PIN")
	}

	// Validate new PIN format
	if err := validatePIN(req.NewPIN); err != nil {
		s.logger.Error("Invalid new PIN format", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	// Hash and update PIN
	newPINHash := hashPIN(req.NewPIN)
	if err := s.repo.UpdatePIN(ctx, userID, newPINHash); err != nil {
		s.logger.Error("Failed to update PIN", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	s.logger.Info("PIN updated successfully", zap.Uint("user_id", userID))
	return nil
}

// VerifyPIN verifies the wallet PIN
func (s *walletService) VerifyPIN(ctx context.Context, userID uint, pin string) error {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", userID))
		return err
	}

	if !verifyPIN(pin, wallet.PINHash) {
		s.logger.Warn("Invalid PIN", zap.Uint("user_id", userID))
		return errors.New("invalid PIN")
	}

	return nil
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

// Helper functions

// validatePIN validates that PIN is exactly 6 numeric digits
func validatePIN(pin string) error {
	if len(pin) != 6 {
		return errors.New("PIN must be exactly 6 digits")
	}

	matched, _ := regexp.MatchString(`^\d{6}$`, pin)
	if !matched {
		return errors.New("PIN must contain only numeric digits")
	}

	return nil
}

// hashPIN creates a SHA256 hash of the PIN
func hashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}

// verifyPIN verifies a PIN against its hash
func verifyPIN(pin string, hash string) bool {
	return hashPIN(pin) == hash
}
