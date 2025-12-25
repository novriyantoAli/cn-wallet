package transfer

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer/service"
	"go.uber.org/fx"
)

var Module = fx.Module("transfer",
	fx.Provide(
		repository.NewTransferRepository,
		service.NewTransferService,
		handler.NewTransferHandler,
	),
)
