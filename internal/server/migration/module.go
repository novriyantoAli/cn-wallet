package migration

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"

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
	// Order matters: Provider before Product (foreign key dependency)
	err := s.db.AutoMigrate(
		&userEntity.User{},
		&entity.Payment{},
		&walletEntity.Wallet{},
		&providerEntity.Provider{},
		&productEntity.Product{},
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

	s.logger.Info("Data seeding completed successfully")
	return nil
}

func (s *Server) DropTables() error {
	s.logger.Warn("Dropping all database tables")

	// Drop in reverse order of creation (Product before Provider due to foreign key)
	err := s.db.Migrator().DropTable(
		&productEntity.Product{},
		&providerEntity.Provider{},
		&userEntity.User{},
		&entity.Payment{},
		&walletEntity.Wallet{},
	)
	if err != nil {
		s.logger.Error("Failed to drop database tables", zap.Error(err))
		return err
	}

	s.logger.Info("Database tables dropped successfully")
	return nil
}
