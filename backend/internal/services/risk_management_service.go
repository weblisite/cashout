package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
)

// RiskManagementService handles all risk management operations
type RiskManagementService struct {
	supabaseService *SupabaseService
	coreBanking     *CoreBankingService
}

// NewRiskManagementService creates a new risk management service
func NewRiskManagementService(supabaseService *SupabaseService, coreBanking *CoreBankingService) *RiskManagementService {
	return &RiskManagementService{
		supabaseService: supabaseService,
		coreBanking:     coreBanking,
	}
}

// RiskScore represents a risk assessment
type RiskScore struct {
	ID            uuid.UUID `json:"id"`
	AccountID     uuid.UUID `json:"account_id"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	Score         float64   `json:"score"` // 0.0 to 1.0
	RiskLevel     string    `json:"risk_level"` // low, medium, high, critical
	Factors       []string  `json:"factors"`
	CreatedAt     time.Time `json:"created_at"`
}

// FraudAlert represents a fraud detection alert
type FraudAlert struct {
	ID            uuid.UUID `json:"id"`
	AlertType     string    `json:"alert_type"`
	Severity      string    `json:"severity"` // low, medium, high, critical
	AccountID     uuid.UUID `json:"account_id"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	Description   string    `json:"description"`
	RiskScore     float64   `json:"risk_score"`
	Status        string    `json:"status"` // active, investigated, resolved, false_positive
	CreatedAt     time.Time `json:"created_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
}

// TransactionLimit represents transaction limits
type TransactionLimit struct {
	ID            uuid.UUID `json:"id"`
	AccountID     uuid.UUID `json:"account_id"`
	LimitType     string    `json:"limit_type"` // daily, weekly, monthly, per_transaction
	Amount        float64   `json:"amount"`
	UsedAmount    float64   `json:"used_amount"`
	ResetDate     time.Time `json:"reset_date"`
	Status        string    `json:"status"` // active, suspended, exceeded
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VelocityMonitor represents velocity monitoring
type VelocityMonitor struct {
	ID            uuid.UUID `json:"id"`
	AccountID     uuid.UUID `json:"account_id"`
	TimeWindow    string    `json:"time_window"` // 1h, 24h, 7d, 30d
	TransactionCount int    `json:"transaction_count"`
	TotalAmount   float64   `json:"total_amount"`
	AverageAmount float64   `json:"average_amount"`
	RiskLevel     string    `json:"risk_level"`
	LastUpdated   time.Time `json:"last_updated"`
	CreatedAt     time.Time `json:"created_at"`
}

// GeographicAlert represents geographic monitoring
type GeographicAlert struct {
	ID            uuid.UUID `json:"id"`
	AccountID     uuid.UUID `json:"account_id"`
	Location      string    `json:"location"`
	RiskType      string    `json:"risk_type"` // unusual_location, high_risk_region
	Description   string    `json:"description"`
	RiskScore     float64   `json:"risk_score"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Risk Management Methods
func (r *RiskManagementService) AssessTransactionRisk(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, location string) (*RiskScore, error) {
	var riskFactors []string
	var totalScore float64

	// Factor 1: Transaction amount risk (0-0.3)
	amountScore := r.calculateAmountRisk(amount)
	totalScore += amountScore
	if amountScore > 0.2 {
		riskFactors = append(riskFactors, "high_transaction_amount")
	}

	// Factor 2: Account history risk (0-0.2)
	historyScore, err := r.calculateHistoryRisk(ctx, fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate history risk: %w", err)
	}
	totalScore += historyScore
	if historyScore > 0.15 {
		riskFactors = append(riskFactors, "suspicious_account_history")
	}

	// Factor 3: Velocity risk (0-0.2)
	velocityScore, err := r.calculateVelocityRisk(ctx, fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate velocity risk: %w", err)
	}
	totalScore += velocityScore
	if velocityScore > 0.15 {
		riskFactors = append(riskFactors, "high_transaction_velocity")
	}

	// Factor 4: Geographic risk (0-0.15)
	geoScore := r.calculateGeographicRisk(location)
	totalScore += geoScore
	if geoScore > 0.1 {
		riskFactors = append(riskFactors, "high_risk_location")
	}

	// Factor 5: Time-based risk (0-0.15)
	timeScore := r.calculateTimeRisk()
	totalScore += timeScore
	if timeScore > 0.1 {
		riskFactors = append(riskFactors, "unusual_transaction_time")
	}

	// Ensure score is between 0 and 1
	totalScore = math.Min(totalScore, 1.0)

	// Determine risk level
	riskLevel := r.determineRiskLevel(totalScore)

	riskScore := &RiskScore{
		ID:        uuid.New(),
		AccountID: fromAccountID,
		Score:     totalScore,
		RiskLevel: riskLevel,
		Factors:   riskFactors,
		CreatedAt: time.Now(),
	}

	if err := r.supabaseService.CreateRiskScore(ctx, riskScore); err != nil {
		return nil, fmt.Errorf("failed to create risk score: %w", err)
	}

	log.Printf("Risk assessment: Score %.2f, Level %s, Factors: %v", totalScore, riskLevel, riskFactors)
	return riskScore, nil
}

func (r *RiskManagementService) calculateAmountRisk(amount float64) float64 {
	// Higher amounts have higher risk
	if amount > 1000000 {
		return 0.3
	} else if amount > 500000 {
		return 0.2
	} else if amount > 100000 {
		return 0.1
	} else if amount > 50000 {
		return 0.05
	}
	return 0.0
}

func (r *RiskManagementService) calculateHistoryRisk(ctx context.Context, accountID uuid.UUID) (float64, error) {
	// Get account transaction history
	transactions, err := r.supabaseService.GetTransactionsByAccount(ctx, accountID, 100, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	if len(transactions) == 0 {
		return 0.1 // New account risk
	}

	var riskScore float64
	var failedTransactions int
	var chargebacks int

	for _, tx := range transactions {
		if tx.Status == "failed" {
			failedTransactions++
		}
		// Check for chargebacks (would need additional logic)
	}

	// Calculate risk based on failed transactions
	if failedTransactions > 0 {
		failureRate := float64(failedTransactions) / float64(len(transactions))
		riskScore = failureRate * 0.2
	}

	return riskScore, nil
}

func (r *RiskManagementService) calculateVelocityRisk(ctx context.Context, accountID uuid.UUID) (float64, error) {
	// Get recent transactions (last 24 hours)
	recentTransactions, err := r.supabaseService.GetTransactionsByAccountAndPeriod(ctx, accountID, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to get recent transactions: %w", err)
	}

	if len(recentTransactions) == 0 {
		return 0.0
	}

	// Calculate velocity metrics
	var totalAmount float64
	for _, tx := range recentTransactions {
		totalAmount += tx.Amount
	}

	averageAmount := totalAmount / float64(len(recentTransactions))

	// Risk based on transaction count and amount
	var riskScore float64
	if len(recentTransactions) > 50 {
		riskScore += 0.1
	}
	if averageAmount > 100000 {
		riskScore += 0.1
	}

	return riskScore, nil
}

func (r *RiskManagementService) calculateGeographicRisk(location string) float64 {
	// High-risk locations
	highRiskLocations := map[string]bool{
		"unknown": true,
		"vpn":     true,
	}

	if highRiskLocations[location] {
		return 0.15
	}

	// Medium-risk locations (would be expanded in production)
	mediumRiskLocations := map[string]bool{
		"international": true,
	}

	if mediumRiskLocations[location] {
		return 0.08
	}

	return 0.0
}

func (r *RiskManagementService) calculateTimeRisk() float64 {
	now := time.Now()
	hour := now.Hour()

	// High risk during unusual hours (2 AM - 6 AM)
	if hour >= 2 && hour <= 6 {
		return 0.15
	}

	// Medium risk during late hours (10 PM - 2 AM)
	if hour >= 22 || hour <= 2 {
		return 0.08
	}

	return 0.0
}

func (r *RiskManagementService) determineRiskLevel(score float64) string {
	if score >= 0.8 {
		return "critical"
	} else if score >= 0.6 {
		return "high"
	} else if score >= 0.4 {
		return "medium"
	} else {
		return "low"
	}
}

// Fraud Detection Methods
func (r *RiskManagementService) DetectFraud(ctx context.Context, accountID uuid.UUID, transactionID *uuid.UUID, riskScore *RiskScore) error {
	// Check if risk score is high enough to trigger fraud alert
	if riskScore.Score < 0.6 {
		return nil // No fraud detected
	}

	// Determine alert type based on risk factors
	alertType := "suspicious_activity"
	if riskScore.Score >= 0.8 {
		alertType = "high_risk_transaction"
	}

	// Determine severity
	severity := "medium"
	if riskScore.Score >= 0.8 {
		severity = "critical"
	} else if riskScore.Score >= 0.7 {
		severity = "high"
	}

	description := fmt.Sprintf("Risk score %.2f with factors: %v", riskScore.Score, riskScore.Factors)

	alert := &FraudAlert{
		ID:            uuid.New(),
		AlertType:     alertType,
		Severity:      severity,
		AccountID:     accountID,
		TransactionID: transactionID,
		Description:   description,
		RiskScore:     riskScore.Score,
		Status:        "active",
		CreatedAt:     time.Now(),
	}

	if err := r.supabaseService.CreateFraudAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to create fraud alert: %w", err)
	}

	log.Printf("🚨 Fraud alert created: %s (Score: %.2f)", alertType, riskScore.Score)
	return nil
}

func (r *RiskManagementService) GetActiveFraudAlerts(ctx context.Context) ([]FraudAlert, error) {
	return r.supabaseService.GetFraudAlertsByStatus(ctx, "active")
}

func (r *RiskManagementService) ResolveFraudAlert(ctx context.Context, alertID uuid.UUID, status string) error {
	alert, err := r.supabaseService.GetFraudAlertByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get fraud alert: %w", err)
	}

	alert.Status = status
	if status != "active" {
		now := time.Now()
		alert.ResolvedAt = &now
	}

	if err := r.supabaseService.UpdateFraudAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to update fraud alert: %w", err)
	}

	log.Printf("Resolved fraud alert: %s", alert.Description)
	return nil
}

// Transaction Limit Management
func (r *RiskManagementService) CreateTransactionLimit(ctx context.Context, accountID uuid.UUID, limitType string, amount float64) (*TransactionLimit, error) {
	limit := &TransactionLimit{
		ID:         uuid.New(),
		AccountID:  accountID,
		LimitType:  limitType,
		Amount:     amount,
		UsedAmount: 0.0,
		ResetDate:  r.calculateResetDate(limitType),
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := r.supabaseService.CreateTransactionLimit(ctx, limit); err != nil {
		return nil, fmt.Errorf("failed to create transaction limit: %w", err)
	}

	log.Printf("Created %s limit for account %s: %.2f", limitType, accountID, amount)
	return limit, nil
}

func (r *RiskManagementService) calculateResetDate(limitType string) time.Time {
	now := time.Now()
	switch limitType {
	case "daily":
		return now.Add(24 * time.Hour)
	case "weekly":
		return now.Add(7 * 24 * time.Hour)
	case "monthly":
		return now.AddDate(0, 1, 0)
	default:
		return now.Add(24 * time.Hour)
	}
}

func (r *RiskManagementService) CheckTransactionLimit(ctx context.Context, accountID uuid.UUID, amount float64) error {
	limits, err := r.supabaseService.GetTransactionLimitsByAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get transaction limits: %w", err)
	}

	for _, limit := range limits {
		if limit.Status != "active" {
			continue
		}

		// Check if limit has reset
		if time.Now().After(limit.ResetDate) {
			limit.UsedAmount = 0.0
			limit.ResetDate = r.calculateResetDate(limit.LimitType)
			limit.UpdatedAt = time.Now()
			r.supabaseService.UpdateTransactionLimit(ctx, limit)
		}

		// Check if transaction would exceed limit
		if limit.UsedAmount+amount > limit.Amount {
			return fmt.Errorf("transaction would exceed %s limit (%.2f/%.2f)", 
				limit.LimitType, limit.UsedAmount+amount, limit.Amount)
		}
	}

	return nil
}

func (r *RiskManagementService) UpdateTransactionLimitUsage(ctx context.Context, accountID uuid.UUID, amount float64) error {
	limits, err := r.supabaseService.GetTransactionLimitsByAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get transaction limits: %w", err)
	}

	for _, limit := range limits {
		if limit.Status != "active" {
			continue
		}

		limit.UsedAmount += amount
		limit.UpdatedAt = time.Now()

		// Update status if limit exceeded
		if limit.UsedAmount >= limit.Amount {
			limit.Status = "exceeded"
		}

		if err := r.supabaseService.UpdateTransactionLimit(ctx, limit); err != nil {
			return fmt.Errorf("failed to update transaction limit: %w", err)
		}
	}

	return nil
}

// Velocity Monitoring
func (r *RiskManagementService) UpdateVelocityMonitor(ctx context.Context, accountID uuid.UUID, amount float64) error {
	// Update velocity for different time windows
	timeWindows := []string{"1h", "24h", "7d", "30d"}

	for _, window := range timeWindows {
		monitor, err := r.supabaseService.GetVelocityMonitor(ctx, accountID, window)
		if err != nil {
			// Create new monitor if doesn't exist
			monitor = &VelocityMonitor{
				ID:              uuid.New(),
				AccountID:       accountID,
				TimeWindow:      window,
				TransactionCount: 0,
				TotalAmount:     0.0,
				AverageAmount:   0.0,
				RiskLevel:       "low",
				LastUpdated:     time.Now(),
				CreatedAt:       time.Now(),
			}
		}

		// Update metrics
		monitor.TransactionCount++
		monitor.TotalAmount += amount
		monitor.AverageAmount = monitor.TotalAmount / float64(monitor.TransactionCount)
		monitor.LastUpdated = time.Now()

		// Determine risk level
		monitor.RiskLevel = r.calculateVelocityRiskLevel(monitor)

		if err := r.supabaseService.UpsertVelocityMonitor(ctx, monitor); err != nil {
			return fmt.Errorf("failed to update velocity monitor: %w", err)
		}
	}

	return nil
}

func (r *RiskManagementService) calculateVelocityRiskLevel(monitor *VelocityMonitor) string {
	// Risk based on transaction count and amount
	if monitor.TransactionCount > 100 || monitor.AverageAmount > 50000 {
		return "high"
	} else if monitor.TransactionCount > 50 || monitor.AverageAmount > 25000 {
		return "medium"
	}
	return "low"
}

// Geographic Monitoring
func (r *RiskManagementService) MonitorGeographicActivity(ctx context.Context, accountID uuid.UUID, location string) error {
	// Check for unusual location activity
	recentLocations, err := r.supabaseService.GetRecentLocations(ctx, accountID, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to get recent locations: %w", err)
	}

	// Check if this is a new location
	isNewLocation := true
	for _, loc := range recentLocations {
		if loc == location {
			isNewLocation = false
			break
		}
	}

	if isNewLocation && len(recentLocations) > 0 {
		// Create geographic alert
		alert := &GeographicAlert{
			ID:          uuid.New(),
			AccountID:   accountID,
			Location:    location,
			RiskType:    "unusual_location",
			Description: fmt.Sprintf("Transaction from new location: %s", location),
			RiskScore:   0.3,
			Status:      "active",
			CreatedAt:   time.Now(),
		}

		if err := r.supabaseService.CreateGeographicAlert(ctx, alert); err != nil {
			return fmt.Errorf("failed to create geographic alert: %w", err)
		}

		log.Printf("⚠️ Geographic alert: New location %s for account %s", location, accountID)
	}

	return nil
}

// Health Check
func (r *RiskManagementService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := r.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("risk management service health check failed: %w", err)
	}

	log.Println("✅ Risk management service health check passed")
	return nil
} 