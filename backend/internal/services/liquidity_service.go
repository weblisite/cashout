package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// LiquidityService handles all liquidity management operations
type LiquidityService struct {
	supabaseService *SupabaseService
	coreBanking     *CoreBankingService
}

// NewLiquidityService creates a new liquidity service
func NewLiquidityService(supabaseService *SupabaseService, coreBanking *CoreBankingService) *LiquidityService {
	return &LiquidityService{
		supabaseService: supabaseService,
		coreBanking:     coreBanking,
	}
}

// LiquiditySnapshot represents a liquidity snapshot
type LiquiditySnapshot struct {
	ID                    uuid.UUID `json:"id"`
	Timestamp             time.Time `json:"timestamp"`
	TotalAssets           float64   `json:"total_assets"`
	TotalLiabilities      float64   `json:"total_liabilities"`
	NetLiquidity          float64   `json:"net_liquidity"`
	ReserveRatio          float64   `json:"reserve_ratio"`
	CashFlowIn            float64   `json:"cash_flow_in"`
	CashFlowOut           float64   `json:"cash_flow_out"`
	NetCashFlow           float64   `json:"net_cash_flow"`
	RiskLevel             string    `json:"risk_level"` // low, medium, high, critical
	AlertThreshold        float64   `json:"alert_threshold"`
	MinimumReserve        float64   `json:"minimum_reserve"`
	CreatedAt             time.Time `json:"created_at"`
}

// CashFlow represents cash flow data
type CashFlow struct {
	ID            uuid.UUID `json:"id"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"` // inflow, outflow
	Category      string    `json:"category"` // transactions, settlements, fees, operational
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Reserve represents reserve management
type Reserve struct {
	ID            uuid.UUID `json:"id"`
	Type          string    `json:"type"` // operational, regulatory, emergency
	Amount        float64   `json:"amount"`
	TargetAmount  float64   `json:"target_amount"`
	Status        string    `json:"status"` // adequate, low, critical
	LastUpdated   time.Time `json:"last_updated"`
	CreatedAt     time.Time `json:"created_at"`
}

// LiquidityAlert represents a liquidity alert
type LiquidityAlert struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"` // low_reserve, high_outflow, risk_threshold
	Severity    string    `json:"severity"` // low, medium, high, critical
	Message     string    `json:"message"`
	CurrentValue float64  `json:"current_value"`
	Threshold   float64   `json:"threshold"`
	Status      string    `json:"status"` // active, resolved
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// Liquidity Methods
func (l *LiquidityService) CreateLiquiditySnapshot(ctx context.Context) (*LiquiditySnapshot, error) {
	// Calculate total assets (all account balances)
	totalAssets, err := l.calculateTotalAssets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total assets: %w", err)
	}

	// Calculate total liabilities (pending settlements + operational costs)
	totalLiabilities, err := l.calculateTotalLiabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total liabilities: %w", err)
	}

	// Calculate net liquidity
	netLiquidity := totalAssets - totalLiabilities

	// Calculate reserve ratio
	reserveRatio := 0.0
	if totalAssets > 0 {
		reserveRatio = netLiquidity / totalAssets
	}

	// Calculate cash flow
	cashFlowIn, cashFlowOut, err := l.calculateCashFlow(ctx, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to calculate cash flow: %w", err)
	}

	netCashFlow := cashFlowIn - cashFlowOut

	// Determine risk level
	riskLevel := l.determineRiskLevel(netLiquidity, reserveRatio, netCashFlow)

	// Get configuration
	alertThreshold := l.getAlertThreshold()
	minimumReserve := l.getMinimumReserve()

	snapshot := &LiquiditySnapshot{
		ID:               uuid.New(),
		Timestamp:        time.Now(),
		TotalAssets:      totalAssets,
		TotalLiabilities: totalLiabilities,
		NetLiquidity:     netLiquidity,
		ReserveRatio:     reserveRatio,
		CashFlowIn:       cashFlowIn,
		CashFlowOut:      cashFlowOut,
		NetCashFlow:      netCashFlow,
		RiskLevel:        riskLevel,
		AlertThreshold:   alertThreshold,
		MinimumReserve:   minimumReserve,
		CreatedAt:        time.Now(),
	}

	if err := l.supabaseService.CreateLiquiditySnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("failed to create liquidity snapshot: %w", err)
	}

	// Check for alerts
	if err := l.checkLiquidityAlerts(ctx, snapshot); err != nil {
		log.Printf("Failed to check liquidity alerts: %v", err)
	}

	log.Printf("Created liquidity snapshot: Net Liquidity: %.2f, Risk Level: %s", 
		netLiquidity, riskLevel)
	return snapshot, nil
}

func (l *LiquidityService) calculateTotalAssets(ctx context.Context) (float64, error) {
	// Get all active accounts and sum their balances
	accounts, err := l.supabaseService.GetAllActiveAccounts(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get accounts: %w", err)
	}

	var totalAssets float64
	for _, account := range accounts {
		totalAssets += account.Balance
	}

	return totalAssets, nil
}

func (l *LiquidityService) calculateTotalLiabilities(ctx context.Context) (float64, error) {
	// Get pending settlements
	settlements, err := l.supabaseService.GetSettlementsByStatus(ctx, "pending")
	if err != nil {
		return 0, fmt.Errorf("failed to get pending settlements: %w", err)
	}

	var totalLiabilities float64
	for _, settlement := range settlements {
		totalLiabilities += settlement.Amount
	}

	// Add operational costs (estimated)
	operationalCosts := l.getOperationalCosts()
	totalLiabilities += operationalCosts

	return totalLiabilities, nil
}

func (l *LiquidityService) calculateCashFlow(ctx context.Context, startTime, endTime time.Time) (float64, float64, error) {
	// Get cash flow records for the period
	cashFlows, err := l.supabaseService.GetCashFlowsByPeriod(ctx, startTime, endTime)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get cash flows: %w", err)
	}

	var cashFlowIn, cashFlowOut float64
	for _, flow := range cashFlows {
		if flow.Type == "inflow" {
			cashFlowIn += flow.Amount
		} else {
			cashFlowOut += flow.Amount
		}
	}

	return cashFlowIn, cashFlowOut, nil
}

func (l *LiquidityService) determineRiskLevel(netLiquidity, reserveRatio, netCashFlow float64) string {
	// Determine risk level based on multiple factors
	if netLiquidity < 0 || reserveRatio < 0.1 {
		return "critical"
	} else if netLiquidity < 1000000 || reserveRatio < 0.2 {
		return "high"
	} else if netLiquidity < 5000000 || reserveRatio < 0.3 {
		return "medium"
	} else {
		return "low"
	}
}

func (l *LiquidityService) checkLiquidityAlerts(ctx context.Context, snapshot *LiquiditySnapshot) error {
	// Check for low reserve alert
	if snapshot.NetLiquidity < snapshot.AlertThreshold {
		alert := &LiquidityAlert{
			ID:           uuid.New(),
			Type:         "low_reserve",
			Severity:     "high",
			Message:      fmt.Sprintf("Low liquidity reserve: %.2f", snapshot.NetLiquidity),
			CurrentValue: snapshot.NetLiquidity,
			Threshold:    snapshot.AlertThreshold,
			Status:       "active",
			CreatedAt:    time.Now(),
		}

		if err := l.supabaseService.CreateLiquidityAlert(ctx, alert); err != nil {
			return fmt.Errorf("failed to create liquidity alert: %w", err)
		}

		log.Printf("⚠️ Liquidity alert: Low reserve - %.2f", snapshot.NetLiquidity)
	}

	// Check for high outflow alert
	if snapshot.CashFlowOut > snapshot.CashFlowIn*2 {
		alert := &LiquidityAlert{
			ID:           uuid.New(),
			Type:         "high_outflow",
			Severity:     "medium",
			Message:      fmt.Sprintf("High cash outflow: %.2f", snapshot.CashFlowOut),
			CurrentValue: snapshot.CashFlowOut,
			Threshold:    snapshot.CashFlowIn * 2,
			Status:       "active",
			CreatedAt:    time.Now(),
		}

		if err := l.supabaseService.CreateLiquidityAlert(ctx, alert); err != nil {
			return fmt.Errorf("failed to create liquidity alert: %w", err)
		}

		log.Printf("⚠️ Liquidity alert: High outflow - %.2f", snapshot.CashFlowOut)
	}

	return nil
}

// Reserve Management Methods
func (l *LiquidityService) CreateReserve(ctx context.Context, reserveType string, targetAmount float64) (*Reserve, error) {
	reserve := &Reserve{
		ID:           uuid.New(),
		Type:         reserveType,
		Amount:       0.0,
		TargetAmount: targetAmount,
		Status:       "low",
		LastUpdated:  time.Now(),
		CreatedAt:    time.Now(),
	}

	if err := l.supabaseService.CreateReserve(ctx, reserve); err != nil {
		return nil, fmt.Errorf("failed to create reserve: %w", err)
	}

	log.Printf("Created %s reserve with target: %.2f", reserveType, targetAmount)
	return reserve, nil
}

func (l *LiquidityService) UpdateReserve(ctx context.Context, reserveID uuid.UUID, amount float64) error {
	reserve, err := l.supabaseService.GetReserveByID(ctx, reserveID)
	if err != nil {
		return fmt.Errorf("failed to get reserve: %w", err)
	}

	reserve.Amount = amount
	reserve.LastUpdated = time.Now()

	// Update status based on amount vs target
	if reserve.Amount >= reserve.TargetAmount {
		reserve.Status = "adequate"
	} else if reserve.Amount >= reserve.TargetAmount*0.5 {
		reserve.Status = "low"
	} else {
		reserve.Status = "critical"
	}

	if err := l.supabaseService.UpdateReserve(ctx, reserve); err != nil {
		return fmt.Errorf("failed to update reserve: %w", err)
	}

	log.Printf("Updated %s reserve: %.2f (%.1f%% of target)", 
		reserve.Type, amount, (amount/reserve.TargetAmount)*100)
	return nil
}

func (l *LiquidityService) GetReserves(ctx context.Context) ([]Reserve, error) {
	return l.supabaseService.GetAllReserves(ctx)
}

// Cash Flow Methods
func (l *LiquidityService) RecordCashFlow(ctx context.Context, flowType, category string, amount float64, description string, accountID, transactionID *uuid.UUID) error {
	cashFlow := &CashFlow{
		ID:            uuid.New(),
		Date:          time.Now(),
		Type:          flowType,
		Category:      category,
		Amount:        amount,
		Description:   description,
		AccountID:     accountID,
		TransactionID: transactionID,
		CreatedAt:     time.Now(),
	}

	if err := l.supabaseService.CreateCashFlow(ctx, cashFlow); err != nil {
		return fmt.Errorf("failed to create cash flow: %w", err)
	}

	log.Printf("Recorded cash flow: %s %s %.2f - %s", flowType, category, amount, description)
	return nil
}

func (l *LiquidityService) GetCashFlowReport(ctx context.Context, startDate, endDate time.Time) ([]CashFlow, error) {
	return l.supabaseService.GetCashFlowsByPeriod(ctx, startDate, endDate)
}

// Alert Management Methods
func (l *LiquidityService) GetActiveAlerts(ctx context.Context) ([]LiquidityAlert, error) {
	return l.supabaseService.GetLiquidityAlertsByStatus(ctx, "active")
}

func (l *LiquidityService) ResolveAlert(ctx context.Context, alertID uuid.UUID) error {
	alert, err := l.supabaseService.GetLiquidityAlertByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}

	now := time.Now()
	alert.Status = "resolved"
	alert.ResolvedAt = &now

	if err := l.supabaseService.UpdateLiquidityAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	log.Printf("Resolved liquidity alert: %s", alert.Message)
	return nil
}

// Emergency Procedures
func (l *LiquidityService) TriggerEmergencyProcedures(ctx context.Context) error {
	log.Println("🚨 Triggering emergency liquidity procedures")

	// 1. Suspend large transactions
	if err := l.suspendLargeTransactions(ctx); err != nil {
		log.Printf("Failed to suspend large transactions: %v", err)
	}

	// 2. Increase reserve requirements
	if err := l.increaseReserveRequirements(ctx); err != nil {
		log.Printf("Failed to increase reserve requirements: %v", err)
	}

	// 3. Notify stakeholders
	if err := l.notifyStakeholders(ctx); err != nil {
		log.Printf("Failed to notify stakeholders: %v", err)
	}

	log.Println("✅ Emergency procedures triggered")
	return nil
}

func (l *LiquidityService) suspendLargeTransactions(ctx context.Context) error {
	// Implementation to suspend transactions above certain threshold
	log.Println("Suspending transactions above 100,000 KES")
	return nil
}

func (l *LiquidityService) increaseReserveRequirements(ctx context.Context) error {
	// Implementation to increase reserve requirements
	log.Println("Increasing reserve requirements to 50%")
	return nil
}

func (l *LiquidityService) notifyStakeholders(ctx context.Context) error {
	// Implementation to notify stakeholders
	log.Println("Notifying stakeholders of liquidity emergency")
	return nil
}

// Configuration Methods
func (l *LiquidityService) getAlertThreshold() float64 {
	// In production, this would come from environment variables
	return 1000000 // 1M KES
}

func (l *LiquidityService) getMinimumReserve() float64 {
	// In production, this would come from environment variables
	return 5000000 // 5M KES
}

func (l *LiquidityService) getOperationalCosts() float64 {
	// In production, this would be calculated from actual operational data
	return 100000 // 100K KES daily
}

// Scheduled Monitoring
func (l *LiquidityService) MonitorLiquidity(ctx context.Context) error {
	// Create liquidity snapshot
	if _, err := l.CreateLiquiditySnapshot(ctx); err != nil {
		return fmt.Errorf("failed to create liquidity snapshot: %w", err)
	}

	// Check for critical alerts
	alerts, err := l.GetActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active alerts: %w", err)
	}

	for _, alert := range alerts {
		if alert.Severity == "critical" {
			if err := l.TriggerEmergencyProcedures(ctx); err != nil {
				log.Printf("Failed to trigger emergency procedures: %v", err)
			}
			break
		}
	}

	return nil
}

// Health Check
func (l *LiquidityService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := l.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("liquidity service health check failed: %w", err)
	}

	log.Println("✅ Liquidity service health check passed")
	return nil
} 