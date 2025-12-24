package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WalletHandler struct {
	service service.WalletService
	jwt     *jwt.JWTManager
	logger  *zap.Logger
}

func NewWalletHandler(service service.WalletService, jwtManager *jwt.JWTManager, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{
		service: service,
		jwt:     jwtManager,
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
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /wallets [post]
func (h *WalletHandler) CreateWallet(ctx *gin.Context) {
	var req dto.CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.CreateWallet(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create wallet", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": wallet})
}

// GetWalletByUserID godoc
// @Summary Get wallet by user ID
// @Description Get a wallet by user ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Wallet details"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Router /wallets/by-user/{user_id} [get]
func (h *WalletHandler) GetWalletByUserID(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	wallet, err := h.service.GetWalletByUserID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("user_id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": wallet})
}

// GetWalletByID godoc
// @Summary Get wallet by wallet ID
// @Description Get a wallet by wallet ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} map[string]interface{} "Wallet details"
// @Failure 400 {object} map[string]interface{} "Invalid wallet ID"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Router /wallets/wallet/{id} [get]
func (h *WalletHandler) GetWalletByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.service.GetWalletByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err), zap.Uint("wallet_id", uint(id)))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": wallet})
}

// AddBalance godoc
// @Summary Add balance to wallet
// @Description Add funds to a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param req body dto.UpdateBalanceRequest true "Add balance request"
// @Success 200 {object} map[string]interface{} "Balance added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /wallets/by-user/{user_id}/add-balance [post]
func (h *WalletHandler) AddBalance(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.AddBalance(ctx.Request.Context(), uint(id), req.Amount, req.Description)
	if err != nil {
		h.logger.Error("Failed to add balance", zap.Error(err), zap.Uint("user_id", uint(id)))
		if err.Error() == "wallet not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add balance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": wallet})
}

// DeductBalance godoc
// @Summary Deduct balance from wallet
// @Description Deduct funds from a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param req body dto.UpdateBalanceRequest true "Deduct balance request"
// @Success 200 {object} map[string]interface{} "Balance deducted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 402 {object} map[string]interface{} "Insufficient balance"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /wallets/by-user/{user_id}/deduct-balance [post]
func (h *WalletHandler) DeductBalance(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.DeductBalance(ctx.Request.Context(), uint(id), req.Amount, req.Description)
	if err != nil {
		h.logger.Error("Failed to deduct balance", zap.Error(err), zap.Uint("user_id", uint(id)))
		if err.Error() == "wallet not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient balance" {
			ctx.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct balance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": wallet})
}

// Transfer godoc
// @Summary Transfer balance between wallets
// @Description Transfer funds from authenticated user's wallet to another wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Security Bearer
// @Param req body dto.TransferRequest true "Transfer request"
// @Success 200 {object} map[string]interface{} "Transfer completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 402 {object} map[string]interface{} "Insufficient balance"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /wallets/transfer [post]
func (h *WalletHandler) Transfer(ctx *gin.Context) {
	token, err := h.extractBearerToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Transfer(ctx.Request.Context(), token, &req)
	if err != nil {
		h.logger.Error("Failed to transfer", zap.Error(err), zap.Uint("to_user_id", req.ToUserID))
		switch err.Error() {
		case "invalid or expired token":
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case "source wallet not found", "destination wallet not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient balance":
			ctx.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
		case "cannot transfer to yourself", "amount must be positive":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// DeleteWallet godoc
// @Summary Delete a wallet
// @Description Delete a wallet by user ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Wallet deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "Wallet not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /wallets/by-user/{user_id} [delete]
func (h *WalletHandler) DeleteWallet(ctx *gin.Context) {
	idStr := ctx.Param("user_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.DeleteWallet(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete wallet", zap.Error(err), zap.Uint("user_id", uint(id)))
		if err.Error() == "wallet not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}

// extractBearerToken extracts and validates the bearer token from Authorization header
func (h *WalletHandler) extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("Missing authorization header")
		return "", gin.Error{Err: http.ErrAbortHandler, Meta: "Authorization header is required"}
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		h.logger.Warn("Invalid authorization header format")
		return "", gin.Error{Err: http.ErrAbortHandler, Meta: "Invalid authorization header format"}
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		h.logger.Warn("Empty token in authorization header")
		return "", gin.Error{Err: http.ErrAbortHandler, Meta: "Invalid authorization header"}
	}

	return token, nil
}

func (h *WalletHandler) RegisterRoutes(api *gin.RouterGroup) {
	wallets := api.Group("/wallets")
	{
		wallets.POST("", h.CreateWallet)
		wallets.GET("/by-user/:user_id", h.GetWalletByUserID)
		wallets.GET("/wallet/:id", h.GetWalletByID)
		wallets.POST("/by-user/:user_id/add-balance", h.AddBalance)
		wallets.POST("/by-user/:user_id/deduct-balance", h.DeductBalance)
		wallets.POST("/transfer", h.Transfer)
		wallets.DELETE("/by-user/:user_id", h.DeleteWallet)
	}
}
