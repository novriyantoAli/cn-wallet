package provider

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider/service"

	"go.uber.org/fx"
)

// Module provides all provider domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewProviderRepository,
		service.NewProviderService,
		handler.NewProviderHandler,
	),
)
