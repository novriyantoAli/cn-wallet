package repository

import (
	"context"
	"errors"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PaylaterRepaymentRepository defines the interface for paylater repayment data access
type PaylaterRepaymentRepository interface {
	Create(ctx context.Context, repayment *entity.PaylaterRepayment) error
	GetByID(ctx context.Context, id uint) (*entity.PaylaterRepayment, error)
	Update(ctx context.Context, repayment *entity.PaylaterRepayment) error
	Delete(ctx context.Context, id uint) error
	GetByLoanID(ctx context.Context, loanID uint) ([]entity.PaylaterRepayment, error)
	GetByUserID(ctx context.Context, userID uint) ([]entity.PaylaterRepayment, error)
}

type paylaterRepaymentRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewPaylaterRepaymentRepository creates a new paylater repayment repository
func NewPaylaterRepaymentRepository(db *gorm.DB, logger *zap.Logger) PaylaterRepaymentRepository {
	return &paylaterRepaymentRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new paylater repayment record
func (r *paylaterRepaymentRepository) Create(ctx context.Context, repayment *entity.PaylaterRepayment) error {
	r.logger.Info("Creating paylater repayment", zap.Uint("loan_id", repayment.PaylaterLoanID), zap.Uint("user_id", repayment.UserID), zap.Int64("amount", repayment.Amount))
	db := database.GetDB(ctx, r.db)
	if err := db.Create(repayment).Error; err != nil {
		r.logger.Error("Failed to create paylater repayment", zap.Error(err), zap.Uint("loan_id", repayment.PaylaterLoanID))
		return err
	}
	return nil
}

// GetByID retrieves a repayment record by its ID
func (r *paylaterRepaymentRepository) GetByID(ctx context.Context, id uint) (*entity.PaylaterRepayment, error) {
	db := database.GetDB(ctx, r.db)
	var repayment entity.PaylaterRepayment
	if err := db.Where("id = ?", id).First(&repayment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("Paylater repayment not found", zap.Uint("id", id))
			return nil, errors.New("repayment not found")
		}
		r.logger.Error("Failed to get paylater repayment", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	return &repayment, nil
}

// Update updates an existing paylater repayment record
func (r *paylaterRepaymentRepository) Update(ctx context.Context, repayment *entity.PaylaterRepayment) error {
	r.logger.Info("Updating paylater repayment", zap.Uint("id", repayment.ID), zap.Int64("amount", repayment.Amount))
	db := database.GetDB(ctx, r.db)
	if err := db.Save(repayment).Error; err != nil {
		r.logger.Error("Failed to update paylater repayment", zap.Error(err), zap.Uint("id", repayment.ID))
		return err
	}
	return nil
}

// Delete deletes a repayment record by its ID
func (r *paylaterRepaymentRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting paylater repayment", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db)
	if err := db.Delete(&entity.PaylaterRepayment{}, id).Error; err != nil {
		r.logger.Error("Failed to delete paylater repayment", zap.Error(err), zap.Uint("id", id))
		return err
	}
	return nil
}

// GetByLoanID retrieves all repayments for a specific loan
func (r *paylaterRepaymentRepository) GetByLoanID(ctx context.Context, loanID uint) ([]entity.PaylaterRepayment, error) {
	r.logger.Info("Getting repayments by loan", zap.Uint("loan_id", loanID))
	db := database.GetDB(ctx, r.db)
	var repayments []entity.PaylaterRepayment
	if err := db.Where("paylater_loan_id = ?", loanID).Order("created_at DESC").Find(&repayments).Error; err != nil {
		r.logger.Error("Failed to get repayments by loan", zap.Error(err), zap.Uint("loan_id", loanID))
		return nil, err
	}
	return repayments, nil
}

// GetByUserID retrieves all repayments for a specific user
func (r *paylaterRepaymentRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.PaylaterRepayment, error) {
	r.logger.Info("Getting repayments by user", zap.Uint("user_id", userID))
	db := database.GetDB(ctx, r.db)
	var repayments []entity.PaylaterRepayment
	if err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&repayments).Error; err != nil {
		r.logger.Error("Failed to get repayments by user", zap.Error(err), zap.Uint("user_id", userID))
		return nil, err
	}
	return repayments, nil
}
