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

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupWalletHandler() (*WalletHandler, *testutil.MockWalletService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockWalletService{}
	jwtManager := testutil.NewMockJWTManager()
	logger := testutil.NewSilentLogger()
	handler := NewWalletHandler(mockService, jwtManager, logger)
	return handler, mockService
}

func TestWalletHandler_CreateWallet(t *testing.T) {
	t.Run("should create wallet successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		req := testutil.CreateWalletRequestFixture()
		response := testutil.CreateWalletFixture()

		mockService.On("CreateWallet", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateWalletRequest) bool {
			return r.UserID == req.UserID && r.PIN == req.PIN
		})).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateWallet(ctx)

		// Then
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, float64(1), data["user_id"])
		assert.Equal(t, 1000.00, data["balance"])
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, _ := setupWalletHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateWallet(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal server error for service failure", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		req := testutil.CreateWalletRequestFixture()
		mockService.On("CreateWallet", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateWalletRequest) bool {
			return r.UserID == req.UserID && r.PIN == req.PIN
		})).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateWallet(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_GetWalletByUserID(t *testing.T) {
	t.Run("should get wallet by user ID successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(1)
		response := &dto.GetWalletResponse{
			ID:        1,
			UserID:    userID,
			Balance:   1000.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetWalletByUserID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/1", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "1"},
		}

		// When
		handler.GetWalletByUserID(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, float64(1), data["user_id"])
	})

	t.Run("should return bad request for invalid user ID", func(t *testing.T) {
		// Setup
		handler, _ := setupWalletHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "invalid"},
		}

		// When
		handler.GetWalletByUserID(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when wallet not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(999)
		mockService.On("GetWalletByUserID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID).Return(nil, errors.New("wallet not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/999", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "999"},
		}

		// When
		handler.GetWalletByUserID(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_GetWalletByID(t *testing.T) {
	t.Run("should get wallet by ID successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		walletID := uint(1)
		response := &dto.GetWalletResponse{
			ID:        walletID,
			UserID:    1,
			Balance:   1000.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetWalletByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), walletID).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/wallet/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.GetWalletByID(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
	})

	t.Run("should return bad request for invalid wallet ID", func(t *testing.T) {
		// Setup
		handler, _ := setupWalletHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/wallet/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.GetWalletByID(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when wallet not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		walletID := uint(999)
		mockService.On("GetWalletByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), walletID).Return(nil, errors.New("wallet not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/wallets/wallet/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.GetWalletByID(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_AddBalance(t *testing.T) {
	t.Run("should add balance successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(1)
		req := &dto.UpdateBalanceRequest{
			Amount:      100.00,
			Description: "test deposit",
		}
		response := &dto.GetWalletResponse{
			ID:        1,
			UserID:    userID,
			Balance:   1100.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("AddBalance", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID, req.Amount, req.Description).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets/1/add-balance", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "1"},
		}

		// When
		handler.AddBalance(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, 1100.00, data["balance"])
	})

	t.Run("should return not found when wallet not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(999)
		req := &dto.UpdateBalanceRequest{
			Amount:      100.00,
			Description: "test deposit",
		}

		mockService.On("AddBalance", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID, req.Amount, req.Description).Return(nil, errors.New("wallet not found"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets/999/add-balance", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "999"},
		}

		// When
		handler.AddBalance(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_DeductBalance(t *testing.T) {
	t.Run("should deduct balance successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(1)
		req := &dto.UpdateBalanceRequest{
			Amount:      100.00,
			Description: "test withdrawal",
		}
		response := &dto.GetWalletResponse{
			ID:        1,
			UserID:    userID,
			Balance:   900.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("DeductBalance", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID, req.Amount, req.Description).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets/1/deduct-balance", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "1"},
		}

		// When
		handler.DeductBalance(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, 900.00, data["balance"])
	})

	t.Run("should return payment required for insufficient balance", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(1)
		req := &dto.UpdateBalanceRequest{
			Amount:      2000.00,
			Description: "test withdrawal",
		}

		mockService.On("DeductBalance", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID, req.Amount, req.Description).Return(nil, errors.New("insufficient balance"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/wallets/1/deduct-balance", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "1"},
		}

		// When
		handler.DeductBalance(ctx)

		// Then
		assert.Equal(t, http.StatusPaymentRequired, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_DeleteWallet(t *testing.T) {
	t.Run("should delete wallet successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(1)
		mockService.On("DeleteWallet", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/wallets/1", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "1"},
		}

		// When
		handler.DeleteWallet(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "message")
		assert.Equal(t, "Wallet deleted successfully", result["message"])
	})

	t.Run("should return bad request for invalid user ID", func(t *testing.T) {
		// Setup
		handler, _ := setupWalletHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/wallets/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "invalid"},
		}

		// When
		handler.DeleteWallet(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when wallet not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupWalletHandler()

		userID := uint(999)
		mockService.On("DeleteWallet", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), userID).Return(errors.New("wallet not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/wallets/999", nil)
		ctx.Params = gin.Params{
			{Key: "user_id", Value: "999"},
		}

		// When
		handler.DeleteWallet(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
