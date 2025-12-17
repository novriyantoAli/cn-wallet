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
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProviderService struct {
	mock.Mock
}

func (m *MockProviderService) CreateProvider(ctx context.Context, req *dto.CreateProviderRequest) (*dto.ProviderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetProviderByID(ctx context.Context, id uint) (*dto.ProviderResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetProviderByCode(ctx context.Context, code string) (*dto.ProviderResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) GetAllProviders(ctx context.Context, filter *dto.ProviderFilter) (*dto.ProviderListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProviderListResponse), args.Error(1)
}

func (m *MockProviderService) UpdateProvider(ctx context.Context, id uint, req *dto.UpdateProviderRequest) (*dto.ProviderResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProviderResponse), args.Error(1)
}

func (m *MockProviderService) DeleteProvider(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupProviderHandler() (*ProviderHandler, *MockProviderService) {
	mockService := new(MockProviderService)
	logger := testutil.NewSilentLogger()
	return NewProviderHandler(mockService, logger), mockService
}

func TestProviderHandler_CreateProvider(t *testing.T) {
	t.Run("should create provider successfully", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		req := &dto.CreateProviderRequest{
			Name: "Telkomsel",
			Code: "TSEL",
			Logo: "https://example.com/logo.png",
		}

		response := &dto.ProviderResponse{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("CreateProvider", mock.Anything, mock.MatchedBy(func(r *dto.CreateProviderRequest) bool {
			return r.Name == req.Name && r.Code == req.Code
		})).Return(response, nil)

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/providers", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.CreateProvider(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 when code already exists", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		req := &dto.CreateProviderRequest{
			Name: "Telkomsel",
			Code: "TSEL",
		}

		mockService.On("CreateProvider", mock.Anything, mock.Anything).Return(nil, errors.New("provider with code already exists"))

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/providers", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.CreateProvider(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProviderHandler_GetProviderByID(t *testing.T) {
	t.Run("should get provider by id successfully", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		response := &dto.ProviderResponse{
			ID:        1,
			Name:      "Telkomsel",
			Code:      "TSEL",
			Logo:      "https://example.com/logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetProviderByID", mock.Anything, uint(1)).Return(response, nil)

		httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/providers/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetProviderByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when provider not found", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		mockService.On("GetProviderByID", mock.Anything, uint(9999)).Return(nil, errors.New("provider not found"))

		httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/providers/9999", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "9999"}}

		handler.GetProviderByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProviderHandler_GetAllProviders(t *testing.T) {
	t.Run("should get all providers", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		response := &dto.ProviderListResponse{
			Data: []dto.ProviderResponse{
				{
					ID:        1,
					Name:      "Provider 1",
					Code:      "PRV1",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			TotalCount: 1,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllProviders", mock.Anything, mock.MatchedBy(func(f *dto.ProviderFilter) bool {
			return f.Page == 1 && f.PageSize == 10
		})).Return(response, nil)

		httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/providers?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.GetAllProviders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProviderHandler_UpdateProvider(t *testing.T) {
	t.Run("should update provider successfully", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		req := &dto.UpdateProviderRequest{
			Name: "Telkomsel Updated",
			Code: "TSEL",
		}

		response := &dto.ProviderResponse{
			ID:        1,
			Name:      "Telkomsel Updated",
			Code:      "TSEL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("UpdateProvider", mock.Anything, uint(1), mock.Anything).Return(response, nil)

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/providers/1", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateProvider(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when provider not found", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		req := &dto.UpdateProviderRequest{
			Name: "Updated",
		}

		mockService.On("UpdateProvider", mock.Anything, uint(999), mock.Anything).Return(nil, errors.New("provider not found"))

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/providers/999", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.UpdateProvider(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProviderHandler_DeleteProvider(t *testing.T) {
	t.Run("should delete provider successfully", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		mockService.On("DeleteProvider", mock.Anything, uint(1)).Return(nil)

		httpReq := httptest.NewRequest(http.MethodDelete, "/api/v1/providers/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteProvider(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when provider not found", func(t *testing.T) {
		handler, mockService := setupProviderHandler()
		gin.SetMode(gin.TestMode)

		mockService.On("DeleteProvider", mock.Anything, uint(999)).Return(errors.New("provider not found"))

		httpReq := httptest.NewRequest(http.MethodDelete, "/api/v1/providers/999", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.DeleteProvider(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
