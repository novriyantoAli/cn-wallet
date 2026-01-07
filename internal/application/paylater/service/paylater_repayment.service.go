package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	ledgerRepository "github.com/novriyantoAli/cn-wallet/internal/application/ledger/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"go.uber.org/zap"
)

// PaylaterRepaymentService defines the interface for paylater repayment business logic
type PaylaterRepaymentService interface {
	// ProcessRepayment processes a repayment for a paylater loan
	ProcessRepayment(ctx context.Context, req *dto.CreatePaylaterRepaymentRequest) (*dto.GetPaylaterRepaymentResponse, error)

	// GetRepaymentByID retrieves a repayment by ID
	GetRepaymentByID(ctx context.Context, id uint) (*dto.GetPaylaterRepaymentResponse, error)

	// ListRepaymentsByLoan retrieves all repayments for a specific loan
	ListRepaymentsByLoan(ctx context.Context, loanID uint) ([]dto.GetPaylaterRepaymentResponse, error)

	// GetRepaymentsByUser retrieves all repayments for a specific user
	GetRepaymentsByUser(ctx context.Context, userID uint) ([]dto.GetPaylaterRepaymentResponse, error)
}

type paylaterRepaymentService struct {
	repaymentRepo repository.PaylaterRepaymentRepository
	loanRepo      repository.PaylaterLoanRepository
	ledgerRepo    ledgerRepository.LedgerRepository
	txManager     database.TransactionManagerI
	logger        *zap.Logger
}

// NewPaylaterRepaymentService creates a new paylater repayment service
func NewPaylaterRepaymentService(
	repaymentRepo repository.PaylaterRepaymentRepository,
	loanRepo repository.PaylaterLoanRepository,
	ledgerRepo ledgerRepository.LedgerRepository,
	txManager database.TransactionManagerI,
	logger *zap.Logger,
) PaylaterRepaymentService {
	return &paylaterRepaymentService{
		repaymentRepo: repaymentRepo,
		loanRepo:      loanRepo,
		ledgerRepo:    ledgerRepo,
		txManager:     txManager,
		logger:        logger,
	}
}

// ProcessRepayment processes a repayment for a paylater loan
// Flow:
// 1. Create repayment record
// 2. Create ledger entries for cash inflow and paylater debt reduction
// 3. Update loan status to 'paid' if total repayments >= loan total
func (s *paylaterRepaymentService) ProcessRepayment(ctx context.Context, req *dto.CreatePaylaterRepaymentRequest) (*dto.GetPaylaterRepaymentResponse, error) {
	var repayment *paylaterEntity.PaylaterRepayment
	var err error

	s.logger.Info("Processing paylater repayment",
		zap.Uint("loan_id", req.PaylaterLoanID),
		zap.Uint("user_id", req.UserID),
		zap.Int64("amount", req.Amount),
		zap.String("payment_source", req.PaymentSource),
	)

	// Execute in transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get loan to verify it exists and check status
		loan, err := s.loanRepo.GetLoanByID(txCtx, req.PaylaterLoanID)
		if err != nil {
			s.logger.Error("failed to get loan", zap.Error(err))
			return errors.New("loan not found")
		}

		// Verify loan belongs to user
		if loan.UserID != req.UserID {
			s.logger.Warn("unauthorized repayment attempt",
				zap.Uint("loan_id", req.PaylaterLoanID),
				zap.Uint("user_id", req.UserID),
				zap.Uint("actual_user_id", loan.UserID),
			)
			return errors.New("unauthorized repayment")
		}

		// Verify loan is active (can only repay active loans)
		if loan.Status != paylaterEntity.PaylaterLoanStatusActive && loan.Status != paylaterEntity.PaylaterLoanStatusOverdue {
			s.logger.Warn("cannot repay loan with status",
				zap.Uint("loan_id", req.PaylaterLoanID),
				zap.String("status", string(loan.Status)),
			)
			return fmt.Errorf("cannot repay loan with status: %s", loan.Status)
		}

		// Validate repayment amount
		if req.Amount <= 0 {
			s.logger.Error("invalid repayment amount", zap.Int64("amount", req.Amount))
			return errors.New("repayment amount must be greater than 0")
		}

		// 1. Create repayment record
		repayment = &paylaterEntity.PaylaterRepayment{
			PaylaterLoanID: req.PaylaterLoanID,
			UserID:         req.UserID,
			Amount:         req.Amount,
			PaymentSource:  req.PaymentSource,
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		if err := s.repaymentRepo.Create(txCtx, repayment); err != nil {
			s.logger.Error("failed to create repayment record", zap.Error(err))
			return fmt.Errorf("failed to create repayment: %w", err)
		}

		s.logger.Info("repayment record created",
			zap.Uint("repayment_id", repayment.ID),
			zap.Uint("loan_id", req.PaylaterLoanID),
		)

		// 2. Create ledger entries
		refID := fmt.Sprintf("%d", repayment.ID)

		// Entry 1: Cash inflow from external/wallet payment to system cash account
		ledgerEntry1 := &entity.LedgerEntry{
			UserID:        uint64(req.UserID),
			ReferenceID:   refID,
			ReferenceType: entity.ReferenceTypeRepayment,
			Debit:         0,
			Credit:        req.Amount,
			AccountType:   entity.AccountTypeMerchantIncome, // system_cash (incoming payment)
			CreatedAt:     time.Now(),
		}

		if err := s.ledgerRepo.CreateEntry(txCtx, ledgerEntry1); err != nil {
			s.logger.Error("failed to create ledger entry for cash inflow", zap.Error(err))
			return fmt.Errorf("failed to create cash inflow ledger entry: %w", err)
		}

		s.logger.Info("ledger entry created for cash inflow",
			zap.Uint64("entry_id", ledgerEntry1.ID),
			zap.Int64("amount", req.Amount),
		)

		// Entry 2: Paylater debt reduction (debit paylater account)
		ledgerEntry2 := &entity.LedgerEntry{
			UserID:        uint64(req.UserID),
			ReferenceID:   refID,
			ReferenceType: entity.ReferenceTypeRepayment,
			Debit:         req.Amount,
			Credit:        0,
			AccountType:   entity.AccountTypePaylater, // paylater debt reduction
			CreatedAt:     time.Now(),
		}

		if err := s.ledgerRepo.CreateEntry(txCtx, ledgerEntry2); err != nil {
			s.logger.Error("failed to create ledger entry for paylater reduction", zap.Error(err))
			return fmt.Errorf("failed to create paylater reduction ledger entry: %w", err)
		}

		s.logger.Info("ledger entry created for paylater debt reduction",
			zap.Uint64("entry_id", ledgerEntry2.ID),
			zap.Int64("amount", req.Amount),
		)

		// 3. Update loan status to 'paid' if total repayments >= loan total
		// This checks if the cumulative paylater debt reduction is >= loan total
		if err := s.loanRepo.UpdateLoanStatus(txCtx, req.PaylaterLoanID, string(paylaterEntity.PaylaterLoanStatusPaid)); err != nil {
			s.logger.Error("failed to update loan status", zap.Error(err))
			return fmt.Errorf("failed to update loan status: %w", err)
		}

		s.logger.Info("loan status updated to paid",
			zap.Uint("loan_id", req.PaylaterLoanID),
		)

		return nil
	})

	if err != nil {
		s.logger.Error("repayment processing failed",
			zap.Error(err),
			zap.Uint("loan_id", req.PaylaterLoanID),
		)
		return nil, err
	}

	s.logger.Info("repayment processed successfully",
		zap.Uint("repayment_id", repayment.ID),
		zap.Uint("loan_id", req.PaylaterLoanID),
	)

	return &dto.GetPaylaterRepaymentResponse{
		ID:             repayment.ID,
		PaylaterLoanID: repayment.PaylaterLoanID,
		UserID:         repayment.UserID,
		Amount:         repayment.Amount,
		PaymentSource:  repayment.PaymentSource,
		Status:         repayment.Status,
		CreatedAt:      repayment.CreatedAt,
	}, nil
}

// GetRepaymentByID retrieves a repayment by ID
func (s *paylaterRepaymentService) GetRepaymentByID(ctx context.Context, id uint) (*dto.GetPaylaterRepaymentResponse, error) {
	repayment, err := s.repaymentRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get repayment", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return &dto.GetPaylaterRepaymentResponse{
		ID:             repayment.ID,
		PaylaterLoanID: repayment.PaylaterLoanID,
		UserID:         repayment.UserID,
		Amount:         repayment.Amount,
		PaymentSource:  repayment.PaymentSource,
		Status:         repayment.Status,
		CreatedAt:      repayment.CreatedAt,
	}, nil
}

// ListRepaymentsByLoan retrieves all repayments for a specific loan
func (s *paylaterRepaymentService) ListRepaymentsByLoan(ctx context.Context, loanID uint) ([]dto.GetPaylaterRepaymentResponse, error) {
	s.logger.Info("Listing repayments by loan", zap.Uint("loan_id", loanID))

	repayments, err := s.repaymentRepo.GetByLoanID(ctx, loanID)
	if err != nil {
		s.logger.Error("failed to get repayments by loan", zap.Error(err), zap.Uint("loan_id", loanID))
		return nil, err
	}

	response := make([]dto.GetPaylaterRepaymentResponse, len(repayments))
	for i, repayment := range repayments {
		response[i] = dto.GetPaylaterRepaymentResponse{
			ID:             repayment.ID,
			PaylaterLoanID: repayment.PaylaterLoanID,
			UserID:         repayment.UserID,
			Amount:         repayment.Amount,
			PaymentSource:  repayment.PaymentSource,
			Status:         repayment.Status,
			CreatedAt:      repayment.CreatedAt,
		}
	}

	s.logger.Info("repayments retrieved",
		zap.Uint("loan_id", loanID),
		zap.Int("count", len(response)),
	)

	return response, nil
}

// GetRepaymentsByUser retrieves all repayments for a specific user
func (s *paylaterRepaymentService) GetRepaymentsByUser(ctx context.Context, userID uint) ([]dto.GetPaylaterRepaymentResponse, error) {
	s.logger.Info("Listing repayments by user", zap.Uint("user_id", userID))

	repayments, err := s.repaymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get repayments by user", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	response := make([]dto.GetPaylaterRepaymentResponse, len(repayments))
	for i, repayment := range repayments {
		response[i] = dto.GetPaylaterRepaymentResponse{
			ID:             repayment.ID,
			PaylaterLoanID: repayment.PaylaterLoanID,
			UserID:         repayment.UserID,
			Amount:         repayment.Amount,
			PaymentSource:  repayment.PaymentSource,
			Status:         repayment.Status,
			CreatedAt:      repayment.CreatedAt,
		}
	}

	s.logger.Info("repayments retrieved",
		zap.Uint("user_id", userID),
		zap.Int("count", len(response)),
	)

	return response, nil
}
