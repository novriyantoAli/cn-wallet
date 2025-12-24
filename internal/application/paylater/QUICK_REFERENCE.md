# Paylater Feature - Quick Reference

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/paylater` | Create paylater account |
| GET | `/api/v1/paylater/by-user/:user_id` | Get account by user ID |
| GET | `/api/v1/paylater/:id` | Get account by ID |
| PUT | `/api/v1/paylater/by-user/:user_id/credit-limit` | Update credit limit |
| PUT | `/api/v1/paylater/by-user/:user_id/status` | Update status |
| POST | `/api/v1/paylater/by-user/:user_id/use-credit` | Use credit |
| POST | `/api/v1/paylater/by-user/:user_id/repayment` | Make repayment |
| DELETE | `/api/v1/paylater/by-user/:user_id` | Delete account |

## Quick Start

### 1. Run Migration
```bash
go run cmd/migration/main.go -action=migrate
```

### 2. Start Server
```bash
go run cmd/api/main.go
```

### 3. Create Account
```bash
curl -X POST http://localhost:8080/api/v1/paylater \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "credit_limit": 1000000
  }'
```

### 4. Use Credit
```bash
curl -X POST http://localhost:8080/api/v1/paylater/by-user/1/use-credit \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000,
    "description": "Purchase"
  }'
```

### 5. Make Repayment
```bash
curl -X POST http://localhost:8080/api/v1/paylater/by-user/1/repayment \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 25000,
    "description": "Payment"
  }'
```

## File Structure
```
internal/application/paylater/
├── README.md
├── module.go
├── dto/
│   └── paylater.dto.go
├── entity/
│   └── paylater.entity.go
├── handler/
│   └── paylater.handler.go
├── repository/
│   └── paylater.repo.go
└── service/
    └── paylater.service.go
```

## Key Constants

### Status Values
- `active` - Account is active
- `suspended` - Account is suspended
- `closed` - Account is closed

## Business Rules

1. ✅ Credit limit must be > 0
2. ✅ Outstanding balance ≥ 0
3. ✅ Available limit = credit_limit - outstanding
4. ✅ One account per user
5. ✅ Must be active to use credit
6. ✅ Cannot delete with outstanding balance
7. ✅ Repayment ≤ outstanding balance

## Status
✅ Feature complete and ready to use!
