package wallet

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"wallet",
	fx.Provide(
		repository.NewWalletRepository,
		service.NewWalletService,
		handler.NewWalletHandler,
	),
)

type WalletHandler struct {
	handler *handler.WalletHandler
}

func NewWalletModule(
	db *gorm.DB,
	logger *zap.Logger,
	h *handler.WalletHandler,
) *WalletHandler {
	return &WalletHandler{
		handler: h,
	}
}

func (m *WalletHandler) RegisterRoutes(api interface{}) {
	// This will be called from the server module
}
