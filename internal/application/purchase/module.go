package purchase

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/service"

	"go.uber.org/fx"
)

// Module provides all purchase domain dependencies
var Module = fx.Options(
	fx.Provide(
		service.NewPurchaseService,
		handler.NewPurchaseHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		service.NewPurchaseService,
	),
)
