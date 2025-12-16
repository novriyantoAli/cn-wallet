package repository

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id uint) (*entity.Payment, error)
	GetAll(ctx context.Context, filter *dto.PaymentFilter) ([]entity.Payment, int64, error)
	Update(ctx context.Context, payment *entity.Payment) error
	Delete(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]entity.Payment, error)
}

type paymentRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPaymentRepository(db *gorm.DB, logger *zap.Logger) PaymentRepository {
	return &paymentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	db := database.GetDB(ctx, r.db)
	r.logger.Info("Creating payment", zap.Uint("user_id", payment.UserID))
	return db.Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uint) (*entity.Payment, error) {
	db := database.GetDB(ctx, r.db)
	var payment entity.Payment
	err := db.First(&payment, id).Error
	if err != nil {
		r.logger.Error("Failed to get payment by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetAll(ctx context.Context, filter *dto.PaymentFilter) ([]entity.Payment, int64, error) {
	db := database.GetDB(ctx, r.db)
	var payments []entity.Payment
	var totalCount int64

	query := db.Model(&entity.Payment{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Currency != "" {
		query = query.Where("currency = ?", filter.Currency)
	}
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	query.Count(&totalCount)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Find(&payments).Error
	if err != nil {
		r.logger.Error("Failed to get payments", zap.Error(err))
		return nil, 0, err
	}

	return payments, totalCount, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	db := database.GetDB(ctx, r.db)
	r.logger.Info("Updating payment", zap.Uint("id", payment.ID))
	return db.Save(payment).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id uint) error {
	db := database.GetDB(ctx, r.db)
	r.logger.Info("Deleting payment", zap.Uint("id", id))
	return db.Delete(&entity.Payment{}, id).Error
}

func (r *paymentRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.Payment, error) {
	db := database.GetDB(ctx, r.db)
	var payments []entity.Payment
	err := db.Where("user_id = ?", userID).Find(&payments).Error
	if err != nil {
		r.logger.Error("Failed to get payments by user ID", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}
	return payments, nil
}
