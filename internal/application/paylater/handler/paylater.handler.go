package handler

import (
	"net/http"
	"strconv"

	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaylaterAccountHandler struct {
	service service.PaylaterAccountService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewPaylaterAccountHandler(service service.PaylaterAccountService, jwtManager *jwt.JWTManager, logger *zap.Logger) *PaylaterAccountHandler {
	return &PaylaterAccountHandler{
		service: service,
		jwt:     jwtManager,
		logger:  logger,
	}
}

// CreateAccount godoc
// @Summary Create a new paylater account
// @Description Create a new paylater account for a user
// @Tags paylater
// @Accept json
// @Produce json
// @Param account body dto.CreatePaylaterAccountRequest true "Paylater account creation request"
// @Success 201 {object} map[string]interface{} "Created paylater account"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater [post]
func (h *PaylaterAccountHandler) CreateAccount(ctx *gin.Context) {
	var req dto.CreatePaylaterAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.CreateAccount(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create paylater account", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": account})
}

// GetAccountByUserID godoc
// @Summary Get paylater account by user ID
// @Description Get a paylater account by user ID
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Paylater account details"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id} [get]
func (h *PaylaterAccountHandler) GetAccountByUserID(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	account, err := h.service.GetAccountByUserID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get paylater account", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": account})
}

// GetAccountByID godoc
// @Summary Get paylater account by ID
// @Description Get a paylater account by account ID
// @Tags paylater
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} map[string]interface{} "Paylater account details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/{id} [get]
func (h *PaylaterAccountHandler) GetAccountByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	account, err := h.service.GetAccountByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get paylater account", zap.Error(err), zap.Uint("id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": account})
}

// UpdateCreditLimit godoc
// @Summary Update paylater account credit limit
// @Description Update the credit limit of a paylater account
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.UpdateCreditLimitRequest true "Update credit limit request"
// @Success 200 {object} map[string]interface{} "Updated account details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id}/credit-limit [put]
func (h *PaylaterAccountHandler) UpdateCreditLimit(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateCreditLimitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.UpdateCreditLimit(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update credit limit", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": account, "message": "Credit limit updated successfully"})
}

// UpdateStatus godoc
// @Summary Update paylater account status
// @Description Update the status of a paylater account
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.UpdateStatusRequest true "Update status request"
// @Success 200 {object} map[string]interface{} "Updated account details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id}/status [put]
func (h *PaylaterAccountHandler) UpdateStatus(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.UpdateStatus(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update status", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": account, "message": "Status updated successfully"})
}

// UseCredit godoc
// @Summary Use credit from paylater account
// @Description Use credit from a paylater account
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.UseCreditRequest true "Use credit request"
// @Success 200 {object} map[string]interface{} "Credit used successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id}/use-credit [post]
func (h *PaylaterAccountHandler) UseCredit(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UseCreditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UseCredit(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to use credit", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// Repayment godoc
// @Summary Repay outstanding balance
// @Description Process repayment for outstanding balance
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.RepaymentRequest true "Repayment request"
// @Success 200 {object} map[string]interface{} "Repayment processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id}/repayment [post]
func (h *PaylaterAccountHandler) Repayment(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.RepaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Repayment(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to process repayment", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// DeleteAccount godoc
// @Summary Delete paylater account
// @Description Delete a paylater account by user ID
// @Tags paylater
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Account deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID or account has outstanding balance"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Router /paylater/by-user/{user_id} [delete]
func (h *PaylaterAccountHandler) DeleteAccount(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.DeleteAccount(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete paylater account", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// RegisterRoutes registers paylater account routes
func (h *PaylaterAccountHandler) RegisterRoutes(api *gin.RouterGroup) {
	paylater := api.Group("/paylater")
	{
		paylater.POST("", h.CreateAccount)
		paylater.GET("/by-user/:user_id", h.GetAccountByUserID)
		paylater.GET("/:id", h.GetAccountByID)
		paylater.PUT("/by-user/:user_id/credit-limit", h.UpdateCreditLimit)
		paylater.PUT("/by-user/:user_id/status", h.UpdateStatus)
		paylater.POST("/by-user/:user_id/use-credit", h.UseCredit)
		paylater.POST("/by-user/:user_id/repayment", h.Repayment)
		paylater.DELETE("/by-user/:user_id", h.DeleteAccount)
	}
}
