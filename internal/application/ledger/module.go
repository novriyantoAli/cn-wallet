package ledger

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/service"

	"go.uber.org/fx"
)

// Module provides all ledger domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewLedgerRepository,
		service.NewLedgerService,
		handler.NewLedgerHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewLedgerRepository,
		service.NewLedgerService,
	),
)
