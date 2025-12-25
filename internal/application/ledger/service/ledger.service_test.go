package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/dto"
	"github.com/novriyantoAli/cn-wallet/internal/application/ledger/entity"
)

type MockLedgerRepository struct {
	mock.Mock
}

func (m *MockLedgerRepository) CreateEntry(ctx context.Context, entry *entity.LedgerEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockLedgerRepository) GetEntryByID(ctx context.Context, id uint64) (*entity.LedgerEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LedgerEntry), args.Error(1)
}

func (m *MockLedgerRepository) GetEntriesByUserID(ctx context.Context, userID uint64) ([]entity.LedgerEntry, error) {
	args := m.Called(ctx, userID)
	var entries []entity.LedgerEntry
	if args.Get(0) != nil {
		entries = args.Get(0).([]entity.LedgerEntry)
	}
	return entries, args.Error(1)
}

func (m *MockLedgerRepository) GetEntriesByReference(ctx context.Context, referenceType string, referenceID string) ([]entity.LedgerEntry, error) {
	args := m.Called(ctx, referenceType, referenceID)
	var entries []entity.LedgerEntry
	if args.Get(0) != nil {
		entries = args.Get(0).([]entity.LedgerEntry)
	}
	return entries, args.Error(1)
}

func (m *MockLedgerRepository) ListEntries(ctx context.Context, req *dto.ListLedgerEntriesRequest) ([]entity.LedgerEntry, int64, error) {
	args := m.Called(ctx, req)
	var entries []entity.LedgerEntry
	if args.Get(0) != nil {
		entries = args.Get(0).([]entity.LedgerEntry)
	}
	return entries, args.Get(1).(int64), args.Error(2)
}

func (m *MockLedgerRepository) GetUserStats(ctx context.Context, userID uint64) (*dto.LedgerStatsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LedgerStatsResponse), args.Error(1)
}

func setupServiceTest() (LedgerService, *MockLedgerRepository) {
	repo := new(MockLedgerRepository)
	logger := zap.NewNop()
	service := NewLedgerService(repo, logger)
	return service, repo
}

func TestCreateEntry(t *testing.T) {
	t.Run("Success with Debit", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		request := &dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
		}

		repo.On("CreateEntry", ctx, mock.MatchedBy(func(e *entity.LedgerEntry) bool {
			return e.UserID == 1 && e.Debit == 50000 && e.Credit == 0
		})).Run(func(args mock.Arguments) {
			entry := args.Get(1).(*entity.LedgerEntry)
			entry.ID = 1
			entry.CreatedAt = time.Now()
		}).Return(nil)

		response, err := service.CreateEntry(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint64(1), response.UserID)
		assert.Equal(t, int64(50000), response.Debit)
		assert.Equal(t, int64(50000), response.Amount)
		repo.AssertExpectations(t)
	})

	t.Run("Success with Credit", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		request := &dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "payment",
			Debit:         0,
			Credit:        50000,
			AccountType:   "wallet",
		}

		repo.On("CreateEntry", ctx, mock.MatchedBy(func(e *entity.LedgerEntry) bool {
			return e.UserID == 1 && e.Debit == 0 && e.Credit == 50000
		})).Run(func(args mock.Arguments) {
			entry := args.Get(1).(*entity.LedgerEntry)
			entry.ID = 2
			entry.CreatedAt = time.Now()
		}).Return(nil)

		response, err := service.CreateEntry(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(50000), response.Credit)
		assert.Equal(t, int64(50000), response.Amount)
		repo.AssertExpectations(t)
	})

	t.Run("Error: Both Debit and Credit", func(t *testing.T) {
		service, _ := setupServiceTest()
		ctx := context.Background()

		request := &dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        50000,
			AccountType:   "wallet",
		}

		response, err := service.CreateEntry(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "either debit or credit must be set, but not both", err.Error())
	})

	t.Run("Error: Neither Debit nor Credit", func(t *testing.T) {
		service, _ := setupServiceTest()
		ctx := context.Background()

		request := &dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         0,
			Credit:        0,
			AccountType:   "wallet",
		}

		response, err := service.CreateEntry(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("Repository Error", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		request := &dto.CreateLedgerEntryRequest{
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: "transfer",
			Debit:         50000,
			Credit:        0,
			AccountType:   "wallet",
		}

		repo.On("CreateEntry", ctx, mock.Anything).Return(errors.New("database error"))

		response, err := service.CreateEntry(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		repo.AssertExpectations(t)
	})
}

func TestGetEntryByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		expectedEntry := &entity.LedgerEntry{
			ID:            1,
			UserID:        1,
			ReferenceID:   "100",
			ReferenceType: entity.ReferenceTypeTransfer,
			Debit:         50000,
			Credit:        0,
			AccountType:   entity.AccountTypeWallet,
			CreatedAt:     time.Now(),
		}

		repo.On("GetEntryByID", ctx, uint64(1)).Return(expectedEntry, nil)

		response, err := service.GetEntryByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint64(1), response.ID)
		assert.Equal(t, int64(50000), response.Amount)
		repo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		repo.On("GetEntryByID", ctx, uint64(999)).Return(nil, errors.New("record not found"))

		response, err := service.GetEntryByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, response)
		repo.AssertExpectations(t)
	})
}

func TestGetEntriesByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		expectedEntries := []entity.LedgerEntry{
			{
				ID:            1,
				UserID:        1,
				ReferenceID:   "100",
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         50000,
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
				CreatedAt:     time.Now(),
			},
			{
				ID:            2,
				UserID:        1,
				ReferenceID:   "101",
				ReferenceType: entity.ReferenceTypePayment,
				Debit:         0,
				Credit:        30000,
				AccountType:   entity.AccountTypeWallet,
				CreatedAt:     time.Now(),
			},
		}

		repo.On("GetEntriesByUserID", ctx, uint64(1)).Return(expectedEntries, nil)

		responses, err := service.GetEntriesByUserID(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, responses, 2)
		assert.Equal(t, "transfer", responses[0].ReferenceType)
		assert.Equal(t, int64(50000), responses[0].Amount)
		repo.AssertExpectations(t)
	})

	t.Run("Empty Results", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		repo.On("GetEntriesByUserID", ctx, uint64(999)).Return([]entity.LedgerEntry{}, nil)

		responses, err := service.GetEntriesByUserID(ctx, 999)

		assert.NoError(t, err)
		assert.Len(t, responses, 0)
		repo.AssertExpectations(t)
	})
}

func TestGetEntriesByReference(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		expectedEntries := []entity.LedgerEntry{
			{
				ID:            1,
				UserID:        1,
				ReferenceID:   "100",
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         50000,
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
				CreatedAt:     time.Now(),
			},
		}

		repo.On("GetEntriesByReference", ctx, "transfer", "100").Return(expectedEntries, nil)

		responses, err := service.GetEntriesByReference(ctx, "transfer", "100")

		assert.NoError(t, err)
		assert.Len(t, responses, 1)
		assert.Equal(t, "transfer", responses[0].ReferenceType)
		repo.AssertExpectations(t)
	})
}

func TestListEntries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		request := &dto.ListLedgerEntriesRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedEntries := []entity.LedgerEntry{
			{
				ID:            1,
				UserID:        1,
				ReferenceID:   "100",
				ReferenceType: entity.ReferenceTypeTransfer,
				Debit:         50000,
				Credit:        0,
				AccountType:   entity.AccountTypeWallet,
				CreatedAt:     time.Now(),
			},
		}

		repo.On("ListEntries", ctx, request).Return(expectedEntries, int64(1), nil)

		response, err := service.ListEntries(ctx, request)

		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.Data, 1)
		assert.Equal(t, int64(1), response.TotalCount)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		repo.AssertExpectations(t)
	})

	t.Run("With Default Pagination", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		request := &dto.ListLedgerEntriesRequest{
			Page:     0,
			PageSize: 0,
		}

		expectedEntries := []entity.LedgerEntry{}

		repo.On("ListEntries", ctx, mock.MatchedBy(func(r *dto.ListLedgerEntriesRequest) bool {
			return r.Page == 1 && r.PageSize == 10
		})).Return(expectedEntries, int64(0), nil)

		response, err := service.ListEntries(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		repo.AssertExpectations(t)
	})
}

func TestGetUserStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, repo := setupServiceTest()
		ctx := context.Background()

		expectedStats := &dto.LedgerStatsResponse{
			UserID:          1,
			TotalDebit:      100000,
			TotalCredit:     50000,
			EntryCount:      5,
			WalletDebit:     80000,
			WalletCredit:    40000,
			PaylaterDebit:   20000,
			PaylaterCredits: 10000,
		}

		repo.On("GetUserStats", ctx, uint64(1)).Return(expectedStats, nil)

		response, err := service.GetUserStats(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint64(1), response.UserID)
		assert.Equal(t, int64(100000), response.TotalDebit)
		assert.Equal(t, int64(50000), response.TotalCredit)
		repo.AssertExpectations(t)
	})
}
