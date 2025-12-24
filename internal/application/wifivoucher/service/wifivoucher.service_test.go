package service

import (
	"context"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockWifiVoucherRepository struct {
	mock.Mock
}

func (m *mockWifiVoucherRepository) Create(ctx context.Context, wv *entity.WifiVoucher) error {
	args := m.Called(ctx, wv)
	return args.Error(0)
}

func (m *mockWifiVoucherRepository) GetByID(ctx context.Context, id uint) (*entity.WifiVoucher, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WifiVoucher), args.Error(1)
}

func (m *mockWifiVoucherRepository) GetByCode(ctx context.Context, code string) (*entity.WifiVoucher, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WifiVoucher), args.Error(1)
}

func (m *mockWifiVoucherRepository) GetAll(ctx context.Context, filter *dto.WifiVoucherFilter) ([]entity.WifiVoucher, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.WifiVoucher), args.Get(1).(int64), args.Error(2)
}

func (m *mockWifiVoucherRepository) Update(ctx context.Context, wv *entity.WifiVoucher) error {
	args := m.Called(ctx, wv)
	return args.Error(0)
}

func (m *mockWifiVoucherRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockWifiVoucherRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *mockWifiVoucherRepository) GetByProviderAndDurationHours(ctx context.Context, providerID uint, durationHours int, filter *dto.WifiVoucherFilter) ([]entity.WifiVoucher, int64, error) {
	args := m.Called(ctx, providerID, durationHours, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.WifiVoucher), args.Get(1).(int64), args.Error(2)
}

func (m *mockWifiVoucherRepository) GetByProviderIDWithDurationHours(ctx context.Context, providerID uint, durationHours int) ([]entity.WifiVoucher, error) {
	args := m.Called(ctx, providerID, durationHours)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.WifiVoucher), args.Error(1)
}

func (m *mockWifiVoucherRepository) GetForUpdate(ctx context.Context, id uint) (*entity.WifiVoucher, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WifiVoucher), args.Error(1)
}

func TestWifiVoucherService_CreateWifiVoucher(t *testing.T) {
	t.Run("should create wifi voucher successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		req := testutil.CreateWifiVoucherRequestFixture()

		mockRepo.On("CodeExists", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), req.Code).Return(false, nil)

		mockRepo.On("Create", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(w *entity.WifiVoucher) bool {
			return w.Code == req.Code
		})).Run(func(args mock.Arguments) {
			w := args.Get(1).(*entity.WifiVoucher)
			w.ID = 1
		}).Return(nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.CreateWifiVoucher(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Code, result.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when code already exists", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		req := testutil.CreateWifiVoucherRequestFixture()

		mockRepo.On("CodeExists", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), req.Code).Return(true, nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.CreateWifiVoucher(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_GetWifiVoucherByID(t *testing.T) {
	t.Run("should get wifi voucher by ID successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetWifiVoucherByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, wifiVoucher.Code, result.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when wifi voucher not found", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetWifiVoucherByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_GetWifiVoucherByCode(t *testing.T) {
	t.Run("should get wifi voucher by code successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1

		mockRepo.On("GetByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), wifiVoucher.Code).Return(wifiVoucher, nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetWifiVoucherByCode(context.Background(), wifiVoucher.Code)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, wifiVoucher.Code, result.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when code not found", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)

		mockRepo.On("GetByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "NONEXISTENT").Return(nil, gorm.ErrRecordNotFound)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetWifiVoucherByCode(context.Background(), "NONEXISTENT")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_GetAllWifiVouchers(t *testing.T) {
	t.Run("should get all wifi vouchers with pagination", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVouchers := []entity.WifiVoucher{
			*testutil.CreateWifiVoucherFixture(),
			*testutil.CreateWifiVoucherFixture(),
		}
		wifiVouchers[0].ID = 1
		wifiVouchers[1].ID = 2
		wifiVouchers[1].Code = "WIFI002"

		filter := &dto.WifiVoucherFilter{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), filter).Return(wifiVouchers, int64(2), nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetAllWifiVouchers(context.Background(), filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.TotalCount)
		assert.Len(t, result.Data, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should filter wifi vouchers by status", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVouchers := []entity.WifiVoucher{*testutil.CreateWifiVoucherFixture()}
		wifiVouchers[0].ID = 1

		filter := &dto.WifiVoucherFilter{
			Status:   entity.StatusAvailable,
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), filter).Return(wifiVouchers, int64(1), nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.GetAllWifiVouchers(context.Background(), filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.GreaterOrEqual(t, result.TotalCount, int64(1))
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_UpdateWifiVoucher(t *testing.T) {
	t.Run("should update wifi voucher successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1

		req := &dto.UpdateWifiVoucherRequest{
			Code:     "WIFI_UPDATED",
			Password: "newpass",
		}

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		mockRepo.On("CodeExists", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), req.Code).Return(false, nil)

		mockRepo.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(w *entity.WifiVoucher) bool {
			return w.ID == 1 && w.Code == req.Code
		})).Return(nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.UpdateWifiVoucher(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when updating non-existent wifi voucher", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)

		req := &dto.UpdateWifiVoucherRequest{
			Code:     "WIFI_UPDATED",
			Password: "newpass",
		}

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.UpdateWifiVoucher(context.Background(), 999, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_DeleteWifiVoucher(t *testing.T) {
	t.Run("should delete wifi voucher successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		mockRepo.On("Delete", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(nil)

		service := NewWifiVoucherService(mockRepo, logger)
		err := service.DeleteWifiVoucher(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when deleting non-existent wifi voucher", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := NewWifiVoucherService(mockRepo, logger)
		err := service.DeleteWifiVoucher(context.Background(), 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_SellWifiVoucher(t *testing.T) {
	t.Run("should sell wifi voucher successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1
		wifiVoucher.Status = entity.StatusAvailable

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		mockRepo.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(w *entity.WifiVoucher) bool {
			return w.ID == 1 && w.Status == entity.StatusSold
		})).Return(nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.SellWifiVoucher(context.Background(), 1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when wifi voucher is not available", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1
		wifiVoucher.Status = entity.StatusSold

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.SellWifiVoucher(context.Background(), 1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWifiVoucherService_UseWifiVoucher(t *testing.T) {
	t.Run("should use wifi voucher successfully", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1
		wifiVoucher.Status = entity.StatusSold

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		mockRepo.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(w *entity.WifiVoucher) bool {
			return w.ID == 1 && w.Status == entity.StatusUsed
		})).Return(nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.UseWifiVoucher(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when wifi voucher is not sold", func(t *testing.T) {
		mockRepo := new(mockWifiVoucherRepository)
		logger := testutil.NewTestLogger(t)
		wifiVoucher := testutil.CreateWifiVoucherFixture()
		wifiVoucher.ID = 1
		wifiVoucher.Status = entity.StatusAvailable

		mockRepo.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(wifiVoucher, nil)

		service := NewWifiVoucherService(mockRepo, logger)
		result, err := service.UseWifiVoucher(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
