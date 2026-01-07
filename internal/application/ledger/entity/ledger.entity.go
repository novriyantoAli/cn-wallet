package entity

import "time"

// ReferenceType represents the type of reference for a ledger entry
type ReferenceType string

const (
	ReferenceTypeTransfer       ReferenceType = "transfer"
	ReferenceTypePaylater       ReferenceType = "paylater"
	ReferenceTypePaylaterLoan   ReferenceType = "paylater_loan"
	ReferenceTypePayment        ReferenceType = "payment"
	ReferenceTypeRepayment      ReferenceType = "repayment"
	ReferenceTypeAdjustment     ReferenceType = "adjustment"
	ReferenceTypePurchase       ReferenceType = "purchase"
	ReferenceTypeResellerTopUp  ReferenceType = "reseller_topup"
	ReferenceTypeWalletTransfer ReferenceType = "wallet_transfer"
)

// AccountType represents the type of account for a ledger entry
type AccountType string

const (
	AccountTypeWallet             AccountType = "wallet"
	AccountTypePaylater           AccountType = "paylater"
	AccountTypeVoucher            AccountType = "voucher"
	AccountTypeMerchantIncome     AccountType = "merchant_income"
	AccountTypePaylaterReceivable AccountType = "paylater_receivable"
)

// LedgerEntry represents a single ledger entry (append-only)
type LedgerEntry struct {
	ID            uint64        `gorm:"primaryKey" json:"id"`
	UserID        uint64        `gorm:"index" json:"user_id"`
	ReferenceID   string        `json:"reference_id"`
	ReferenceType ReferenceType `json:"reference_type"`
	Debit         int64         `json:"debit"`
	Credit        int64         `json:"credit"`
	AccountType   AccountType   `json:"account_type"`
	CreatedAt     time.Time     `gorm:"index" json:"created_at"`
}

// TableName specifies the table name for LedgerEntry
func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

// Amount returns the transaction amount (debit or credit)
func (le *LedgerEntry) Amount() int64 {
	if le.Debit > 0 {
		return le.Debit
	}
	return le.Credit
}

// IsDebit checks if this is a debit entry
func (le *LedgerEntry) IsDebit() bool {
	return le.Debit > 0
}

// IsCredit checks if this is a credit entry
func (le *LedgerEntry) IsCredit() bool {
	return le.Credit > 0
}
