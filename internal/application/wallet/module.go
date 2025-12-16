package wallet

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"wallet",
	fx.Provide(
		repository.NewWalletRepository,
		service.NewWalletService,
		handler.NewWalletHandler,
	),
)
