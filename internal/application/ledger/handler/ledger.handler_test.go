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

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/config"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"
)

func setupHandlerTest() (*LedgerHandler, *testutil.MockLedgerService) {
	gin.SetMode(gin.TestMode)
	mockService := new(testutil.MockLedgerService)
	logger := zap.NewNop()

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret",
			Expiry:    24 * time.Hour,
		},
	}
	jwtManager := jwt.NewJWTManager(cfg)
	handler := NewLedgerHandler(mockService, jwtManager, logger)
	return handler, mockService
}

func TestCreateEntry(t *testing.T) {
	t.Run("Success with Debit", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
		}
		expectedResponse := &dto.GetLedgerEntryResponse{
			ID:            1,
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
			Amount:        50000,
			CreatedAt:     time.Now(),
		}

		mockService.On("CreateEntry", mock.Anything, mock.AnythingOfType("*dto.CreateLedgerEntryRequest")).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/ledger", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEntry(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.Equal(t, "Ledger entry created successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Success with Credit", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "payment",
			Debit:         0,
			Credit:        50000,
			AccountType:   "wallet",
		}
		expectedResponse := &dto.GetLedgerEntryResponse{
			ID:            2,
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "payment",
			Debit:         0,
			Credit:        50000,
			AccountType:   "wallet",
			Amount:        50000,
			CreatedAt:     time.Now(),
		}

		mockService.On("CreateEntry", mock.Anything, mock.AnythingOfType("*dto.CreateLedgerEntryRequest")).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/ledger", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEntry(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/ledger", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEntry(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		request := dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
		}

		mockService.On("CreateEntry", mock.Anything, mock.AnythingOfType("*dto.CreateLedgerEntryRequest")).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/ledger", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEntry(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetEntryByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.GetLedgerEntryResponse{
			ID:            1,
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
			Amount:        50000,
			CreatedAt:     time.Now(),
		}

		mockService.On("GetEntryByID", mock.Anything, uint64(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ledger/1", nil)

		handler.GetEntryByID(c)

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
		c.Request = httptest.NewRequest("GET", "/ledger/invalid", nil)

		handler.GetEntryByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("EntryNotFound", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetEntryByID", mock.Anything, uint64(999)).Return(nil, errors.New("entry not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("GET", "/ledger/999", nil)

		handler.GetEntryByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetEntriesByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedEntries := []dto.GetLedgerEntryResponse{
			{
				ID:            1,
				UserID:        1,
				ReferenceID:   "100",
				ReferenceType: "transfer",
				Debit:         50000,
				Credit:        0,
				AccountType:   "wallet",
				Amount:        50000,
			},
			{
				ID:            2,
				UserID:        1,
				ReferenceID:   "101",
				ReferenceType: "payment",
				Debit:         0,
				Credit:        30000,
				AccountType:   "wallet",
				Amount:        30000,
			},
		}

		mockService.On("GetEntriesByUserID", mock.Anything, uint64(1)).Return(expectedEntries, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ledger/user/1", nil)

		handler.GetEntriesByUserID(c)

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
		c.Request = httptest.NewRequest("GET", "/ledger/user/invalid", nil)

		handler.GetEntriesByUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetEntriesByUserID", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ledger/user/1", nil)

		handler.GetEntriesByUserID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetEntriesByReference(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedEntries := []dto.GetLedgerEntryResponse{
			{
				ID:            1,
				UserID:        1,
				ReferenceID:   "100",
				ReferenceType: "transfer",
				Debit:         50000,
				Credit:        0,
				AccountType:   "wallet",
				Amount:        50000,
			},
		}

		mockService.On("GetEntriesByReference", mock.Anything, "transfer", "100").Return(expectedEntries, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "reference_type", Value: "transfer"}, {Key: "reference_id", Value: "100"}}
		c.Request = httptest.NewRequest("GET", "/ledger/reference/transfer/100", nil)

		handler.GetEntriesByReference(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidReferenceID", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetEntriesByReference", mock.Anything, "transfer", "invalid").Return([]dto.GetLedgerEntryResponse{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "reference_type", Value: "transfer"}, {Key: "reference_id", Value: "invalid"}}
		c.Request = httptest.NewRequest("GET", "/ledger/reference/transfer/invalid", nil)

		handler.GetEntriesByReference(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListEntries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.ListLedgerEntriesResponse{
			Data: []dto.GetLedgerEntryResponse{
				{ID: 1, UserID: 1, ReferenceID: "100", ReferenceType: "transfer", Debit: 50000, Credit: 0, AccountType: "wallet", Amount: 50000},
				{ID: 2, UserID: 1, ReferenceID: "101", ReferenceType: "payment", Debit: 0, Credit: 30000, AccountType: "wallet", Amount: 30000},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListEntries", mock.Anything, mock.AnythingOfType("*dto.ListLedgerEntriesRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ledger?page=1&page_size=10", nil)

		handler.ListEntries(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListLedgerEntriesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("WithFilters", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.ListLedgerEntriesResponse{
			Data: []dto.GetLedgerEntryResponse{
				{ID: 1, UserID: 1, ReferenceID: "100", ReferenceType: "transfer", Debit: 50000, Credit: 0, AccountType: "wallet", Amount: 50000},
			},
			TotalCount: 1,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("ListEntries", mock.Anything, mock.AnythingOfType("*dto.ListLedgerEntriesRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ledger?page=1&page_size=10&user_id=1&account_type=wallet", nil)

		handler.ListEntries(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidQuery", func(t *testing.T) {
		handler, _ := setupHandlerTest()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ledger?page=invalid", nil)

		handler.ListEntries(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("ListEntries", mock.Anything, mock.AnythingOfType("*dto.ListLedgerEntriesRequest")).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ledger?page=1&page_size=10", nil)

		handler.ListEntries(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		expectedResponse := &dto.LedgerStatsResponse{
			UserID:          1,
			TotalDebit:      100000,
			TotalCredit:     50000,
			EntryCount:      5,
			WalletDebit:     80000,
			WalletCredit:    40000,
			PaylaterDebit:   20000,
			PaylaterCredits: 10000,
		}

		mockService.On("GetUserStats", mock.Anything, uint64(1)).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ledger/stats/1", nil)

		handler.GetUserStats(c)

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
		c.Request = httptest.NewRequest("GET", "/ledger/stats/invalid", nil)

		handler.GetUserStats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := setupHandlerTest()

		mockService.On("GetUserStats", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "user_id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ledger/stats/1", nil)

		handler.GetUserStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
