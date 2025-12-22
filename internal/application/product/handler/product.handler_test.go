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

	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupProductHandler() (*ProductHandler, *testutil.MockProductService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockProductService{}
	logger := testutil.NewSilentLogger()
	handler := NewProductHandler(mockService, logger)
	return handler, mockService
}

func TestProductHandler_CreateProduct(t *testing.T) {
	t.Run("should create product successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		req := testutil.CreateProductRequestFixture()
		response := &dto.ProductResponse{
			ID:         1,
			ProviderID: req.ProviderID,
			Name:       req.Name,
			Code:       req.Code,
			Category:   req.Category,
			PriceBasic: req.PriceBasic,
			PriceSell:  req.PriceSell,
			Provider: &dto.ProviderInfo{
				ID:   1,
				Name: "Telkomsel",
				Code: "TSEL",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("CreateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateProductRequest) bool {
			return r.Code == req.Code && r.Name == req.Name
		})).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, req.Code, result.Code)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return bad request when provider not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		req := testutil.CreateProductRequestFixture()
		mockService.On("CreateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateProductRequest) bool {
			return r.Code == req.Code
		})).Return(nil, errors.New("provider not found"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request when code already exists", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		req := testutil.CreateProductRequestFixture()
		mockService.On("CreateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateProductRequest) bool {
			return r.Code == req.Code
		})).Return(nil, errors.New("code already exists"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		req := testutil.CreateProductRequestFixture()
		mockService.On("CreateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(r *dto.CreateProductRequest) bool {
			return r.Code == req.Code
		})).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetProductByID(t *testing.T) {
	t.Run("should get product by ID successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		response := &dto.ProductResponse{
			ID:         productID,
			ProviderID: 1,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Provider: &dto.ProviderInfo{
				ID:   1,
				Name: "Telkomsel",
				Code: "TSEL",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetProductByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.GetProductByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, uint(1), result.ID)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.GetProductByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(999)
		mockService.On("GetProductByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(nil, errors.New("product not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		handler.GetProductByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		mockService.On("GetProductByID", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.GetProductByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetProductByCode(t *testing.T) {
	t.Run("should get product by code successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		response := &dto.ProductResponse{
			ID:         1,
			ProviderID: 1,
			Name:       "Pulsa 10K",
			Code:       "PULSA10K",
			Category:   "Mobile",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
			Provider: &dto.ProviderInfo{
				ID:   1,
				Name: "Telkomsel",
				Code: "TSEL",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetProductByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "PULSA10K").Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/code/PULSA10K", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: "PULSA10K"},
		}

		handler.GetProductByCode(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, "PULSA10K", result.Code)
	})

	t.Run("should return bad request when code is empty", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/code/", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: ""},
		}

		handler.GetProductByCode(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		mockService.On("GetProductByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "NONEXISTENT").Return(nil, errors.New("product not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/code/NONEXISTENT", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: "NONEXISTENT"},
		}

		handler.GetProductByCode(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		mockService.On("GetProductByCode", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), "PROD1").Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products/code/PROD1", nil)
		ctx.Params = gin.Params{
			{Key: "code", Value: "PROD1"},
		}

		handler.GetProductByCode(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	t.Run("should get all products successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		response := &dto.ProductListResponse{
			Data: []dto.ProductResponse{
				{ID: 1, Name: "Product 1", Code: "PROD1", Category: "Mobile", PriceBasic: 9000, PriceSell: 10000},
				{ID: 2, Name: "Product 2", Code: "PROD2", Category: "Internet", PriceBasic: 18000, PriceSell: 20000},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllProducts", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.ProductFilter) bool {
			return true
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products?page=1&page_size=10", nil)

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductListResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int64(2), result.TotalCount)
	})

	t.Run("should return bad request for invalid query parameters", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products?page=invalid", nil)

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when provider not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		mockService.On("GetAllProducts", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.ProductFilter) bool {
			return true
		})).Return(nil, errors.New("provider not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products?provider_id=999", nil)

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		mockService.On("GetAllProducts", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.ProductFilter) bool {
			return true
		})).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products", nil)

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should get products with filters successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		response := &dto.ProductListResponse{
			Data: []dto.ProductResponse{
				{ID: 1, Name: "Product 1", Code: "PULSA10K", Category: "Mobile", PriceBasic: 9000, PriceSell: 10000},
			},
			TotalCount: 1,
			Page:       1,
			PageSize:   10,
		}

		mockService.On("GetAllProducts", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), mock.MatchedBy(func(f *dto.ProductFilter) bool {
			return f.ProviderID == 1
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/products?provider_id=1&type=PULSA&page=1&page_size=10", nil)

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductListResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result.Data, 1)
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	t.Run("should update product successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		req := testutil.CreateUpdateProductRequestFixture()
		response := &dto.ProductResponse{
			ID:         productID,
			ProviderID: 1,
			Name:       req.Name,
			Code:       "PULSA10K",
			Category:   req.Category,
			PriceBasic: req.PriceBasic,
			PriceSell:  req.PriceSell,
			IsActive:   req.IsActive,
			IconURL:    req.IconURL,
			Provider: &dto.ProviderInfo{
				ID:   1,
				Name: "Telkomsel",
				Code: "TSEL",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("UpdateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID, mock.MatchedBy(func(r *dto.UpdateProductRequest) bool {
			return r.Name == req.Name && r.PriceSell == req.PriceSell
		})).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, uint(1), result.ID)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, _ := setupProductHandler()

		req := testutil.CreateUpdateProductRequestFixture()
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/products/invalid", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(999)
		req := testutil.CreateUpdateProductRequestFixture()
		mockService.On("UpdateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID, mock.MatchedBy(func(r *dto.UpdateProductRequest) bool {
			return true
		})).Return(nil, errors.New("product not found"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/products/999", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		req := testutil.CreateUpdateProductRequestFixture()
		mockService.On("UpdateProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID, mock.MatchedBy(func(r *dto.UpdateProductRequest) bool {
			return true
		})).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	t.Run("should delete product successfully", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		mockService.On("DeleteProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/products/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.DeleteProduct(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "message")
		assert.Equal(t, "product deleted successfully", result["message"])
	})

	t.Run("should return not found when product not found", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(999)
		mockService.On("DeleteProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(errors.New("product not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/products/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		handler.DeleteProduct(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		handler, _ := setupProductHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/products/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.DeleteProduct(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		handler, mockService := setupProductHandler()

		productID := uint(1)
		mockService.On("DeleteProduct", mock.MatchedBy(func(ctx context.Context) bool {
			return true
		}), productID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/products/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		handler.DeleteProduct(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_RegisterRoutes(t *testing.T) {
	t.Run("should register all routes correctly", func(t *testing.T) {
		handler, _ := setupProductHandler()
		router := gin.New()
		api := router.Group("/api/v1")

		handler.RegisterRoutes(api)

		routes := router.Routes()
		expectedRoutes := []string{
			"POST /api/v1/products",
			"GET /api/v1/products",
			"GET /api/v1/products/:id",
			"GET /api/v1/products/code/:code",
			"PUT /api/v1/products/:id",
			"DELETE /api/v1/products/:id",
		}

		for _, expectedRoute := range expectedRoutes {
			found := false
			for _, route := range routes {
				if route.Method+" "+route.Path == expectedRoute {
					found = true
					break
				}
			}
			assert.True(t, found, "Route %s not found", expectedRoute)
		}
	})
}
