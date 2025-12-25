# Paylater Tests Quick Reference

## Running Tests

### All Paylater Repository Tests
```bash
cd /home/rhein/projects/backend/cn-wallet
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -v
```

### Paylater Account Repository Only
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/paylater.repo_test.go ./internal/application/paylater/repository/paylater.repo.go -v
```

### Paylater Loan Repository Only
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/paylater_loan.repo_test.go ./internal/application/paylater/repository/paylater_loan.repo.go -v
```

### With Coverage Report
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -cover
```

### Detailed Coverage HTML Report
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html in browser
```

### Run Specific Test
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -run TestPaylaterLoanRepository_CreateLoan -v
```

### Clean Test Cache
```bash
go clean -testcache
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -v
```

## Test Results

✅ **Paylater Account Repository**: 16 tests passing  
✅ **Paylater Loan Repository**: 21 tests passing  
✅ **Total**: 37 tests passing  
✅ **Coverage**: 84.1%  
✅ **Execution Time**: ~0.2s  

## Prerequisites

### Install GCC (if not already installed)
```bash
sudo apt-get update
sudo apt-get install -y gcc
```

### Verify Installation
```bash
gcc --version
```

## Common Issues

### Issue: "CGO_ENABLED=0"
**Solution**: Ensure CGO is enabled:
```bash
CGO_ENABLED=1 go test ...
```

### Issue: "gcc: not found"
**Solution**: Install gcc (see Prerequisites above)

### Issue: Test cache preventing fresh runs
**Solution**: Use `-count=1` flag:
```bash
CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -count=1
```

## CI/CD Integration

### Makefile Target
```makefile
.PHONY: test-paylater
test-paylater:
	CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -v -cover

.PHONY: test-paylater-coverage
test-paylater-coverage:
	CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
```

### GitHub Actions
```yaml
name: Paylater Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y gcc
      
      - name: Run Paylater Repository Tests
        run: CGO_ENABLED=1 go test ./internal/application/paylater/repository/... -v -cover
```

## Test Output Example

```
=== RUN   TestPaylaterLoanRepository_CreateLoan
=== RUN   TestPaylaterLoanRepository_CreateLoan/should_create_paylater_loan_successfully
=== RUN   TestPaylaterLoanRepository_CreateLoan/should_create_multiple_loans_for_same_user
--- PASS: TestPaylaterLoanRepository_CreateLoan (0.01s)
    --- PASS: TestPaylaterLoanRepository_CreateLoan/should_create_paylater_loan_successfully (0.01s)
    --- PASS: TestPaylaterLoanRepository_CreateLoan/should_create_multiple_loans_for_same_user (0.01s)
...
PASS
ok      github.com/novriyantoAli/cn-wallet/internal/application/paylater/repository     0.201s  coverage: 84.1% of statements
```

## Files Created

1. **paylater_loan.repo_test.go** (829 lines)
   - 11 test groups
   - 21 individual test cases
   - Comprehensive coverage of all repository methods

2. **Updated: testutil/database.go**
   - Added PaylaterLoan entity to auto-migration

## Test Statistics

| Repository | Test Groups | Test Cases | Coverage |
|------------|-------------|------------|----------|
| Paylater Account | 8 | 16 | ~85% |
| Paylater Loan | 11 | 21 | ~84% |
| **Total** | **19** | **37** | **84.1%** |

---

**Status**: ✅ All Tests Passing  
**Last Run**: December 25, 2025  
**Environment**: Go 1.21+ with CGO enabled
