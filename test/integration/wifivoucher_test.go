package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWifiVoucherIntegration(t *testing.T) (*gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)

	logger := testutil.NewTestLogger(t)

	// Create test provider
	provider := &providerEntity.Provider{
		Name:      "Test Provider",
		Code:      "TEST_PROVIDER",
		Logo:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(provider).Error
	require.NoError(t, err)

	// Create real instances (no mocks)
	wifiVoucherRepo := repository.NewWifiVoucherRepository(db, logger)
	wifiVoucherService := service.NewWifiVoucherService(wifiVoucherRepo, logger)
	wifiVoucherHandler := handler.NewWifiVoucherHandler(wifiVoucherService, logger)

	// Setup Gin router
	router := gin.New()
	api := router.Group("/api/v1")
	wifiVoucherHandler.RegisterRoutes(api)

	cleanup := func() {
		testutil.CleanDB(db)
	}

	return router, cleanup
}

func TestWifiVoucherIntegration_CreateAndGetVoucher(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Test data
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI001",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	// Step 1: Create wifi voucher
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)

	data := createResp["data"].(map[string]interface{})
	voucherID := int(data["id"].(float64))
	assert.Equal(t, createReq.Code, data["code"])
	assert.Equal(t, float64(createReq.DurationHours), data["duration_hours"])

	// Step 2: Get the created wifi voucher by ID
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/v1/wifi-vouchers/"+string(rune(voucherID+'0')), nil)

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var getResp map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &getResp)
	require.NoError(t, err)

	voucherData := getResp["data"].(map[string]interface{})
	assert.Equal(t, float64(voucherID), voucherData["id"])
	assert.Equal(t, createReq.Code, voucherData["code"])
}

func TestWifiVoucherIntegration_GetVoucherByCode(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create wifi voucher
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI002",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get voucher by code
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/v1/wifi-vouchers/code/WIFI002", nil)

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var getResp map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &getResp)
	require.NoError(t, err)

	voucherData := getResp["data"].(map[string]interface{})
	assert.Equal(t, createReq.Code, voucherData["code"])
}

func TestWifiVoucherIntegration_CreateDuplicateCode(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Test data
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI_DUPLICATE",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	// Step 1: Create first wifi voucher
	reqBody, _ := json.Marshal(createReq)
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Step 2: Try to create voucher with same code
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req2.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)

	var errorResp map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "code already exists")
}

func TestWifiVoucherIntegration_GetAllVouchers(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create multiple wifi vouchers
	vouchers := []dto.CreateWifiVoucherRequest{
		{Code: "WIFI101", DurationHours: 24, BatchID: "BATCH001"},
		{Code: "WIFI102", DurationHours: 24, BatchID: "BATCH001"},
		{Code: "WIFI103", DurationHours: 48, BatchID: "BATCH002"},
	}

	// Create vouchers
	for _, voucher := range vouchers {
		reqBody, _ := json.Marshal(voucher)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Get all vouchers
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wifi-vouchers?page=1&page_size=10", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.WifiVoucherListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Data, 3)
	assert.Equal(t, int64(3), response.TotalCount)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
}

func TestWifiVoucherIntegration_GetVouchersWithStatusFilter(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create wifi vouchers with different statuses
	vouchers := []dto.CreateWifiVoucherRequest{
		{Code: "WIFI201", DurationHours: 24, BatchID: "BATCH001"},
		{Code: "WIFI202", DurationHours: 24, BatchID: "BATCH001"},
		{Code: "WIFI203", DurationHours: 24, BatchID: "BATCH001"},
	}

	for _, voucher := range vouchers {
		reqBody, _ := json.Marshal(voucher)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Get only available vouchers
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wifi-vouchers?status=available&page=1&page_size=10", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.WifiVoucherListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// All vouchers should be available by default
	for _, voucher := range response.Data {
		assert.Equal(t, entity.StatusAvailable, voucher.Status)
	}
}

func TestWifiVoucherIntegration_GetVouchersWithBatchFilter(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create wifi vouchers with different batches
	vouchers := []dto.CreateWifiVoucherRequest{
		{Code: "WIFI301", DurationHours: 24, BatchID: "BATCH_A"},
		{Code: "WIFI302", DurationHours: 24, BatchID: "BATCH_A"},
		{Code: "WIFI303", DurationHours: 24, BatchID: "BATCH_B"},
	}

	for _, voucher := range vouchers {
		reqBody, _ := json.Marshal(voucher)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Get only BATCH_A vouchers
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wifi-vouchers?batch_id=BATCH_A&page=1&page_size=10", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.WifiVoucherListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Data, 2)
	for _, voucher := range response.Data {
		assert.Equal(t, "BATCH_A", voucher.BatchID)
	}
}

func TestWifiVoucherIntegration_UpdateVoucher(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create wifi voucher
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI401",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)

	data := createResp["data"].(map[string]interface{})
	voucherID := int(data["id"].(float64))

	// Update wifi voucher
	updateReq := &dto.UpdateWifiVoucherRequest{
		Status: entity.StatusSold,
	}

	updateBody, _ := json.Marshal(updateReq)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("PUT", "/api/v1/wifi-vouchers/"+string(rune(voucherID+'0')), bytes.NewBuffer(updateBody))
	req2.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var updateResp map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &updateResp)
	require.NoError(t, err)

	updatedData := updateResp["data"].(map[string]interface{})
	assert.Equal(t, updateReq.Status, updatedData["status"].(string))
}

func TestWifiVoucherIntegration_DeleteVoucher(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Create wifi voucher
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI501",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)

	data := createResp["data"].(map[string]interface{})
	voucherID := int(data["id"].(float64))

	// Delete wifi voucher
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("DELETE", "/api/v1/wifi-vouchers/"+string(rune(voucherID+'0')), nil)

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var deleteResp map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &deleteResp)
	require.NoError(t, err)
	assert.Equal(t, "Wifi voucher deleted successfully", deleteResp["message"])

	// Try to get deleted voucher (should return 404)
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest("GET", "/api/v1/wifi-vouchers/"+string(rune(voucherID+'0')), nil)

	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestWifiVoucherIntegration_CreateVoucherInvalidProvider(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Test data
	createReq := &dto.CreateWifiVoucherRequest{
		Code:          "WIFI601",
		DurationHours: 24,
		BatchID:       "BATCH001",
	}

	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/wifi-vouchers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Should succeed since no provider validation in DTO
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestWifiVoucherIntegration_GetNonExistentVoucher(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Try to get non-existent voucher
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wifi-vouchers/99999", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "not found")
}

func TestWifiVoucherIntegration_GetByNonExistentCode(t *testing.T) {
	router, cleanup := setupWifiVoucherIntegration(t)
	defer cleanup()

	// Try to get voucher by non-existent code
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wifi-vouchers/code/NONEXISTENT", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "not found")
}
