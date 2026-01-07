package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
)

type TransferHandler struct {
	service service.TransferService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewTransferHandler(service service.TransferService, jwtManager *jwt.JWTManager, logger *zap.Logger) *TransferHandler {
	return &TransferHandler{
		service: service,
		jwt:     jwtManager,
		logger:  logger,
	}
}

// CreateTransfer godoc
// @Summary Create a new transfer
// @Description Create a new transfer from wallet or paylater to another user
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body dto.CreateTransferRequest true "Transfer creation request"
// @Success 201 {object} map[string]interface{} "Created transfer"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transfers [post]
func (h *TransferHandler) CreateTransfer(ctx *gin.Context) {
	var req dto.CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer, err := h.service.CreateTransfer(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create transfer", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": transfer, "message": "Transfer created successfully"})
}

// GetTransferByID godoc
// @Summary Get transfer by ID
// @Description Get a transfer by transfer ID
// @Tags transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} map[string]interface{} "Transfer details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Transfer not found"
// @Router /transfers/{id} [get]
func (h *TransferHandler) GetTransferByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	transfer, err := h.service.GetTransferByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get transfer", zap.Error(err), zap.Uint("id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transfer not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": transfer})
}

// GetTransfersByUserID godoc
// @Summary Get transfers by user ID (sent)
// @Description Get all transfers sent by a user
// @Tags transfers
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "List of transfers"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /transfers/by-user/{user_id} [get]
func (h *TransferHandler) GetTransfersByUserID(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	transfers, err := h.service.GetTransfersByUserID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get transfers by user ID", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": transfers})
}

// GetTransfersByTargetUserID godoc
// @Summary Get transfers by target user ID (received)
// @Description Get all transfers received by a user
// @Tags transfers
// @Accept json
// @Produce json
// @Param target_user_id path int true "Target User ID"
// @Success 200 {object} map[string]interface{} "List of transfers"
// @Failure 400 {object} map[string]interface{} "Invalid target user ID"
// @Router /transfers/by-target/{target_user_id} [get]
func (h *TransferHandler) GetTransfersByTargetUserID(ctx *gin.Context) {
	idStr := ctx.Param("target_user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	transfers, err := h.service.GetTransfersByTargetUserID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get transfers by target user ID", zap.Error(err), zap.Uint("target_user_id", uint(id)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": transfers})
}

// ListTransfers godoc
// @Summary List transfers with filters
// @Description List transfers with pagination and filters
// @Tags transfers
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by sender user ID"
// @Param target_user_id query int false "Filter by receiver user ID"
// @Param source query string false "Filter by source (wallet/paylater)"
// @Param status query string false "Filter by status"
// @Param from_date query string false "Filter by from date (YYYY-MM-DD)"
// @Param to_date query string false "Filter by to date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{} "Paginated transfer list"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /transfers [get]
func (h *TransferHandler) ListTransfers(ctx *gin.Context) {
	var req dto.TransferFilter
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ListTransfers(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list transfers", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateTransferStatus godoc
// @Summary Update transfer status
// @Description Update the status of a transfer
// @Tags transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param request body dto.UpdateTransferStatusRequest true "Update status request"
// @Success 200 {object} map[string]interface{} "Updated transfer details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Transfer not found"
// @Router /transfers/{id}/status [put]
func (h *TransferHandler) UpdateTransferStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateTransferStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer, err := h.service.UpdateTransferStatus(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update transfer status", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": transfer, "message": "Transfer status updated successfully"})
}

// GetUserTransferStats godoc
// @Summary Get user transfer statistics
// @Description Get transfer statistics for a specific user
// @Tags transfers
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User transfer statistics"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /transfers/stats/{user_id} [get]
func (h *TransferHandler) GetUserTransferStats(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.service.GetUserTransferStats(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get transfer stats", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}

// CancelTransfer godoc
// @Summary Cancel transfer
// @Description Cancel a pending transfer
// @Tags transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} map[string]interface{} "Transfer cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID or transfer cannot be cancelled"
// @Failure 404 {object} map[string]interface{} "Transfer not found"
// @Router /transfers/{id}/cancel [post]
func (h *TransferHandler) CancelTransfer(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.CancelTransfer(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to cancel transfer", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transfer cancelled successfully"})
}

// RegisterRoutes registers transfer routes
func (h *TransferHandler) RegisterRoutes(api *gin.RouterGroup) {
	transfers := api.Group("/transfers")
	{
		transfers.POST("", h.CreateTransfer)
		transfers.GET("", h.ListTransfers)
		transfers.GET("/:id", h.GetTransferByID)
		transfers.GET("/by-user/:user_id", h.GetTransfersByUserID)
		transfers.GET("/by-target/:target_user_id", h.GetTransfersByTargetUserID)
		transfers.GET("/stats/:user_id", h.GetUserTransferStats)
		transfers.PUT("/:id/status", h.UpdateTransferStatus)
		transfers.POST("/:id/cancel", h.CancelTransfer)
	}
}
