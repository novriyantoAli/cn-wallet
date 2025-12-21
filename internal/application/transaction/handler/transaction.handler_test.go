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
	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTransactionHandler() (*TransactionHandler, *testutil.MockTransactionService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockTransactionService{}
	handler := NewTransactionHandler(mockService)
	return handler, mockService
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should create transaction successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		req := testutil.CreateTransactionRequestFixture()

		response := &dto.TransactionResponse{
			ID:                 uuid.New(),
			WalletID:           req.WalletID,
			Type:               req.Type,
			Amount:             req.Amount,
			Status:             entity.StatusPending,
			Description:        req.Description,
			PaymentMethod:      req.PaymentMethod,
			PaymentProviderRef: req.PaymentProviderRef,
			ProductID:          req.ProductID,
			TargetNumber:       req.TargetNumber,
			SerialNumber:       req.SerialNumber,
			RelatedWalletID:    req.RelatedWalletID,
			CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
		}

		mockService.On("CreateTransaction", ctx, mock.MatchedBy(func(r *dto.CreateTransactionRequest) bool {
			return r.WalletID == req.WalletID && r.Amount == req.Amount
		})).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, _ := setupTransactionHandler()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("invalid")))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		req := testutil.CreateTransactionRequestFixture()

		mockService.On("CreateTransaction", ctx, mock.MatchedBy(func(r *dto.CreateTransactionRequest) bool {
			return true
		})).Return(nil, errors.New("service error"))

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetTransactionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get transaction by ID successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()

		response := &dto.TransactionResponse{
			ID:        transactionID,
			WalletID:  1,
			Type:      entity.TypeTopup,
			Amount:    100.00,
			Status:    entity.StatusPending,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		mockService.On("GetTransactionByID", ctx, transactionID).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions/"+transactionID.String(), nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.GetTransactionByID(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID format", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions/invalid-id", nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: "invalid-id"},
		}

		// When
		handler.GetTransactionByID(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when transaction not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()

		mockService.On("GetTransactionByID", ctx, transactionID).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions/"+transactionID.String(), nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.GetTransactionByID(ginCtx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetAllTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all transactions successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		response := &dto.TransactionListResponse{
			Data: []dto.TransactionResponse{
				{ID: uuid.New(), WalletID: 1, Type: entity.TypeTopup, Amount: 100.00, Status: entity.StatusPending},
				{ID: uuid.New(), WalletID: 1, Type: entity.TypePurchase, Amount: 50.00, Status: entity.StatusSuccess},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllTransactions", ctx, mock.MatchedBy(func(filter *dto.TransactionFilter) bool {
			return filter.Page == 1 && filter.PageSize == 10
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions", nil)

		// When
		handler.GetAllTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should handle pagination parameters", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		response := &dto.TransactionListResponse{
			Data:       []dto.TransactionResponse{},
			TotalCount: 0,
			Page:       2,
			PageSize:   5,
		}

		mockService.On("GetAllTransactions", ctx, mock.MatchedBy(func(filter *dto.TransactionFilter) bool {
			return filter.Page == 2 && filter.PageSize == 5
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions?page=2&page_size=5", nil)
		ginCtx.Request.URL.RawQuery = "page=2&page_size=5"

		// When
		handler.GetAllTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should handle filter parameters", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		response := &dto.TransactionListResponse{
			Data:       []dto.TransactionResponse{},
			TotalCount: 0,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllTransactions", ctx, mock.MatchedBy(func(filter *dto.TransactionFilter) bool {
			return filter.Type == entity.TypeTopup && filter.Status == entity.StatusPending
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions?type=topup&status=pending", nil)
		ginCtx.Request.URL.RawQuery = "type=topup&status=pending"

		// When
		handler.GetAllTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		mockService.On("GetAllTransactions", ctx, mock.MatchedBy(func(filter *dto.TransactionFilter) bool {
			return true
		})).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/transactions", nil)

		// When
		handler.GetAllTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_UpdateTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should update transaction successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()
		req := testutil.CreateUpdateTransactionRequestFixture()

		response := &dto.TransactionResponse{
			ID:        transactionID,
			WalletID:  1,
			Type:      entity.TypeTopup,
			Amount:    100.00,
			Status:    req.Status,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		mockService.On("UpdateTransaction", ctx, transactionID, mock.MatchedBy(func(r *dto.UpdateTransactionRequest) bool {
			return r.Status == req.Status
		})).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("PUT", "/transactions/"+transactionID.String(), bytes.NewBuffer(reqBody))
		ginCtx.Request.Header.Set("Content-Type", "application/json")
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.UpdateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID format", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		req := testutil.CreateUpdateTransactionRequestFixture()

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("PUT", "/transactions/invalid-id", bytes.NewBuffer(reqBody))
		ginCtx.Request.Header.Set("Content-Type", "application/json")
		ginCtx.Params = gin.Params{
			{Key: "id", Value: "invalid-id"},
		}

		// When
		handler.UpdateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("PUT", "/transactions/"+transactionID.String(), bytes.NewBuffer([]byte("invalid")))
		ginCtx.Request.Header.Set("Content-Type", "application/json")
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.UpdateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()
		req := testutil.CreateUpdateTransactionRequestFixture()

		mockService.On("UpdateTransaction", ctx, transactionID, mock.MatchedBy(func(r *dto.UpdateTransactionRequest) bool {
			return true
		})).Return(nil, errors.New("service error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("PUT", "/transactions/"+transactionID.String(), bytes.NewBuffer(reqBody))
		ginCtx.Request.Header.Set("Content-Type", "application/json")
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.UpdateTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete transaction successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()

		mockService.On("DeleteTransaction", ctx, transactionID).Return(nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("DELETE", "/transactions/"+transactionID.String(), nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.DeleteTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID format", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("DELETE", "/transactions/invalid-id", nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: "invalid-id"},
		}

		// When
		handler.DeleteTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()
		transactionID := uuid.New()

		mockService.On("DeleteTransaction", ctx, transactionID).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("DELETE", "/transactions/"+transactionID.String(), nil)
		ginCtx.Params = gin.Params{
			{Key: "id", Value: transactionID.String()},
		}

		// When
		handler.DeleteTransaction(ginCtx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetWalletTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("should get wallet transactions successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		response := &dto.TransactionListResponse{
			Data: []dto.TransactionResponse{
				{ID: uuid.New(), WalletID: 1, Type: entity.TypeTopup, Amount: 100.00, Status: entity.StatusPending},
				{ID: uuid.New(), WalletID: 1, Type: entity.TypePurchase, Amount: 50.00, Status: entity.StatusSuccess},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetWalletTransactions", ctx, uint(1), 1, 10).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/wallets/1/transactions", nil)
		ginCtx.Params = gin.Params{
			{Key: "wallet_id", Value: "1"},
		}

		// When
		handler.GetWalletTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid wallet ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/wallets/invalid/transactions", nil)
		ginCtx.Params = gin.Params{
			{Key: "wallet_id", Value: "invalid"},
		}

		// When
		handler.GetWalletTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should handle pagination parameters", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		response := &dto.TransactionListResponse{
			Data:       []dto.TransactionResponse{},
			TotalCount: 0,
			Page:       2,
			PageSize:   5,
		}

		mockService.On("GetWalletTransactions", ctx, uint(1), 2, 5).Return(response, nil)

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/wallets/1/transactions?page=2&page_size=5", nil)
		ginCtx.Request.URL.RawQuery = "page=2&page_size=5"
		ginCtx.Params = gin.Params{
			{Key: "wallet_id", Value: "1"},
		}

		// When
		handler.GetWalletTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupTransactionHandler()

		mockService.On("GetWalletTransactions", ctx, uint(1), 1, 10).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest("GET", "/wallets/1/transactions", nil)
		ginCtx.Params = gin.Params{
			{Key: "wallet_id", Value: "1"},
		}

		// When
		handler.GetWalletTransactions(ginCtx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
