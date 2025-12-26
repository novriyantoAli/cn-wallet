package worker

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment"
	"github.com/novriyantoAli/cn-wallet/internal/application/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include domain worker modules
	ledger.WorkerModule,
	payment.WorkerModule,
	user.WorkerModule,

	// Worker api
	fx.Provide(NewServer),
)
