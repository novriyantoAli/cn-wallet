package repository

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
)

// LedgerRepository defines the interface for ledger data access
type LedgerRepository interface {
	// CreateEntry creates a new ledger entry (append-only)
	CreateEntry(ctx context.Context, entry *entity.LedgerEntry) error

	// GetEntryByID retrieves a ledger entry by ID
	GetEntryByID(ctx context.Context, id uint64) (*entity.LedgerEntry, error)

	// GetEntriesByUserID retrieves all ledger entries for a user
	GetEntriesByUserID(ctx context.Context, userID uint64) ([]entity.LedgerEntry, error)

	// GetEntriesByReference retrieves ledger entries by reference
	GetEntriesByReference(ctx context.Context, referenceType string, referenceID uint64) ([]entity.LedgerEntry, error)

	// ListEntries retrieves paginated ledger entries with filters
	ListEntries(ctx context.Context, req *dto.ListLedgerEntriesRequest) ([]entity.LedgerEntry, int64, error)

	// GetUserStats retrieves ledger statistics for a user
	GetUserStats(ctx context.Context, userID uint64) (*dto.LedgerStatsResponse, error)
}
