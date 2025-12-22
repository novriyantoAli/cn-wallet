package usersecurity

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/user-security/service"

	"go.uber.org/fx"
)

// Module provides all user-security domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewUserSecurityRepository,
		service.NewUserSecurityService,
		handler.NewUserSecurityHandler,
	),
)
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewUserSecurityRepository,
		service.NewUserSecurityService,
	),
)
