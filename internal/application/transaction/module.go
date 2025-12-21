package transaction

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction/service"

	"go.uber.org/fx"
)

// Module provides all transaction domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewTransactionRepository,
		service.NewTransactionService,
		handler.NewTransactionHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewTransactionRepository,
		service.NewTransactionService,
	),
)
