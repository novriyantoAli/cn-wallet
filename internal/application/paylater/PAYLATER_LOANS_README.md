# Paylater Loans Feature

## Overview
The Paylater Loans feature allows users to take loans from their paylater accounts. This feature integrates with the existing paylater accounts system and tracks loan history, repayments, and overdue status.

## Database Schema

### paylater_loans Table
```sql
CREATE TABLE paylater_loans (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  amount BIGINT NOT NULL CHECK (amount > 0),
  interest BIGINT NOT NULL CHECK (interest >= 0),
  total BIGINT NOT NULL CHECK (total = amount + interest),
  due_date DATE NOT NULL,
  source VARCHAR(20) NOT NULL CHECK (source IN ('transfer', 'checkout')),
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Features

### 1. **Loan Creation**
- Creates a new loan for a user
- Automatically deducts from available credit limit
- Validates sufficient credit availability
- Supports two sources: `transfer` and `checkout`
- Calculates total including interest

### 2. **Loan Management**
- View loan details by ID
- List all loans for a user
- Filter loans by status, source, and date range
- Paginated loan listing
- Update loan status

### 3. **Loan Repayment**
- Mark loan as paid
- Automatically restores credit to the account
- Updates available credit limit

### 4. **Overdue Loan Processing**
- Automatically identifies overdue loans
- Updates status to overdue
- Can be run periodically via cron or manually

### 5. **Loan Statistics**
- Total loans count
- Active loans count
- Total borrowed amount
- Total paid amount
- Total outstanding amount
- Overdue loans count

## API Endpoints

### Create Loan
```
POST /api/v1/paylater/loans
```
**Request Body:**
```json
{
  "user_id": 1,
  "amount": 100000,
  "interest": 5000,
  "due_date": "2026-01-25",
  "source": "checkout"
}
```

### Get Loan by ID
```
GET /api/v1/paylater/loans/:id
```

### Get Loans by User ID
```
GET /api/v1/paylater/loans/by-user/:user_id
```

### List Loans (with filters)
```
GET /api/v1/paylater/loans?user_id=1&status=active&page=1&page_size=10
```

**Query Parameters:**
- `user_id` (optional): Filter by user
- `status` (optional): Filter by status (pending, active, paid, overdue, defaulted)
- `source` (optional): Filter by source (transfer, checkout)
- `from_date` (optional): Filter by date range (YYYY-MM-DD)
- `to_date` (optional): Filter by date range (YYYY-MM-DD)
- `page` (required): Page number
- `page_size` (required): Number of items per page

### Update Loan Status
```
PUT /api/v1/paylater/loans/:id/status
```
**Request Body:**
```json
{
  "status": "paid"
}
```

### Mark Loan as Paid
```
POST /api/v1/paylater/loans/:id/mark-paid
```

### Process Overdue Loans
```
POST /api/v1/paylater/loans/process-overdue
```

### Get User Loan Statistics
```
GET /api/v1/paylater/loans/stats/:user_id
```

### Delete Loan
```
DELETE /api/v1/paylater/loans/:id
```
*Note: Only paid loans can be deleted*

## Loan Status Flow

```
pending → active → paid
           ↓
        overdue → defaulted
```

- **pending**: Loan is created but not yet active
- **active**: Loan is active and due
- **paid**: Loan has been repaid
- **overdue**: Loan is past due date
- **defaulted**: Loan has been defaulted

## Loan Sources

- **transfer**: Loan created from a transfer operation
- **checkout**: Loan created from a checkout/purchase operation

## Business Rules

1. **Credit Limit Validation**: Loan total (amount + interest) must not exceed available credit
2. **Account Status**: Only active paylater accounts can create loans
3. **Repayment**: Marking a loan as paid restores credit to the account
4. **Deletion**: Only paid loans can be deleted
5. **Overdue Detection**: Loans past their due date are automatically marked as overdue

## Integration with Paylater Accounts

The loans feature is tightly integrated with paylater accounts:

1. When a loan is created:
   - Validates account exists and is active
   - Deducts loan total from available credit
   - Increases outstanding balance

2. When a loan is marked as paid:
   - Decreases outstanding balance
   - Increases available credit

## Running the Migration

To create the paylater_loans table:

```bash
psql -U your_username -d your_database -f scripts/migrations/create_paylater_loans_table.sql
```

Or use your migration tool to run:
```bash
make migrate-up
```

## Example Usage

### Creating a Loan
```bash
curl -X POST http://localhost:8080/api/v1/paylater/loans \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 100000,
    "interest": 5000,
    "due_date": "2026-01-25",
    "source": "checkout"
  }'
```

### Listing User's Loans
```bash
curl http://localhost:8080/api/v1/paylater/loans/by-user/1
```

### Getting Loan Statistics
```bash
curl http://localhost:8080/api/v1/paylater/loans/stats/1
```

### Marking Loan as Paid
```bash
curl -X POST http://localhost:8080/api/v1/paylater/loans/5/mark-paid
```

## Testing

The feature includes comprehensive test coverage:
- Repository tests for data access
- Service tests for business logic
- Handler tests for HTTP endpoints

Run tests with:
```bash
go test ./internal/application/paylater/...
```

## Notes

- All monetary amounts are stored in the smallest currency unit (e.g., cents)
- Dates are stored in UTC timezone
- The feature uses database transactions to ensure data consistency
- Row-level locking is used to prevent race conditions

## Future Enhancements

- Automatic interest calculation based on loan amount and duration
- Partial payment support
- Payment reminders
- Late payment penalties
- Loan refinancing options
