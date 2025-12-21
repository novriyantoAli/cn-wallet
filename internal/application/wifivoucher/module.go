package wifivoucher

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/service"

	"go.uber.org/fx"
)

// Module provides all wifi voucher domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewWifiVoucherRepository,
		service.NewWifiVoucherService,
		handler.NewWifiVoucherHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewWifiVoucherRepository,
		service.NewWifiVoucherService,
	),
)
