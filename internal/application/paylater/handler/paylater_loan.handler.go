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

type PaylaterLoanHandler struct {
	service service.PaylaterLoanService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewPaylaterLoanHandler(service service.PaylaterLoanService, jwtManager *jwt.JWTManager, logger *zap.Logger) *PaylaterLoanHandler {
	return &PaylaterLoanHandler{
		service: service,
		jwt:     jwtManager,
		logger:  logger,
	}
}

// CreateLoan godoc
// @Summary Create a new paylater loan
// @Description Create a new paylater loan for a user
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param loan body dto.CreatePaylaterLoanRequest true "Loan creation request"
// @Success 201 {object} map[string]interface{} "Created loan"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/loans [post]
func (h *PaylaterLoanHandler) CreateLoan(ctx *gin.Context) {
	var req dto.CreatePaylaterLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.service.CreateLoan(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create paylater loan", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": loan})
}

// GetLoanByID godoc
// @Summary Get paylater loan by ID
// @Description Get a paylater loan by loan ID
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{} "Loan details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Router /paylater/loans/{id} [get]
func (h *PaylaterLoanHandler) GetLoanByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	loan, err := h.service.GetLoanByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get paylater loan", zap.Error(err), zap.Uint("id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loan})
}

// GetLoansByUserID godoc
// @Summary Get loans by user ID
// @Description Get all loans for a specific user
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "List of loans"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /paylater/loans/by-user/{user_id} [get]
func (h *PaylaterLoanHandler) GetLoansByUserID(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	loans, err := h.service.GetLoansByUserID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get loans by user ID", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loans})
}

// ListLoans godoc
// @Summary List paylater loans with filters
// @Description List paylater loans with pagination and filters
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Param source query string false "Filter by source"
// @Param from_date query string false "Filter by from date (YYYY-MM-DD)"
// @Param to_date query string false "Filter by to date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{} "Paginated loan list"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /paylater/loans [get]
func (h *PaylaterLoanHandler) ListLoans(ctx *gin.Context) {
	var req dto.ListPaylaterLoansRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ListLoans(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list loans", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateLoanStatus godoc
// @Summary Update loan status
// @Description Update the status of a paylater loan
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param request body dto.UpdateLoanStatusRequest true "Update status request"
// @Success 200 {object} map[string]interface{} "Updated loan details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Router /paylater/loans/{id}/status [put]
func (h *PaylaterLoanHandler) UpdateLoanStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateLoanStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.service.UpdateLoanStatus(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update loan status", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loan, "message": "Loan status updated successfully"})
}

// MarkLoanAsPaid godoc
// @Summary Mark loan as paid
// @Description Mark a loan as paid and restore credit to the account
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{} "Loan marked as paid"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Router /paylater/loans/{id}/mark-paid [post]
func (h *PaylaterLoanHandler) MarkLoanAsPaid(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	loan, err := h.service.MarkLoanAsPaid(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to mark loan as paid", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loan, "message": "Loan marked as paid successfully"})
}

// ProcessOverdueLoans godoc
// @Summary Process overdue loans
// @Description Process all overdue loans and update their status
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Number of loans processed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /paylater/loans/process-overdue [post]
func (h *PaylaterLoanHandler) ProcessOverdueLoans(ctx *gin.Context) {
	count, err := h.service.ProcessOverdueLoans(ctx.Request.Context())
	if err != nil {
		h.logger.Error("Failed to process overdue loans", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Overdue loans processed successfully",
		"count":   count,
	})
}

// GetUserLoanStats godoc
// @Summary Get user loan statistics
// @Description Get loan statistics for a specific user
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User loan statistics"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /paylater/loans/stats/{user_id} [get]
func (h *PaylaterLoanHandler) GetUserLoanStats(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.service.GetUserLoanStats(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get loan stats", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}

// DeleteLoan godoc
// @Summary Delete paylater loan
// @Description Delete a paylater loan by ID
// @Tags paylater-loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{} "Loan deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID or loan cannot be deleted"
// @Failure 404 {object} map[string]interface{} "Loan not found"
// @Router /paylater/loans/{id} [delete]
func (h *PaylaterLoanHandler) DeleteLoan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteLoan(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete paylater loan", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}

// RegisterRoutes registers paylater loan routes
func (h *PaylaterLoanHandler) RegisterRoutes(api *gin.RouterGroup) {
	loans := api.Group("/paylater/loans")
	{
		loans.POST("", h.CreateLoan)
		loans.GET("", h.ListLoans)
		loans.GET("/:id", h.GetLoanByID)
		loans.GET("/by-user/:user_id", h.GetLoansByUserID)
		loans.GET("/stats/:user_id", h.GetUserLoanStats)
		loans.PUT("/:id/status", h.UpdateLoanStatus)
		loans.POST("/:id/mark-paid", h.MarkLoanAsPaid)
		loans.POST("/process-overdue", h.ProcessOverdueLoans)
		loans.DELETE("/:id", h.DeleteLoan)
	}
}
