package handler

import (
	"net/http"
	"strconv"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/service"
	"github.com/shopspring/decimal"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WalletHandler struct {
	service service.WalletService
	logger  *zap.Logger
}

func NewWalletHandler(service service.WalletService, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{
		service: service,
		logger:  logger,
	}
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet for a user
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body dto.CreateWalletRequest true "Wallet creation request"
// @Success 201 {object} map[string]interface{} "Created wallet"
// @Router /wallets [post]
func (h *WalletHandler) CreateWallet(ctx *gin.Context) {
	var req dto.CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.CreateWallet(&req)
	if err != nil {
		h.logger.Error("Failed to create wallet", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": wallet})
}

// GetWallet godoc
// @Summary Get wallet by ID
// @Description Get a wallet by its ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Router /wallets/{id} [get]
func (h *WalletHandler) GetWallet(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.service.GetWalletByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": wallet})
}

// GetWallets godoc
// @Summary Get all wallets
// @Description Get a list of wallets with pagination
// @Tags wallets
// @Accept json
// @Produce json
// @Router /wallets [get]
func (h *WalletHandler) GetWallets(ctx *gin.Context) {
	var filter dto.WalletFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallets, err := h.service.GetWallets(&filter)
	if err != nil {
		h.logger.Error("Failed to get wallets", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallets"})
		return
	}

	ctx.JSON(http.StatusOK, wallets)
}

// AddBalance godoc
// @Summary Add balance to wallet
// @Description Add balance to a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Router /wallets/{id}/add-balance [post]
func (h *WalletHandler) AddBalance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	var req dto.AddBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	err = h.service.AddBalance(uint(id), amount)
	if err != nil {
		h.logger.Error("Failed to add balance", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add balance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Balance added successfully"})
}

// WithdrawBalance godoc
// @Summary Withdraw balance from wallet
// @Description Withdraw balance from a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Router /wallets/{id}/withdraw [post]
func (h *WalletHandler) WithdrawBalance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	var req dto.WithdrawBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	err = h.service.WithdrawBalance(uint(id), amount, req.PIN)
	if err != nil {
		h.logger.Error("Failed to withdraw balance", zap.Error(err))
		if err.Error() == "insufficient balance" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to withdraw balance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Balance withdrawn successfully"})
}

// RegisterRoutes registers wallet routes
func (h *WalletHandler) RegisterRoutes(api *gin.RouterGroup) {
	wallets := api.Group("/wallets")
	{
		wallets.POST("", h.CreateWallet)
		wallets.GET("", h.GetWallets)
		wallets.GET("/:id", h.GetWallet)
		wallets.POST("/:id/add-balance", h.AddBalance)
		wallets.POST("/:id/withdraw", h.WithdrawBalance)
	}
}
