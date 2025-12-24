package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	transactionRepo "github.com/novriyantoAli/cn-wallet/internal/application/transaction/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"

	"go.uber.org/zap"
)

// WalletService defines the interface for wallet business logic
type WalletService interface {
	CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*walletEntity.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uint) (*dto.GetWalletResponse, error)
	GetWalletByID(ctx context.Context, id uint) (*dto.GetWalletResponse, error)
	AddBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error)
	DeductBalance(ctx context.Context, userID uint, amount float64, description string) (*dto.GetWalletResponse, error)
	Transfer(ctx context.Context, token string, req *dto.TransferRequest) (*dto.TransferResponse, error)
	DeleteWallet(ctx context.Context, userID uint) error
}

type walletService struct {
	repo            repository.WalletRepository
	transactionRepo transactionRepo.TransactionRepository
	jwtManager      *jwt.JWTManager
	txManager       database.TransactionManagerI
	logger          *zap.Logger
}

// NewWalletService creates a new wallet service
func NewWalletService(
	repo repository.WalletRepository,
	transactionRepo transactionRepo.TransactionRepository,
	jwtManager *jwt.JWTManager,
	txManager database.TransactionManagerI,
	logger *zap.Logger,
) WalletService {
	return &walletService{
		repo:            repo,
		transactionRepo: transactionRepo,
		jwtManager:      jwtManager,
		txManager:       txManager,
		logger:          logger,
	}
}

// CreateWallet creates a new wallet for a user
func (s *walletService) CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*walletEntity.Wallet, error) {
	// Create wallet entity
	wallet := &walletEntity.Wallet{
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

// Transfer transfers funds from one wallet to another
func (s *walletService) Transfer(ctx context.Context, token string, req *dto.TransferRequest) (*dto.TransferResponse, error) {
	// Verify token and extract user ID
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		s.logger.Error("Invalid token", zap.Error(err))
		return nil, errors.New("invalid or expired token")
	}

	fromUserID := claims.UserID

	// Validate that user is not transferring to themselves
	if fromUserID == req.ToUserID {
		s.logger.Warn("Cannot transfer to self", zap.Uint("user_id", fromUserID))
		return nil, errors.New("cannot transfer to yourself")
	}

	// Validate amount
	if req.Amount <= 0 {
		s.logger.Warn("Invalid transfer amount", zap.Uint("from_user_id", fromUserID), zap.Float64("amount", req.Amount))
		return nil, errors.New("amount must be positive")
	}

	// Execute transfer within transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get source wallet
		fromWallet, err := s.repo.GetForUpdate(txCtx, fromUserID)
		if err != nil {
			s.logger.Error("Failed to get source wallet", zap.Error(err), zap.Uint("user_id", fromUserID))
			return errors.New("source wallet not found")
		}

		// Check sufficient balance
		if fromWallet.Balance < req.Amount {
			s.logger.Warn("Insufficient balance for transfer",
				zap.Uint("from_user_id", fromUserID),
				zap.Float64("balance", fromWallet.Balance),
				zap.Float64("amount", req.Amount))
			return errors.New("insufficient balance")
		}

		// Deduct from source wallet
		fromWallet.Balance -= req.Amount
		if err := s.repo.UpdateWallet(txCtx, fromWallet); err != nil {
			s.logger.Error("Failed to deduct from source wallet", zap.Error(err), zap.Uint("from_user_id", fromUserID))
			return err
		}

		// Create transfer_out transaction for sender
		txOutID := uuid.New()
		transferOutTxn := &transactionEntity.Transaction{
			ID:          txOutID,
			WalletID:    fromWallet.ID,
			Type:        transactionEntity.TypeTransferOut,
			Amount:      req.Amount,
			Status:      transactionEntity.StatusPending,
			Description: req.Description,
			// RelatedWalletID: &toWallet.ID, // will set after getting toWallet
			CreatedAt: time.Now(),
		}

		if err := s.transactionRepo.Create(txCtx, transferOutTxn); err != nil {
			s.logger.Error("Failed to create transfer_out transaction", zap.Error(err))
			return errors.New("failed to create transfer transaction")
		}

		// Get destination wallet
		toWallet, err := s.repo.GetForUpdate(txCtx, req.ToUserID)
		if err != nil {
			s.logger.Error("Failed to get destination wallet", zap.Error(err), zap.Uint("user_id", req.ToUserID))
			return errors.New("destination wallet not found")
		}

		// Add to destination wallet
		toWallet.Balance += req.Amount
		if err := s.repo.UpdateWallet(txCtx, toWallet); err != nil {
			s.logger.Error("Failed to add to destination wallet", zap.Error(err), zap.Uint("to_user_id", req.ToUserID))
			return err
		}

		// update
		transferOutTxn.RelatedWalletID = &toWallet.ID
		transferOutTxn.Status = transactionEntity.StatusSuccess
		if err := s.transactionRepo.Update(txCtx, transferOutTxn); err != nil {
			s.logger.Error("Failed to update transfer_out transaction", zap.Error(err))
			return errors.New("failed to update transfer transaction")
		}

		// Create transfer_in transaction for receiver
		txInID := uuid.New()
		transferInTxn := &transactionEntity.Transaction{
			ID:              txInID,
			WalletID:        toWallet.ID,
			Type:            transactionEntity.TypeTransferIn,
			Amount:          req.Amount,
			Status:          transactionEntity.StatusSuccess,
			Description:     req.Description,
			RelatedWalletID: &fromWallet.ID,
			CreatedAt:       time.Now(),
		}

		if err := s.transactionRepo.Create(txCtx, transferInTxn); err != nil {
			s.logger.Error("Failed to create transfer_in transaction", zap.Error(err))
			return errors.New("failed to create transfer transaction")
		}

		s.logger.Info("Transfer transactions created",
			zap.String("transfer_out_id", txOutID.String()),
			zap.String("transfer_in_id", txInID.String()))

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info("Transfer completed successfully",
		zap.Uint("from_user_id", fromUserID),
		zap.Uint("to_user_id", req.ToUserID),
		zap.Float64("amount", req.Amount),
		zap.String("description", req.Description))

	return &dto.TransferResponse{
		FromUserID:  fromUserID,
		ToUserID:    req.ToUserID,
		Amount:      req.Amount,
		Description: req.Description,
		Message:     "Transfer completed successfully",
	}, nil
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
func (s *walletService) entityToResponse(wallet *walletEntity.Wallet) *dto.GetWalletResponse {
	return &dto.GetWalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
