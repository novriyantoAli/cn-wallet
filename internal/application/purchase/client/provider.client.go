package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// ProviderClient handles external provider API calls for purchases
type ProviderClient struct {
	logger *zap.Logger
}

// NewProviderClient creates a new provider client instance
func NewProviderClient(logger *zap.Logger) *ProviderClient {
	return &ProviderClient{
		logger: logger,
	}
}

// ProcessPurchase calls the provider API to process a purchase
// Returns the serial number or error
func (c *ProviderClient) ProcessPurchase(ctx context.Context, productCode string, phone string, reference string) (string, error) {
	// TODO: Implement actual provider API integration
	// For now, this is a stub that simulates success
	c.logger.Info("Processing purchase with provider",
		zap.String("product_code", productCode),
		zap.String("phone", phone),
		zap.String("reference", reference))

	// Simulate serial number generation
	serialNumber := fmt.Sprintf("SN-%s-%s", productCode, reference[:8])

	return serialNumber, nil
}
