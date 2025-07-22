package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// EnhancedFraudDetectionService handles advanced fraud detection
type EnhancedFraudDetectionService struct {
	supabaseService *SupabaseService
	riskManagement  *RiskManagementService
	kycAML          *KYCAMLService
}

// NewEnhancedFraudDetectionService creates a new enhanced fraud detection service
func NewEnhancedFraudDetectionService(supabaseService *SupabaseService, riskManagement *RiskManagementService, kycAML *KYCAMLService) *EnhancedFraudDetectionService {
	return &EnhancedFraudDetectionService{
		supabaseService: supabaseService,
		riskManagement:  riskManagement,
		kycAML:          kycAML,
	}
}

// FraudModel represents a fraud detection model
type FraudModel struct {
	ID            uuid.UUID `json:"id"`
	ModelName     string    `json:"model_name"`
	ModelVersion  string    `json:"model_version"`
	ModelType     string    `json:"model_type"` // ml, rule_based, hybrid
	Accuracy      float64   `json:"accuracy"`
	Precision     float64   `json:"precision"`
	Recall        float64   `json:"recall"`
	F1Score       float64   `json:"f1_score"`
	IsActive      bool      `json:"is_active"`
	LastUpdated   time.Time `json:"last_updated"`
	CreatedAt     time.Time `json:"created_at"`
}

// BehavioralProfile represents user behavioral patterns
type BehavioralProfile struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	TransactionPatterns   TransactionPatterns `json:"transaction_patterns"`
	DevicePatterns        DevicePatterns `json:"device_patterns"`
	LocationPatterns      LocationPatterns `json:"location_patterns"`
	TimePatterns          TimePatterns `json:"time_patterns"`
	RiskScore             float64   `json:"risk_score"`
	LastUpdated           time.Time `json:"last_updated"`
	CreatedAt             time.Time `json:"created_at"`
}

// TransactionPatterns represents transaction behavior patterns
type TransactionPatterns struct {
	AverageAmount        float64 `json:"average_amount"`
	MaxAmount            float64 `json:"max_amount"`
	MinAmount            float64 `json:"min_amount"`
	TransactionFrequency float64 `json:"transaction_frequency"` // per day
	PreferredAmounts     []float64 `json:"preferred_amounts"`
	RecipientPatterns    map[string]int `json:"recipient_patterns"`
	AmountVariance       float64 `json:"amount_variance"`
}

// DevicePatterns represents device usage patterns
type DevicePatterns struct {
	PrimaryDevice        string `json:"primary_device"`
	DeviceFingerprint    string `json:"device_fingerprint"`
	IPAddresses          []string `json:"ip_addresses"`
	UserAgents           []string `json:"user_agents"`
	DeviceChanges        int     `json:"device_changes"`
	LastDeviceChange     *time.Time `json:"last_device_change,omitempty"`
}

// LocationPatterns represents location behavior patterns
type LocationPatterns struct {
	PrimaryLocation      string `json:"primary_location"`
	FrequentLocations    []string `json:"frequent_locations"`
	LocationChanges      int     `json:"location_changes"`
	LastLocationChange   *time.Time `json:"last_location_change,omitempty"`
	TravelPatterns       []TravelPattern `json:"travel_patterns"`
}

// TimePatterns represents time-based behavior patterns
type TimePatterns struct {
	PreferredHours       []int `json:"preferred_hours"`
	PreferredDays        []int `json:"preferred_days"`
	TransactionTimes     []time.Time `json:"transaction_times"`
	TimeVariance         float64 `json:"time_variance"`
}

// TravelPattern represents travel behavior
type TravelPattern struct {
	FromLocation         string    `json:"from_location"`
	ToLocation           string    `json:"to_location"`
	Frequency            int       `json:"frequency"`
	LastTravel           time.Time `json:"last_travel"`
}

// FraudAlert represents an enhanced fraud alert
type EnhancedFraudAlert struct {
	ID              uuid.UUID `json:"id"`
	AlertType       string    `json:"alert_type"`
	Severity        string    `json:"severity"`
	UserID          uuid.UUID `json:"user_id"`
	TransactionID   *uuid.UUID `json:"transaction_id,omitempty"`
	Description     string    `json:"description"`
	RiskScore       float64   `json:"risk_score"`
	Confidence      float64   `json:"confidence"`
	ModelUsed       string    `json:"model_used"`
	Features        map[string]interface{} `json:"features"`
	Status          string    `json:"status"`
	InvestigationNotes string `json:"investigation_notes"`
	CreatedAt       time.Time `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

// Enhanced Fraud Detection Methods
func (e *EnhancedFraudDetectionService) AnalyzeTransaction(ctx context.Context, userID, transactionID uuid.UUID, amount float64, location, deviceInfo string) (*EnhancedFraudAlert, error) {
	// Get or create behavioral profile
	profile, err := e.getOrCreateBehavioralProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get behavioral profile: %w", err)
	}

	// Extract features for fraud detection
	features := e.extractFeatures(profile, amount, location, deviceInfo, time.Now())

	// Calculate risk score using multiple models
	riskScore, confidence, modelUsed := e.calculateRiskScore(features)

	// Check if alert should be triggered
	if riskScore > 0.7 {
		alert := &EnhancedFraudAlert{
			ID:            uuid.New(),
			AlertType:     "suspicious_transaction",
			Severity:      e.determineSeverity(riskScore),
			UserID:        userID,
			TransactionID: &transactionID,
			Description:   e.generateAlertDescription(features, riskScore),
			RiskScore:     riskScore,
			Confidence:    confidence,
			ModelUsed:     modelUsed,
			Features:      features,
			Status:        "active",
			CreatedAt:     time.Now(),
		}

		if err := e.supabaseService.CreateEnhancedFraudAlert(ctx, alert); err != nil {
			return nil, fmt.Errorf("failed to create fraud alert: %w", err)
		}

		log.Printf("🚨 Enhanced fraud alert: Score %.2f, Confidence %.2f, Model: %s", 
			riskScore, confidence, modelUsed)
		return alert, nil
	}

	// Update behavioral profile
	if err := e.updateBehavioralProfile(ctx, profile, amount, location, deviceInfo); err != nil {
		return nil, fmt.Errorf("failed to update behavioral profile: %w", err)
	}

	return nil, nil
}

func (e *EnhancedFraudDetectionService) getOrCreateBehavioralProfile(ctx context.Context, userID uuid.UUID) (*BehavioralProfile, error) {
	profile, err := e.supabaseService.GetBehavioralProfileByUserID(ctx, userID)
	if err != nil {
		// Create new profile
		profile = &BehavioralProfile{
			ID:          uuid.New(),
			UserID:      userID,
			RiskScore:   0.0,
			LastUpdated: time.Now(),
			CreatedAt:   time.Now(),
		}

		if err := e.supabaseService.CreateBehavioralProfile(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to create behavioral profile: %w", err)
		}
	}

	return profile, nil
}

func (e *EnhancedFraudDetectionService) extractFeatures(profile *BehavioralProfile, amount float64, location, deviceInfo string, timestamp time.Time) map[string]interface{} {
	features := make(map[string]interface{})

	// Transaction amount features
	features["amount"] = amount
	features["amount_deviation"] = e.calculateAmountDeviation(profile.TransactionPatterns.AverageAmount, amount)
	features["amount_ratio"] = amount / profile.TransactionPatterns.AverageAmount

	// Time-based features
	hour := timestamp.Hour()
	dayOfWeek := int(timestamp.Weekday())
	features["hour"] = hour
	features["day_of_week"] = dayOfWeek
	features["is_preferred_hour"] = e.isPreferredHour(profile.TimePatterns.PreferredHours, hour)
	features["is_preferred_day"] = e.isPreferredDay(profile.TimePatterns.PreferredDays, dayOfWeek)

	// Location features
	features["location"] = location
	features["is_primary_location"] = location == profile.LocationPatterns.PrimaryLocation
	features["is_frequent_location"] = e.isFrequentLocation(profile.LocationPatterns.FrequentLocations, location)
	features["location_change"] = e.calculateLocationChange(profile, location)

	// Device features
	features["device_info"] = deviceInfo
	features["device_change"] = e.calculateDeviceChange(profile, deviceInfo)

	// Behavioral features
	features["transaction_frequency"] = profile.TransactionPatterns.TransactionFrequency
	features["amount_variance"] = profile.TransactionPatterns.AmountVariance

	return features
}

func (e *EnhancedFraudDetectionService) calculateAmountDeviation(average, current float64) float64 {
	if average == 0 {
		return 1.0
	}
	return math.Abs(current-average) / average
}

func (e *EnhancedFraudDetectionService) isPreferredHour(preferredHours []int, currentHour int) bool {
	for _, hour := range preferredHours {
		if hour == currentHour {
			return true
		}
	}
	return false
}

func (e *EnhancedFraudDetectionService) isPreferredDay(preferredDays []int, currentDay int) bool {
	for _, day := range preferredDays {
		if day == currentDay {
			return true
		}
	}
	return false
}

func (e *EnhancedFraudDetectionService) isFrequentLocation(frequentLocations []string, currentLocation string) bool {
	for _, location := range frequentLocations {
		if location == currentLocation {
			return true
		}
	}
	return false
}

func (e *EnhancedFraudDetectionService) calculateLocationChange(profile *BehavioralProfile, currentLocation string) float64 {
	if profile.LocationPatterns.PrimaryLocation == "" {
		return 0.0
	}

	if currentLocation != profile.LocationPatterns.PrimaryLocation {
		return 1.0
	}

	return 0.0
}

func (e *EnhancedFraudDetectionService) calculateDeviceChange(profile *BehavioralProfile, currentDevice string) float64 {
	if profile.DevicePatterns.PrimaryDevice == "" {
		return 0.0
	}

	if currentDevice != profile.DevicePatterns.PrimaryDevice {
		return 1.0
	}

	return 0.0
}

func (e *EnhancedFraudDetectionService) calculateRiskScore(features map[string]interface{}) (float64, float64, string) {
	// Multi-model approach
	scores := make(map[string]float64)
	confidences := make(map[string]float64)

	// Rule-based model
	scores["rule_based"], confidences["rule_based"] = e.ruleBasedModel(features)

	// ML model (simulated)
	scores["ml_model"], confidences["ml_model"] = e.mlModel(features)

	// Hybrid model
	scores["hybrid"], confidences["hybrid"] = e.hybridModel(features)

	// Select best model based on confidence
	bestModel := "rule_based"
	bestScore := scores["rule_based"]
	bestConfidence := confidences["rule_based"]

	for model, score := range scores {
		if confidences[model] > bestConfidence {
			bestModel = model
			bestScore = score
			bestConfidence = confidences[model]
		}
	}

	return bestScore, bestConfidence, bestModel
}

func (e *EnhancedFraudDetectionService) ruleBasedModel(features map[string]interface{}) (float64, float64) {
	score := 0.0

	// Amount deviation
	if amountDeviation, ok := features["amount_deviation"].(float64); ok {
		if amountDeviation > 2.0 {
			score += 0.3
		} else if amountDeviation > 1.5 {
			score += 0.2
		}
	}

	// Location change
	if locationChange, ok := features["location_change"].(float64); ok {
		if locationChange > 0 {
			score += 0.2
		}
	}

	// Device change
	if deviceChange, ok := features["device_change"].(float64); ok {
		if deviceChange > 0 {
			score += 0.2
		}
	}

	// Time patterns
	if isPreferredHour, ok := features["is_preferred_hour"].(bool); ok {
		if !isPreferredHour {
			score += 0.1
		}
	}

	// Transaction frequency
	if frequency, ok := features["transaction_frequency"].(float64); ok {
		if frequency > 10 {
			score += 0.2
		}
	}

	return math.Min(score, 1.0), 0.8
}

func (e *EnhancedFraudDetectionService) mlModel(features map[string]interface{}) (float64, float64) {
	// Simulated ML model
	// In production, this would use a trained machine learning model
	
	// Feature engineering
	featureVector := e.createFeatureVector(features)
	
	// Simulated prediction
	score := e.simulateMLPrediction(featureVector)
	
	return score, 0.9
}

func (e *EnhancedFraudDetectionService) createFeatureVector(features map[string]interface{}) []float64 {
	vector := make([]float64, 0)
	
	// Normalize features
	if amount, ok := features["amount"].(float64); ok {
		vector = append(vector, amount/1000000) // Normalize to millions
	}
	
	if deviation, ok := features["amount_deviation"].(float64); ok {
		vector = append(vector, deviation)
	}
	
	if locationChange, ok := features["location_change"].(float64); ok {
		vector = append(vector, locationChange)
	}
	
	if deviceChange, ok := features["device_change"].(float64); ok {
		vector = append(vector, deviceChange)
	}
	
	return vector
}

func (e *EnhancedFraudDetectionService) simulateMLPrediction(featureVector []float64) float64 {
	// Simulated ML prediction using weighted sum
	weights := []float64{0.3, 0.25, 0.25, 0.2}
	
	score := 0.0
	for i, feature := range featureVector {
		if i < len(weights) {
			score += feature * weights[i]
		}
	}
	
	// Add some randomness to simulate ML uncertainty
	score += (rand.Float64() - 0.5) * 0.1
	
	return math.Max(0.0, math.Min(score, 1.0))
}

func (e *EnhancedFraudDetectionService) hybridModel(features map[string]interface{}) (float64, float64) {
	ruleScore, ruleConfidence := e.ruleBasedModel(features)
	mlScore, mlConfidence := e.mlModel(features)
	
	// Weighted combination
	hybridScore := (ruleScore*ruleConfidence + mlScore*mlConfidence) / (ruleConfidence + mlConfidence)
	hybridConfidence := (ruleConfidence + mlConfidence) / 2
	
	return hybridScore, hybridConfidence
}

func (e *EnhancedFraudDetectionService) determineSeverity(riskScore float64) string {
	if riskScore >= 0.9 {
		return "critical"
	} else if riskScore >= 0.7 {
		return "high"
	} else if riskScore >= 0.5 {
		return "medium"
	}
	return "low"
}

func (e *EnhancedFraudDetectionService) generateAlertDescription(features map[string]interface{}, riskScore float64) string {
	description := fmt.Sprintf("Suspicious transaction detected (Risk Score: %.2f). ", riskScore)
	
	if deviation, ok := features["amount_deviation"].(float64); ok && deviation > 1.5 {
		description += "Unusual transaction amount. "
	}
	
	if locationChange, ok := features["location_change"].(float64); ok && locationChange > 0 {
		description += "Transaction from new location. "
	}
	
	if deviceChange, ok := features["device_change"].(float64); ok && deviceChange > 0 {
		description += "Transaction from new device. "
	}
	
	return description
}

func (e *EnhancedFraudDetectionService) updateBehavioralProfile(ctx context.Context, profile *BehavioralProfile, amount float64, location, deviceInfo string) error {
	// Update transaction patterns
	profile.TransactionPatterns.AverageAmount = e.updateAverageAmount(profile.TransactionPatterns.AverageAmount, amount)
	if amount > profile.TransactionPatterns.MaxAmount {
		profile.TransactionPatterns.MaxAmount = amount
	}
	if amount < profile.TransactionPatterns.MinAmount || profile.TransactionPatterns.MinAmount == 0 {
		profile.TransactionPatterns.MinAmount = amount
	}

	// Update location patterns
	if profile.LocationPatterns.PrimaryLocation == "" {
		profile.LocationPatterns.PrimaryLocation = location
	} else if location != profile.LocationPatterns.PrimaryLocation {
		profile.LocationPatterns.LocationChanges++
		now := time.Now()
		profile.LocationPatterns.LastLocationChange = &now
	}

	// Update device patterns
	if profile.DevicePatterns.PrimaryDevice == "" {
		profile.DevicePatterns.PrimaryDevice = deviceInfo
	} else if deviceInfo != profile.DevicePatterns.PrimaryDevice {
		profile.DevicePatterns.DeviceChanges++
		now := time.Now()
		profile.DevicePatterns.LastDeviceChange = &now
	}

	// Update time patterns
	now := time.Now()
	profile.TimePatterns.TransactionTimes = append(profile.TimePatterns.TransactionTimes, now)

	profile.LastUpdated = time.Now()

	if err := e.supabaseService.UpdateBehavioralProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update behavioral profile: %w", err)
	}

	return nil
}

func (e *EnhancedFraudDetectionService) updateAverageAmount(currentAverage, newAmount float64) float64 {
	// Simple moving average (in production, would use more sophisticated methods)
	if currentAverage == 0 {
		return newAmount
	}
	return (currentAverage + newAmount) / 2
}

func (e *EnhancedFraudDetectionService) GetActiveFraudAlerts(ctx context.Context) ([]EnhancedFraudAlert, error) {
	return e.supabaseService.GetEnhancedFraudAlertsByStatus(ctx, "active")
}

func (e *EnhancedFraudDetectionService) ResolveFraudAlert(ctx context.Context, alertID uuid.UUID, status string, notes string) error {
	alert, err := e.supabaseService.GetEnhancedFraudAlertByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get fraud alert: %w", err)
	}

	alert.Status = status
	alert.InvestigationNotes = notes
	if status != "active" {
		now := time.Now()
		alert.ResolvedAt = &now
	}

	if err := e.supabaseService.UpdateEnhancedFraudAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to update fraud alert: %w", err)
	}

	log.Printf("Resolved enhanced fraud alert: %s", alert.Description)
	return nil
}

// Health Check
func (e *EnhancedFraudDetectionService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := e.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("enhanced fraud detection service health check failed: %w", err)
	}

	log.Println("✅ Enhanced fraud detection service health check passed")
	return nil
} 