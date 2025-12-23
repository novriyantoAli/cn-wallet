package purchase

import (
	"context"

	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/service"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module provides all purchase domain dependencies
var Module = fx.Options(
	fx.Provide(
		// Provide ProviderClient implementation
		NewProviderClient,
		// Provide TransactionManager
		fx.Annotate(
			func(db *gorm.DB) database.TransactionManagerI {
				return database.NewTransactionManager(db)
			},
			fx.As(new(database.TransactionManagerI)),
		),
		service.NewPurchaseService,
		handler.NewPurchaseHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		// Provide ProviderClient implementation
		NewProviderClient,
		// Provide TransactionManager
		fx.Annotate(
			func(db *gorm.DB) database.TransactionManagerI {
				return database.NewTransactionManager(db)
			},
			fx.As(new(database.TransactionManagerI)),
		),
		service.NewPurchaseService,
	),
)

// NewProviderClient creates a new provider client
// This is a mock implementation that should be replaced with actual API calls
func NewProviderClient() service.ProviderClient {
	return &ProviderClientImpl{}
}

// ProviderClientImpl is a mock implementation of ProviderClient
type ProviderClientImpl struct{}

// ProcessPurchase processes a purchase with the provider API
func (c *ProviderClientImpl) ProcessPurchase(ctx context.Context, productCode string, phone string, reference string) (serialNumber string, err error) {
	// This is a mock implementation
	// In production, this would call the actual provider API
	return "SN-" + reference, nil
}
