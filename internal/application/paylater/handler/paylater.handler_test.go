package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func TestPaylaterAccountHandler_CreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Create paylater account", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterAccountRequest{
			UserID:      1,
			CreditLimit: 1000000,
		}

		account := &entity.PaylaterAccount{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    0,
			AvailableLimit: 1000000,
			Status:         "active",
			CreatedAt:      time.Now(),
		}

		mockService.On("CreateAccount", mock.Anything, req).Return(account, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.CreateAccount(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid request body", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.CreateAccount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Service returns error", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterAccountRequest{
			UserID:      1,
			CreditLimit: 1000000,
		}

		mockService.On("CreateAccount", mock.Anything, req).Return(nil, errors.New("account already exists"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.CreateAccount(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterAccountHandler_GetAccountByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Get account by user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		response := &dto.GetPaylaterAccountResponse{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         "active",
			CreatedAt:      time.Now(),
		}

		mockService.On("GetAccountByUserID", mock.Anything, uint(1)).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/by-user/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.GetAccountByUserID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/by-user/invalid", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.GetAccountByUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		mockService.On("GetAccountByUserID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/by-user/999", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "999"}}

		handler.GetAccountByUserID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterAccountHandler_GetAccountByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Get account by ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		response := &dto.GetPaylaterAccountResponse{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Status:         "active",
			CreatedAt:      time.Now(),
		}

		mockService.On("GetAccountByID", mock.Anything, uint(1)).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/invalid", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Account not found", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		mockService.On("GetAccountByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/999", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterAccountHandler_UpdateCreditLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Update credit limit", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.UpdateCreditLimitRequest{
			CreditLimit: 2000000,
		}

		response := &dto.GetPaylaterAccountResponse{
			ID:             1,
			UserID:         1,
			CreditLimit:    2000000,
			Outstanding:    500000,
			AvailableLimit: 1500000,
			Status:         "active",
			CreatedAt:      time.Now(),
		}

		mockService.On("UpdateCreditLimit", mock.Anything, uint(1), req).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/paylater/by-user/1/credit-limit", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.UpdateCreditLimit(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/paylater/by-user/invalid/credit-limit", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.UpdateCreditLimit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Invalid request body", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/paylater/by-user/1/credit-limit", bytes.NewBuffer([]byte("invalid")))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.UpdateCreditLimit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPaylaterAccountHandler_UpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Update status", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.UpdateStatusRequest{
			Status: string(entity.PaylaterStatusSuspended),
		}

		response := &dto.GetPaylaterAccountResponse{
			ID:             1,
			UserID:         1,
			CreditLimit:    1000000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Status:         string(entity.PaylaterStatusSuspended),
			CreatedAt:      time.Now(),
		}

		mockService.On("UpdateStatus", mock.Anything, uint(1), req).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/paylater/by-user/1/status", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.UpdateStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/paylater/by-user/invalid/status", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.UpdateStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPaylaterAccountHandler_UseCredit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Use credit", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.UseCreditRequest{
			Amount: 200000,
		}

		response := &dto.UseCreditResponse{
			UserID:         1,
			Amount:         200000,
			Outstanding:    500000,
			AvailableLimit: 500000,
			Message:        "Credit used successfully",
		}

		mockService.On("UseCredit", mock.Anything, uint(1), req).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/1/use-credit", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.UseCredit(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/invalid/use-credit", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.UseCredit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Insufficient credit", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.UseCreditRequest{
			Amount: 800000,
		}

		mockService.On("UseCredit", mock.Anything, uint(1), req).Return(nil, errors.New("insufficient credit limit"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/1/use-credit", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.UseCredit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterAccountHandler_Repayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Process repayment", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.RepaymentRequest{
			Amount: 200000,
		}

		response := &dto.RepaymentResponse{
			UserID:         1,
			Amount:         200000,
			Outstanding:    300000,
			AvailableLimit: 700000,
			Message:        "Repayment processed successfully",
		}

		mockService.On("Repayment", mock.Anything, uint(1), req).Return(response, nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/1/repayment", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.Repayment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/invalid/repayment", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.Repayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Repayment exceeds outstanding", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		req := &dto.RepaymentRequest{
			Amount: 600000,
		}

		mockService.On("Repayment", mock.Anything, uint(1), req).Return(nil, errors.New("repayment amount cannot exceed outstanding balance"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/by-user/1/repayment", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.Repayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterAccountHandler_DeleteAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Delete account", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		mockService.On("DeleteAccount", mock.Anything, uint(1)).Return(nil)

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/paylater/by-user/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.DeleteAccount(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/paylater/by-user/invalid", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.DeleteAccount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Account has outstanding balance", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterAccountService)
		logger := zap.NewNop()

		mockService.On("DeleteAccount", mock.Anything, uint(1)).Return(errors.New("cannot delete account with outstanding balance"))

		handler := NewPaylaterAccountHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/paylater/by-user/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.DeleteAccount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}
