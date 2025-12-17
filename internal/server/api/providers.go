package api

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/oauth"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider"
	"github.com/novriyantoAli/cn-wallet/internal/application/user"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include all domain modules
	oauth.Module,
	user.Module,
	payment.Module,
	provider.Module,
	wallet.Module,

	// API api
	fx.Provide(NewServer),
)
