package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
)

// PaylaterLoanService defines the interface for paylater loan business logic
type PaylaterLoanService interface {
	CreateLoan(ctx context.Context, req *dto.CreatePaylaterLoanRequest) (*entity.PaylaterLoan, error)
	GetLoanByID(ctx context.Context, id uint) (*dto.GetPaylaterLoanResponse, error)
	GetLoansByUserID(ctx context.Context, userID uint) ([]dto.GetPaylaterLoanResponse, error)
	ListLoans(ctx context.Context, req *dto.ListPaylaterLoansRequest) (*dto.ListPaylaterLoansResponse, error)
	UpdateLoanStatus(ctx context.Context, id uint, req *dto.UpdateLoanStatusRequest) (*dto.GetPaylaterLoanResponse, error)
	MarkLoanAsPaid(ctx context.Context, id uint) (*dto.GetPaylaterLoanResponse, error)
	ProcessOverdueLoans(ctx context.Context) (int, error)
	GetUserLoanStats(ctx context.Context, userID uint) (*dto.PaylaterLoanStatsResponse, error)
	DeleteLoan(ctx context.Context, id uint) error
}

type paylaterLoanService struct {
	loanRepo    repository.PaylaterLoanRepository
	accountRepo repository.PaylaterAccountRepository
	txManager   database.TransactionManagerI
	logger      *zap.Logger
}

// NewPaylaterLoanService creates a new paylater loan service
func NewPaylaterLoanService(
	loanRepo repository.PaylaterLoanRepository,
	accountRepo repository.PaylaterAccountRepository,
	txManager database.TransactionManagerI,
	logger *zap.Logger,
) PaylaterLoanService {
	return &paylaterLoanService{
		loanRepo:    loanRepo,
		accountRepo: accountRepo,
		txManager:   txManager,
		logger:      logger,
	}
}

// CreateLoan creates a new paylater loan and deducts from available credit
func (s *paylaterLoanService) CreateLoan(ctx context.Context, req *dto.CreatePaylaterLoanRequest) (*entity.PaylaterLoan, error) {
	var loan *entity.PaylaterLoan
	var err error

	// Execute in transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Verify paylater account exists and is active
		account, err := s.accountRepo.GetForUpdate(txCtx, req.UserID)
		if err != nil {
			return errors.New("paylater account not found or inactive")
		}

		if account.Status != entity.PaylaterStatusActive {
			return errors.New("paylater account is not active")
		}

		// Calculate total
		total := req.Amount + req.Interest

		// Check if user has sufficient credit
		if account.AvailableLimit < total {
			return errors.New("insufficient credit limit for this loan")
		}

		// Create loan entity
		loan = &entity.PaylaterLoan{
			UserID:   req.UserID,
			Amount:   req.Amount,
			Interest: req.Interest,
			Total:    total,
			DueDate:  req.DueDate,
			Source:   req.Source,
			Status:   entity.PaylaterLoanStatusActive,
		}

		// Save loan to repository
		if err := s.loanRepo.CreateLoan(txCtx, loan); err != nil {
			return err
		}

		// Update account outstanding and available limit
		account.Outstanding += total
		account.AvailableLimit = account.CreditLimit - account.Outstanding

		if err := s.accountRepo.UpdateAccount(txCtx, account); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to create paylater loan", zap.Error(err), zap.Uint("user_id", req.UserID))
		return nil, err
	}

	s.logger.Info("Paylater loan created successfully",
		zap.Uint("user_id", req.UserID),
		zap.Uint("loan_id", loan.ID),
		zap.Int64("amount", loan.Amount))

	return loan, nil
}

// GetLoanByID retrieves a paylater loan by ID
func (s *paylaterLoanService) GetLoanByID(ctx context.Context, id uint) (*dto.GetPaylaterLoanResponse, error) {
	loan, err := s.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get paylater loan", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	return &dto.GetPaylaterLoanResponse{
		ID:        loan.ID,
		UserID:    loan.UserID,
		Amount:    loan.Amount,
		Interest:  loan.Interest,
		Total:     loan.Total,
		DueDate:   loan.DueDate,
		Source:    loan.Source,
		Status:    loan.Status,
		CreatedAt: loan.CreatedAt,
	}, nil
}

// GetLoansByUserID retrieves all loans for a specific user
func (s *paylaterLoanService) GetLoansByUserID(ctx context.Context, userID uint) ([]dto.GetPaylaterLoanResponse, error) {
	loans, err := s.loanRepo.GetLoansByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get loans by user ID", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	responses := make([]dto.GetPaylaterLoanResponse, 0, len(loans))
	for _, loan := range loans {
		responses = append(responses, dto.GetPaylaterLoanResponse{
			ID:        loan.ID,
			UserID:    loan.UserID,
			Amount:    loan.Amount,
			Interest:  loan.Interest,
			Total:     loan.Total,
			DueDate:   loan.DueDate,
			Source:    loan.Source,
			Status:    loan.Status,
			CreatedAt: loan.CreatedAt,
		})
	}

	return responses, nil
}

// ListLoans retrieves loans with filters and pagination
func (s *paylaterLoanService) ListLoans(ctx context.Context, req *dto.ListPaylaterLoansRequest) (*dto.ListPaylaterLoansResponse, error) {
	// Set default pagination values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.Source != nil {
		filters["source"] = *req.Source
	}
	if req.FromDate != nil {
		fromDate, err := time.Parse("2006-01-02", *req.FromDate)
		if err == nil {
			filters["from_date"] = fromDate
		}
	}
	if req.ToDate != nil {
		toDate, err := time.Parse("2006-01-02", *req.ToDate)
		if err == nil {
			filters["to_date"] = toDate
		}
	}

	loans, totalCount, err := s.loanRepo.ListLoans(ctx, filters, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list loans", zap.Error(err))
		return nil, err
	}

	// Convert to response DTOs
	loanResponses := make([]dto.GetPaylaterLoanResponse, 0, len(loans))
	for _, loan := range loans {
		loanResponses = append(loanResponses, dto.GetPaylaterLoanResponse{
			ID:        loan.ID,
			UserID:    loan.UserID,
			Amount:    loan.Amount,
			Interest:  loan.Interest,
			Total:     loan.Total,
			DueDate:   loan.DueDate,
			Source:    loan.Source,
			Status:    loan.Status,
			CreatedAt: loan.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	return &dto.ListPaylaterLoansResponse{
		Data:       loanResponses,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateLoanStatus updates the status of a loan
func (s *paylaterLoanService) UpdateLoanStatus(ctx context.Context, id uint, req *dto.UpdateLoanStatusRequest) (*dto.GetPaylaterLoanResponse, error) {
	// Get the loan first
	loan, err := s.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status
	if err := s.loanRepo.UpdateLoanStatus(ctx, id, req.Status); err != nil {
		s.logger.Error("Failed to update loan status", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	loan.Status = req.Status

	s.logger.Info("Loan status updated", zap.Uint("id", id), zap.String("status", req.Status))

	return &dto.GetPaylaterLoanResponse{
		ID:        loan.ID,
		UserID:    loan.UserID,
		Amount:    loan.Amount,
		Interest:  loan.Interest,
		Total:     loan.Total,
		DueDate:   loan.DueDate,
		Source:    loan.Source,
		Status:    loan.Status,
		CreatedAt: loan.CreatedAt,
	}, nil
}

// MarkLoanAsPaid marks a loan as paid and restores credit to the account
func (s *paylaterLoanService) MarkLoanAsPaid(ctx context.Context, id uint) (*dto.GetPaylaterLoanResponse, error) {
	var loan *entity.PaylaterLoan
	var err error

	// Execute in transaction to ensure atomicity
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Get loan
		loan, err = s.loanRepo.GetLoanByID(txCtx, id)
		if err != nil {
			return err
		}

		// Check if already paid
		if loan.Status == entity.PaylaterLoanStatusPaid {
			return errors.New("loan is already marked as paid")
		}

		// Get account with lock
		account, err := s.accountRepo.GetForUpdate(txCtx, loan.UserID)
		if err != nil {
			return err
		}

		// Update loan status to paid
		loan.Status = entity.PaylaterLoanStatusPaid
		if err := s.loanRepo.UpdateLoan(txCtx, loan); err != nil {
			return err
		}

		// Restore credit to account
		account.Outstanding -= loan.Total
		if account.Outstanding < 0 {
			account.Outstanding = 0
		}
		account.AvailableLimit = account.CreditLimit - account.Outstanding

		if err := s.accountRepo.UpdateAccount(txCtx, account); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to mark loan as paid", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	s.logger.Info("Loan marked as paid",
		zap.Uint("id", id),
		zap.Uint("user_id", loan.UserID),
		zap.Int64("amount", loan.Total))

	return &dto.GetPaylaterLoanResponse{
		ID:        loan.ID,
		UserID:    loan.UserID,
		Amount:    loan.Amount,
		Interest:  loan.Interest,
		Total:     loan.Total,
		DueDate:   loan.DueDate,
		Source:    loan.Source,
		Status:    loan.Status,
		CreatedAt: loan.CreatedAt,
	}, nil
}

// ProcessOverdueLoans processes all overdue loans and updates their status
func (s *paylaterLoanService) ProcessOverdueLoans(ctx context.Context) (int, error) {
	loans, err := s.loanRepo.GetOverdueLoans(ctx, time.Now())
	if err != nil {
		s.logger.Error("Failed to get overdue loans", zap.Error(err))
		return 0, err
	}

	count := 0
	for _, loan := range loans {
		if err := s.loanRepo.UpdateLoanStatus(ctx, loan.ID, entity.PaylaterLoanStatusOverdue); err != nil {
			s.logger.Error("Failed to mark loan as overdue",
				zap.Error(err),
				zap.Uint("loan_id", loan.ID))
			continue
		}
		count++
	}

	s.logger.Info("Processed overdue loans", zap.Int("count", count))
	return count, nil
}

// GetUserLoanStats retrieves loan statistics for a user
func (s *paylaterLoanService) GetUserLoanStats(ctx context.Context, userID uint) (*dto.PaylaterLoanStatsResponse, error) {
	stats, err := s.loanRepo.GetLoanStats(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get loan stats", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}

	return &dto.PaylaterLoanStatsResponse{
		UserID:           userID,
		TotalLoans:       stats["total_loans"].(int64),
		ActiveLoans:      stats["active_loans"].(int64),
		TotalBorrowed:    stats["total_borrowed"].(int64),
		TotalPaid:        stats["total_paid"].(int64),
		TotalOutstanding: stats["total_outstanding"].(int64),
		OverdueLoans:     stats["overdue_loans"].(int64),
	}, nil
}

// DeleteLoan deletes a paylater loan
func (s *paylaterLoanService) DeleteLoan(ctx context.Context, id uint) error {
	// Get loan first to check status
	loan, err := s.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deletion of paid or cancelled loans
	if loan.Status != entity.PaylaterLoanStatusPaid {
		return errors.New("can only delete paid loans")
	}

	if err := s.loanRepo.DeleteLoan(ctx, id); err != nil {
		s.logger.Error("Failed to delete paylater loan", zap.Error(err), zap.Uint("id", id))
		return err
	}

	s.logger.Info("Paylater loan deleted successfully", zap.Uint("id", id))
	return nil
}
