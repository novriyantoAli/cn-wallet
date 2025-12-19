package wifivoucher

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module represents the wifi voucher module.
type Module struct {
	Handler handler.WifiVoucherHandler
	Service service.WifiVoucherService
	Repo    repository.WifiVoucherRepository
}

// NewModule creates and initializes a new wifi voucher module.
func NewModule(db *gorm.DB, logger *zap.Logger) Module {
	repo := repository.NewWifiVoucherRepository(db, logger)
	svc := service.NewWifiVoucherService(repo, logger)
	hdlr := handler.NewWifiVoucherHandler(svc)

	return Module{
		Handler: *hdlr,
		Service: svc,
		Repo:    repo,
	}
}
