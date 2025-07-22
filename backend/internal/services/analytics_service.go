package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
)

// AnalyticsService handles all analytics and reporting operations
type AnalyticsService struct {
	supabaseService *SupabaseService
	coreBanking     *CoreBankingService
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(supabaseService *SupabaseService, coreBanking *CoreBankingService) *AnalyticsService {
	return &AnalyticsService{
		supabaseService: supabaseService,
		coreBanking:     coreBanking,
	}
}

// AnalyticsReport represents an analytics report
type AnalyticsReport struct {
	ID              uuid.UUID `json:"id"`
	ReportType      string    `json:"report_type"`
	ReportName      string    `json:"report_name"`
	Period          string    `json:"period"` // daily, weekly, monthly, quarterly, yearly
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Data            map[string]interface{} `json:"data"`
	Insights        []Insight `json:"insights"`
	GeneratedAt     time.Time `json:"generated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// Insight represents a business insight
type Insight struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"` // trend, anomaly, recommendation, alert
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // low, medium, high, critical
	Confidence  float64   `json:"confidence"`
	Data        map[string]interface{} `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
}

// BusinessMetrics represents key business metrics
type BusinessMetrics struct {
	ID                    uuid.UUID `json:"id"`
	Date                  time.Time `json:"date"`
	TotalTransactions     int       `json:"total_transactions"`
	TotalVolume           float64   `json:"total_volume"`
	TotalRevenue          float64   `json:"total_revenue"`
	TotalFees             float64   `json:"total_fees"`
	ActiveUsers           int       `json:"active_users"`
	ActiveAgents          int       `json:"active_agents"`
	ActiveBusinesses      int       `json:"active_businesses"`
	AverageTransactionSize float64  `json:"average_transaction_size"`
	TransactionSuccessRate float64  `json:"transaction_success_rate"`
	CustomerSatisfaction  float64   `json:"customer_satisfaction"`
	FraudRate             float64   `json:"fraud_rate"`
	CreatedAt             time.Time `json:"created_at"`
}

// UserAnalytics represents user behavior analytics
type UserAnalytics struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	TotalTransactions     int       `json:"total_transactions"`
	TotalVolume           float64   `json:"total_volume"`
	AverageTransactionSize float64  `json:"average_transaction_size"`
	LastTransactionDate   *time.Time `json:"last_transaction_date,omitempty"`
	DaysSinceLastActivity int       `json:"days_since_last_activity"`
	RiskScore             float64   `json:"risk_score"`
	UserSegment           string    `json:"user_segment"` // low_value, medium_value, high_value, vip
	LifetimeValue         float64   `json:"lifetime_value"`
	RetentionRate         float64   `json:"retention_rate"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// AgentAnalytics represents agent performance analytics
type AgentAnalytics struct {
	ID                    uuid.UUID `json:"id"`
	AgentID               uuid.UUID `json:"agent_id"`
	TotalTransactions     int       `json:"total_transactions"`
	TotalVolume           float64   `json:"total_volume"`
	CommissionEarned      float64   `json:"commission_earned"`
	AverageFloat          float64   `json:"average_float"`
	FloatUtilization      float64   `json:"float_utilization"`
	TransactionSuccessRate float64  `json:"transaction_success_rate"`
	CustomerSatisfaction  float64   `json:"customer_satisfaction"`
	PerformanceRating     string    `json:"performance_rating"` // excellent, good, average, poor
	LastActivityDate      *time.Time `json:"last_activity_date,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Analytics Methods
func (a *AnalyticsService) GenerateBusinessReport(ctx context.Context, reportType, period string, startDate, endDate time.Time) (*AnalyticsReport, error) {
	report := &AnalyticsReport{
		ID:          uuid.New(),
		ReportType:  reportType,
		ReportName:  fmt.Sprintf("%s_%s_Report", reportType, period),
		Period:      period,
		StartDate:   startDate,
		EndDate:     endDate,
		Data:        make(map[string]interface{}),
		Insights:    make([]Insight, 0),
		GeneratedAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	// Generate report data based on type
	switch reportType {
	case "financial":
		a.generateFinancialReport(ctx, report)
	case "operational":
		a.generateOperationalReport(ctx, report)
	case "risk":
		a.generateRiskReport(ctx, report)
	case "user":
		a.generateUserReport(ctx, report)
	case "agent":
		a.generateAgentReport(ctx, report)
	default:
		return nil, fmt.Errorf("unknown report type: %s", reportType)
	}

	// Generate insights
	a.generateInsights(report)

	// Save report
	if err := a.supabaseService.CreateAnalyticsReport(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create analytics report: %w", err)
	}

	log.Printf("Generated %s report for period %s", reportType, period)
	return report, nil
}

func (a *AnalyticsService) generateFinancialReport(ctx context.Context, report *AnalyticsReport) {
	// Get financial metrics
	metrics, err := a.getBusinessMetrics(ctx, report.StartDate, report.EndDate)
	if err != nil {
		log.Printf("Failed to get business metrics: %v", err)
		return
	}

	report.Data["total_revenue"] = metrics.TotalRevenue
	report.Data["total_volume"] = metrics.TotalVolume
	report.Data["total_fees"] = metrics.TotalFees
	report.Data["average_transaction_size"] = metrics.AverageTransactionSize
	report.Data["transaction_success_rate"] = metrics.TransactionSuccessRate

	// Calculate growth rates
	previousMetrics, _ := a.getBusinessMetrics(ctx, report.StartDate.AddDate(0, 0, -30), report.StartDate)
	if previousMetrics.TotalRevenue > 0 {
		revenueGrowth := ((metrics.TotalRevenue - previousMetrics.TotalRevenue) / previousMetrics.TotalRevenue) * 100
		report.Data["revenue_growth"] = revenueGrowth
	}

	// Calculate profitability metrics
	report.Data["profit_margin"] = (metrics.TotalFees / metrics.TotalRevenue) * 100
	report.Data["revenue_per_user"] = metrics.TotalRevenue / float64(metrics.ActiveUsers)
}

func (a *AnalyticsService) generateOperationalReport(ctx context.Context, report *AnalyticsReport) {
	// Get operational metrics
	metrics, err := a.getBusinessMetrics(ctx, report.StartDate, report.EndDate)
	if err != nil {
		log.Printf("Failed to get business metrics: %v", err)
		return
	}

	report.Data["total_transactions"] = metrics.TotalTransactions
	report.Data["active_users"] = metrics.ActiveUsers
	report.Data["active_agents"] = metrics.ActiveAgents
	report.Data["active_businesses"] = metrics.ActiveBusinesses
	report.Data["customer_satisfaction"] = metrics.CustomerSatisfaction

	// Calculate operational efficiency
	report.Data["transactions_per_user"] = float64(metrics.TotalTransactions) / float64(metrics.ActiveUsers)
	report.Data["transactions_per_agent"] = float64(metrics.TotalTransactions) / float64(metrics.ActiveAgents)
	report.Data["volume_per_agent"] = metrics.TotalVolume / float64(metrics.ActiveAgents)
}

func (a *AnalyticsService) generateRiskReport(ctx context.Context, report *AnalyticsReport) {
	// Get risk metrics
	metrics, err := a.getBusinessMetrics(ctx, report.StartDate, report.EndDate)
	if err != nil {
		log.Printf("Failed to get business metrics: %v", err)
		return
	}

	report.Data["fraud_rate"] = metrics.FraudRate
	report.Data["total_transactions"] = metrics.TotalTransactions

	// Calculate risk indicators
	fraudulentTransactions := int(metrics.FraudRate * float64(metrics.TotalTransactions) / 100)
	report.Data["fraudulent_transactions"] = fraudulentTransactions
	report.Data["fraud_loss"] = float64(fraudulentTransactions) * metrics.AverageTransactionSize

	// Get high-risk users
	highRiskUsers, _ := a.getHighRiskUsers(ctx)
	report.Data["high_risk_users"] = len(highRiskUsers)
	report.Data["high_risk_volume"] = a.calculateHighRiskVolume(highRiskUsers)
}

func (a *AnalyticsService) generateUserReport(ctx context.Context, report *AnalyticsReport) {
	// Get user analytics
	userAnalytics, err := a.getAllUserAnalytics(ctx)
	if err != nil {
		log.Printf("Failed to get user analytics: %v", err)
		return
	}

	report.Data["total_users"] = len(userAnalytics)
	report.Data["active_users"] = a.countActiveUsers(userAnalytics)
	report.Data["average_lifetime_value"] = a.calculateAverageLTV(userAnalytics)
	report.Data["retention_rate"] = a.calculateRetentionRate(userAnalytics)

	// User segmentation
	segments := a.segmentUsers(userAnalytics)
	report.Data["user_segments"] = segments
}

func (a *AnalyticsService) generateAgentReport(ctx context.Context, report *AnalyticsReport) {
	// Get agent analytics
	agentAnalytics, err := a.getAllAgentAnalytics(ctx)
	if err != nil {
		log.Printf("Failed to get agent analytics: %v", err)
		return
	}

	report.Data["total_agents"] = len(agentAnalytics)
	report.Data["active_agents"] = a.countActiveAgents(agentAnalytics)
	report.Data["total_commission_paid"] = a.calculateTotalCommission(agentAnalytics)
	report.Data["average_agent_performance"] = a.calculateAverageAgentPerformance(agentAnalytics)

	// Agent performance distribution
	performanceDistribution := a.getAgentPerformanceDistribution(agentAnalytics)
	report.Data["performance_distribution"] = performanceDistribution
}

func (a *AnalyticsService) generateInsights(report *AnalyticsReport) {
	// Generate insights based on report data
	insights := make([]Insight, 0)

	// Revenue insights
	if revenue, ok := report.Data["revenue_growth"].(float64); ok {
		if revenue > 10 {
			insights = append(insights, Insight{
				ID:          uuid.New(),
				Type:        "trend",
				Title:       "Strong Revenue Growth",
				Description: fmt.Sprintf("Revenue grew by %.1f%% in the reporting period", revenue),
				Severity:    "low",
				Confidence:  0.9,
				Data:        map[string]interface{}{"growth_rate": revenue},
				CreatedAt:   time.Now(),
			})
		} else if revenue < -5 {
			insights = append(insights, Insight{
				ID:          uuid.New(),
				Type:        "alert",
				Title:       "Revenue Decline Detected",
				Description: fmt.Sprintf("Revenue declined by %.1f%% in the reporting period", math.Abs(revenue)),
				Severity:    "high",
				Confidence:  0.8,
				Data:        map[string]interface{}{"decline_rate": math.Abs(revenue)},
				CreatedAt:   time.Now(),
			})
		}
	}

	// Fraud insights
	if fraudRate, ok := report.Data["fraud_rate"].(float64); ok {
		if fraudRate > 2.0 {
			insights = append(insights, Insight{
				ID:          uuid.New(),
				Type:        "alert",
				Title:       "High Fraud Rate",
				Description: fmt.Sprintf("Fraud rate of %.2f%% is above acceptable threshold", fraudRate),
				Severity:    "critical",
				Confidence:  0.95,
				Data:        map[string]interface{}{"fraud_rate": fraudRate},
				CreatedAt:   time.Now(),
			})
		}
	}

	// User engagement insights
	if retentionRate, ok := report.Data["retention_rate"].(float64); ok {
		if retentionRate < 70 {
			insights = append(insights, Insight{
				ID:          uuid.New(),
				Type:        "recommendation",
				Title:       "Low User Retention",
				Description: fmt.Sprintf("User retention rate of %.1f%% needs improvement", retentionRate),
				Severity:    "medium",
				Confidence:  0.85,
				Data:        map[string]interface{}{"retention_rate": retentionRate},
				CreatedAt:   time.Now(),
			})
		}
	}

	report.Insights = insights
}

// Business Metrics Methods
func (a *AnalyticsService) getBusinessMetrics(ctx context.Context, startDate, endDate time.Time) (*BusinessMetrics, error) {
	// Get transactions for the period
	transactions, err := a.supabaseService.GetTransactionsByPeriod(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	metrics := &BusinessMetrics{
		ID:                uuid.New(),
		Date:              startDate,
		TotalTransactions: len(transactions),
		CreatedAt:         time.Now(),
	}

	// Calculate metrics from transactions
	var totalVolume, totalFees float64
	var successfulTransactions int

	for _, tx := range transactions {
		totalVolume += tx.Amount
		totalFees += tx.Fee

		if tx.Status == "completed" {
			successfulTransactions++
		}
	}

	metrics.TotalVolume = totalVolume
	metrics.TotalFees = totalFees
	metrics.TotalRevenue = totalFees // Assuming fees are the revenue

	if len(transactions) > 0 {
		metrics.AverageTransactionSize = totalVolume / float64(len(transactions))
		metrics.TransactionSuccessRate = float64(successfulTransactions) / float64(len(transactions)) * 100
	}

	// Get active users, agents, businesses
	activeUsers, _ := a.supabaseService.GetActiveUsers(ctx)
	activeAgents, _ := a.supabaseService.GetActiveAgents(ctx)
	activeBusinesses, _ := a.supabaseService.GetActiveBusinesses(ctx)

	metrics.ActiveUsers = len(activeUsers)
	metrics.ActiveAgents = len(activeAgents)
	metrics.ActiveBusinesses = len(activeBusinesses)

	// Simulate other metrics
	metrics.CustomerSatisfaction = 85.0 + float64(time.Now().UnixNano()%15)
	metrics.FraudRate = 0.5 + float64(time.Now().UnixNano()%10)/100

	return metrics, nil
}

func (a *AnalyticsService) getHighRiskUsers(ctx context.Context) ([]UserAnalytics, error) {
	// Get users with high risk scores
	userAnalytics, err := a.getAllUserAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}

	var highRiskUsers []UserAnalytics
	for _, user := range userAnalytics {
		if user.RiskScore > 0.7 {
			highRiskUsers = append(highRiskUsers, user)
		}
	}

	return highRiskUsers, nil
}

func (a *AnalyticsService) calculateHighRiskVolume(users []UserAnalytics) float64 {
	var totalVolume float64
	for _, user := range users {
		totalVolume += user.TotalVolume
	}
	return totalVolume
}

func (a *AnalyticsService) getAllUserAnalytics(ctx context.Context) ([]UserAnalytics, error) {
	// Get all users and calculate analytics
	users, err := a.supabaseService.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var userAnalytics []UserAnalytics
	for _, user := range users {
		analytics := a.calculateUserAnalytics(ctx, user.ID)
		userAnalytics = append(userAnalytics, analytics)
	}

	return userAnalytics, nil
}

func (a *AnalyticsService) calculateUserAnalytics(ctx context.Context, userID uuid.UUID) UserAnalytics {
	// Get user transactions
	transactions, _ := a.supabaseService.GetTransactionsByUser(ctx, userID, 1000, 0)

	analytics := UserAnalytics{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate metrics
	analytics.TotalTransactions = len(transactions)
	
	var totalVolume float64
	for _, tx := range transactions {
		totalVolume += tx.Amount
	}
	analytics.TotalVolume = totalVolume

	if len(transactions) > 0 {
		analytics.AverageTransactionSize = totalVolume / float64(len(transactions))
		if len(transactions) > 0 {
			lastTx := transactions[len(transactions)-1]
			analytics.LastTransactionDate = &lastTx.CreatedAt
		}
	}

	// Calculate user segment
	if analytics.TotalVolume > 1000000 {
		analytics.UserSegment = "vip"
	} else if analytics.TotalVolume > 100000 {
		analytics.UserSegment = "high_value"
	} else if analytics.TotalVolume > 10000 {
		analytics.UserSegment = "medium_value"
	} else {
		analytics.UserSegment = "low_value"
	}

	analytics.LifetimeValue = analytics.TotalVolume
	analytics.RiskScore = 0.1 + float64(time.Now().UnixNano()%90)/100 // Simulated risk score

	return analytics
}

func (a *AnalyticsService) countActiveUsers(users []UserAnalytics) int {
	activeCount := 0
	for _, user := range users {
		if user.DaysSinceLastActivity <= 30 {
			activeCount++
		}
	}
	return activeCount
}

func (a *AnalyticsService) calculateAverageLTV(users []UserAnalytics) float64 {
	if len(users) == 0 {
		return 0
	}

	var totalLTV float64
	for _, user := range users {
		totalLTV += user.LifetimeValue
	}
	return totalLTV / float64(len(users))
}

func (a *AnalyticsService) calculateRetentionRate(users []UserAnalytics) float64 {
	if len(users) == 0 {
		return 0
	}

	activeUsers := a.countActiveUsers(users)
	return float64(activeUsers) / float64(len(users)) * 100
}

func (a *AnalyticsService) segmentUsers(users []UserAnalytics) map[string]int {
	segments := make(map[string]int)
	for _, user := range users {
		segments[user.UserSegment]++
	}
	return segments
}

func (a *AnalyticsService) getAllAgentAnalytics(ctx context.Context) ([]AgentAnalytics, error) {
	// Get all agents and calculate analytics
	agents, err := a.supabaseService.GetAllAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	var agentAnalytics []AgentAnalytics
	for _, agent := range agents {
		analytics := a.calculateAgentAnalytics(ctx, agent.ID)
		agentAnalytics = append(agentAnalytics, analytics)
	}

	return agentAnalytics, nil
}

func (a *AnalyticsService) calculateAgentAnalytics(ctx context.Context, agentID uuid.UUID) AgentAnalytics {
	// Get agent transactions
	transactions, _ := a.supabaseService.GetTransactionsByAgent(ctx, agentID, 1000, 0)

	analytics := AgentAnalytics{
		ID:        uuid.New(),
		AgentID:   agentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate metrics
	analytics.TotalTransactions = len(transactions)
	
	var totalVolume float64
	var successfulTransactions int
	for _, tx := range transactions {
		totalVolume += tx.Amount
		if tx.Status == "completed" {
			successfulTransactions++
		}
	}
	analytics.TotalVolume = totalVolume

	if len(transactions) > 0 {
		analytics.TransactionSuccessRate = float64(successfulTransactions) / float64(len(transactions)) * 100
		if len(transactions) > 0 {
			lastTx := transactions[len(transactions)-1]
			analytics.LastActivityDate = &lastTx.CreatedAt
		}
	}

	// Calculate commission (assuming 0.25% commission rate)
	analytics.CommissionEarned = totalVolume * 0.0025

	// Simulate other metrics
	analytics.AverageFloat = 50000 + float64(time.Now().UnixNano()%100000)
	analytics.FloatUtilization = 60.0 + float64(time.Now().UnixNano()%40)
	analytics.CustomerSatisfaction = 80.0 + float64(time.Now().UnixNano()%20)

	// Determine performance rating
	if analytics.TransactionSuccessRate > 95 && analytics.TotalVolume > 1000000 {
		analytics.PerformanceRating = "excellent"
	} else if analytics.TransactionSuccessRate > 90 && analytics.TotalVolume > 500000 {
		analytics.PerformanceRating = "good"
	} else if analytics.TransactionSuccessRate > 80 {
		analytics.PerformanceRating = "average"
	} else {
		analytics.PerformanceRating = "poor"
	}

	return analytics
}

func (a *AnalyticsService) countActiveAgents(agents []AgentAnalytics) int {
	activeCount := 0
	for _, agent := range agents {
		if agent.LastActivityDate != nil && time.Since(*agent.LastActivityDate) <= 7*24*time.Hour {
			activeCount++
		}
	}
	return activeCount
}

func (a *AnalyticsService) calculateTotalCommission(agents []AgentAnalytics) float64 {
	var totalCommission float64
	for _, agent := range agents {
		totalCommission += agent.CommissionEarned
	}
	return totalCommission
}

func (a *AnalyticsService) calculateAverageAgentPerformance(agents []AgentAnalytics) float64 {
	if len(agents) == 0 {
		return 0
	}

	var totalPerformance float64
	for _, agent := range agents {
		totalPerformance += agent.TransactionSuccessRate
	}
	return totalPerformance / float64(len(agents))
}

func (a *AnalyticsService) getAgentPerformanceDistribution(agents []AgentAnalytics) map[string]int {
	distribution := make(map[string]int)
	for _, agent := range agents {
		distribution[agent.PerformanceRating]++
	}
	return distribution
}

// Health Check
func (a *AnalyticsService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := a.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("analytics service health check failed: %w", err)
	}

	log.Println("✅ Analytics service health check passed")
	return nil
} 