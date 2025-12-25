package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	paylaterRepo "github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/repository"
	walletRepo "github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"go.uber.org/zap"
)

type TransferService interface {
	CreateTransfer(ctx context.Context, req *dto.CreateTransferRequest) (*entity.Transfer, error)
	GetTransferByID(ctx context.Context, id uint) (*dto.GetTransferResponse, error)
	GetTransfersByUserID(ctx context.Context, userID uint) ([]dto.GetTransferResponse, error)
	GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]dto.GetTransferResponse, error)
	ListTransfers(ctx context.Context, req *dto.ListTransfersRequest) (*dto.ListTransfersResponse, error)
	UpdateTransferStatus(ctx context.Context, id uint, req *dto.UpdateTransferStatusRequest) (*dto.GetTransferResponse, error)
	GetUserTransferStats(ctx context.Context, userID uint) (*dto.TransferStatsResponse, error)
	CancelTransfer(ctx context.Context, id uint) error
}

type transferService struct {
	transferRepo        repository.TransferRepository
	walletRepo          walletRepo.WalletRepository
	paylaterAccountRepo paylaterRepo.PaylaterAccountRepository
	txManager           database.TransactionManagerI
	logger              *zap.Logger
}

func NewTransferService(
	transferRepo repository.TransferRepository,
	walletRepo walletRepo.WalletRepository,
	paylaterAccountRepo paylaterRepo.PaylaterAccountRepository,
	txManager database.TransactionManagerI,
	logger *zap.Logger,
) TransferService {
	return &transferService{
		transferRepo:        transferRepo,
		walletRepo:          walletRepo,
		paylaterAccountRepo: paylaterAccountRepo,
		txManager:           txManager,
		logger:              logger,
	}
}

// CreateTransfer creates a new transfer with balance validation and transaction management
func (s *transferService) CreateTransfer(ctx context.Context, req *dto.CreateTransferRequest) (*entity.Transfer, error) {
	// Validate sender and receiver are different
	if req.UserID == req.TargetUserID {
		return nil, errors.New("cannot transfer to yourself")
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("transfer amount must be greater than 0")
	}

	transfer := &entity.Transfer{
		UserID:       req.UserID,
		TargetUserID: req.TargetUserID,
		Amount:       req.Amount,
		Source:       entity.TransferSource(req.Source),
		Status:       entity.TransferStatusPending,
	}

	// Execute transfer in transaction
	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Verify and deduct from source
		if entity.TransferSource(req.Source) == entity.TransferSourceWallet {
			// Check wallet balance
			senderWallet, err := s.walletRepo.GetWalletByUserID(txCtx, req.UserID)
			if err != nil {
				return fmt.Errorf("failed to get sender wallet: %w", err)
			}
			if senderWallet.Balance < float64(req.Amount) {
				return errors.New("insufficient wallet balance")
			}

			// Calculate new balance for sender
			newSenderBalance := senderWallet.Balance - float64(req.Amount)
			if err := s.walletRepo.UpdateBalance(txCtx, req.UserID, newSenderBalance); err != nil {
				return fmt.Errorf("failed to deduct from sender wallet: %w", err)
			}

			// Add to receiver wallet
			receiverWallet, err := s.walletRepo.GetWalletByUserID(txCtx, req.TargetUserID)
			if err != nil {
				return fmt.Errorf("failed to get receiver wallet: %w", err)
			}
			if receiverWallet == nil {
				return errors.New("receiver wallet not found")
			}

			newReceiverBalance := receiverWallet.Balance + float64(req.Amount)
			if err := s.walletRepo.UpdateBalance(txCtx, req.TargetUserID, newReceiverBalance); err != nil {
				return fmt.Errorf("failed to add to receiver wallet: %w", err)
			}
		} else if entity.TransferSource(req.Source) == entity.TransferSourcePaylater {
			// Check paylater credit
			senderAccount, err := s.paylaterAccountRepo.GetForUpdate(txCtx, req.UserID)
			if err != nil {
				return fmt.Errorf("failed to get sender paylater account: %w", err)
			}
			if senderAccount.AvailableLimit < req.Amount {
				return errors.New("insufficient paylater credit")
			}

			// Increase outstanding balance and update available limit
			senderAccount.Outstanding += req.Amount
			senderAccount.AvailableLimit = senderAccount.CreditLimit - senderAccount.Outstanding
			if err := s.paylaterAccountRepo.UpdateAccount(txCtx, senderAccount); err != nil {
				return fmt.Errorf("failed to update paylater outstanding: %w", err)
			}

			// Add to receiver wallet
			receiverWallet, err := s.walletRepo.GetWalletByUserID(txCtx, req.TargetUserID)
			if err != nil {
				return fmt.Errorf("failed to get receiver wallet: %w", err)
			}
			if receiverWallet == nil {
				return errors.New("receiver wallet not found")
			}

			newReceiverBalance := receiverWallet.Balance + float64(req.Amount)
			if err := s.walletRepo.UpdateBalance(txCtx, req.TargetUserID, newReceiverBalance); err != nil {
				return fmt.Errorf("failed to add to receiver wallet: %w", err)
			}
		}

		// Mark transfer as completed
		transfer.Status = entity.TransferStatusCompleted

		// Create transfer record
		if err := s.transferRepo.CreateTransfer(txCtx, transfer); err != nil {
			return fmt.Errorf("failed to create transfer record: %w", err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to create transfer", zap.Error(err))
		transfer.Status = entity.TransferStatusFailed
		// Try to save failed transfer record (best effort)
		_ = s.transferRepo.CreateTransfer(ctx, transfer)
		return nil, err
	}

	return transfer, nil
}

// GetTransferByID retrieves a transfer by ID
func (s *transferService) GetTransferByID(ctx context.Context, id uint) (*dto.GetTransferResponse, error) {
	transfer, err := s.transferRepo.GetTransferByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.mapToTransferResponse(transfer), nil
}

// GetTransfersByUserID retrieves all transfers sent by a user
func (s *transferService) GetTransfersByUserID(ctx context.Context, userID uint) ([]dto.GetTransferResponse, error) {
	transfers, err := s.transferRepo.GetTransfersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.mapToTransferResponses(transfers), nil
}

// GetTransfersByTargetUserID retrieves all transfers received by a user
func (s *transferService) GetTransfersByTargetUserID(ctx context.Context, targetUserID uint) ([]dto.GetTransferResponse, error) {
	transfers, err := s.transferRepo.GetTransfersByTargetUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	return s.mapToTransferResponses(transfers), nil
}

// ListTransfers retrieves transfers with filters and pagination
func (s *transferService) ListTransfers(ctx context.Context, req *dto.ListTransfersRequest) (*dto.ListTransfersResponse, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.TargetUserID != nil {
		filters["target_user_id"] = *req.TargetUserID
	}
	if req.Source != nil {
		filters["source"] = *req.Source
	}
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.FromDate != nil {
		if fromDate, err := time.Parse("2006-01-02", *req.FromDate); err == nil {
			filters["from_date"] = fromDate
		}
	}
	if req.ToDate != nil {
		if toDate, err := time.Parse("2006-01-02", *req.ToDate); err == nil {
			filters["to_date"] = toDate
		}
	}

	transfers, totalCount, err := s.transferRepo.ListTransfers(ctx, filters, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / req.PageSize
	if int(totalCount)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.ListTransfersResponse{
		Data:       s.mapToTransferResponses(transfers),
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateTransferStatus updates the status of a transfer
func (s *transferService) UpdateTransferStatus(ctx context.Context, id uint, req *dto.UpdateTransferStatusRequest) (*dto.GetTransferResponse, error) {
	// Verify transfer exists
	transfer, err := s.transferRepo.GetTransferByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status
	if err := s.transferRepo.UpdateTransferStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}

	// Fetch updated transfer
	transfer.Status = entity.TransferStatus(req.Status)
	return s.mapToTransferResponse(transfer), nil
}

// GetUserTransferStats retrieves transfer statistics for a user
func (s *transferService) GetUserTransferStats(ctx context.Context, userID uint) (*dto.TransferStatsResponse, error) {
	stats, err := s.transferRepo.GetUserTransferStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.TransferStatsResponse{
		UserID:        userID,
		TotalSent:     stats["total_sent"].(int64),
		TotalReceived: stats["total_received"].(int64),
		CountSent:     stats["count_sent"].(int64),
		CountReceived: stats["count_received"].(int64),
		WalletSent:    stats["wallet_sent"].(int64),
		PaylaterSent:  stats["paylater_sent"].(int64),
		CompletedSent: stats["completed_sent"].(int64),
		PendingSent:   stats["pending_sent"].(int64),
		FailedSent:    stats["failed_sent"].(int64),
	}, nil
}

// CancelTransfer cancels a pending transfer
func (s *transferService) CancelTransfer(ctx context.Context, id uint) error {
	// Verify transfer exists and is pending
	transfer, err := s.transferRepo.GetTransferByID(ctx, id)
	if err != nil {
		return err
	}

	if transfer.Status != entity.TransferStatusPending {
		return errors.New("can only cancel pending transfers")
	}

	return s.transferRepo.CancelTransfer(ctx, id)
}

// Helper methods
func (s *transferService) mapToTransferResponse(transfer *entity.Transfer) *dto.GetTransferResponse {
	return &dto.GetTransferResponse{
		ID:           transfer.ID,
		UserID:       transfer.UserID,
		TargetUserID: transfer.TargetUserID,
		Amount:       transfer.Amount,
		Source:       string(transfer.Source),
		Status:       string(transfer.Status),
		CreatedAt:    transfer.CreatedAt,
		UpdatedAt:    transfer.UpdatedAt,
	}
}

func (s *transferService) mapToTransferResponses(transfers []entity.Transfer) []dto.GetTransferResponse {
	responses := make([]dto.GetTransferResponse, len(transfers))
	for i, transfer := range transfers {
		responses[i] = *s.mapToTransferResponse(&transfer)
	}
	return responses
}
