package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// CoreBankingService handles all core banking operations
type CoreBankingService struct {
	supabaseService *SupabaseService
}

// NewCoreBankingService creates a new core banking service
func NewCoreBankingService(supabaseService *SupabaseService) *CoreBankingService {
	return &CoreBankingService{
		supabaseService: supabaseService,
	}
}

// Account represents a banking account
type Account struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccountType   string    `json:"account_type"` // user, agent, business
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // active, suspended, closed
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Transaction represents a banking transaction
type BankingTransaction struct {
	ID              uuid.UUID `json:"id"`
	TransactionID   string    `json:"transaction_id"`
	FromAccountID   uuid.UUID `json:"from_account_id"`
	ToAccountID     uuid.UUID `json:"to_account_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionType string    `json:"transaction_type"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

// LedgerEntry represents a ledger entry
type LedgerEntry struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     uuid.UUID `json:"account_id"`
	EntryType     string    `json:"entry_type"` // debit, credit
	Amount        float64   `json:"amount"`
	Balance       float64   `json:"balance"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

// Account Management Methods
func (c *CoreBankingService) CreateAccount(ctx context.Context, userID uuid.UUID, accountType string) (*Account, error) {
	account := &Account{
		ID:            uuid.New(),
		UserID:        userID,
		AccountType:   accountType,
		AccountNumber: c.generateAccountNumber(accountType),
		Balance:       0.0,
		Currency:      "KES",
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store account in database
	if err := c.supabaseService.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	log.Printf("Created account %s for user %s", account.AccountNumber, userID)
	return account, nil
}

func (c *CoreBankingService) GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	return c.supabaseService.GetAccountByID(ctx, accountID)
}

func (c *CoreBankingService) GetAccountByUserID(ctx context.Context, userID uuid.UUID) (*Account, error) {
	return c.supabaseService.GetAccountByUserID(ctx, userID)
}

func (c *CoreBankingService) UpdateAccountBalance(ctx context.Context, accountID uuid.UUID, newBalance float64) error {
	account, err := c.GetAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	account.Balance = newBalance
	account.UpdatedAt = time.Now()

	return c.supabaseService.UpdateAccount(ctx, account)
}

// Transaction Processing Methods
func (c *CoreBankingService) ProcessTransaction(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, transactionType, description string) (*BankingTransaction, error) {
	// Validate accounts
	fromAccount, err := c.GetAccount(ctx, fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}

	toAccount, err := c.GetAccount(ctx, toAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get to account: %w", err)
	}

	// Validate transaction
	if err := c.validateTransaction(fromAccount, toAccount, amount); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	// Create transaction
	transaction := &BankingTransaction{
		ID:              uuid.New(),
		TransactionID:   c.generateTransactionID(),
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		Amount:          amount,
		Currency:        "KES",
		TransactionType: transactionType,
		Status:          "processing",
		Description:     description,
		CreatedAt:       time.Now(),
	}

	// Process transaction atomically
	if err := c.processTransactionAtomically(ctx, transaction, fromAccount, toAccount); err != nil {
		return nil, fmt.Errorf("failed to process transaction: %w", err)
	}

	log.Printf("Processed transaction %s: %s -> %s, Amount: %.2f", 
		transaction.TransactionID, fromAccount.AccountNumber, toAccount.AccountNumber, amount)

	return transaction, nil
}

func (c *CoreBankingService) validateTransaction(fromAccount, toAccount *Account, amount float64) error {
	// Check if accounts are active
	if fromAccount.Status != "active" {
		return fmt.Errorf("from account is not active")
	}
	if toAccount.Status != "active" {
		return fmt.Errorf("to account is not active")
	}

	// Check if from account has sufficient balance
	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// Check transaction limits
	if amount < 50 {
		return fmt.Errorf("minimum transaction amount is 50")
	}
	if amount > 1000000 {
		return fmt.Errorf("maximum transaction amount is 1,000,000")
	}

	return nil
}

func (c *CoreBankingService) processTransactionAtomically(ctx context.Context, transaction *BankingTransaction, fromAccount, toAccount *Account) error {
	// Update account balances
	fromAccount.Balance -= transaction.Amount
	fromAccount.UpdatedAt = time.Now()

	toAccount.Balance += transaction.Amount
	toAccount.UpdatedAt = time.Now()

	// Update transaction status
	transaction.Status = "completed"

	// Save all changes atomically
	if err := c.supabaseService.ProcessTransactionAtomically(ctx, transaction, fromAccount, toAccount); err != nil {
		return fmt.Errorf("failed to save transaction atomically: %w", err)
	}

	return nil
}

// Ledger Management Methods
func (c *CoreBankingService) CreateLedgerEntry(ctx context.Context, transactionID, accountID uuid.UUID, entryType string, amount, balance float64, description string) error {
	entry := &LedgerEntry{
		ID:            uuid.New(),
		TransactionID: transactionID,
		AccountID:     accountID,
		EntryType:     entryType,
		Amount:        amount,
		Balance:       balance,
		Description:   description,
		CreatedAt:     time.Now(),
	}

	return c.supabaseService.CreateLedgerEntry(ctx, entry)
}

func (c *CoreBankingService) GetAccountLedger(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]LedgerEntry, error) {
	return c.supabaseService.GetAccountLedger(ctx, accountID, limit, offset)
}

// Utility Methods
func (c *CoreBankingService) generateAccountNumber(accountType string) string {
	// Generate unique account number based on type
	prefix := "CA"
	switch accountType {
	case "user":
		prefix = "CU"
	case "agent":
		prefix = "CA"
	case "business":
		prefix = "CB"
	}

	// Generate 8-digit number
	number := fmt.Sprintf("%08d", time.Now().UnixNano()%100000000)
	return prefix + number
}

func (c *CoreBankingService) generateTransactionID() string {
	// Generate unique transaction ID
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	return "TXN" + timestamp + random
}

// Balance Management Methods
func (c *CoreBankingService) GetAccountBalance(ctx context.Context, accountID uuid.UUID) (float64, error) {
	account, err := c.GetAccount(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}
	return account.Balance, nil
}

func (c *CoreBankingService) AddFunds(ctx context.Context, accountID uuid.UUID, amount float64, description string) error {
	account, err := c.GetAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Create internal transaction (system to account)
	systemAccountID := uuid.MustParse("00000000-0000-0000-0000-000000000000") // System account
	transaction := &BankingTransaction{
		ID:              uuid.New(),
		TransactionID:   c.generateTransactionID(),
		FromAccountID:   systemAccountID,
		ToAccountID:     accountID,
		Amount:          amount,
		Currency:        "KES",
		TransactionType: "fund_addition",
		Status:          "completed",
		Description:     description,
		CreatedAt:       time.Now(),
	}

	// Update account balance
	account.Balance += amount
	account.UpdatedAt = time.Now()

	// Save transaction and account update
	return c.supabaseService.ProcessTransactionAtomically(ctx, transaction, nil, account)
}

func (c *CoreBankingService) DeductFunds(ctx context.Context, accountID uuid.UUID, amount float64, description string) error {
	account, err := c.GetAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Validate sufficient balance
	if account.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// Create internal transaction (account to system)
	systemAccountID := uuid.MustParse("00000000-0000-0000-0000-000000000000") // System account
	transaction := &BankingTransaction{
		ID:              uuid.New(),
		TransactionID:   c.generateTransactionID(),
		FromAccountID:   accountID,
		ToAccountID:     systemAccountID,
		Amount:          amount,
		Currency:        "KES",
		TransactionType: "fund_deduction",
		Status:          "completed",
		Description:     description,
		CreatedAt:       time.Now(),
	}

	// Update account balance
	account.Balance -= amount
	account.UpdatedAt = time.Now()

	// Save transaction and account update
	return c.supabaseService.ProcessTransactionAtomically(ctx, transaction, account, nil)
}

// Health Check
func (c *CoreBankingService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := c.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("core banking health check failed: %w", err)
	}

	log.Println("✅ Core banking service health check passed")
	return nil
} 