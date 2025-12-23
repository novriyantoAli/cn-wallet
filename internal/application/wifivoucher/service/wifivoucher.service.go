package service

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WifiVoucherService defines the interface for wifi voucher business logic.
type WifiVoucherService interface {
	CreateWifiVoucher(ctx context.Context, req *dto.CreateWifiVoucherRequest) (*dto.WifiVoucherResponse, error)
	GetWifiVoucherByID(ctx context.Context, id uint) (*dto.WifiVoucherResponse, error)
	GetWifiVoucherByCode(ctx context.Context, code string) (*dto.WifiVoucherResponse, error)
	GetAllWifiVouchers(ctx context.Context, filter *dto.WifiVoucherFilter) (*dto.WifiVoucherListResponse, error)
	UpdateWifiVoucher(ctx context.Context, id uint, req *dto.UpdateWifiVoucherRequest) (*dto.WifiVoucherResponse, error)
	DeleteWifiVoucher(ctx context.Context, id uint) error
	SellWifiVoucher(ctx context.Context, voucherID uint, userID uint) (*dto.WifiVoucherResponse, error)
	UseWifiVoucher(ctx context.Context, voucherID uint) (*dto.WifiVoucherResponse, error)
}

type wifiVoucherService struct {
	repo   repository.WifiVoucherRepository
	logger *zap.Logger
}

// NewWifiVoucherService creates a new instance of WifiVoucherService.
func NewWifiVoucherService(repo repository.WifiVoucherRepository, logger *zap.Logger) WifiVoucherService {
	return &wifiVoucherService{
		repo:   repo,
		logger: logger,
	}
}

// CreateWifiVoucher creates a new wifi voucher.
func (s *wifiVoucherService) CreateWifiVoucher(ctx context.Context, req *dto.CreateWifiVoucherRequest) (*dto.WifiVoucherResponse, error) {
	// Check if code already exists
	exists, err := s.repo.CodeExists(ctx, req.Code)
	if err != nil {
		s.logger.Error("Failed to check code existence", zap.String("code", req.Code), zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("voucher code already exists")
	}

	wifiVoucher := &entity.WifiVoucher{
		Code:          req.Code,
		Password:      req.Password,
		DurationHours: req.DurationHours,
		BatchID:       req.BatchID,
		ProviderID:    req.ProviderID,
		Status:        entity.StatusAvailable,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.repo.Create(ctx, wifiVoucher)
	if err != nil {
		s.logger.Error("Failed to create wifi voucher", zap.String("code", req.Code), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(wifiVoucher), nil
}

// GetWifiVoucherByID retrieves a wifi voucher by ID.
func (s *wifiVoucherService) GetWifiVoucherByID(ctx context.Context, id uint) (*dto.WifiVoucherResponse, error) {
	wifiVoucher, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	if wifiVoucher == nil {
		return nil, errors.New("wifi voucher not found")
	}

	return s.entityToResponse(wifiVoucher), nil
}

// GetWifiVoucherByCode retrieves a wifi voucher by code.
func (s *wifiVoucherService) GetWifiVoucherByCode(ctx context.Context, code string) (*dto.WifiVoucherResponse, error) {
	wifiVoucher, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}
	if wifiVoucher == nil {
		return nil, errors.New("wifi voucher not found")
	}

	return s.entityToResponse(wifiVoucher), nil
}

// GetAllWifiVouchers retrieves all wifi vouchers with filters and pagination.
func (s *wifiVoucherService) GetAllWifiVouchers(ctx context.Context, filter *dto.WifiVoucherFilter) (*dto.WifiVoucherListResponse, error) {
	wifiVouchers, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get wifi vouchers", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.WifiVoucherResponse, 0)
	for _, wv := range wifiVouchers {
		responses = append(responses, *s.entityToResponse(&wv))
	}

	return &dto.WifiVoucherListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

// UpdateWifiVoucher updates an existing wifi voucher.
func (s *wifiVoucherService) UpdateWifiVoucher(ctx context.Context, id uint, req *dto.UpdateWifiVoucherRequest) (*dto.WifiVoucherResponse, error) {
	wifiVoucher, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	if wifiVoucher == nil {
		return nil, errors.New("wifi voucher not found")
	}

	// Check if new code already exists (if code is being changed)
	if req.Code != "" && req.Code != wifiVoucher.Code {
		exists, err := s.repo.CodeExists(ctx, req.Code)
		if err != nil {
			s.logger.Error("Failed to check code existence", zap.String("code", req.Code), zap.Error(err))
			return nil, err
		}
		if exists {
			return nil, errors.New("voucher code already exists")
		}
		wifiVoucher.Code = req.Code
	}

	if req.Password != "" {
		wifiVoucher.Password = req.Password
	}
	if req.DurationHours > 0 {
		wifiVoucher.DurationHours = req.DurationHours
	}
	if req.BatchID != "" {
		wifiVoucher.BatchID = req.BatchID
	}
	if req.ProviderID > 0 {
		wifiVoucher.ProviderID = req.ProviderID
	}
	if req.Status != "" {
		wifiVoucher.Status = req.Status
	}
	if req.SoldToUserID != nil {
		wifiVoucher.SoldToUserID = req.SoldToUserID
	}
	if req.SoldAt != nil {
		wifiVoucher.SoldAt = req.SoldAt
	}
	if req.UsedAt != nil {
		wifiVoucher.UsedAt = req.UsedAt
	}

	wifiVoucher.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, wifiVoucher)
	if err != nil {
		s.logger.Error("Failed to update wifi voucher", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(wifiVoucher), nil
}

// DeleteWifiVoucher deletes a wifi voucher.
func (s *wifiVoucherService) DeleteWifiVoucher(ctx context.Context, id uint) error {
	wifiVoucher, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher", zap.Uint("id", id), zap.Error(err))
		return err
	}
	if wifiVoucher == nil {
		return errors.New("wifi voucher not found")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete wifi voucher", zap.Uint("id", id), zap.Error(err))
		return err
	}

	return nil
}

// SellWifiVoucher marks a wifi voucher as sold to a user.
func (s *wifiVoucherService) SellWifiVoucher(ctx context.Context, voucherID uint, userID uint) (*dto.WifiVoucherResponse, error) {
	wifiVoucher, err := s.repo.GetByID(ctx, voucherID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher", zap.Uint("id", voucherID), zap.Error(err))
		return nil, err
	}
	if wifiVoucher == nil {
		return nil, errors.New("wifi voucher not found")
	}

	if wifiVoucher.Status != entity.StatusAvailable {
		s.logger.Warn("Cannot sell wifi voucher with status", zap.String("status", wifiVoucher.Status))
		return nil, errors.New("wifi voucher cannot be sold, current status: " + wifiVoucher.Status)
	}

	now := time.Now()
	wifiVoucher.Status = entity.StatusSold
	wifiVoucher.SoldToUserID = &userID
	wifiVoucher.SoldAt = &now
	wifiVoucher.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, wifiVoucher)
	if err != nil {
		s.logger.Error("Failed to update wifi voucher", zap.Uint("id", voucherID), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(wifiVoucher), nil
}

// UseWifiVoucher marks a wifi voucher as used.
func (s *wifiVoucherService) UseWifiVoucher(ctx context.Context, voucherID uint) (*dto.WifiVoucherResponse, error) {
	wifiVoucher, err := s.repo.GetByID(ctx, voucherID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wifi voucher not found")
		}
		s.logger.Error("Failed to get wifi voucher", zap.Uint("id", voucherID), zap.Error(err))
		return nil, err
	}
	if wifiVoucher == nil {
		return nil, errors.New("wifi voucher not found")
	}

	if wifiVoucher.Status != entity.StatusSold {
		s.logger.Warn("Cannot use wifi voucher with status", zap.String("status", wifiVoucher.Status))
		return nil, errors.New("wifi voucher cannot be used, current status: " + wifiVoucher.Status)
	}

	now := time.Now()
	wifiVoucher.Status = entity.StatusUsed
	wifiVoucher.UsedAt = &now
	wifiVoucher.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, wifiVoucher)
	if err != nil {
		s.logger.Error("Failed to update wifi voucher", zap.Uint("id", voucherID), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(wifiVoucher), nil
}

// entityToResponse converts a WifiVoucher entity to response DTO.
func (s *wifiVoucherService) entityToResponse(wv *entity.WifiVoucher) *dto.WifiVoucherResponse {
	return &dto.WifiVoucherResponse{
		ID:            wv.ID,
		Code:          wv.Code,
		Password:      wv.Password,
		DurationHours: wv.DurationHours,
		BatchID:       wv.BatchID,
		ProviderID:    wv.ProviderID,
		Status:        wv.Status,
		SoldToUserID:  wv.SoldToUserID,
		SoldAt:        wv.SoldAt,
		UsedAt:        wv.UsedAt,
		CreatedAt:     wv.CreatedAt,
		UpdatedAt:     wv.UpdatedAt,
	}
}
