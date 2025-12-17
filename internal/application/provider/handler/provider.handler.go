package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/service"
	"go.uber.org/zap"
)

type ProviderHandler struct {
	service service.ProviderService
	logger  *zap.Logger
}

func NewProviderHandler(service service.ProviderService, logger *zap.Logger) *ProviderHandler {
	return &ProviderHandler{
		service: service,
		logger:  logger,
	}
}

// CreateProvider godoc
// @Summary Create a new provider
// @Description Create a new provider with name, code, and optional logo
// @Tags Provider
// @Accept json
// @Produce json
// @Param body body dto.CreateProviderRequest true "Create provider request"
// @Success 201 {object} dto.ProviderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers [post]
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req dto.CreateProviderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.logger.Warn("validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateProvider(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "provider with code already exists" {
			h.logger.Warn("provider creation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to create provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create provider"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProviderByID godoc
// @Summary Get provider by ID
// @Description Get a provider by its ID with all associated products
// @Tags Provider
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} dto.ProviderResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers/{id} [get]
func (h *ProviderHandler) GetProviderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("invalid provider id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	response, err := h.service.GetProviderByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "provider not found" {
			h.logger.Warn("provider not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}
		h.logger.Error("failed to get provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get provider"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProviderByCode godoc
// @Summary Get provider by code
// @Description Get a provider by its code with all associated products
// @Tags Provider
// @Accept json
// @Produce json
// @Param code path string true "Provider Code"
// @Success 200 {object} dto.ProviderResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers/code/{code} [get]
func (h *ProviderHandler) GetProviderByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		h.logger.Warn("provider code is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider code is required"})
		return
	}

	response, err := h.service.GetProviderByCode(c.Request.Context(), code)
	if err != nil {
		if err.Error() == "provider not found" {
			h.logger.Warn("provider not found", zap.String("code", code))
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}
		h.logger.Error("failed to get provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get provider"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllProviders godoc
// @Summary Get all providers
// @Description Get all providers with pagination (without products for performance)
// @Tags Provider
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{} "List of providers with pagination"
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers [get]
func (h *ProviderHandler) GetAllProviders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.logger.Warn("Invalid page number", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		h.logger.Warn("Invalid page size", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
		return
	}

	filter := &dto.ProviderFilter{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := h.service.GetAllProviders(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get all providers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get providers"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProvider godoc
// @Summary Update a provider
// @Description Update an existing provider
// @Tags Provider
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Param body body dto.UpdateProviderRequest true "Update provider request"
// @Success 200 {object} dto.ProviderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers/{id} [put]
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid provider ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	var req dto.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateProvider(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "provider not found" {
			h.logger.Warn("Provider not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}
		h.logger.Error("Failed to update provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update provider"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProvider godoc
// @Summary Delete a provider
// @Description Delete a provider by its ID
// @Tags Provider
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/providers/{id} [delete]
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("invalid provider id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	err = h.service.DeleteProvider(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "provider not found" {
			h.logger.Warn("provider not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}
		h.logger.Error("failed to delete provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider deleted successfully"})
}

// RegisterRoutes registers provider routes
func (h *ProviderHandler) RegisterRoutes(api *gin.RouterGroup) {
	providers := api.Group("/providers")
	{
		providers.POST("", h.CreateProvider)
		providers.GET("", h.GetAllProviders)
		providers.GET("/:id", h.GetProviderByID)
		providers.GET("/code/:code", h.GetProviderByCode)
		providers.PUT("/:id", h.UpdateProvider)
		providers.DELETE("/:id", h.DeleteProvider)
	}
}
