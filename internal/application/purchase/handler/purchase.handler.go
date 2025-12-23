package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/service"
	"github.com/novriyantoAli/cn-wallet/internal/middleware"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"

	userSecurityService "github.com/novriyantoAli/cn-wallet/internal/application/user-security/service"

	"go.uber.org/zap"
)

type PurchaseHandler struct {
	service             service.PurchaseService
	userSecurityService userSecurityService.UserSecurityService
	jwt                 *jwt.JWTManager
	logger              *zap.Logger
}

func NewPurchaseHandler(
	svc service.PurchaseService,
	usrSecuritySvc userSecurityService.UserSecurityService,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) *PurchaseHandler {
	return &PurchaseHandler{
		service:             svc,
		userSecurityService: usrSecuritySvc,
		jwt:                 jwtManager,
		logger:              logger,
	}
}

// ProcessPurchase godoc
// @Summary Process a purchase
// @Description Process a purchase transaction with product ID and phone number
// @Tags purchases
// @Accept json
// @Produce json
// @Security Bearer
// @Param purchase body dto.PurchaseRequest true "Purchase request"
// @Success 200 {object} map[string]interface{} "Purchase processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient balance"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product or user not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /purchases [post]
func (h *PurchaseHandler) ProcessPurchase(ctx *gin.Context) {
	var req dto.PurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.extractBearerToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ProcessPurchase(ctx.Request.Context(), token, &req)
	if err != nil {
		h.logger.Error("Failed to process purchase", zap.Error(err))
		h.handlePurchaseError(ctx, err)
		return
	}

	statusCode := http.StatusOK
	if response.Status == "failed" {
		statusCode = http.StatusBadRequest
	}

	ctx.JSON(statusCode, gin.H{"data": response})
}

// ProcessPurchaseWifi godoc
// @Summary Process a wifi voucher purchase
// @Description Process a wifi voucher purchase transaction with provider ID and duration hours
// @Tags purchases
// @Accept json
// @Produce json
// @Security Bearer
// @Param purchase body dto.PurchaseWifiRequest true "Wifi purchase request"
// @Success 200 {object} map[string]interface{} "Wifi purchase processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient balance"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product or user not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /purchases/wifi [post]
func (h *PurchaseHandler) ProcessPurchaseWifi(ctx *gin.Context) {
	var req dto.PurchaseWifiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.extractBearerToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ProcessPurchaseWifi(ctx.Request.Context(), token, &req)
	if err != nil {
		h.logger.Error("Failed to process wifi purchase", zap.Error(err))
		h.handlePurchaseError(ctx, err)
		return
	}

	statusCode := http.StatusOK
	if response.Status == "failed" {
		statusCode = http.StatusBadRequest
	}

	ctx.JSON(statusCode, gin.H{"data": response})
}

// GetPurchaseHistory godoc
// @Summary Get purchase history
// @Description Get purchase history for a wallet with pagination
// @Tags purchases
// @Accept json
// @Produce json
// @Param wallet_id query int true "Wallet ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} dto.PurchaseHistoryList "Purchase history"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /purchases/history [get]
func (h *PurchaseHandler) GetPurchaseHistory(ctx *gin.Context) {
	walletID, err := h.parseWalletID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, pageSize := h.parsePagination(ctx)

	filter := &dto.PurchaseHistoryFilter{
		WalletID: walletID,
		Page:     page,
		PageSize: pageSize,
	}

	response, err := h.service.GetPurchaseHistory(ctx.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get purchase history", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get purchase history"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// extractBearerToken extracts and validates the bearer token from Authorization header
func (h *PurchaseHandler) extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("Missing authorization header")
		return "", &HandlerError{
			Message: "Authorization header is required",
			Code:    "MISSING_AUTH",
		}
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		h.logger.Warn("Invalid authorization header format")
		return "", &HandlerError{
			Message: "Invalid authorization header format",
			Code:    "INVALID_AUTH_FORMAT",
		}
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		h.logger.Warn("Empty token in authorization header")
		return "", &HandlerError{
			Message: "Invalid authorization header",
			Code:    "EMPTY_TOKEN",
		}
	}

	return token, nil
}

// parseWalletID parses and validates the wallet_id query parameter
func (h *PurchaseHandler) parseWalletID(ctx *gin.Context) (uint, error) {
	walletID := ctx.Query("wallet_id")
	walletIDUint, err := strconv.ParseUint(walletID, 10, 32)
	if err != nil || walletIDUint == 0 {
		h.logger.Error("Invalid wallet_id", zap.String("wallet_id", walletID), zap.Error(err))
		return 0, &HandlerError{
			Message: "Invalid wallet_id",
			Code:    "INVALID_WALLET_ID",
		}
	}
	return uint(walletIDUint), nil
}

// parsePagination parses page and page_size query parameters with defaults
func (h *PurchaseHandler) parsePagination(ctx *gin.Context) (int, int) {
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		pageSizeInt = 10
	}

	return pageInt, pageSizeInt
}

// handlePurchaseError handles purchase service errors with appropriate HTTP status codes
func (h *PurchaseHandler) handlePurchaseError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	switch {
	case errMsg == "invalid or expired token":
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
	case errMsg == "user not found":
		ctx.JSON(http.StatusNotFound, gin.H{"error": errMsg})
	case errMsg == "product not found":
		ctx.JSON(http.StatusNotFound, gin.H{"error": errMsg})
	case errMsg == "wallet not found":
		ctx.JSON(http.StatusNotFound, gin.H{"error": errMsg})
	case errMsg == "insufficient wallet balance":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	case errMsg == "product_id is required" || errMsg == "phone is required":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process purchase"})
	}
}

func (h *PurchaseHandler) RegisterRoutes(api *gin.RouterGroup) {
	purchases := api.Group("/purchases")
	{
		purchases.GET("/history", h.GetPurchaseHistory)
	}
	purchases.Use(middleware.JWTMiddleware(h.jwt, h.logger))

	purchase := api.Group("/purchase")
	{
		purchase.POST("", h.ProcessPurchase)
		purchase.POST("/wifi", h.ProcessPurchaseWifi)
	}
	purchase.Use(middleware.PINMiddleware(h.userSecurityService, h.jwt, h.logger))
}

// HandlerError represents a handler-level error with code and message
type HandlerError struct {
	Message string
	Code    string
}

func (e *HandlerError) Error() string {
	return e.Message
}
