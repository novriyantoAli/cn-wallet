package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
)

type ledgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) CreateEntry(ctx context.Context, entry *entity.LedgerEntry) error {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	db := database.GetDB(ctx, r.db)
	return db.Create(entry).Error
}

func (r *ledgerRepository) GetEntryByID(ctx context.Context, id uint64) (*entity.LedgerEntry, error) {
	var entry entity.LedgerEntry
	db := database.GetDB(ctx, r.db)
	err := db.Where("id = ?", id).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *ledgerRepository) GetEntriesByUserID(ctx context.Context, userID uint64) ([]entity.LedgerEntry, error) {
	var entries []entity.LedgerEntry
	db := database.GetDB(ctx, r.db)
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&entries).Error
	return entries, err
}

func (r *ledgerRepository) GetEntriesByReference(ctx context.Context, referenceType string, referenceID uint64) ([]entity.LedgerEntry, error) {
	var entries []entity.LedgerEntry
	db := database.GetDB(ctx, r.db)
	err := db.
		Where("reference_type = ? AND reference_id = ?", referenceType, referenceID).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}

func (r *ledgerRepository) ListEntries(ctx context.Context, req *dto.ListLedgerEntriesRequest) ([]entity.LedgerEntry, int64, error) {
	var entries []entity.LedgerEntry
	query := database.GetDB(ctx, r.db)

	// Apply filters
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.ReferenceType != nil {
		query = query.Where("reference_type = ?", *req.ReferenceType)
	}
	if req.AccountType != nil {
		query = query.Where("account_type = ?", *req.AccountType)
	}
	if req.FromDate != nil {
		fromDate, err := time.Parse("2006-01-02", *req.FromDate)
		if err == nil {
			query = query.Where("created_at >= ?", fromDate)
		}
	}
	if req.ToDate != nil {
		toDate, err := time.Parse("2006-01-02", *req.ToDate)
		if err == nil {
			// Add 1 day to include the entire end date
			toDate = toDate.Add(24 * time.Hour)
			query = query.Where("created_at < ?", toDate)
		}
	}

	// Get total count
	var totalCount int64
	if err := query.Model(&entity.LedgerEntry{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get entries
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, totalCount, nil
}

func (r *ledgerRepository) GetUserStats(ctx context.Context, userID uint64) (*dto.LedgerStatsResponse, error) {
	stats := &dto.LedgerStatsResponse{
		UserID: userID,
	}

	// Get aggregated stats
	type StatsResult struct {
		TotalDebit     int64
		TotalCredit    int64
		EntryCount     int64
		WalletDebit    int64
		WalletCredit   int64
		PaylaterDebit  int64
		PaylaterCredit int64
	}

	var result StatsResult
	db := database.GetDB(ctx, r.db)

	err := db.Model(&entity.LedgerEntry{}).
		Select(
			"COALESCE(SUM(CASE WHEN debit > 0 THEN debit ELSE 0 END), 0) as total_debit",
			"COALESCE(SUM(CASE WHEN credit > 0 THEN credit ELSE 0 END), 0) as total_credit",
			"COUNT(*) as entry_count",
			"COALESCE(SUM(CASE WHEN account_type = 'wallet' AND debit > 0 THEN debit ELSE 0 END), 0) as wallet_debit",
			"COALESCE(SUM(CASE WHEN account_type = 'wallet' AND credit > 0 THEN credit ELSE 0 END), 0) as wallet_credit",
			"COALESCE(SUM(CASE WHEN account_type = 'paylater' AND debit > 0 THEN debit ELSE 0 END), 0) as paylater_debit",
			"COALESCE(SUM(CASE WHEN account_type = 'paylater' AND credit > 0 THEN credit ELSE 0 END), 0) as paylater_credit",
		).
		Where("user_id = ?", userID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.TotalDebit = result.TotalDebit
	stats.TotalCredit = result.TotalCredit
	stats.EntryCount = result.EntryCount
	stats.WalletDebit = result.WalletDebit
	stats.WalletCredit = result.WalletCredit
	stats.PaylaterDebit = result.PaylaterDebit
	stats.PaylaterCredits = result.PaylaterCredit

	return stats, nil
}
