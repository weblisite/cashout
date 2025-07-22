package services

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// FraudDetectionService handles fraud detection operations
type FraudDetectionService struct {
	userPatterns map[string]*UserPattern
	blacklist    map[string]bool
	mutex        sync.RWMutex
	config       FraudDetectionConfig
}

// UserPattern represents user behavior patterns
type UserPattern struct {
	UserID           string    `json:"user_id"`
	TransactionCount int       `json:"transaction_count"`
	TotalAmount      float64   `json:"total_amount"`
	AverageAmount    float64   `json:"average_amount"`
	LastTransaction  time.Time `json:"last_transaction"`
	DeviceFingerprint string   `json:"device_fingerprint"`
	LocationHistory  []string  `json:"location_history"`
	RiskScore        float64   `json:"risk_score"`
	LastUpdated      time.Time `json:"last_updated"`
}

// FraudDetectionConfig represents fraud detection configuration
type FraudDetectionConfig struct {
	MaxTransactionAmount    float64 `json:"max_transaction_amount"`
	MaxDailyTransactions    int     `json:"max_daily_transactions"`
	MaxDailyAmount          float64 `json:"max_daily_amount"`
	VelocityThreshold       int     `json:"velocity_threshold"`
	AmountDeviationThreshold float64 `json:"amount_deviation_threshold"`
	LocationChangeThreshold int     `json:"location_change_threshold"`
	RiskScoreThreshold      float64 `json:"risk_score_threshold"`
}

// TransactionData represents transaction data for fraud detection
type TransactionData struct {
	UserID           string    `json:"user_id"`
	Amount           float64   `json:"amount"`
	TransactionType  string    `json:"transaction_type"`
	RecipientID      string    `json:"recipient_id"`
	DeviceFingerprint string   `json:"device_fingerprint"`
	IPAddress        string    `json:"ip_address"`
	Location         string    `json:"location"`
	Timestamp        time.Time `json:"timestamp"`
}

// FraudDetectionResult represents the result of fraud detection
type FraudDetectionResult struct {
	IsFraudulent     bool      `json:"is_fraudulent"`
	RiskScore        float64   `json:"risk_score"`
	RiskFactors      []string  `json:"risk_factors"`
	Recommendation   string    `json:"recommendation"`
	BlockTransaction bool      `json:"block_transaction"`
	Timestamp        time.Time `json:"timestamp"`
}

// NewFraudDetectionService creates a new fraud detection service
func NewFraudDetectionService() *FraudDetectionService {
	return &FraudDetectionService{
		userPatterns: make(map[string]*UserPattern),
		blacklist:    make(map[string]bool),
		config: FraudDetectionConfig{
			MaxTransactionAmount:      1000000.0,
			MaxDailyTransactions:      50,
			MaxDailyAmount:            500000.0,
			VelocityThreshold:         10,
			AmountDeviationThreshold:  3.0,
			LocationChangeThreshold:   3,
			RiskScoreThreshold:        0.7,
		},
	}
}

// AnalyzeTransaction analyzes a transaction for potential fraud
func (f *FraudDetectionService) AnalyzeTransaction(data TransactionData) *FraudDetectionResult {
	result := &FraudDetectionResult{
		Timestamp: time.Now(),
		RiskFactors: []string{},
	}

	// Check blacklist
	if f.isBlacklisted(data.UserID) {
		result.IsFraudulent = true
		result.RiskScore = 1.0
		result.RiskFactors = append(result.RiskFactors, "User is blacklisted")
		result.Recommendation = "Block transaction - user is blacklisted"
		result.BlockTransaction = true
		return result
	}

	// Get or create user pattern
	pattern := f.getUserPattern(data.UserID)
	if pattern == nil {
		pattern = f.createUserPattern(data.UserID, data.DeviceFingerprint)
	}

	// Perform various fraud checks
	riskScore := 0.0
	riskFactors := []string{}

	// Check transaction amount
	if data.Amount > f.config.MaxTransactionAmount {
		riskScore += 0.3
		riskFactors = append(riskFactors, "Transaction amount exceeds limit")
	}

	// Check velocity (transaction frequency)
	if f.checkVelocity(data.UserID, data.Timestamp) {
		riskScore += 0.25
		riskFactors = append(riskFactors, "High transaction velocity detected")
	}

	// Check amount deviation
	if f.checkAmountDeviation(data.UserID, data.Amount) {
		riskScore += 0.2
		riskFactors = append(riskFactors, "Unusual transaction amount")
	}

	// Check location change
	if f.checkLocationChange(data.UserID, data.Location) {
		riskScore += 0.15
		riskFactors = append(riskFactors, "Suspicious location change")
	}

	// Check device fingerprint
	if f.checkDeviceFingerprint(data.UserID, data.DeviceFingerprint) {
		riskScore += 0.1
		riskFactors = append(riskFactors, "New device detected")
	}

	// Check daily limits
	if f.checkDailyLimits(data.UserID, data.Amount) {
		riskScore += 0.2
		riskFactors = append(riskFactors, "Daily limits exceeded")
	}

	// Update user pattern
	f.updateUserPattern(data.UserID, data)

	// Set result
	result.RiskScore = math.Min(riskScore, 1.0)
	result.RiskFactors = riskFactors
	result.IsFraudulent = riskScore >= f.config.RiskScoreThreshold

	if result.IsFraudulent {
		result.Recommendation = "Review transaction manually"
		result.BlockTransaction = riskScore >= 0.8
	} else {
		result.Recommendation = "Transaction appears legitimate"
	}

	log.Info().
		Str("user_id", data.UserID).
		Float64("amount", data.Amount).
		Float64("risk_score", result.RiskScore).
		Bool("is_fraudulent", result.IsFraudulent).
		Msg("Fraud detection analysis completed")

	return result
}

// AddToBlacklist adds a user to the blacklist
func (f *FraudDetectionService) AddToBlacklist(userID string, reason string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.blacklist[userID] = true
	log.Warn().
		Str("user_id", userID).
		Str("reason", reason).
		Msg("User added to fraud blacklist")
}

// RemoveFromBlacklist removes a user from the blacklist
func (f *FraudDetectionService) RemoveFromBlacklist(userID string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	delete(f.blacklist, userID)
	log.Info().
		Str("user_id", userID).
		Msg("User removed from fraud blacklist")
}

// GetUserRiskProfile gets the risk profile for a user
func (f *FraudDetectionService) GetUserRiskProfile(userID string) *UserPattern {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	if pattern, exists := f.userPatterns[userID]; exists {
		return pattern
	}
	return nil
}

// UpdateConfiguration updates fraud detection configuration
func (f *FraudDetectionService) UpdateConfiguration(config FraudDetectionConfig) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.config = config
	log.Info().Msg("Fraud detection configuration updated")
}

// GetStatistics gets fraud detection statistics
func (f *FraudDetectionService) GetStatistics() map[string]interface{} {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_users":           len(f.userPatterns),
		"blacklisted_users":     len(f.blacklist),
		"high_risk_users":       0,
		"average_risk_score":    0.0,
	}

	totalRiskScore := 0.0
	highRiskCount := 0

	for _, pattern := range f.userPatterns {
		totalRiskScore += pattern.RiskScore
		if pattern.RiskScore > 0.5 {
			highRiskCount++
		}
	}

	if len(f.userPatterns) > 0 {
		stats["average_risk_score"] = totalRiskScore / float64(len(f.userPatterns))
	}
	stats["high_risk_users"] = highRiskCount

	return stats
}

// isBlacklisted checks if a user is blacklisted
func (f *FraudDetectionService) isBlacklisted(userID string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.blacklist[userID]
}

// getUserPattern gets user pattern
func (f *FraudDetectionService) getUserPattern(userID string) *UserPattern {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.userPatterns[userID]
}

// createUserPattern creates a new user pattern
func (f *FraudDetectionService) createUserPattern(userID, deviceFingerprint string) *UserPattern {
	pattern := &UserPattern{
		UserID:           userID,
		TransactionCount: 0,
		TotalAmount:      0.0,
		AverageAmount:    0.0,
		LastTransaction:  time.Now(),
		DeviceFingerprint: deviceFingerprint,
		LocationHistory:  []string{},
		RiskScore:        0.0,
		LastUpdated:      time.Now(),
	}

	f.mutex.Lock()
	f.userPatterns[userID] = pattern
	f.mutex.Unlock()

	return pattern
}

// updateUserPattern updates user pattern with new transaction data
func (f *FraudDetectionService) updateUserPattern(userID string, data TransactionData) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	pattern, exists := f.userPatterns[userID]
	if !exists {
		return
	}

	// Update transaction count and amounts
	pattern.TransactionCount++
	pattern.TotalAmount += data.Amount
	pattern.AverageAmount = pattern.TotalAmount / float64(pattern.TransactionCount)
	pattern.LastTransaction = data.Timestamp
	pattern.LastUpdated = time.Now()

	// Update location history
	if data.Location != "" {
		pattern.LocationHistory = append(pattern.LocationHistory, data.Location)
		// Keep only last 10 locations
		if len(pattern.LocationHistory) > 10 {
			pattern.LocationHistory = pattern.LocationHistory[len(pattern.LocationHistory)-10:]
		}
	}

	// Update device fingerprint if changed
	if data.DeviceFingerprint != "" && data.DeviceFingerprint != pattern.DeviceFingerprint {
		pattern.DeviceFingerprint = data.DeviceFingerprint
	}
}

// checkVelocity checks for high transaction velocity
func (f *FraudDetectionService) checkVelocity(userID string, timestamp time.Time) bool {
	pattern := f.getUserPattern(userID)
	if pattern == nil {
		return false
	}

	// Check if multiple transactions in short time
	timeDiff := timestamp.Sub(pattern.LastTransaction)
	if timeDiff < time.Minute*5 && pattern.TransactionCount > f.config.VelocityThreshold {
		return true
	}

	return false
}

// checkAmountDeviation checks for unusual transaction amounts
func (f *FraudDetectionService) checkAmountDeviation(userID string, amount float64) bool {
	pattern := f.getUserPattern(userID)
	if pattern == nil || pattern.AverageAmount == 0 {
		return false
	}

	// Check if amount deviates significantly from average
	deviation := math.Abs(amount - pattern.AverageAmount) / pattern.AverageAmount
	return deviation > f.config.AmountDeviationThreshold
}

// checkLocationChange checks for suspicious location changes
func (f *FraudDetectionService) checkLocationChange(userID, newLocation string) bool {
	pattern := f.getUserPattern(userID)
	if pattern == nil || len(pattern.LocationHistory) == 0 {
		return false
	}

	// Check if location changed recently
	lastLocation := pattern.LocationHistory[len(pattern.LocationHistory)-1]
	if newLocation != lastLocation {
		// Simple check - in production, you'd use geolocation services
		return len(pattern.LocationHistory) >= f.config.LocationChangeThreshold
	}

	return false
}

// checkDeviceFingerprint checks for new device usage
func (f *FraudDetectionService) checkDeviceFingerprint(userID, newFingerprint string) bool {
	pattern := f.getUserPattern(userID)
	if pattern == nil {
		return false
	}

	return pattern.DeviceFingerprint != "" && pattern.DeviceFingerprint != newFingerprint
}

// checkDailyLimits checks if daily limits are exceeded
func (f *FraudDetectionService) checkDailyLimits(userID string, amount float64) bool {
	pattern := f.getUserPattern(userID)
	if pattern == nil {
		return false
	}

	// Check if last transaction was today
	today := time.Now().Truncate(24 * time.Hour)
	lastTransactionDay := pattern.LastTransaction.Truncate(24 * time.Hour)

	if today.Equal(lastTransactionDay) {
		// Simple daily limit check - in production, you'd track daily totals
		return amount > f.config.MaxDailyAmount || pattern.TransactionCount > f.config.MaxDailyTransactions
	}

	return false
}

// ExportUserPatterns exports user patterns for analysis
func (f *FraudDetectionService) ExportUserPatterns() ([]byte, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return json.Marshal(f.userPatterns)
}

// ImportUserPatterns imports user patterns
func (f *FraudDetectionService) ImportUserPatterns(data []byte) error {
	var patterns map[string]*UserPattern
	if err := json.Unmarshal(data, &patterns); err != nil {
		return err
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.userPatterns = patterns
	return nil
} 