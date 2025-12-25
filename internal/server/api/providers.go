package api

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/oauth"
	"github.com/novriyantoAli/cn-wallet/internal/application/paylater"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment"
	"github.com/novriyantoAli/cn-wallet/internal/application/product"
	"github.com/novriyantoAli/cn-wallet/internal/application/provider"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase"
	"github.com/novriyantoAli/cn-wallet/internal/application/transaction"
	"github.com/novriyantoAli/cn-wallet/internal/application/transfer"
	"github.com/novriyantoAli/cn-wallet/internal/application/user"
	usersecurity "github.com/novriyantoAli/cn-wallet/internal/application/user-security"
	"github.com/novriyantoAli/cn-wallet/internal/application/wallet"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include all domain modules
	oauth.Module,
	user.Module,
	usersecurity.Module,
	payment.Module,
	paylater.Module,
	product.Module,
	provider.Module,
	purchase.Module,
	wallet.Module,
	wifivoucher.Module,
	transaction.Module,
	transfer.Module,

	// API api
	fx.Provide(NewServer),
)
