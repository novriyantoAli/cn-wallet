package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWifiVoucherService struct {
	mock.Mock
}

func (m *mockWifiVoucherService) CreateWifiVoucher(ctx context.Context, req *dto.CreateWifiVoucherRequest) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func (m *mockWifiVoucherService) GetWifiVoucherByID(ctx context.Context, id uint) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func (m *mockWifiVoucherService) GetWifiVoucherByCode(ctx context.Context, code string) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func (m *mockWifiVoucherService) GetAllWifiVouchers(ctx context.Context, filter *dto.WifiVoucherFilter) (*dto.WifiVoucherListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherListResponse), args.Error(1)
}

func (m *mockWifiVoucherService) UpdateWifiVoucher(ctx context.Context, id uint, req *dto.UpdateWifiVoucherRequest) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func (m *mockWifiVoucherService) DeleteWifiVoucher(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockWifiVoucherService) SellWifiVoucher(ctx context.Context, voucherID uint, userID uint) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, voucherID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func (m *mockWifiVoucherService) UseWifiVoucher(ctx context.Context, voucherID uint) (*dto.WifiVoucherResponse, error) {
	args := m.Called(ctx, voucherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WifiVoucherResponse), args.Error(1)
}

func setupWifiVoucherHandler() (*WifiVoucherHandler, *mockWifiVoucherService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockWifiVoucherService)
	logger := testutil.NewSilentLogger()
	handler := NewWifiVoucherHandler(mockService, logger)
	return handler, mockService
}

func TestWifiVoucherHandler_CreateWifiVoucher(t *testing.T) {
	t.Run("should create wifi voucher successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()
		req := testutil.CreateWifiVoucherRequestFixture()

		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            req.Code,
			Password:        req.Password,
			DurationHours: req.DurationHours,
			BatchID:         req.BatchID,
			Status:          entity.StatusAvailable,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockService.On("CreateWifiVoucher", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateWifiVoucherRequest) bool {
			return r.Code == req.Code
		})).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wifi-vouchers", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateWifiVoucher(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wifi-vouchers", bytes.NewBuffer([]byte("invalid")))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateWifiVoucher(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_GetWifiVoucherByID(t *testing.T) {
	t.Run("should get wifi voucher by ID successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            "WIFI001",
			Password:        "pass123",
			DurationHours: 24,
			BatchID:         "BATCH001",
			Status:          entity.StatusAvailable,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockService.On("GetWifiVoucherByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.GetWifiVoucherByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.GetWifiVoucherByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_GetWifiVoucherByCode(t *testing.T) {
	t.Run("should get wifi voucher by code successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            "WIFI001",
			Password:        "pass123",
			DurationHours: 24,
			BatchID:         "BATCH001",
			Status:          entity.StatusAvailable,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockService.On("GetWifiVoucherByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "WIFI001").Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers/code/WIFI001", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: "WIFI001"},
		}

		handler.GetWifiVoucherByCode(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request when code is empty", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers/code/", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: ""},
		}

		handler.GetWifiVoucherByCode(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_GetAllWifiVouchers(t *testing.T) {
	t.Run("should get all wifi vouchers successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		response := &dto.WifiVoucherListResponse{
			Data: []dto.WifiVoucherResponse{
				{ID: 1, Code: "WIFI001", Status: entity.StatusAvailable, DurationHours: 24},
				{ID: 2, Code: "WIFI002", Status: entity.StatusAvailable, DurationHours: 24},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllWifiVouchers", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.WifiVoucherFilter) bool {
			return true
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers", nil)

		handler.GetAllWifiVouchers(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var result dto.WifiVoucherListResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result.Data, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		mockService.On("GetAllWifiVouchers", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.WifiVoucherFilter) bool {
			return true
		})).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wifi-vouchers", nil)

		handler.GetAllWifiVouchers(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_UpdateWifiVoucher(t *testing.T) {
	t.Run("should update wifi voucher successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		req := testutil.CreateUpdateWifiVoucherRequestFixture()
		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            req.Code,
			Password:        req.Password,
			DurationHours: req.DurationHours,
			BatchID:         req.BatchID,
			Status:          req.Status,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockService.On("UpdateWifiVoucher", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1), mock.MatchedBy(func(r *dto.UpdateWifiVoucherRequest) bool {
			return r.Code == req.Code
		})).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/wifi-vouchers/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.UpdateWifiVoucher(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/wifi-vouchers/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.UpdateWifiVoucher(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_DeleteWifiVoucher(t *testing.T) {
	t.Run("should delete wifi voucher successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		mockService.On("DeleteWifiVoucher", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/wifi-vouchers/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.DeleteWifiVoucher(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/wifi-vouchers/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.DeleteWifiVoucher(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_SellWifiVoucher(t *testing.T) {
	t.Run("should sell wifi voucher successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            "WIFI001",
			Status:          entity.StatusSold,
			DurationHours: 24,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		sellReq := struct {
			UserID uint `json:"user_id"`
		}{UserID: 2}

		mockService.On("SellWifiVoucher", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1), uint(2)).Return(response, nil)

		reqBody, _ := json.Marshal(sellReq)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wifi-vouchers/1/sell", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.SellWifiVoucher(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_UseWifiVoucher(t *testing.T) {
	t.Run("should use wifi voucher successfully", func(t *testing.T) {
		handler, mockService := setupWifiVoucherHandler()

		response := &dto.WifiVoucherResponse{
			ID:              1,
			Code:            "WIFI001",
			Status:          entity.StatusUsed,
			DurationHours: 24,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockService.On("UseWifiVoucher", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), uint(1)).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wifi-vouchers/1/use", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.UseWifiVoucher(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWifiVoucherHandler_RegisterRoutes(t *testing.T) {
	t.Run("should register all routes correctly", func(t *testing.T) {
		handler, _ := setupWifiVoucherHandler()
		router := gin.New()
		api := router.Group("/api/v1")

		handler.RegisterRoutes(api)

		routes := router.Routes()
		assert.NotNil(t, routes)
		assert.Greater(t, len(routes), 0)
	})
}
