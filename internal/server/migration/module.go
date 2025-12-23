package migration

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	userSecurityEntity "github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewServer(db *gorm.DB, logger *zap.Logger) *Server {
	return &Server{
		db:     db,
		logger: logger,
	}
}

func (s *Server) RunMigrations() error {
	s.logger.Info("Starting database migrations")

	// Run auto migrations for all entities
	// Order matters: dependencies should be created before dependent tables
	err := s.db.AutoMigrate(
		&userEntity.User{},
		&userSecurityEntity.UserSecurity{},
		&entity.Payment{},
		&walletEntity.Wallet{},
		&providerEntity.Provider{},
		&productEntity.Product{},
		&transactionEntity.Transaction{},
		&wifiVoucherEntity.WifiVoucher{},
	)
	if err != nil {
		s.logger.Error("Failed to run database migrations", zap.Error(err))
		return err
	}

	s.logger.Info("Database migrations completed successfully")
	return nil
}

func (s *Server) SeedData() error {
	s.logger.Info("Starting data seeding")

	// Add any initial data seeding here
	// Example: Create default admin user, initial payment statuses, etc.

	// Seed data for provider table
	providers := []providerEntity.Provider{
		{
			Name: "CN Hotspot",
			Code: "CNHOTSPOT",
		},
	}
	err := s.db.Create(&providers).Error
	if err != nil {
		s.logger.Error("Failed to seed data for provider table", zap.Error(err))
		return err
	}

	// Seed data for products table
	products := []productEntity.Product{
		{
			ProviderID: 1,
			Name:       "CN Hotspot 1GB",
			Code:       "CNHOTSPOT1GB",
			Category:   "wifi",
			PriceBasic: 9000.00,
			PriceSell:  10000.00,
		},
	}
	err = s.db.Create(&products).Error
	if err != nil {
		s.logger.Error("Failed to seed data for products table", zap.Error(err))
		return err
	}
	s.logger.Info("Data seeding completed successfully")
	return nil
}

func (s *Server) DropTables() error {
	s.logger.Warn("Dropping all database tables")

	// Drop in reverse order of creation to handle foreign key constraints
	err := s.db.Migrator().DropTable(
		&productEntity.Product{},
		&providerEntity.Provider{},
		&walletEntity.Wallet{},
		&entity.Payment{},
		&userEntity.User{},
		&userSecurityEntity.UserSecurity{},
		&transactionEntity.Transaction{},
		&wifiVoucherEntity.WifiVoucher{},
	)
	if err != nil {
		s.logger.Error("Failed to drop database tables", zap.Error(err))
		return err
	}

	s.logger.Info("Database tables dropped successfully")
	return nil
}
