package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	productEntity "github.com/novriyantoAli/cn-wallet/internal/application/product/entity"
	productRepo "github.com/novriyantoAli/cn-wallet/internal/application/product/repository"
	"github.com/novriyantoAli/cn-wallet/internal/application/purchase/dto"
	transactionDto "github.com/novriyantoAli/cn-wallet/internal/application/transaction/dto"
	transactionEntity "github.com/novriyantoAli/cn-wallet/internal/application/transaction/entity"
	transactionRepo "github.com/novriyantoAli/cn-wallet/internal/application/transaction/repository"
	userEntity "github.com/novriyantoAli/cn-wallet/internal/application/user/entity"
	userRepo "github.com/novriyantoAli/cn-wallet/internal/application/user/repository"
	walletRepo "github.com/novriyantoAli/cn-wallet/internal/application/wallet/repository"
	wifiVoucherDto "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/dto"
	wifiVoucherEntity "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/entity"
	wifiVoucherRepo "github.com/novriyantoAli/cn-wallet/internal/application/wifivoucher/repository"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/database"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProviderClient defines the interface for provider API calls
type ProviderClient interface {
	ProcessPurchase(ctx context.Context, productCode string, phone string, reference string) (serialNumber string, err error)
}

// PurchaseService defines the interface for purchase business logic
type PurchaseService interface {
	// ProcessPurchase handles the complete purchase flow:
	// Step 1: Validate product exists
	// Step 2: Check wallet balance >= Product Price
	// Step 3: Create pending transaction
	// Step 4: Deduct balance from wallet
	// Step 5: Call Provider API
	// Success: Update Transaction to success, save serial_number
	// Failed: Update Transaction to failed, Refund balance
	ProcessPurchase(ctx context.Context, token string, req *dto.PurchaseRequest) (*dto.PurchaseResponse, error)

	// ProcessWifiPurchase handles the wifi voucher purchase flow:
	// Step 1: Validate product exists and is wifi type
	// Step 2: Check wallet balance >= Product Price
	// Step 3: Get available wifi voucher
	// Step 4: Create pending transaction
	// Step 5: Deduct balance from wallet
	// Step 6: Mark voucher as sold and assign to user
	ProcessWifiPurchase(ctx context.Context, token string, req *dto.PurchaseWifiRequest) (*dto.PurchaseWifiResponse, error)

	// GetPurchaseHistory retrieves purchase history for a wallet
	GetPurchaseHistory(ctx context.Context, filter *dto.PurchaseHistoryFilter) (*dto.PurchaseHistoryList, error)
}

type purchaseService struct {
	productRepo     productRepo.ProductRepository
	walletRepo      walletRepo.WalletRepository
	transactionRepo transactionRepo.TransactionRepository
	userRepo        userRepo.UserRepository
	wifiVoucherRepo wifiVoucherRepo.WifiVoucherRepository
	providerClient  ProviderClient
	txManager       database.TransactionManagerI
	jwtManager      *jwt.JWTManager
	logger          *zap.Logger
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(
	productRepo productRepo.ProductRepository,
	walletRepo walletRepo.WalletRepository,
	transactionRepo transactionRepo.TransactionRepository,
	userRepo userRepo.UserRepository,
	wifiVoucherRepo wifiVoucherRepo.WifiVoucherRepository,
	providerClient ProviderClient,
	txManager database.TransactionManagerI,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) PurchaseService {
	return &purchaseService{
		productRepo:     productRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		wifiVoucherRepo: wifiVoucherRepo,
		providerClient:  providerClient,
		txManager:       txManager,
		jwtManager:      jwtManager,
		logger:          logger,
	}
}

// ProcessPurchase processes a complete purchase transaction
func (s *purchaseService) ProcessPurchase(ctx context.Context, token string, req *dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	// Validate request
	if req.ProductID == 0 {
		return nil, errors.New("product_id is required")
	}
	if req.Phone == "" {
		return nil, errors.New("phone is required")
	}

	// Verify token and get user
	user, err := s.getUserFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Validate product exists
	product, err := s.getProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	txID := uuid.New()

	// Execute database operations within a transaction
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.processTransaction(txCtx, user.ID, req.Phone, product, txID)
	})

	if err != nil {
		return nil, err
	}

	// Call Provider API (outside transaction to avoid long locks)
	serialNumber, err := s.providerClient.ProcessPurchase(ctx, product.Code, req.Phone, txID.String())

	if err != nil {
		s.logger.Error("Provider API call failed", zap.String("tx_id", txID.String()), zap.Error(err))
		return s.handlePurchaseFailure(ctx, txID, user.ID, product.PriceSell)
	}

	// Update transaction to success
	err = s.updateTransactionSuccess(ctx, txID, serialNumber)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Purchase completed successfully", zap.String("tx_id", txID.String()), zap.String("serial_number", serialNumber))

	return &dto.PurchaseResponse{
		TransactionID: txID,
		Status:        transactionEntity.StatusSuccess,
		SerialNumber:  serialNumber,
		Message:       "Purchase successful",
	}, nil
}

// getUserFromToken verifies the token and retrieves the user
func (s *purchaseService) getUserFromToken(ctx context.Context, token string) (*userEntity.User, error) {
	if s.jwtManager == nil {
		s.logger.Warn("JWT manager not initialized")
		return nil, errors.New("invalid or expired token")
	}

	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		s.logger.Warn("Invalid token provided", zap.Error(err))
		return nil, errors.New("invalid or expired token")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.Error("Failed to retrieve user", zap.Uint("user_id", claims.UserID), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// getProduct retrieves a product by ID
func (s *purchaseService) getProduct(ctx context.Context, productID uint) (*productEntity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		s.logger.Error("Failed to get product", zap.Uint("product_id", productID), zap.Error(err))
		return nil, err
	}
	return product, nil
}

// processTransaction handles the transaction creation and balance deduction
func (s *purchaseService) processTransaction(ctx context.Context, userID uint, phone string, product *productEntity.Product, txID uuid.UUID) error {
	// Get wallet with FOR UPDATE lock
	wallet, err := s.walletRepo.GetForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		s.logger.Error("Failed to get wallet", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}

	s.logger.Info("Starting purchase process", zap.String("tx_id", txID.String()), zap.Uint("wallet_id", wallet.ID), zap.Float64("balance", wallet.Balance), zap.Float64("price", product.PriceSell))

	// Check wallet balance
	if wallet.Balance < product.PriceSell {
		s.logger.Warn("Insufficient balance", zap.Uint("wallet_id", wallet.ID), zap.Float64("balance", wallet.Balance), zap.Float64("price", product.PriceSell))
		return errors.New("insufficient wallet balance")
	}

	// Create pending transaction
	txn := &transactionEntity.Transaction{
		ID:           txID,
		WalletID:     wallet.ID,
		Type:         transactionEntity.TypePurchase,
		Amount:       product.PriceSell,
		Status:       transactionEntity.StatusPending,
		Description:  "Product purchase",
		ProductID:    &product.ID,
		TargetNumber: phone,
		CreatedAt:    time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, txn); err != nil {
		s.logger.Error("Failed to create transaction", zap.String("tx_id", txID.String()), zap.Error(err))
		return errors.New("failed to create transaction")
	}

	s.logger.Info("Transaction created", zap.String("tx_id", txID.String()), zap.Float64("amount", product.PriceSell))

	// Deduct balance from wallet
	newBalance := wallet.Balance - product.PriceSell
	if err := s.walletRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		s.logger.Error("Failed to update wallet balance", zap.Uint("user_id", userID), zap.Error(err))
		return errors.New("failed to update wallet balance")
	}

	s.logger.Info("Wallet balance updated", zap.Uint("wallet_id", wallet.ID), zap.Float64("new_balance", newBalance))

	return nil
}

// handlePurchaseFailure handles failed purchases with refund
func (s *purchaseService) handlePurchaseFailure(ctx context.Context, txID uuid.UUID, userID uint, refundAmount float64) (*dto.PurchaseResponse, error) {
	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Update transaction status to failed
		txn := &transactionEntity.Transaction{
			ID:     txID,
			Status: transactionEntity.StatusFailed,
		}
		if err := s.transactionRepo.Update(txCtx, txn); err != nil {
			s.logger.Error("Failed to update transaction status", zap.String("tx_id", txID.String()), zap.Error(err))
			return err
		}

		// Refund balance
		wallet, err := s.walletRepo.GetForUpdate(txCtx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found for refund")
			}
			s.logger.Error("Failed to get wallet for refund", zap.Uint("user_id", userID), zap.Error(err))
			return err
		}

		if err := s.walletRepo.UpdateBalance(txCtx, userID, wallet.Balance+refundAmount); err != nil {
			s.logger.Error("Failed to refund wallet balance", zap.Uint("wallet_id", wallet.ID), zap.Error(err))
			return err
		}

		s.logger.Info("Purchase refunded", zap.String("tx_id", txID.String()), zap.Float64("refund_amount", refundAmount))
		return nil
	})

	if err != nil {
		s.logger.Error("Failed to handle purchase failure", zap.Error(err))
	}

	return &dto.PurchaseResponse{
		TransactionID: txID,
		Status:        transactionEntity.StatusFailed,
		Message:       "Purchase failed, balance refunded",
	}, nil
}

// updateTransactionSuccess updates transaction to success with serial number
func (s *purchaseService) updateTransactionSuccess(ctx context.Context, txID uuid.UUID, serialNumber string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		txn := &transactionEntity.Transaction{
			ID:           txID,
			Status:       transactionEntity.StatusSuccess,
			SerialNumber: serialNumber,
		}

		if err := s.transactionRepo.Update(txCtx, txn); err != nil {
			s.logger.Error("Failed to update transaction to success", zap.String("tx_id", txID.String()), zap.Error(err))
			return errors.New("failed to update transaction status")
		}

		return nil
	})
}

// GetPurchaseHistory retrieves purchase history for a wallet
func (s *purchaseService) GetPurchaseHistory(ctx context.Context, filter *dto.PurchaseHistoryFilter) (*dto.PurchaseHistoryList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}

	s.logger.Info("Getting purchase history", zap.Uint("wallet_id", filter.WalletID), zap.Int("page", filter.Page), zap.Int("page_size", filter.PageSize))

	// Retrieve purchase transactions
	transactionFilter := &transactionDto.TransactionFilter{
		WalletID: filter.WalletID,
		Type:     transactionEntity.TypePurchase,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	transactions, total, err := s.transactionRepo.GetAll(ctx, transactionFilter)
	if err != nil {
		s.logger.Error("Failed to get transactions", zap.Error(err))
		return nil, err
	}

	// Convert transactions to purchase history DTOs
	histories := s.transactionsToHistories(transactions)

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &dto.PurchaseHistoryList{
		Data:      histories,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPages,
	}, nil
}

// transactionsToHistories converts transaction entities to purchase history DTOs
func (s *purchaseService) transactionsToHistories(transactions []transactionEntity.Transaction) []dto.PurchaseHistory {
	histories := make([]dto.PurchaseHistory, 0, len(transactions))
	for _, txn := range transactions {
		hist := dto.PurchaseHistory{
			ID:           txn.ID,
			WalletID:     txn.WalletID,
			Amount:       txn.Amount,
			Phone:        txn.TargetNumber,
			Status:       txn.Status,
			SerialNumber: txn.SerialNumber,
			ProviderRef:  txn.PaymentProviderRef,
			CreatedAt:    txn.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if txn.ProductID != nil {
			hist.ProductID = *txn.ProductID
		}
		histories = append(histories, hist)
	}
	return histories
}

// ProcessWifiPurchase processes a wifi voucher purchase
func (s *purchaseService) ProcessWifiPurchase(ctx context.Context, token string, req *dto.PurchaseWifiRequest) (*dto.PurchaseWifiResponse, error) {
	// Validate request
	if req.ProductID == 0 {
		return nil, errors.New("product_id is required")
	}

	// Verify token and get user
	user, err := s.getUserFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Validate product exists
	product, err := s.getProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Check if product is active
	if !product.IsActive {
		return nil, errors.New("product is not active")
	}

	txID := uuid.New()
	var selectedVoucher *wifiVoucherEntity.WifiVoucher

	// Execute database operations within a transaction
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		selectedVoucher, err = s.processWifiTransaction(txCtx, user.ID, product, txID)
		return err
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info("WiFi purchase completed successfully", zap.String("tx_id", txID.String()), zap.Uint("voucher_id", selectedVoucher.ID))

	return &dto.PurchaseWifiResponse{
		TransactionID:   txID,
		VoucherID:       selectedVoucher.ID,
		VoucherCode:     selectedVoucher.Code,
		VoucherPassword: selectedVoucher.Password,
		Status:          transactionEntity.StatusSuccess,
		Message:         "WiFi voucher purchase successful",
	}, nil
}

// processWifiTransaction handles the wifi voucher transaction and assignment
func (s *purchaseService) processWifiTransaction(ctx context.Context, userID uint, product *productEntity.Product, txID uuid.UUID) (*wifiVoucherEntity.WifiVoucher, error) {
	// Get wallet with FOR UPDATE lock
	wallet, err := s.walletRepo.GetForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		s.logger.Error("Failed to get wallet", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Starting wifi purchase process", zap.String("tx_id", txID.String()), zap.Uint("wallet_id", wallet.ID), zap.Float64("balance", wallet.Balance), zap.Float64("price", product.PriceSell))

	// Check wallet balance
	if wallet.Balance < product.PriceSell {
		s.logger.Warn("Insufficient balance", zap.Uint("wallet_id", wallet.ID), zap.Float64("balance", wallet.Balance), zap.Float64("price", product.PriceSell))
		return nil, errors.New("insufficient wallet balance")
	}

	// Get available wifi voucher
	filter := &wifiVoucherDto.WifiVoucherFilter{
		Status:   wifiVoucherEntity.StatusAvailable,
		Page:     1,
		PageSize: 1,
	}
	vouchers, _, err := s.wifiVoucherRepo.GetByProviderAndDurationHours(ctx, product.ProviderID, product.DurationHours, filter)
	// vouchers, _, err := s.wifiVoucherRepo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get available wifi vouchers", zap.Error(err))
		return nil, errors.New("failed to get available vouchers")
	}

	if len(vouchers) == 0 {
		s.logger.Warn("No available wifi vouchers", zap.Uint("product_id", product.ID))
		return nil, errors.New("no available wifi vouchers in stock")
	}

	selectedVoucher := &vouchers[0]

	// Create pending transaction
	voucherID := selectedVoucher.ID
	txn := &transactionEntity.Transaction{
		ID:            txID,
		WalletID:      wallet.ID,
		Type:          transactionEntity.TypePurchase,
		Amount:        product.PriceSell,
		Status:        transactionEntity.StatusSuccess, // Direct success for wifi vouchers
		Description:   "WiFi voucher purchase",
		ProductID:     &product.ID,
		WifiVoucherID: &voucherID,
		TargetNumber:  "", // No phone for wifi purchases
		CreatedAt:     time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, txn); err != nil {
		s.logger.Error("Failed to create transaction", zap.String("tx_id", txID.String()), zap.Error(err))
		return nil, errors.New("failed to create transaction")
	}

	s.logger.Info("Transaction created for wifi voucher", zap.String("tx_id", txID.String()), zap.Float64("amount", product.PriceSell))

	// Update wifi voucher: mark as sold and assign to user
	now := time.Now()
	selectedVoucher.Status = wifiVoucherEntity.StatusSold
	selectedVoucher.SoldToUserID = &userID
	selectedVoucher.SoldAt = &now
	selectedVoucher.UpdatedAt = now

	if err := s.wifiVoucherRepo.Update(ctx, selectedVoucher); err != nil {
		s.logger.Error("Failed to update wifi voucher", zap.Uint("voucher_id", selectedVoucher.ID), zap.Error(err))
		return nil, errors.New("failed to assign voucher")
	}

	s.logger.Info("WiFi voucher assigned to user", zap.Uint("voucher_id", selectedVoucher.ID), zap.Uint("user_id", userID))

	// Deduct balance from wallet
	newBalance := wallet.Balance - product.PriceSell
	if err := s.walletRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		s.logger.Error("Failed to update wallet balance", zap.Uint("user_id", userID), zap.Error(err))
		return nil, errors.New("failed to update wallet balance")
	}

	s.logger.Info("Wallet balance updated for wifi purchase", zap.Uint("wallet_id", wallet.ID), zap.Float64("new_balance", newBalance))

	return selectedVoucher, nil
}
