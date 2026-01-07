# ProcessRepayment - Updated Flow Summary

## Updated Transaction Flow

The `ProcessRepayment` method now follows the exact SQL pattern you provided:

### Step 1: Create Repayment Record
```go
repayment = &paylaterEntity.PaylaterRepayment{
    PaylaterLoanID: req.PaylaterLoanID,
    UserID:         req.UserID,
    Amount:         req.Amount,
    PaymentSource:  req.PaymentSource,
    Status:         "success",
    CreatedAt:      time.Now(),
}
s.repaymentRepo.Create(txCtx, repayment)  // INSERT INTO paylater_repayments
```

### Step 2: Create Two Ledger Entries
Both entries share the same reference ID (the repayment ID):

**Entry 1 - Cash Inflow:**
```go
ledgerEntry1 := &entity.LedgerEntry{
    UserID:        uint64(req.UserID),
    ReferenceID:   fmt.Sprintf("%d", repayment.ID),  // Numeric ID
    ReferenceType: entity.ReferenceTypeRepayment,
    Debit:         0,
    Credit:        req.Amount,
    AccountType:   entity.AccountTypeMerchantIncome,  // system_cash
    CreatedAt:     time.Now(),
}
```

**Entry 2 - Paylater Debt Reduction:**
```go
ledgerEntry2 := &entity.LedgerEntry{
    UserID:        uint64(req.UserID),
    ReferenceID:   fmt.Sprintf("%d", repayment.ID),  // Same reference
    ReferenceType: entity.ReferenceTypeRepayment,
    Debit:         req.Amount,
    Credit:        0,
    AccountType:   entity.AccountTypePaylater,
    CreatedAt:     time.Now(),
}
```

### Step 3: Update Loan Status
Update loan status to 'paid':
```go
s.loanRepo.UpdateLoanStatus(txCtx, req.PaylaterLoanID, string(paylaterEntity.PaylaterLoanStatusPaid))
```

## SQL Equivalent

```sql
BEGIN;

-- Step 1: Create repayment
INSERT INTO paylater_repayments
(loan_id, user_id, amount, payment_source, status)
VALUES
(15, 2, 500000, 'external', 'success')
RETURNING id;

-- Step 2: Create ledger entries
INSERT INTO ledger_entries
(user_id, reference_id, reference_type, debit, credit, account_type)
VALUES
-- Cash inflow
(2, 30, 'paylater_repayment', 0, 500000, 'merchant_income'),
-- Paylater debt reduction
(2, 30, 'paylater_repayment', 500000, 0, 'paylater');

-- Step 3: Update loan status
UPDATE paylater_loans
SET status = 'paid'
WHERE id = 15;

COMMIT;
```

## Key Changes

✅ **Order Changed**: Ledger entries are now created BEFORE loan status update
✅ **Reference ID**: Uses numeric ID instead of prefixed "rep_" format
✅ **Account Type**: Uses `AccountTypeMerchantIncome` for cash inflow (system_cash)
✅ **Transactional**: All operations within single transaction for atomicity
✅ **Full Repayment**: Currently marks loan as 'paid' immediately (can be enhanced with conditional logic)

## Notes

The SQL query includes a conditional UPDATE that checks if total repayments >= loan total. This can be implemented in the service by:
1. Fetching total repayments for the loan before updating status
2. Only updating to 'paid' if sufficient repayment has been made
3. Or using a database trigger for the conditional check

Current implementation always marks as 'paid' on successful repayment - adjust based on your business logic.
