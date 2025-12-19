package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/service"
)

// WifiVoucherHandler handles HTTP requests for wifi vouchers.
type WifiVoucherHandler struct {
	service service.WifiVoucherService
}

// NewWifiVoucherHandler creates a new instance of WifiVoucherHandler.
func NewWifiVoucherHandler(svc service.WifiVoucherService) *WifiVoucherHandler {
	return &WifiVoucherHandler{service: svc}
}

// CreateWifiVoucher handles POST /wifi-vouchers
func (h *WifiVoucherHandler) CreateWifiVoucher(c *gin.Context) {
	var req dto.CreateWifiVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateWifiVoucher(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetWifiVoucherByID handles GET /wifi-vouchers/:id
func (h *WifiVoucherHandler) GetWifiVoucherByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := h.service.GetWifiVoucherByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetWifiVoucherByCode handles GET /wifi-vouchers/code/:code
func (h *WifiVoucherHandler) GetWifiVoucherByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	response, err := h.service.GetWifiVoucherByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllWifiVouchers handles GET /wifi-vouchers
func (h *WifiVoucherHandler) GetAllWifiVouchers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	code := c.Query("code")
	batchID := c.Query("batch_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	filter := &dto.WifiVoucherFilter{
		Status:   status,
		Code:     code,
		BatchID:  batchID,
		Page:     page,
		PageSize: pageSize,
	}

	response, err := h.service.GetAllWifiVouchers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateWifiVoucher handles PUT /wifi-vouchers/:id
func (h *WifiVoucherHandler) UpdateWifiVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req dto.UpdateWifiVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateWifiVoucher(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteWifiVoucher handles DELETE /wifi-vouchers/:id
func (h *WifiVoucherHandler) DeleteWifiVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.service.DeleteWifiVoucher(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wifi voucher deleted successfully"})
}

// SellWifiVoucher handles POST /wifi-vouchers/:id/sell
func (h *WifiVoucherHandler) SellWifiVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SellWifiVoucher(c.Request.Context(), uint(id), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UseWifiVoucher handles POST /wifi-vouchers/:id/use
func (h *WifiVoucherHandler) UseWifiVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := h.service.UseWifiVoucher(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all wifi voucher routes.
func (h *WifiVoucherHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/wifi-vouchers")
	{
		group.POST("", h.CreateWifiVoucher)
		group.GET("", h.GetAllWifiVouchers)
		group.GET("/:id", h.GetWifiVoucherByID)
		group.GET("/code/:code", h.GetWifiVoucherByCode)
		group.PUT("/:id", h.UpdateWifiVoucher)
		group.DELETE("/:id", h.DeleteWifiVoucher)
		group.POST("/:id/sell", h.SellWifiVoucher)
		group.POST("/:id/use", h.UseWifiVoucher)
	}
}
