package service

import (
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type WalletService interface {
	CreateWallet(req *dto.CreateWalletRequest) (*dto.WalletResponse, error)
	GetWalletByID(id uint) (*dto.WalletResponse, error)
	GetWalletByUserID(userID uint) (*dto.WalletResponse, error)
	GetWallets(filter *dto.WalletFilter) (*dto.WalletListResponse, error)
	UpdateWalletPIN(id uint, req *dto.UpdateWalletPINRequest) error
	AddBalance(id uint, amount decimal.Decimal) error
	WithdrawBalance(id uint, amount decimal.Decimal, pin string) error
	TransferBalance(fromID, toUserID uint, amount decimal.Decimal, pin string) error
	Delete(id uint) error
}

type walletService struct {
	repo   repository.WalletRepository
	logger *zap.Logger
}

func NewWalletService(repo repository.WalletRepository, logger *zap.Logger) WalletService {
	return &walletService{
		repo:   repo,
		logger: logger,
	}
}

func (s *walletService) CreateWallet(req *dto.CreateWalletRequest) (*dto.WalletResponse, error) {
	// Hash the PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash PIN", zap.Error(err))
		return nil, err
	}

	wallet := &entity.Wallet{
		UserID:    req.UserID,
		Balance:   decimal.RequireFromString("0.00"),
		PinHash:   string(hashedPIN),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(wallet)
	if err != nil {
		s.logger.Error("Failed to create wallet", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(wallet), nil
}

func (s *walletService) GetWalletByID(id uint) (*dto.WalletResponse, error) {
	wallet, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return s.entityToResponse(wallet), nil
}

func (s *walletService) GetWalletByUserID(userID uint) (*dto.WalletResponse, error) {
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return s.entityToResponse(wallet), nil
}

func (s *walletService) GetWallets(filter *dto.WalletFilter) (*dto.WalletListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	wallets, totalCount, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.WalletResponse, 0, len(wallets))
	for _, wallet := range wallets {
		responses = append(responses, *s.entityToResponse(&wallet))
	}

	return &dto.WalletListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (s *walletService) UpdateWalletPIN(id uint, req *dto.UpdateWalletPINRequest) error {
	wallet, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		return err
	}

	// Verify current PIN
	err = bcrypt.CompareHashAndPassword([]byte(wallet.PinHash), []byte(req.CurrentPIN))
	if err != nil {
		return errors.New("current PIN is incorrect")
	}

	// Hash new PIN
	hashedNewPIN, err := bcrypt.GenerateFromPassword([]byte(req.NewPIN), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash new PIN", zap.Error(err))
		return err
	}

	wallet.PinHash = string(hashedNewPIN)
	wallet.UpdatedAt = time.Now()

	return s.repo.Update(wallet)
}

func (s *walletService) AddBalance(id uint, amount decimal.Decimal) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		return err
	}

	if amount.IsNegative() {
		return errors.New("amount must be positive")
	}

	s.logger.Info("Adding balance to wallet", zap.Uint("wallet_id", id), zap.String("amount", amount.String()))
	return s.repo.AddBalance(id, amount)
}

func (s *walletService) WithdrawBalance(id uint, amount decimal.Decimal, pin string) error {
	wallet, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		return err
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(wallet.PinHash), []byte(pin))
	if err != nil {
		return errors.New("PIN is incorrect")
	}

	if amount.IsNegative() {
		return errors.New("amount must be positive")
	}

	// Check balance
	if wallet.Balance.LessThan(amount) {
		return errors.New("insufficient balance")
	}

	s.logger.Info("Withdrawing balance from wallet", zap.Uint("wallet_id", id), zap.String("amount", amount.String()))
	return s.repo.SubtractBalance(id, amount)
}

func (s *walletService) TransferBalance(fromID, toUserID uint, amount decimal.Decimal, pin string) error {
	fromWallet, err := s.repo.GetByID(fromID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("source wallet not found")
		}
		return err
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(fromWallet.PinHash), []byte(pin))
	if err != nil {
		return errors.New("PIN is incorrect")
	}

	// Get recipient wallet
	toWallet, err := s.repo.GetByUserID(toUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recipient wallet not found")
		}
		return err
	}

	if amount.IsNegative() {
		return errors.New("amount must be positive")
	}

	// Check balance
	if fromWallet.Balance.LessThan(amount) {
		return errors.New("insufficient balance")
	}

	s.logger.Info("Transferring balance", 
		zap.Uint("from_wallet_id", fromID), 
		zap.Uint("to_wallet_id", toWallet.ID), 
		zap.String("amount", amount.String()))

	// Subtract from source
	if err := s.repo.SubtractBalance(fromID, amount); err != nil {
		return err
	}

	// Add to recipient
	if err := s.repo.AddBalance(toWallet.ID, amount); err != nil {
		return err
	}

	return nil
}

func (s *walletService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		return err
	}

	s.logger.Info("Deleting wallet", zap.Uint("id", id))
	return s.repo.Delete(id)
}

func (s *walletService) entityToResponse(wallet *entity.Wallet) *dto.WalletResponse {
	return &dto.WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance.String(),
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
