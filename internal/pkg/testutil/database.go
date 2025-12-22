package testutil

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	securityEntity "github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all entities
	err = db.AutoMigrate(
		&userEntity.User{},
		&entity.Payment{},
		&walletEntity.Wallet{},
		&providerEntity.Provider{},
		&productEntity.Product{},
		&wifiVoucherEntity.WifiVoucher{},
		&transactionEntity.Transaction{},
		&securityEntity.UserSecurity{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CleanDB cleans all data from test database
func CleanDB(db *gorm.DB) error {
	// Delete in reverse order of dependencies (Product before Provider)
	if err := db.Exec("DELETE FROM wifi_vouchers").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM products").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wallets").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM payments").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM providers").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	return nil
}
