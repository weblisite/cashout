package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// SettlementService handles all settlement operations
type SettlementService struct {
	supabaseService *SupabaseService
	coreBanking     *CoreBankingService
}

// NewSettlementService creates a new settlement service
func NewSettlementService(supabaseService *SupabaseService, coreBanking *CoreBankingService) *SettlementService {
	return &SettlementService{
		supabaseService: supabaseService,
		coreBanking:     coreBanking,
	}
}

// Settlement represents a settlement record
type Settlement struct {
	ID              uuid.UUID `json:"id"`
	SettlementID    string    `json:"settlement_id"`
	AccountID       uuid.UUID `json:"account_id"`
	AccountType     string    `json:"account_type"` // agent, business
	SettlementType  string    `json:"settlement_type"` // daily, weekly, monthly
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"` // pending, processing, completed, failed
	BankAccount     string    `json:"bank_account"`
	BankName        string    `json:"bank_name"`
	Reference       string    `json:"reference"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SettlementTransaction represents a transaction in a settlement
type SettlementTransaction struct {
	ID           uuid.UUID `json:"id"`
	SettlementID uuid.UUID `json:"settlement_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Amount       float64   `json:"amount"`
	Commission   float64   `json:"commission"`
	NetAmount    float64   `json:"net_amount"`
	CreatedAt    time.Time `json:"created_at"`
}

// Float represents agent float management
type Float struct {
	ID            uuid.UUID `json:"id"`
	AgentID       uuid.UUID `json:"agent_id"`
	CurrentFloat  float64   `json:"current_float"`
	MaxFloat      float64   `json:"max_float"`
	MinThreshold  float64   `json:"min_threshold"`
	LastReplenish time.Time `json:"last_replenish"`
	Status        string    `json:"status"` // active, suspended, low
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Settlement Methods
func (s *SettlementService) CreateSettlement(ctx context.Context, accountID uuid.UUID, accountType, settlementType string, amount float64, bankAccount, bankName string) (*Settlement, error) {
	settlement := &Settlement{
		ID:             uuid.New(),
		SettlementID:   s.generateSettlementID(),
		AccountID:      accountID,
		AccountType:    accountType,
		SettlementType: settlementType,
		Amount:         amount,
		Currency:       "KES",
		Status:         "pending",
		BankAccount:    bankAccount,
		BankName:       bankName,
		Reference:      s.generateReference(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.supabaseService.CreateSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("failed to create settlement: %w", err)
	}

	log.Printf("Created settlement %s for account %s, Amount: %.2f", 
		settlement.SettlementID, accountID, amount)
	return settlement, nil
}

func (s *SettlementService) ProcessSettlement(ctx context.Context, settlementID uuid.UUID) error {
	settlement, err := s.supabaseService.GetSettlementByID(ctx, settlementID)
	if err != nil {
		return fmt.Errorf("failed to get settlement: %w", err)
	}

	if settlement.Status != "pending" {
		return fmt.Errorf("settlement is not in pending status")
	}

	// Update settlement status to processing
	settlement.Status = "processing"
	settlement.UpdatedAt = time.Now()

	if err := s.supabaseService.UpdateSettlement(ctx, settlement); err != nil {
		return fmt.Errorf("failed to update settlement status: %w", err)
	}

	// Process the actual bank transfer (simulated)
	if err := s.processBankTransfer(ctx, settlement); err != nil {
		settlement.Status = "failed"
		settlement.UpdatedAt = time.Now()
		s.supabaseService.UpdateSettlement(ctx, settlement)
		return fmt.Errorf("failed to process bank transfer: %w", err)
	}

	// Update settlement status to completed
	now := time.Now()
	settlement.Status = "completed"
	settlement.ProcessedAt = &now
	settlement.UpdatedAt = now

	if err := s.supabaseService.UpdateSettlement(ctx, settlement); err != nil {
		return fmt.Errorf("failed to update settlement: %w", err)
	}

	// Deduct funds from account
	if err := s.coreBanking.DeductFunds(ctx, settlement.AccountID, settlement.Amount, 
		fmt.Sprintf("Settlement %s", settlement.SettlementID)); err != nil {
		return fmt.Errorf("failed to deduct settlement funds: %w", err)
	}

	log.Printf("Processed settlement %s, Amount: %.2f", settlement.SettlementID, settlement.Amount)
	return nil
}

func (s *SettlementService) processBankTransfer(ctx context.Context, settlement *Settlement) error {
	// Simulate bank transfer processing
	// In production, this would integrate with actual banking APIs
	log.Printf("Processing bank transfer: %s -> %s (%s), Amount: %.2f", 
		settlement.BankName, settlement.BankAccount, settlement.Reference, settlement.Amount)
	
	// Simulate processing time
	time.Sleep(2 * time.Second)
	
	return nil
}

func (s *SettlementService) GetSettlementsByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]Settlement, error) {
	return s.supabaseService.GetSettlementsByAccount(ctx, accountID, limit, offset)
}

func (s *SettlementService) GetPendingSettlements(ctx context.Context) ([]Settlement, error) {
	return s.supabaseService.GetSettlementsByStatus(ctx, "pending")
}

// Float Management Methods
func (s *SettlementService) CreateFloat(ctx context.Context, agentID uuid.UUID, maxFloat, minThreshold float64) (*Float, error) {
	float := &Float{
		ID:            uuid.New(),
		AgentID:       agentID,
		CurrentFloat:  0.0,
		MaxFloat:      maxFloat,
		MinThreshold:  minThreshold,
		LastReplenish: time.Now(),
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.supabaseService.CreateFloat(ctx, float); err != nil {
		return nil, fmt.Errorf("failed to create float: %w", err)
	}

	log.Printf("Created float for agent %s, Max: %.2f, Min: %.2f", agentID, maxFloat, minThreshold)
	return float, nil
}

func (s *SettlementService) UpdateFloat(ctx context.Context, agentID uuid.UUID, amount float64, operation string) error {
	float, err := s.supabaseService.GetFloatByAgentID(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get float: %w", err)
	}

	switch operation {
	case "add":
		float.CurrentFloat += amount
	case "deduct":
		if float.CurrentFloat < amount {
			return fmt.Errorf("insufficient float balance")
		}
		float.CurrentFloat -= amount
	default:
		return fmt.Errorf("invalid operation: %s", operation)
	}

	// Update float status based on current balance
	if float.CurrentFloat <= float.MinThreshold {
		float.Status = "low"
	} else if float.CurrentFloat >= float.MaxFloat {
		float.Status = "suspended"
	} else {
		float.Status = "active"
	}

	float.UpdatedAt = time.Now()

	if err := s.supabaseService.UpdateFloat(ctx, float); err != nil {
		return fmt.Errorf("failed to update float: %w", err)
	}

	log.Printf("Updated float for agent %s: %s %.2f, New balance: %.2f", 
		agentID, operation, amount, float.CurrentFloat)
	return nil
}

func (s *SettlementService) ReplenishFloat(ctx context.Context, agentID uuid.UUID, amount float64) error {
	float, err := s.supabaseService.GetFloatByAgentID(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get float: %w", err)
	}

	// Add funds to agent's account
	if err := s.coreBanking.AddFunds(ctx, agentID, amount, "Float replenishment"); err != nil {
		return fmt.Errorf("failed to add funds to agent account: %w", err)
	}

	// Update float
	float.CurrentFloat += amount
	float.LastReplenish = time.Now()
	float.Status = "active"
	float.UpdatedAt = time.Now()

	if err := s.supabaseService.UpdateFloat(ctx, float); err != nil {
		return fmt.Errorf("failed to update float: %w", err)
	}

	log.Printf("Replenished float for agent %s: %.2f, New balance: %.2f", 
		agentID, amount, float.CurrentFloat)
	return nil
}

func (s *SettlementService) GetFloatByAgentID(ctx context.Context, agentID uuid.UUID) (*Float, error) {
	return s.supabaseService.GetFloatByAgentID(ctx, agentID)
}

func (s *SettlementService) GetLowFloatAgents(ctx context.Context) ([]Float, error) {
	return s.supabaseService.GetFloatsByStatus(ctx, "low")
}

// Reconciliation Methods
func (s *SettlementService) ReconcileTransactions(ctx context.Context, settlementID uuid.UUID) error {
	settlement, err := s.supabaseService.GetSettlementByID(ctx, settlementID)
	if err != nil {
		return fmt.Errorf("failed to get settlement: %w", err)
	}

	// Get all transactions for this settlement period
	transactions, err := s.supabaseService.GetTransactionsForSettlement(ctx, settlement.AccountID, settlement.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	var totalAmount float64
	var totalCommission float64

	for _, transaction := range transactions {
		// Create settlement transaction record
		settlementTx := &SettlementTransaction{
			ID:           uuid.New(),
			SettlementID: settlementID,
			TransactionID: transaction.ID,
			Amount:       transaction.Amount,
			Commission:   transaction.Fee,
			NetAmount:    transaction.Amount - transaction.Fee,
			CreatedAt:    time.Now(),
		}

		if err := s.supabaseService.CreateSettlementTransaction(ctx, settlementTx); err != nil {
			return fmt.Errorf("failed to create settlement transaction: %w", err)
		}

		totalAmount += transaction.Amount
		totalCommission += transaction.Fee
	}

	// Verify reconciliation
	if totalAmount != settlement.Amount {
		return fmt.Errorf("reconciliation failed: expected %.2f, got %.2f", 
			settlement.Amount, totalAmount)
	}

	log.Printf("Reconciled settlement %s: %d transactions, Total: %.2f, Commission: %.2f", 
		settlement.SettlementID, len(transactions), totalAmount, totalCommission)
	return nil
}

// Utility Methods
func (s *SettlementService) generateSettlementID() string {
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	return "SET" + timestamp + random
}

func (s *SettlementService) generateReference() string {
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
	return "REF" + timestamp + random
}

// Scheduled Settlement Processing
func (s *SettlementService) ProcessScheduledSettlements(ctx context.Context) error {
	// Get all pending settlements
	settlements, err := s.GetPendingSettlements(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending settlements: %w", err)
	}

	for _, settlement := range settlements {
		if err := s.ProcessSettlement(ctx, settlement.ID); err != nil {
			log.Printf("Failed to process settlement %s: %v", settlement.SettlementID, err)
			continue
		}
	}

	log.Printf("Processed %d scheduled settlements", len(settlements))
	return nil
}

// Health Check
func (s *SettlementService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := s.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("settlement service health check failed: %w", err)
	}

	log.Println("✅ Settlement service health check passed")
	return nil
} 