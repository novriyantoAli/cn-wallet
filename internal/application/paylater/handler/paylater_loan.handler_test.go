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

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/config"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
)

func setupHandlerTest() (*PaylaterLoanHandler, *testutil.MockPaylaterLoanService) {
	gin.SetMode(gin.TestMode)
	mockService := new(testutil.MockPaylaterLoanService)
	logger := zap.NewNop()

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret",
			Expiry:    24 * time.Hour,
		},
	}
	jwtManager := jwt.NewJWTManager(cfg)
	handler := NewPaylaterLoanHandler(mockService, jwtManager, logger)
	return handler, mockService
}

func TestCreateLoan(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		now := time.Now()
		dueDate := now.Add(30 * 24 * time.Hour)
		request := dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   500000,
			Interest: 50000,
			DueDate:  dueDate,
			Source:   "transfer",
		}
		expectedLoan := &entity.PaylaterLoan{
			UserID:   1,
			Amount:   500000,
			Interest: 50000,
			Status:   entity.PaylaterLoanStatusActive,
			DueDate:  dueDate,
			Source:   entity.PaylaterLoanSourceTransfer,
		}
		expectedLoan.ID = 1

		mockService.On("CreateLoan", mock.Anything, mock.AnythingOfType("*dto.CreatePaylaterLoanRequest")).Return(expectedLoan, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/paylater/loans", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLoan(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/paylater/loans", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLoan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		now := time.Now()
		dueDate := now.Add(30 * 24 * time.Hour)
		request := dto.CreatePaylaterLoanRequest{
			UserID:   1,
			Amount:   500000,
			Interest: 50000,
			DueDate:  dueDate,
			Source:   "transfer",
		}

		mockService.On("CreateLoan", mock.Anything, mock.AnythingOfType("*dto.CreatePaylaterLoanRequest")).Return(nil, errors.New("insufficient credit limit"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/paylater/loans", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLoan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetLoanByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.GetPaylaterLoanResponse{
			ID:       1,
			UserID:   1,
			Amount:   500000,
			Interest: 50000,
			Total:    550000,
			Status:   "active",
			DueDate:  time.Now().Add(30 * 24 * time.Hour),
			Source:   "transfer",
		}

		mockService.On("GetLoanByID", mock.Anything, uint(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/1", nil)

		handler.GetLoanByID(c)

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
		c.Request = httptest.NewRequest("GET", "/paylater/loans/invalid", nil)

		handler.GetLoanByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("LoanNotFound", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetLoanByID", mock.Anything, uint(999)).Return(nil, errors.New("loan not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/999", nil)

		handler.GetLoanByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetLoansByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedLoans := []dto.GetPaylaterLoanResponse{
			{
				ID:       1,
				UserID:   1,
				Amount:   500000,
				Interest: 50000,
				Status:   "active",
			},
			{
				ID:       2,
				UserID:   1,
				Amount:   300000,
				Interest: 30000,
				Status:   "paid",
			},
		}

		mockService.On("GetLoansByUserID", mock.Anything, uint(1)).Return(expectedLoans, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/by-user/1", nil)

		handler.GetLoansByUserID(c)

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
		c.Request = httptest.NewRequest("GET", "/paylater/loans/by-user/invalid", nil)

		handler.GetLoansByUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetLoansByUserID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/by-user/1", nil)

		handler.GetLoansByUserID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListLoans(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.ListPaylaterLoansResponse{
			Data: []dto.GetPaylaterLoanResponse{
				{ID: 1, UserID: 1, Amount: 500000, Status: "active"},
				{ID: 2, UserID: 2, Amount: 300000, Status: "paid"},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListLoans", mock.Anything, mock.AnythingOfType("*dto.ListPaylaterLoansRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/paylater/loans?page=1&page_size=10", nil)

		handler.ListLoans(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListPaylaterLoansResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("WithFilters", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.ListPaylaterLoansResponse{
			Data: []dto.GetPaylaterLoanResponse{
				{ID: 1, UserID: 1, Amount: 500000, Status: "active"},
			},
			TotalCount: 1,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListLoans", mock.Anything, mock.AnythingOfType("*dto.ListPaylaterLoansRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/paylater/loans?page=1&page_size=10&user_id=1&status=active", nil)

		handler.ListLoans(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidQuery", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/paylater/loans?page=invalid", nil)

		handler.ListLoans(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("ListLoans", mock.Anything, mock.AnythingOfType("*dto.ListPaylaterLoansRequest")).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/paylater/loans?page=1&page_size=10", nil)

		handler.ListLoans(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateLoanStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.UpdateLoanStatusRequest{
			Status: "overdue",
		}
		expectedResponse := &dto.GetPaylaterLoanResponse{
			ID:     1,
			UserID: 1,
			Amount: 500000,
			Status: "overdue",
		}

		mockService.On("UpdateLoanStatus", mock.Anything, uint(1), mock.AnythingOfType("*dto.UpdateLoanStatusRequest")).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/paylater/loans/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLoanStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.Equal(t, "Loan status updated successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		request := dto.UpdateLoanStatusRequest{Status: "overdue"}
		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/paylater/loans/invalid/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLoanStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		c.Request = httptest.NewRequest("PUT", "/paylater/loans/1/status", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLoanStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.UpdateLoanStatusRequest{
			Status: string(entity.PaylaterLoanStatusOverdue),
		}

		mockService.On("UpdateLoanStatus", mock.Anything, uint(1), mock.AnythingOfType("*dto.UpdateLoanStatusRequest")).Return(nil, errors.New("loan not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("PUT", "/paylater/loans/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLoanStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestMarkLoanAsPaid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.GetPaylaterLoanResponse{
			ID:     1,
			UserID: 1,
			Amount: 500000,
			Status: "paid",
		}

		mockService.On("MarkLoanAsPaid", mock.Anything, uint(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/paylater/loans/1/mark-paid", nil)

		handler.MarkLoanAsPaid(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.Equal(t, "Loan marked as paid successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("POST", "/paylater/loans/invalid/mark-paid", nil)

		handler.MarkLoanAsPaid(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("MarkLoanAsPaid", mock.Anything, uint(1)).Return(nil, errors.New("loan already paid"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/paylater/loans/1/mark-paid", nil)

		handler.MarkLoanAsPaid(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProcessOverdueLoans(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("ProcessOverdueLoans", mock.Anything).Return(5, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/paylater/loans/process-overdue", nil)

		handler.ProcessOverdueLoans(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Overdue loans processed successfully", response["message"])
		assert.Equal(t, float64(5), response["count"])
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("ProcessOverdueLoans", mock.Anything).Return(0, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/paylater/loans/process-overdue", nil)

		handler.ProcessOverdueLoans(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserLoanStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedStats := &dto.PaylaterLoanStatsResponse{
			UserID:           1,
			TotalLoans:       10,
			ActiveLoans:      3,
			TotalBorrowed:    5000000,
			TotalPaid:        3500000,
			TotalOutstanding: 1500000,
			OverdueLoans:     1,
		}

		mockService.On("GetUserLoanStats", mock.Anything, uint(1)).Return(expectedStats, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/stats/1", nil)

		handler.GetUserLoanStats(c)

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
		c.Request = httptest.NewRequest("GET", "/paylater/loans/stats/invalid", nil)

		handler.GetUserLoanStats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetUserLoanStats", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/paylater/loans/stats/1", nil)

		handler.GetUserLoanStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteLoan(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("DeleteLoan", mock.Anything, uint(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("DELETE", "/paylater/loans/1", nil)

		handler.DeleteLoan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Loan deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("DELETE", "/paylater/loans/invalid", nil)

		handler.DeleteLoan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("DeleteLoan", mock.Anything, uint(1)).Return(errors.New("loan not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("DELETE", "/paylater/loans/1", nil)

		handler.DeleteLoan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}
