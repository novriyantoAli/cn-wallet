package purchase

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/client"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/fx"
)

// Module provides all purchase domain dependencies
var Module = fx.Options(
	fx.Provide(
		client.NewProviderClient,
		// Provide the ProviderClient interface
		func(c *client.ProviderClient) service.ProviderClient {
			return c
		},
		database.NewTransactionManager,
		service.NewPurchaseService,
		handler.NewPurchaseHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		client.NewProviderClient,
		func(c *client.ProviderClient) service.ProviderClient {
			return c
		},
		database.NewTransactionManager,
		service.NewPurchaseService,
	),
)
