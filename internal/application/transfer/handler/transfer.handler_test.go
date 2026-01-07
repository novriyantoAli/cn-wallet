package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	"github.com/novriyantoAli/cn-wallet/internal/config"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
)

func setupHandlerTest() (*TransferHandler, *testutil.MockTransferService) {
	gin.SetMode(gin.TestMode)
	mockService := new(testutil.MockTransferService)
	logger := zap.NewNop()

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret",
			Expiry:    24 * time.Hour,
		},
	}
	jwtManager := jwt.NewJWTManager(cfg)
	handler := NewTransferHandler(mockService, jwtManager, logger)
	return handler, mockService
}

func TestCreateTransfer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       "wallet",
		}
		expectedTransfer := &entity.Transfer{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Status:       entity.TransferStatusPending,
			Source:       entity.TransferSourceWallet,
		}
		expectedTransfer.ID = 1

		mockService.On("CreateTransfer", mock.Anything, mock.AnythingOfType("*dto.CreateTransferRequest")).Return(expectedTransfer, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/transfers", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateTransfer(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.Equal(t, "Transfer created successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/transfers", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateTransfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.CreateTransferRequest{
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       "wallet",
		}

		mockService.On("CreateTransfer", mock.Anything, mock.AnythingOfType("*dto.CreateTransferRequest")).Return(nil, errors.New("insufficient balance"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/transfers", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateTransfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTransferByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.TransferResponse{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       "wallet",
			Status:       "pending",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockService.On("GetTransferByID", mock.Anything, uint(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/transfers/1", nil)

		handler.GetTransferByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("GET", "/transfers/invalid", nil)

		handler.GetTransferByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TransferNotFound", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetTransferByID", mock.Anything, uint(999)).Return(nil, errors.New("transfer not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("GET", "/transfers/999", nil)

		handler.GetTransferByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTransfersByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedTransfers := []dto.TransferResponse{
			{
				ID:           1,
				UserID:       1,
				TargetUserID: 2,
				Amount:       100000,
				Source:       "wallet",
				Status:       "completed",
			},
			{
				ID:           2,
				UserID:       1,
				TargetUserID: 3,
				Amount:       50000,
				Source:       "paylater",
				Status:       "pending",
			},
		}

		mockService.On("GetTransfersByUserID", mock.Anything, uint(1)).Return(expectedTransfers, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-user/1", nil)

		handler.GetTransfersByUserID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-user/invalid", nil)

		handler.GetTransfersByUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetTransfersByUserID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-user/1", nil)

		handler.GetTransfersByUserID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTransfersByTargetUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedTransfers := []dto.TransferResponse{
			{
				ID:           3,
				UserID:       1,
				TargetUserID: 2,
				Amount:       100000,
				Source:       "wallet",
				Status:       "completed",
			},
			{
				ID:           4,
				UserID:       3,
				TargetUserID: 2,
				Amount:       50000,
				Source:       "paylater",
				Status:       "completed",
			},
		}

		mockService.On("GetTransfersByTargetUserID", mock.Anything, uint(2)).Return(expectedTransfers, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "target_user_id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-target/2", nil)

		handler.GetTransfersByTargetUserID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidTargetUserID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "target_user_id", Value: "invalid"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-target/invalid", nil)

		handler.GetTransfersByTargetUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetTransfersByTargetUserID", mock.Anything, uint(2)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "target_user_id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/transfers/by-target/2", nil)

		handler.GetTransfersByTargetUserID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListTransfers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.TransferListResponse{
			Data: []dto.TransferResponse{
				{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: "wallet", Status: "completed"},
				{ID: 2, UserID: 1, TargetUserID: 3, Amount: 50000, Source: "paylater", Status: "pending"},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListTransfers", mock.Anything, mock.AnythingOfType("*dto.TransferFilter")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/transfers?page=1&page_size=10", nil)

		handler.ListTransfers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.TransferListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("WithFilters", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.TransferListResponse{
			Data: []dto.TransferResponse{
				{ID: 1, UserID: 1, TargetUserID: 2, Amount: 100000, Source: "wallet", Status: "completed"},
			},
			TotalCount: 1,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListTransfers", mock.Anything, mock.AnythingOfType("*dto.TransferFilter")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/transfers?page=1&page_size=10&user_id=1&status=completed", nil)

		handler.ListTransfers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidQuery", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/transfers?page=invalid", nil)

		handler.ListTransfers(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("ListTransfers", mock.Anything, mock.AnythingOfType("*dto.TransferFilter")).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/transfers?page=1&page_size=10", nil)

		handler.ListTransfers(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateTransferStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.UpdateTransferStatusRequest{
			Status: "completed",
		}
		expectedResponse := &dto.TransferResponse{
			ID:           1,
			UserID:       1,
			TargetUserID: 2,
			Amount:       100000,
			Source:       "wallet",
			Status:       "completed",
		}

		mockService.On("UpdateTransferStatus", mock.Anything, uint(1), mock.AnythingOfType("*dto.UpdateTransferStatusRequest")).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/transfers/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateTransferStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.Equal(t, "Transfer status updated successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		request := dto.UpdateTransferStatusRequest{Status: "completed"}
		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/transfers/invalid/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateTransferStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		c.Request = httptest.NewRequest("PUT", "/transfers/1/status", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateTransferStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.UpdateTransferStatusRequest{
			Status: "completed",
		}

		mockService.On("UpdateTransferStatus", mock.Anything, uint(1), mock.AnythingOfType("*dto.UpdateTransferStatusRequest")).Return(nil, errors.New("transfer not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/transfers/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateTransferStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserTransferStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.TransferStatsResponse{
			UserID:        1,
			TotalSent:     250000,
			TotalReceived: 100000,
			CountSent:     3,
			CountReceived: 2,
			WalletSent:    150000,
			PaylaterSent:  100000,
			CompletedSent: 200000,
		}

		mockService.On("GetUserTransferStats", mock.Anything, uint(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/transfers/stats/1", nil)

		handler.GetUserTransferStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}
		c.Request = httptest.NewRequest("GET", "/transfers/stats/invalid", nil)

		handler.GetUserTransferStats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetUserTransferStats", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/transfers/stats/1", nil)

		handler.GetUserTransferStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCancelTransfer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("CancelTransfer", mock.Anything, uint(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/transfers/1/cancel", nil)

		handler.CancelTransfer(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer cancelled successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("POST", "/transfers/invalid/cancel", nil)

		handler.CancelTransfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("CancelTransfer", mock.Anything, uint(999)).Return(errors.New("transfer not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("POST", "/transfers/999/cancel", nil)

		handler.CancelTransfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}
