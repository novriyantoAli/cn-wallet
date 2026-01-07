# Paylater Repayment Implementation Summary

## Overview
Implemented the complete paylater repayment feature with service, DTOs, and entity models following the transactional SQL flow provided.

## Files Created/Modified

### 1. Entity - `paylater_repayment.entity.go`
**Updated Fields:**
- `ID` - Primary key
- `PaylaterLoanID` - Foreign key to loan
- `UserID` - Foreign key to user
- `Amount` - Repayment amount (must be > 0)
- `PaymentSource` - `wallet` or `external_payment`
- `Status` - Repayment status (pending, success, failed)
- `CreatedAt` - Timestamp

### 2. DTO - `paylater_repayment.dto.go`
**Request DTOs:**
- `CreatePaylaterRepaymentRequest` - Repayment creation
- `ListPaylaterRepaymentsRequest` - List with filters

**Response DTOs:**
- `GetPaylaterRepaymentResponse` - Single repayment details
- `ListPaylaterRepaymentsResponse` - Paginated repayment list

### 3. Service - `paylater_repayment.service.go`
**Interface Methods:**
1. `ProcessRepayment(ctx, req)` - Main business logic
2. `GetRepaymentByID(ctx, id)` - Retrieve single repayment
3. `ListRepaymentsByLoan(ctx, loanID)` - Get repayments by loan
4. `GetRepaymentsByUser(ctx, userID)` - Get user's repayments

## ProcessRepayment Flow (Transactional)

The main method implements a 3-step transactional flow:

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
```

### Step 2: Update Loan Status
- Validates loan exists and belongs to user
- Checks loan status is `active` or `overdue`
- Calculates new outstanding amount: `loan.Total - newPaidAmount`
- If outstanding = 0, sets loan status to `paid`
- Updates via `s.loanRepo.UpdateLoanStatus()`

### Step 3: Create Ledger Entries
Two entries are created with reference type `ReferenceTypeRepayment`:

**Entry 1 - Cash Inflow:**
- `Debit: 0`
- `Credit: amount` (money received)
- `AccountType: AccountTypeWallet` (system cash)

**Entry 2 - Paylater Debt Reduction:**
- `Debit: amount` (debt reduced)
- `Credit: 0`
- `AccountType: AccountTypePaylater` (paylater receivable)

## SQL Equivalent

```sql
BEGIN;

-- Step 1: Create repayment
INSERT INTO paylater_repayments (loan_id, user_id, amount, payment_source, status)
VALUES (15, 2, 500000, 'external_payment', 'success')
RETURNING id;

-- Step 2: Update loan
UPDATE paylater_loans
SET
  paid_amount = paid_amount + 500000,
  outstanding_amount = outstanding_amount - 500000,
  status = CASE
      WHEN outstanding_amount - 500000 = 0 THEN 'paid'
      ELSE status
  END
WHERE id = 15;

-- Step 3: Create ledger entries
INSERT INTO ledger_entries (user_id, reference_id, reference_type, debit, credit, account_type)
VALUES
-- Cash received
(2, 'rep_20', 'repayment', 0, 500000, 'wallet'),
-- Paylater debt reduction
(2, 'rep_20', 'repayment', 500000, 0, 'paylater');

COMMIT;
```

## Dependencies

### Repositories Used:
1. **PaylaterRepaymentRepository** - Create, GetByID operations
2. **PaylaterLoanRepository** - GetLoanByID, UpdateLoanStatus
3. **LedgerRepository** - CreateEntry for ledger records

### Infrastructure:
1. **TransactionManager** - Ensures ACID compliance
2. **Logger** - Structured logging with zap
3. **Database** - Context-aware DB connection via `database.GetDB()`

## Error Handling

The service handles:
- Loan not found
- Unauthorized repayment (loan belongs to different user)
- Invalid loan status (cannot repay paid/defaulted loans)
- Invalid repayment amount (must be > 0)
- Database errors with proper logging
- Rollback on any failure (via transaction manager)

## Logging

All operations logged with structured fields:
- `loan_id` - Loan identifier
- `user_id` - User identifier
- `amount` - Repayment amount
- `payment_source` - Payment method
- `repayment_id` - Created repayment ID
- Error messages with full context

## Notes

1. **PaidAmount Field**: The current `PaylaterLoan` entity doesn't have a `PaidAmount` field. This should be added for complete tracking:
   ```go
   PaidAmount        int64
   OutstandingAmount int64  // Calculated as Total - PaidAmount
   ```

2. **List Methods**: `ListRepaymentsByLoan` and `GetRepaymentsByUser` are placeholders. They require additional repository methods like:
   - `GetByLoanID(ctx, loanID)` 
   - `GetByUserID(ctx, userID)`

3. **Validation**: For enhanced validation, consider:
   - Checking if repayment amount exceeds outstanding
   - Validating payment source availability
   - Check user's wallet balance if `payment_source == "wallet"`

4. **Transaction Isolation**: Uses transaction manager for ACID compliance, suitable for high-concurrency scenarios.
