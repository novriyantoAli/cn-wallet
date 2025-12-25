package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/repository"
)

// LedgerService defines the interface for ledger business logic
type LedgerService interface {
	// CreateEntry creates a new ledger entry
	CreateEntry(ctx context.Context, req *dto.CreateLedgerEntryRequest) (*dto.GetLedgerEntryResponse, error)

	// GetEntryByID retrieves a ledger entry by ID
	GetEntryByID(ctx context.Context, id uint64) (*dto.GetLedgerEntryResponse, error)

	// GetEntriesByUserID retrieves all ledger entries for a user
	GetEntriesByUserID(ctx context.Context, userID uint64) ([]dto.GetLedgerEntryResponse, error)

	// GetEntriesByReference retrieves ledger entries by reference
	GetEntriesByReference(ctx context.Context, referenceType string, referenceID uint64) ([]dto.GetLedgerEntryResponse, error)

	// ListEntries retrieves paginated ledger entries with filters
	ListEntries(ctx context.Context, req *dto.ListLedgerEntriesRequest) (*dto.ListLedgerEntriesResponse, error)

	// GetUserStats retrieves ledger statistics for a user
	GetUserStats(ctx context.Context, userID uint64) (*dto.LedgerStatsResponse, error)
}

type ledgerService struct {
	repo   repository.LedgerRepository
	logger *zap.Logger
}

func NewLedgerService(repo repository.LedgerRepository, logger *zap.Logger) LedgerService {
	return &ledgerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ledgerService) CreateEntry(ctx context.Context, req *dto.CreateLedgerEntryRequest) (*dto.GetLedgerEntryResponse, error) {
	// Validate that either debit or credit is set, but not both
	if (req.Debit > 0 && req.Credit > 0) || (req.Debit == 0 && req.Credit == 0) {
		return nil, errors.New("either debit or credit must be set, but not both")
	}

	entry := &entity.LedgerEntry{
		UserID:        req.UserID,
		ReferenceID:   req.ReferenceID,
		ReferenceType: entity.ReferenceType(req.ReferenceType),
		Debit:         req.Debit,
		Credit:        req.Credit,
		AccountType:   entity.AccountType(req.AccountType),
	}

	if err := s.repo.CreateEntry(ctx, entry); err != nil {
		s.logger.Error("failed to create ledger entry", zap.Error(err))
		return nil, err
	}

	return &dto.GetLedgerEntryResponse{
		ID:            entry.ID,
		UserID:        entry.UserID,
		ReferenceID:   entry.ReferenceID,
		ReferenceType: string(entry.ReferenceType),
		Debit:         entry.Debit,
		Credit:        entry.Credit,
		AccountType:   string(entry.AccountType),
		Amount:        entry.Amount(),
		CreatedAt:     entry.CreatedAt,
	}, nil
}

func (s *ledgerService) GetEntryByID(ctx context.Context, id uint64) (*dto.GetLedgerEntryResponse, error) {
	entry, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.GetLedgerEntryResponse{
		ID:            entry.ID,
		UserID:        entry.UserID,
		ReferenceID:   entry.ReferenceID,
		ReferenceType: string(entry.ReferenceType),
		Debit:         entry.Debit,
		Credit:        entry.Credit,
		AccountType:   string(entry.AccountType),
		Amount:        entry.Amount(),
		CreatedAt:     entry.CreatedAt,
	}, nil
}

func (s *ledgerService) GetEntriesByUserID(ctx context.Context, userID uint64) ([]dto.GetLedgerEntryResponse, error) {
	entries, err := s.repo.GetEntriesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GetLedgerEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dto.GetLedgerEntryResponse{
			ID:            entry.ID,
			UserID:        entry.UserID,
			ReferenceID:   entry.ReferenceID,
			ReferenceType: string(entry.ReferenceType),
			Debit:         entry.Debit,
			Credit:        entry.Credit,
			AccountType:   string(entry.AccountType),
			Amount:        entry.Amount(),
			CreatedAt:     entry.CreatedAt,
		}
	}

	return responses, nil
}

func (s *ledgerService) GetEntriesByReference(ctx context.Context, referenceType string, referenceID uint64) ([]dto.GetLedgerEntryResponse, error) {
	entries, err := s.repo.GetEntriesByReference(ctx, referenceType, referenceID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GetLedgerEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dto.GetLedgerEntryResponse{
			ID:            entry.ID,
			UserID:        entry.UserID,
			ReferenceID:   entry.ReferenceID,
			ReferenceType: string(entry.ReferenceType),
			Debit:         entry.Debit,
			Credit:        entry.Credit,
			AccountType:   string(entry.AccountType),
			Amount:        entry.Amount(),
			CreatedAt:     entry.CreatedAt,
		}
	}

	return responses, nil
}

func (s *ledgerService) ListEntries(ctx context.Context, req *dto.ListLedgerEntriesRequest) (*dto.ListLedgerEntriesResponse, error) {
	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	entries, totalCount, err := s.repo.ListEntries(ctx, req)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GetLedgerEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dto.GetLedgerEntryResponse{
			ID:            entry.ID,
			UserID:        entry.UserID,
			ReferenceID:   entry.ReferenceID,
			ReferenceType: string(entry.ReferenceType),
			Debit:         entry.Debit,
			Credit:        entry.Credit,
			AccountType:   string(entry.AccountType),
			Amount:        entry.Amount(),
			CreatedAt:     entry.CreatedAt,
		}
	}

	totalPages := (totalCount + int64(req.PageSize) - 1) / int64(req.PageSize)

	return &dto.ListLedgerEntriesResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (s *ledgerService) GetUserStats(ctx context.Context, userID uint64) (*dto.LedgerStatsResponse, error) {
	return s.repo.GetUserStats(ctx, userID)
}
