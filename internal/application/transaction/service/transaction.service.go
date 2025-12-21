package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransactionService defines the interface for transaction business logic.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
	GetAllTransactions(ctx context.Context, filter *dto.TransactionFilter) (*dto.TransactionListResponse, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetWalletTransactions(ctx context.Context, walletID uint, page, pageSize int) (*dto.TransactionListResponse, error)
}

type transactionService struct {
	repo   repository.TransactionRepository
	logger *zap.Logger
}

// NewTransactionService creates a new instance of TransactionService.
func NewTransactionService(repo repository.TransactionRepository, logger *zap.Logger) TransactionService {
	return &transactionService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTransaction creates a new transaction.
func (s *transactionService) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	transaction := &entity.Transaction{
		ID:                 uuid.New(),
		WalletID:           req.WalletID,
		Type:               req.Type,
		Amount:             req.Amount,
		Status:             entity.StatusPending,
		Description:        req.Description,
		PaymentMethod:      req.PaymentMethod,
		PaymentProviderRef: req.PaymentProviderRef,
		ProductID:          req.ProductID,
		TargetNumber:       req.TargetNumber,
		SerialNumber:       req.SerialNumber,
		RelatedWalletID:    req.RelatedWalletID,
	}

	if err := s.repo.Create(ctx, transaction); err != nil {
		s.logger.Error("Failed to create transaction", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(transaction), nil
}

// GetTransactionByID retrieves a transaction by ID.
func (s *transactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Transaction not found", zap.String("id", id.String()))
			return nil, err
		}
		s.logger.Error("Failed to get transaction", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(transaction), nil
}

// GetAllTransactions retrieves all transactions with filters and pagination.
func (s *transactionService) GetAllTransactions(ctx context.Context, filter *dto.TransactionFilter) (*dto.TransactionListResponse, error) {
	transactions, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get transactions", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *s.entityToResponse(&t))
	}

	return &dto.TransactionListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

// UpdateTransaction updates an existing transaction.
func (s *transactionService) UpdateTransaction(ctx context.Context, id uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get transaction", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	transaction.Status = req.Status
	if req.Description != "" {
		transaction.Description = req.Description
	}

	if err := s.repo.Update(ctx, transaction); err != nil {
		s.logger.Error("Failed to update transaction", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(transaction), nil
}

// DeleteTransaction deletes a transaction.
func (s *transactionService) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get transaction", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	if err := s.repo.Delete(ctx, transaction.ID); err != nil {
		s.logger.Error("Failed to delete transaction", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.logger.Info("Transaction deleted successfully", zap.String("id", id.String()))
	return nil
}

// GetWalletTransactions retrieves all transactions for a specific wallet.
func (s *transactionService) GetWalletTransactions(ctx context.Context, walletID uint, page, pageSize int) (*dto.TransactionListResponse, error) {
	transactions, totalCount, err := s.repo.GetByWalletID(ctx, walletID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get wallet transactions", zap.Uint("wallet_id", walletID), zap.Error(err))
		return nil, err
	}

	responses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *s.entityToResponse(&t))
	}

	return &dto.TransactionListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// entityToResponse converts a transaction entity to a response DTO.
func (s *transactionService) entityToResponse(transaction *entity.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:                 transaction.ID,
		WalletID:           transaction.WalletID,
		Type:               transaction.Type,
		Amount:             transaction.Amount,
		Status:             transaction.Status,
		Description:        transaction.Description,
		PaymentMethod:      transaction.PaymentMethod,
		PaymentProviderRef: transaction.PaymentProviderRef,
		ProductID:          transaction.ProductID,
		TargetNumber:       transaction.TargetNumber,
		SerialNumber:       transaction.SerialNumber,
		RelatedWalletID:    transaction.RelatedWalletID,
		CreatedAt:          transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
