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
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func TestPaylaterRepaymentHandler_ProcessRepayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Process repayment with valid loan", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		response := &dto.GetPaylaterRepaymentResponse{
			ID:             1,
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockService.On("ProcessRepayment", mock.Anything, req).Return(response, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid request body", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Loan not found", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 999,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		mockService.On("ProcessRepayment", mock.Anything, req).Return(nil, errors.New("loan not found"))

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Unauthorized repayment", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         2,
			Amount:         500000,
			PaymentSource:  "wallet",
		}

		mockService.On("ProcessRepayment", mock.Anything, req).Return(nil, errors.New("unauthorized repayment"))

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid repayment amount", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		// The DTO has binding:"gt=0" which validates the amount client-side
		// So we need to send a zero amount to trigger binding validation error
		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         0,
			PaymentSource:  "wallet",
		}

		// Do NOT expect service call - binding validation should fail first
		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "ProcessRepayment")
	})

	t.Run("Success: Process repayment with external payment source", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		req := &dto.CreatePaylaterRepaymentRequest{
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         250000,
			PaymentSource:  "external_payment",
		}

		response := &dto.GetPaylaterRepaymentResponse{
			ID:             2,
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         250000,
			PaymentSource:  "external_payment",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockService.On("ProcessRepayment", mock.Anything, req).Return(response, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/paylater/repayments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler.ProcessRepayment(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentHandler_GetRepaymentByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Get repayment by ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		response := &dto.GetPaylaterRepaymentResponse{
			ID:             1,
			PaylaterLoanID: 1,
			UserID:         1,
			Amount:         500000,
			PaymentSource:  "wallet",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockService.On("GetRepaymentByID", mock.Anything, uint(1)).Return(response, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/repayments/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetRepaymentByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid repayment ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/repayments/invalid", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetRepaymentByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error: Repayment not found", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		mockService.On("GetRepaymentByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/repayments/999", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetRepaymentByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success: Get different repayment by ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		response := &dto.GetPaylaterRepaymentResponse{
			ID:             5,
			PaylaterLoanID: 2,
			UserID:         3,
			Amount:         1000000,
			PaymentSource:  "external_payment",
			Status:         "success",
			CreatedAt:      time.Now(),
		}

		mockService.On("GetRepaymentByID", mock.Anything, uint(5)).Return(response, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/repayments/5", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "id", Value: "5"}}

		handler.GetRepaymentByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentHandler_ListRepaymentsByLoan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: List all repayments for a loan", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		responses := []dto.GetPaylaterRepaymentResponse{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         300000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-1 * time.Hour),
			},
		}

		mockService.On("ListRepaymentsByLoan", mock.Anything, uint(1)).Return(responses, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/loans/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "loan_id", Value: "1"}}

		handler.ListRepaymentsByLoan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid loan ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/loans/invalid/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "loan_id", Value: "invalid"}}

		handler.ListRepaymentsByLoan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success: Empty list for loan with no repayments", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		mockService.On("ListRepaymentsByLoan", mock.Anything, uint(999)).Return([]dto.GetPaylaterRepaymentResponse{}, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/loans/999/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "loan_id", Value: "999"}}

		handler.ListRepaymentsByLoan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Database error when listing repayments", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		mockService.On("ListRepaymentsByLoan", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/loans/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "loan_id", Value: "1"}}

		handler.ListRepaymentsByLoan(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success: List multiple repayments for the same loan", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		responses := []dto.GetPaylaterRepaymentResponse{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-2 * time.Hour),
			},
			{
				ID:             3,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         1000000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-4 * time.Hour),
			},
		}

		mockService.On("ListRepaymentsByLoan", mock.Anything, uint(1)).Return(responses, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/loans/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "loan_id", Value: "1"}}

		handler.ListRepaymentsByLoan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestPaylaterRepaymentHandler_GetRepaymentsByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success: Get all repayments for a user", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		responses := []dto.GetPaylaterRepaymentResponse{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         300000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 2,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-1 * time.Hour),
			},
		}

		mockService.On("GetRepaymentsByUser", mock.Anything, uint(1)).Return(responses, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Invalid user ID", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/invalid/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "invalid"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success: Empty list for user with no repayments", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		mockService.On("GetRepaymentsByUser", mock.Anything, uint(999)).Return([]dto.GetPaylaterRepaymentResponse{}, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/999/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "999"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error: Database error when getting repayments", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		mockService.On("GetRepaymentsByUser", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success: User with repayments from multiple loans", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		responses := []dto.GetPaylaterRepaymentResponse{
			{
				ID:             1,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				PaylaterLoanID: 1,
				UserID:         1,
				Amount:         500000,
				PaymentSource:  "external_payment",
				Status:         "success",
				CreatedAt:      time.Now().Add(-2 * time.Hour),
			},
			{
				ID:             3,
				PaylaterLoanID: 2,
				UserID:         1,
				Amount:         1000000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-4 * time.Hour),
			},
			{
				ID:             4,
				PaylaterLoanID: 3,
				UserID:         1,
				Amount:         200000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now().Add(-6 * time.Hour),
			},
		}

		mockService.On("GetRepaymentsByUser", mock.Anything, uint(1)).Return(responses, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/1/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success: Get repayments for different user", func(t *testing.T) {
		mockService := new(testutil.MockPaylaterRepaymentService)
		logger := zap.NewNop()

		responses := []dto.GetPaylaterRepaymentResponse{
			{
				ID:             10,
				PaylaterLoanID: 5,
				UserID:         2,
				Amount:         750000,
				PaymentSource:  "wallet",
				Status:         "success",
				CreatedAt:      time.Now(),
			},
		}

		mockService.On("GetRepaymentsByUser", mock.Anything, uint(2)).Return(responses, nil)

		handler := NewPaylaterRepaymentHandler(mockService, nil, logger)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/paylater/users/2/repayments", nil)

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		c.Params = gin.Params{{Key: "user_id", Value: "2"}}

		handler.GetRepaymentsByUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
