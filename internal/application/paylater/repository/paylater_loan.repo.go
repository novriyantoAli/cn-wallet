package repository

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PaylaterLoanRepository defines the interface for paylater loan data access
type PaylaterLoanRepository interface {
	CreateLoan(ctx context.Context, loan *entity.PaylaterLoan) error
	GetLoanByID(ctx context.Context, id uint) (*entity.PaylaterLoan, error)
	GetLoansByUserID(ctx context.Context, userID uint) ([]entity.PaylaterLoan, error)
	ListLoans(ctx context.Context, filters *dto.ListPaylaterLoansRequest) ([]entity.PaylaterLoan, int64, error)
	UpdateLoanStatus(ctx context.Context, id uint, status string) error
	UpdateLoan(ctx context.Context, loan *entity.PaylaterLoan) error
	GetOverdueLoans(ctx context.Context, asOf time.Time) ([]entity.PaylaterLoan, error)
	GetLoanStats(ctx context.Context, userID uint) (map[string]interface{}, error)
	DeleteLoan(ctx context.Context, id uint) error
}

type paylaterLoanRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewPaylaterLoanRepository creates a new paylater loan repository
func NewPaylaterLoanRepository(db *gorm.DB, logger *zap.Logger) PaylaterLoanRepository {
	return &paylaterLoanRepository{
		db:     db,
		logger: logger,
	}
}

// CreateLoan creates a new paylater loan record
func (r *paylaterLoanRepository) CreateLoan(ctx context.Context, loan *entity.PaylaterLoan) error {
	r.logger.Info("Creating paylater loan", zap.Uint("user_id", loan.UserID))
	db := database.GetDB(ctx, r.db)
	return db.Create(loan).Error
}

// GetLoanByID retrieves a paylater loan by ID
func (r *paylaterLoanRepository) GetLoanByID(ctx context.Context, id uint) (*entity.PaylaterLoan, error) {
	var loan entity.PaylaterLoan
	db := database.GetDB(ctx, r.db)
	result := db.First(&loan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("Paylater loan not found", zap.Uint("id", id))
			return nil, errors.New("paylater loan not found")
		}
		r.logger.Error("Failed to get paylater loan", zap.Error(result.Error), zap.Uint("id", id))
		return nil, result.Error
	}
	return &loan, nil
}

// GetLoansByUserID retrieves all loans for a specific user
func (r *paylaterLoanRepository) GetLoansByUserID(ctx context.Context, userID uint) ([]entity.PaylaterLoan, error) {
	var loans []entity.PaylaterLoan
	db := database.GetDB(ctx, r.db)
	result := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&loans)
	if result.Error != nil {
		r.logger.Error("Failed to get loans by user ID", zap.Error(result.Error), zap.Uint("user_id", userID))
		return nil, result.Error
	}
	return loans, nil
}

// ListLoans retrieves loans with filters and pagination
func (r *paylaterLoanRepository) ListLoans(ctx context.Context, req *dto.ListPaylaterLoansRequest) ([]entity.PaylaterLoan, int64, error) {
	var loans []entity.PaylaterLoan
	var totalCount int64

	db := database.GetDB(ctx, r.db)
	query := db.Model(&entity.PaylaterLoan{})

	// Apply filters
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Source != nil {
		query = query.Where("source = ?", *req.Source)
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
			query = query.Where("created_at <= ?", toDate)
		}
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to count loans", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&loans)
	if result.Error != nil {
		r.logger.Error("Failed to list loans", zap.Error(result.Error))
		return nil, 0, result.Error
	}

	return loans, totalCount, nil
}

// UpdateLoanStatus updates the status of a loan
func (r *paylaterLoanRepository) UpdateLoanStatus(ctx context.Context, id uint, status string) error {
	db := database.GetDB(ctx, r.db)
	result := db.Model(&entity.PaylaterLoan{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		r.logger.Error("Failed to update loan status", zap.Error(result.Error), zap.Uint("id", id))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("loan not found")
	}
	return nil
}

// UpdateLoan updates a paylater loan record
func (r *paylaterLoanRepository) UpdateLoan(ctx context.Context, loan *entity.PaylaterLoan) error {
	db := database.GetDB(ctx, r.db)
	return db.Save(loan).Error
}

// GetOverdueLoans retrieves all overdue loans as of a specific date
func (r *paylaterLoanRepository) GetOverdueLoans(ctx context.Context, asOf time.Time) ([]entity.PaylaterLoan, error) {
	var loans []entity.PaylaterLoan
	db := database.GetDB(ctx, r.db)
	result := db.Where("due_date < ? AND status IN (?)", asOf, []string{
		string(entity.PaylaterLoanStatusActive),
		string(entity.PaylaterLoanStatusPending),
	}).Find(&loans)
	if result.Error != nil {
		r.logger.Error("Failed to get overdue loans", zap.Error(result.Error))
		return nil, result.Error
	}
	return loans, nil
}

// GetLoanStats retrieves loan statistics for a user
func (r *paylaterLoanRepository) GetLoanStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	db := database.GetDB(ctx, r.db)
	stats := make(map[string]interface{})

	// Total loans count
	var totalLoans int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ?", userID).Count(&totalLoans)
	stats["total_loans"] = totalLoans

	// Active loans count
	var activeLoans int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ? AND status = ?", userID, entity.PaylaterLoanStatusActive).Count(&activeLoans)
	stats["active_loans"] = activeLoans

	// Overdue loans count
	var overdueLoans int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ? AND status = ?", userID, entity.PaylaterLoanStatusOverdue).Count(&overdueLoans)
	stats["overdue_loans"] = overdueLoans

	// Total borrowed (sum of all loan totals)
	var totalBorrowed int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ?", userID).Select("COALESCE(SUM(total), 0)").Scan(&totalBorrowed)
	stats["total_borrowed"] = totalBorrowed

	// Total paid (sum of paid loan totals)
	var totalPaid int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ? AND status = ?", userID, entity.PaylaterLoanStatusPaid).Select("COALESCE(SUM(total), 0)").Scan(&totalPaid)
	stats["total_paid"] = totalPaid

	// Total outstanding (sum of active and overdue loan totals)
	var totalOutstanding int64
	db.Model(&entity.PaylaterLoan{}).Where("user_id = ? AND status IN (?)", userID, []string{
		string(entity.PaylaterLoanStatusActive),
		string(entity.PaylaterLoanStatusOverdue),
		string(entity.PaylaterLoanStatusPending),
	}).Select("COALESCE(SUM(total), 0)").Scan(&totalOutstanding)
	stats["total_outstanding"] = totalOutstanding

	return stats, nil
}

// DeleteLoan deletes a paylater loan record
func (r *paylaterLoanRepository) DeleteLoan(ctx context.Context, id uint) error {
	db := database.GetDB(ctx, r.db)
	result := db.Delete(&entity.PaylaterLoan{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete loan", zap.Error(result.Error), zap.Uint("id", id))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("loan not found")
	}
	return nil
}
