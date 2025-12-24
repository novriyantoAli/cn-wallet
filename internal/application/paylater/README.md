# Paylater Account Feature

## Overview
The Paylater Account feature allows users to have a credit account where they can use credit up to a specified limit and repay later. This feature tracks credit limits, outstanding balances, and available credit.

## Database Schema

### Table: `paylater_accounts`
```sql
CREATE TABLE paylater_accounts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
  credit_limit BIGINT NOT NULL CHECK (credit_limit > 0),
  outstanding BIGINT NOT NULL DEFAULT 0 CHECK (outstanding >= 0),
  available_limit BIGINT NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CHECK (available_limit = credit_limit - outstanding)
);
```

### Constraints
- `user_id` must be unique (one paylater account per user)
- `credit_limit` must be greater than 0
- `outstanding` must be greater than or equal to 0
- `available_limit` must equal `credit_limit - outstanding`

### Status Values
- `active`: Account is active and can be used
- `suspended`: Account is temporarily suspended
- `closed`: Account is permanently closed

## API Endpoints

### 1. Create Paylater Account
**POST** `/api/v1/paylater`

Creates a new paylater account for a user.

**Request Body:**
```json
{
  "user_id": 1,
  "credit_limit": 1000000
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "credit_limit": 1000000,
    "outstanding": 0,
    "available_limit": 1000000,
    "status": "active",
    "created_at": "2025-12-24T10:00:00Z"
  }
}
```

### 2. Get Account by User ID
**GET** `/api/v1/paylater/by-user/:user_id`

Retrieves paylater account information for a specific user.

**Response:**
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "credit_limit": 1000000,
    "outstanding": 250000,
    "available_limit": 750000,
    "status": "active",
    "created_at": "2025-12-24T10:00:00Z"
  }
}
```

### 3. Get Account by ID
**GET** `/api/v1/paylater/:id`

Retrieves paylater account information by account ID.

### 4. Update Credit Limit
**PUT** `/api/v1/paylater/by-user/:user_id/credit-limit`

Updates the credit limit for a paylater account.

**Request Body:**
```json
{
  "credit_limit": 2000000
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "credit_limit": 2000000,
    "outstanding": 250000,
    "available_limit": 1750000,
    "status": "active",
    "created_at": "2025-12-24T10:00:00Z"
  },
  "message": "Credit limit updated successfully"
}
```

### 5. Update Status
**PUT** `/api/v1/paylater/by-user/:user_id/status`

Updates the status of a paylater account.

**Request Body:**
```json
{
  "status": "suspended"
}
```

**Valid status values:** `active`, `suspended`, `closed`

### 6. Use Credit
**POST** `/api/v1/paylater/by-user/:user_id/use-credit`

Uses credit from the paylater account. This increases the outstanding balance and decreases the available limit.

**Request Body:**
```json
{
  "amount": 50000,
  "description": "Purchase product X"
}
```

**Response:**
```json
{
  "data": {
    "user_id": 1,
    "amount": 50000,
    "outstanding": 300000,
    "available_limit": 700000,
    "message": "Credit used successfully"
  }
}
```

**Business Rules:**
- Account must be in `active` status
- Amount must not exceed available limit
- Transaction is atomic (uses database transaction)

### 7. Repayment
**POST** `/api/v1/paylater/by-user/:user_id/repayment`

Processes a repayment to reduce outstanding balance.

**Request Body:**
```json
{
  "amount": 100000,
  "description": "Monthly payment"
}
```

**Response:**
```json
{
  "data": {
    "user_id": 1,
    "amount": 100000,
    "outstanding": 200000,
    "available_limit": 800000,
    "message": "Repayment processed successfully"
  }
}
```

**Business Rules:**
- Repayment amount cannot exceed outstanding balance
- Transaction is atomic (uses database transaction)

### 8. Delete Account
**DELETE** `/api/v1/paylater/by-user/:user_id`

Deletes a paylater account.

**Response:**
```json
{
  "message": "Account deleted successfully"
}
```

**Business Rules:**
- Account cannot be deleted if it has outstanding balance
- Outstanding balance must be 0

## Architecture

The paylater feature follows the clean architecture pattern:

```
internal/application/paylater/
├── module.go                    # Dependency injection module
├── dto/
│   └── paylater.dto.go         # Data transfer objects
├── entity/
│   └── paylater.entity.go      # Database entity
├── handler/
│   └── paylater.handler.go     # HTTP handlers
├── repository/
│   └── paylater.repo.go        # Database operations
└── service/
    └── paylater.service.go     # Business logic
```

## Running Migrations

To create the paylater_accounts table, run:

```bash
go run cmd/migration/main.go -action=migrate
```

This will automatically create the table using GORM AutoMigrate.

## Testing

### Manual Testing with cURL

1. **Create Account:**
```bash
curl -X POST http://localhost:8080/api/v1/paylater \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "credit_limit": 1000000
  }'
```

2. **Get Account:**
```bash
curl http://localhost:8080/api/v1/paylater/by-user/1
```

3. **Use Credit:**
```bash
curl -X POST http://localhost:8080/api/v1/paylater/by-user/1/use-credit \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000,
    "description": "Test purchase"
  }'
```

4. **Repayment:**
```bash
curl -X POST http://localhost:8080/api/v1/paylater/by-user/1/repayment \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 25000,
    "description": "Test repayment"
  }'
```

## Integration with Existing Features

The paylater feature can be integrated with:
- **Payment Module**: Process repayments through payment gateway
- **Transaction Module**: Track all credit usage and repayments as transactions
- **User Module**: Link paylater accounts to user profiles
- **Wallet Module**: Transfer funds between wallet and paylater accounts

## Future Enhancements

1. **Interest Calculation**: Add interest on outstanding balance
2. **Payment Due Dates**: Track due dates for repayments
3. **Credit History**: Store transaction history for credit usage
4. **Credit Score**: Calculate credit score based on repayment behavior
5. **Auto-repayment**: Automatic deduction from wallet on due date
6. **Notifications**: Send alerts for due payments
7. **Installment Plans**: Support for installment-based repayments
