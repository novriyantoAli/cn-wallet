package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/service"
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	svc service.TransactionService
}

// NewTransactionHandler creates a new instance of TransactionHandler.
func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		svc: svc,
	}
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetTransactionByID handles GET /transactions/:id
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	result, err := h.svc.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAllTransactions handles GET /transactions
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	filter := &dto.TransactionFilter{
		WalletID: 0,
		Type:     c.Query("type"),
		Status:   c.Query("status"),
		Page:     page,
		PageSize: pageSize,
	}

	if w := c.Query("wallet_id"); w != "" {
		if parsed, err := strconv.ParseUint(w, 10, 32); err == nil {
			filter.WalletID = uint(parsed)
		}
	}

	result, err := h.svc.GetAllTransactions(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateTransaction handles PUT /transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.UpdateTransaction(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteTransaction handles DELETE /transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	if err := h.svc.DeleteTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// GetWalletTransactions handles GET /wallets/:wallet_id/transactions
func (h *TransactionHandler) GetWalletTransactions(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("wallet_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	result, err := h.svc.GetWalletTransactions(c.Request.Context(), uint(walletID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RegisterRoutes registers transaction routes.
func (h *TransactionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/transactions", h.CreateTransaction)
	rg.GET("/transactions", h.GetAllTransactions)
	rg.GET("/transactions/:id", h.GetTransactionByID)
	rg.PUT("/transactions/:id", h.UpdateTransaction)
	rg.DELETE("/transactions/:id", h.DeleteTransaction)
	rg.GET("/wallets/:wallet_id/transactions", h.GetWalletTransactions)
}
