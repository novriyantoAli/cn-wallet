package grpc

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment"
	paymentHandler "github.com/novriyantoAli/cn-wallet/internal/application/payment/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/user"
	userHandler "github.com/novriyantoAli/cn-wallet/internal/application/user/handler"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include domain modules
	ledger.Module,
	user.Module,
	payment.Module,

	// gRPC handlers
	fx.Provide(
		userHandler.NewUserGrpcHandler,
		paymentHandler.NewPaymentGrpcHandler,
		NewServer,
	),
)
