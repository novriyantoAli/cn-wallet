package product

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/product/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/product/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewProductRepository,
		service.NewProductService,
		handler.NewProductHandler,
	),
)

var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewProductRepository,
		service.NewProductService,
	),
)
