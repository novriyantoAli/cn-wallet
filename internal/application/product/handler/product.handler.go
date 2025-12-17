package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/service"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service service.ProductService
	logger  *zap.Logger
}

func NewProductHandler(service service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for a provider
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create product request"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "provider not found" || err.Error() == "code already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	response, err := h.service.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to get product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductByCode godoc
// @Summary Get product by code
// @Description Get a product by its code
// @Tags products
// @Produce json
// @Param code path string true "Product code"
// @Success 200 {object} dto.ProductResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/code/{code} [get]
func (h *ProductHandler) GetProductByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product code is required"})
		return
	}

	response, err := h.service.GetProductByCode(c.Request.Context(), code)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to get product by code", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllProducts godoc
// @Summary List products
// @Description Get all products with optional filtering and pagination
// @Tags products
// @Produce json
// @Param provider_id query int false "Filter by provider ID"
// @Param type query string false "Filter by product type (PULSA|DATA|GAME|PLN)"
// @Param code query string false "Filter by product code (partial match)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	var filter dto.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.GetAllProducts(c.Request.Context(), &filter)
	if err != nil {
		if err.Error() == "provider not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to get products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (partial update allowed)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Update product request"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateProduct(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	err = h.service.DeleteProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// RegisterRoutes registers all product routes
func (h *ProductHandler) RegisterRoutes(group *gin.RouterGroup) {
	products := group.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.GetAllProducts)
		products.GET("/:id", h.GetProductByID)
		products.GET("/code/:code", h.GetProductByCode)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}
