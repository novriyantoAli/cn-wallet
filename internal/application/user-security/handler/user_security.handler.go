package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/service"

	"go.uber.org/zap"
)

type UserSecurityHandler struct {
	service service.UserSecurityService
	logger  *zap.Logger
}

func NewUserSecurityHandler(service service.UserSecurityService, logger *zap.Logger) *UserSecurityHandler {
	return &UserSecurityHandler{
		service: service,
		logger:  logger,
	}
}

// SetPIN godoc
// @Summary Set PIN for user
// @Description Set or update PIN for user account
// @Tags user-security
// @Accept json
// @Produce json
// @Param request body dto.SetPINRequest true "Set PIN request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-security/pin [post]
func (h *UserSecurityHandler) SetPIN(c *gin.Context) {
	var req dto.SetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetPIN(c.Request.Context(), &req); err != nil {
		h.logger.Error("Failed to set PIN", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN set successfully"})
}

// VerifyPIN godoc
// @Summary Verify PIN for user
// @Description Verify PIN for user account
// @Tags user-security
// @Accept json
// @Produce json
// @Param request body dto.VerifyPINRequest true "Verify PIN request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-security/pin/verify [post]
func (h *UserSecurityHandler) VerifyPIN(c *gin.Context) {
	var req dto.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verified, err := h.service.VerifyPIN(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to verify PIN", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": verified})
}

// GetSecurity godoc
// @Summary Get user security info
// @Description Get user security information including lock status
// @Tags user-security
// @Produce json
// @Param user_id query uint true "User ID"
// @Success 200 {object} dto.UserSecurityResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-security [get]
func (h *UserSecurityHandler) GetSecurity(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || userIDUint == 0 {
		h.logger.Error("Invalid user_id", zap.String("user_id", userIDStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	resp, err := h.service.GetSecurity(c.Request.Context(), uint(userIDUint))
	if err != nil {
		h.logger.Error("Failed to get security info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get security info"})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "user not set security"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserSecurityHandler) RegisterRoutes(api *gin.RouterGroup) {
	security := api.Group("/user-security")
	{
		security.POST("/pin", h.SetPIN)
		security.POST("/pin/verify", h.VerifyPIN)
		security.GET("", h.GetSecurity)
	}
}
