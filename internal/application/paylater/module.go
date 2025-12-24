package paylater

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater/service"

	"go.uber.org/fx"
)

// Module provides all paylater domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewPaylaterAccountRepository,
		service.NewPaylaterAccountService,
		handler.NewPaylaterAccountHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewPaylaterAccountRepository,
		service.NewPaylaterAccountService,
	),
)
