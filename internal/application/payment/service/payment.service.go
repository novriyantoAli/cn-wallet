package service

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/payment/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/user/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetPaymentByID(ctx context.Context, id uint) (*dto.PaymentResponse, error)
	GetPayments(ctx context.Context, filter *dto.PaymentFilter) (*dto.PaymentListResponse, error)
	UpdatePayment(ctx context.Context, id uint, req *dto.UpdatePaymentRequest) (*dto.PaymentResponse, error)
	DeletePayment(ctx context.Context, id uint) error
	GetPaymentsByUser(ctx context.Context, userID uint) ([]dto.PaymentResponse, error)
}

type paymentService struct {
	repo        repository.PaymentRepository
	userService service.UserService
	logger      *zap.Logger
}

func NewPaymentService(
	repo repository.PaymentRepository,
	userService service.UserService,
	logger *zap.Logger,
) PaymentService {
	return &paymentService{
		repo:        repo,
		userService: userService,
		logger:      logger,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Validate that user exists before creating payment
	_, err := s.userService.GetUserByID(ctx, req.UserID)
	if err != nil {
		s.logger.Error("User not found for payment creation", zap.Uint("user_id", req.UserID), zap.Error(err))
		return nil, errors.New("user not found")
	}

	payment := &entity.Payment{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      entity.PaymentStatusPending,
		Description: req.Description,
		UserID:      req.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(payment), nil
}

func (s *paymentService) GetPaymentByID(ctx context.Context, id uint) (*dto.PaymentResponse, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return s.entityToResponse(payment), nil
}

func (s *paymentService) GetPayments(ctx context.Context, filter *dto.PaymentFilter) (*dto.PaymentListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	payments, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, *s.entityToResponse(&payment))
	}

	return &dto.PaymentListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (s *paymentService) UpdatePayment(ctx context.Context, id uint, req *dto.UpdatePaymentRequest) (*dto.PaymentResponse, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	status := entity.PaymentStatus(req.Status)
	if !status.IsValid() {
		return nil, errors.New("invalid payment status")
	}

	payment.Status = status
	if req.Description != "" {
		payment.Description = req.Description
	}
	payment.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to update payment", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(payment), nil
}

func (s *paymentService) DeletePayment(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *paymentService) GetPaymentsByUser(ctx context.Context, userID uint) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, *s.entityToResponse(&payment))
	}

	return responses, nil
}

func (s *paymentService) entityToResponse(payment *entity.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:          payment.ID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		Status:      payment.Status.String(),
		Description: payment.Description,
		UserID:      payment.UserID,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}
}
