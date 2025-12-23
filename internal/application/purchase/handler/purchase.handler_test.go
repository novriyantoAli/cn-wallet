package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/dto"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPurchaseHandler() (*PurchaseHandler, *testutil.MockPurchaseService, *testutil.MockUserSecurityService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockPurchaseService{}
	mockUserSecurityService := &testutil.MockUserSecurityService{}
	jwtManager := &jwt.JWTManager{}
	logger := testutil.NewSilentLogger()
	handler := NewPurchaseHandler(mockService, mockUserSecurityService, jwtManager, logger)
	return handler, mockService, mockUserSecurityService
}

func TestPurchaseHandler_ProcessPurchase(t *testing.T) {
	t.Run("should process purchase successfully", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		txID := uuid.New()
		response := &dto.PurchaseResponse{
			TransactionID: txID,
			Status:        transactionEntity.StatusSuccess,
			SerialNumber:  "SN123456",
			Message:       "Purchase successful",
		}

		mockService.On("ProcessPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
	})

	t.Run("should return unauthorized when authorization header is missing", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		// No Authorization header

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Authorization header is required", result["error"])
	})

	t.Run("should return unauthorized when authorization header format is invalid", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "InvalidFormat token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Invalid authorization header format", result["error"])
	})

	t.Run("should return unauthorized when token is empty", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer ")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Invalid authorization header", result["error"])
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when user not found", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		mockService.On("ProcessPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("user not found"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "user not found", result["error"])
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 999,
			Phone:     "08123456789",
		}

		mockService.On("ProcessPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("product not found"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "product not found", result["error"])
	})

	t.Run("should return bad request when insufficient balance", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseRequest{
			ProductID: 1,
			Phone:     "08123456789",
		}

		mockService.On("ProcessPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("insufficient wallet balance"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "insufficient wallet balance", result["error"])
	})
}

func TestPurchaseHandler_GetPurchaseHistory(t *testing.T) {
	t.Run("should get purchase history successfully", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		filter := &dto.PurchaseHistoryFilter{
			WalletID: 1,
			Page:     1,
			PageSize: 10,
		}

		response := &dto.PurchaseHistoryList{
			Data:      []dto.PurchaseHistory{},
			Total:     0,
			Page:      1,
			PageSize:  10,
			TotalPage: 0,
		}

		mockService.On("GetPurchaseHistory", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), filter).Return(response, nil)

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=1", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.NotNil(t, result)
	})

	t.Run("should return bad request when wallet_id is missing", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Invalid wallet_id", result["error"])
	})

	t.Run("should return bad request when wallet_id is invalid", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=invalid", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Invalid wallet_id", result["error"])
	})

	t.Run("should return bad request when wallet_id is zero", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=0", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Invalid wallet_id", result["error"])
	})

	t.Run("should use default pagination values", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		expectedFilter := &dto.PurchaseHistoryFilter{
			WalletID: 1,
			Page:     1,
			PageSize: 10,
		}

		response := &dto.PurchaseHistoryList{
			Data:      []dto.PurchaseHistory{},
			Total:     0,
			Page:      1,
			PageSize:  10,
			TotalPage: 0,
		}

		mockService.On("GetPurchaseHistory", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), expectedFilter).Return(response, nil)

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=1", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "GetPurchaseHistory", mock.Anything, expectedFilter)
	})

	t.Run("should accept custom pagination values", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		expectedFilter := &dto.PurchaseHistoryFilter{
			WalletID: 1,
			Page:     2,
			PageSize: 20,
		}

		response := &dto.PurchaseHistoryList{
			Data:      []dto.PurchaseHistory{},
			Total:     0,
			Page:      2,
			PageSize:  20,
			TotalPage: 0,
		}

		mockService.On("GetPurchaseHistory", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), expectedFilter).Return(response, nil)

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=1&page=2&page_size=20", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "GetPurchaseHistory", mock.Anything, expectedFilter)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		filter := &dto.PurchaseHistoryFilter{
			WalletID: 1,
			Page:     1,
			PageSize: 10,
		}

		mockService.On("GetPurchaseHistory", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), filter).Return(nil, errors.New("database error"))

		// Prepare request
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/purchases/history?wallet_id=1", nil)

		// When
		handler.GetPurchaseHistory(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Failed to get purchase history", result["error"])
	})
}
func TestPurchaseHandler_ProcessWifiPurchase(t *testing.T) {
	t.Run("should process wifi purchase successfully", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		txID := uuid.New()
		response := &dto.PurchaseWifiResponse{
			TransactionID:   txID,
			VoucherID:       1,
			VoucherCode:     "WIFI001",
			VoucherPassword: "pass123",
			Status:          transactionEntity.StatusSuccess,
			Message:         "WiFi voucher purchase successful",
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.NotNil(t, result["data"])
	})

	t.Run("should return bad request when product_id is missing", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 0,
		}

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return unauthorized when authorization header missing", func(t *testing.T) {
		// Setup
		handler, _, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Authorization header is required", result["error"])
	})

	t.Run("should return unauthorized when token is invalid", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "invalid-token", req).Return(nil, errors.New("invalid or expired token"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "invalid or expired token", result["error"])
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 999,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("product not found"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "product not found", result["error"])
	})

	t.Run("should return bad request when product is not active", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("product is not active"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "product is not active", result["error"])
	})

	t.Run("should return bad request when insufficient balance", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("insufficient wallet balance"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "insufficient wallet balance", result["error"])
	})

	t.Run("should return bad request when no available vouchers", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("no available wifi vouchers in stock"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "no available wifi vouchers in stock", result["error"])
	})

	t.Run("should return internal server error when service fails with unknown error", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(nil, errors.New("unexpected database error"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "Failed to process WiFi purchase", result["error"])
	})

	t.Run("should return bad request when wifi purchase fails", func(t *testing.T) {
		// Setup
		handler, mockService, _ := setupPurchaseHandler()

		req := &dto.PurchaseWifiRequest{
			ProductID: 1,
		}

		txID := uuid.New()
		response := &dto.PurchaseWifiResponse{
			TransactionID: txID,
			Status:        transactionEntity.StatusFailed,
			Message:       "WiFi purchase failed",
		}

		mockService.On("ProcessWifiPurchase", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "valid-token", req).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/purchases/wifi", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.Header.Set("Authorization", "Bearer valid-token")

		// When
		handler.ProcessWifiPurchase(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.NotNil(t, result["data"])
	})
}
