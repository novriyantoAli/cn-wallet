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

type PaylaterRepaymentHandler struct {
	service service.PaylaterRepaymentService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewPaylaterRepaymentHandler(service service.PaylaterRepaymentService, jwtManager *jwt.JWTManager, logger *zap.Logger) *PaylaterRepaymentHandler {
	return &PaylaterRepaymentHandler{
		service: service,
		jwt:     jwtManager,
		logger:  logger,
	}
}

// ProcessRepayment godoc
// @Summary Process a paylater loan repayment
// @Description Process a repayment for a paylater loan with ledger entries and loan status update
// @Tags paylater-repayment
// @Accept json
// @Produce json
// @Param request body dto.CreatePaylaterRepaymentRequest true "Repayment request"
// @Success 201 {object} map[string]interface{} "Repayment processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or business logic error"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/repayments [post]
func (h *PaylaterRepaymentHandler) ProcessRepayment(ctx *gin.Context) {
	var req dto.CreatePaylaterRepaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ProcessRepayment(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to process repayment", zap.Error(err), zap.Uint("loan_id", req.PaylaterLoanID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    result,
		"message": "Repayment processed successfully",
	})
}

// GetRepaymentByID godoc
// @Summary Get a repayment by ID
// @Description Retrieve a specific repayment record by its ID
// @Tags paylater-repayment
// @Accept json
// @Produce json
// @Param id path int true "Repayment ID"
// @Success 200 {object} map[string]interface{} "Repayment details"
// @Failure 400 {object} map[string]interface{} "Invalid repayment ID"
// @Failure 404 {object} map[string]interface{} "Repayment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/repayments/{id} [get]
func (h *PaylaterRepaymentHandler) GetRepaymentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repayment ID"})
		return
	}

	result, err := h.service.GetRepaymentByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get repayment", zap.Error(err), zap.Uint("id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Repayment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// ListRepaymentsByLoan godoc
// @Summary List all repayments for a loan
// @Description Retrieve all repayment records for a specific paylater loan
// @Tags paylater-repayment
// @Accept json
// @Produce json
// @Param loan_id path int true "Loan ID"
// @Success 200 {object} map[string]interface{} "List of repayments"
// @Failure 400 {object} map[string]interface{} "Invalid loan ID"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/loans/{loan_id}/repayments [get]
func (h *PaylaterRepaymentHandler) ListRepaymentsByLoan(ctx *gin.Context) {
	loanIDStr := ctx.Param("loan_id")
	loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	results, err := h.service.ListRepaymentsByLoan(ctx.Request.Context(), uint(loanID))
	if err != nil {
		h.logger.Error("Failed to list repayments for loan", zap.Error(err), zap.Uint("loan_id", uint(loanID)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"count": len(results),
	})
}

// GetRepaymentsByUser godoc
// @Summary Get all repayments for a user
// @Description Retrieve all repayment records made by a specific user
// @Tags paylater-repayment
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "List of user's repayments"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/users/{user_id}/repayments [get]
func (h *PaylaterRepaymentHandler) GetRepaymentsByUser(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	results, err := h.service.GetRepaymentsByUser(ctx.Request.Context(), uint(userID))
	if err != nil {
		h.logger.Error("Failed to get repayments for user", zap.Error(err), zap.Uint("user_id", uint(userID)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"count": len(results),
	})
}

// RegisterRoutes registers paylater repayment routes
func (h *PaylaterRepaymentHandler) RegisterRoutes(api *gin.RouterGroup) {
	repayments := api.Group("/paylater/repayments")
	{
		repayments.POST("", h.ProcessRepayment)
		repayments.GET("/:id", h.GetRepaymentByID)
	}

	loans := api.Group("/paylater/loans")
	{
		loans.GET("/:loan_id/repayments", h.ListRepaymentsByLoan)
	}

	users := api.Group("/paylater/users")
	{
		users.GET("/:user_id/repayments", h.GetRepaymentsByUser)
	}
}
