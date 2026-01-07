package testutil

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	ledgerEntity "github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
	paylaterEntity "github.com/novriyantoAli/cn-wallet/internal/application/paylater/entity"
	"github.com/novriyantoAli/cn-wallet/internal/application/payment/entity"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	providerEntity "github.com/novriyantoAli/cn-wallet/internal/application/provider/entity"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	transferEntity "github.com/novriyantoAli/cn-wallet/internal/application/transfer/entity"
	securityEntity "github.com/novriyantoAli/cn-wallet/internal/application/user-security/entity"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	walletEntity "github.com/novriyantoAli/cn-wallet/internal/application/wallet/entity"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
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
		&paylaterEntity.PaylaterAccount{},
		&paylaterEntity.PaylaterLoan{},
		&paylaterEntity.PaylaterRepayment{},
		&wifiVoucherEntity.WifiVoucher{},
		&transactionEntity.Transaction{},
		&securityEntity.UserSecurity{},
		&transferEntity.Transfer{},
		&ledgerEntity.LedgerEntry{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CleanDB cleans all data from test database
func CleanDB(db *gorm.DB) error {
	// Delete in reverse order of dependencies
	if err := db.Exec("DELETE FROM ledger_entries").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM transfers").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM paylater_repayments").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM paylater_loans").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM paylater_accounts").Error; err != nil {
		return err
	}
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
