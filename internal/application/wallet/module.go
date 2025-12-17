package wallet

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet/service"

	"go.uber.org/fx"
)

// Module provides all wallet domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewWalletRepository,
		service.NewWalletService,
		handler.NewWalletHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewWalletRepository,
		service.NewWalletService,
	),
)
