package api

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/payment"
	"github.com/novriyantoAli/cn-wallet/internal/application/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include all domain modules
	user.Module,
	payment.Module,

	// API api
	fx.Provide(NewServer),
)
