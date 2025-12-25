package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
)

type LedgerHandler struct {
	service service.LedgerService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewLedgerHandler(service service.LedgerService, jwtManager *jwt.JWTManager, logger *zap.Logger) *LedgerHandler {
	return &LedgerHandler{
		service: service,
		jwt:     jwtManager,
		logger:  logger,
	}
}

// CreateEntry godoc
// @Summary Create a new ledger entry
// @Description Create a new append-only ledger entry
// @Tags ledger
// @Accept json
// @Produce json
// @Param entry body dto.CreateLedgerEntryRequest true "Ledger entry creation request"
// @Success 201 {object} map[string]interface{} "Created ledger entry"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ledger [post]
func (h *LedgerHandler) CreateEntry(ctx *gin.Context) {
	var req dto.CreateLedgerEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.CreateEntry(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create ledger entry", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": entry, "message": "Ledger entry created successfully"})
}

// GetEntryByID godoc
// @Summary Get ledger entry by ID
// @Description Get a ledger entry by entry ID
// @Tags ledger
// @Accept json
// @Produce json
// @Param id path int true "Entry ID"
// @Success 200 {object} map[string]interface{} "Ledger entry details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Entry not found"
// @Router /ledger/{id} [get]
func (h *LedgerHandler) GetEntryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	entry, err := h.service.GetEntryByID(ctx.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get ledger entry", zap.Error(err), zap.Uint64("id", id))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": entry})
}

// GetEntriesByUserID godoc
// @Summary Get ledger entries by user ID
// @Description Get all ledger entries for a specific user
// @Tags ledger
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "List of ledger entries"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /ledger/user/{user_id} [get]
func (h *LedgerHandler) GetEntriesByUserID(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	entries, err := h.service.GetEntriesByUserID(ctx.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get ledger entries by user ID", zap.Error(err), zap.Uint64("user_id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": entries})
}

// GetEntriesByReference godoc
// @Summary Get ledger entries by reference
// @Description Get all ledger entries for a specific reference
// @Tags ledger
// @Accept json
// @Produce json
// @Param reference_type path string true "Reference type (transfer/paylater/payment/repayment/adjustment)"
// @Param reference_id path int true "Reference ID"
// @Success 200 {object} map[string]interface{} "List of ledger entries"
// @Failure 400 {object} map[string]interface{} "Invalid reference"
// @Router /ledger/reference/{reference_type}/{reference_id} [get]
func (h *LedgerHandler) GetEntriesByReference(ctx *gin.Context) {
	refType := ctx.Param("reference_type")
	idStr := ctx.Param("reference_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reference ID"})
		return
	}

	entries, err := h.service.GetEntriesByReference(ctx.Request.Context(), refType, id)
	if err != nil {
		h.logger.Error("Failed to get ledger entries by reference", zap.Error(err), zap.String("reference_type", refType), zap.Uint64("reference_id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": entries})
}

// ListEntries godoc
// @Summary List ledger entries with filters
// @Description List ledger entries with pagination and filters
// @Tags ledger
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Param reference_type query string false "Filter by reference type"
// @Param account_type query string false "Filter by account type"
// @Param from_date query string false "Filter by from date (YYYY-MM-DD)"
// @Param to_date query string false "Filter by to date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{} "Paginated ledger entries"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /ledger [get]
func (h *LedgerHandler) ListEntries(ctx *gin.Context) {
	var req dto.ListLedgerEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ListEntries(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list ledger entries", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetUserStats godoc
// @Summary Get user ledger statistics
// @Description Get ledger statistics for a specific user
// @Tags ledger
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User ledger statistics"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Router /ledger/stats/{user_id} [get]
func (h *LedgerHandler) GetUserStats(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.service.GetUserStats(ctx.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get ledger stats", zap.Error(err), zap.Uint64("user_id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}

// RegisterRoutes registers ledger routes
func (h *LedgerHandler) RegisterRoutes(api *gin.RouterGroup) {
	ledger := api.Group("/ledger")
	{
		ledger.POST("", h.CreateEntry)
		ledger.GET("", h.ListEntries)
		ledger.GET("/:id", h.GetEntryByID)
		ledger.GET("/user/:user_id", h.GetEntriesByUserID)
		ledger.GET("/stats/:user_id", h.GetUserStats)
		ledger.GET("/reference/:reference_type/:reference_id", h.GetEntriesByReference)
	}
}
