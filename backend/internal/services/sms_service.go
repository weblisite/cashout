package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// SMSService handles SMS operations using Africa's Talking
type SMSService struct {
	apiKey     string
	username   string
	senderID   string
	baseURL    string
	httpClient *http.Client
}

// SMSRequest represents an SMS request to Africa's Talking
type SMSRequest struct {
	Username string   `json:"username"`
	To       []string `json:"to"`
	Message  string   `json:"message"`
	From     string   `json:"from"`
}

// SMSResponse represents the response from Africa's Talking
type SMSResponse struct {
	SMS struct {
		Message string `json:"Message"`
		Recipients []struct {
			Number    string `json:"number"`
			Status    string `json:"status"`
			MessageID string `json:"messageId"`
			Cost      string `json:"cost"`
		} `json:"Recipients"`
	} `json:"SMS"`
}

// NewSMSService creates a new SMS service instance
func NewSMSService() *SMSService {
	return &SMSService{
		apiKey:     getEnv("AT_API_KEY", ""),
		username:   getEnv("AT_USERNAME", ""),
		senderID:   getEnv("AT_SENDER_ID", "CASHOUT"),
		baseURL:    "https://api.africastalking.com/version1/messaging",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendOTP sends an OTP via SMS
func (s *SMSService) SendOTP(phoneNumber, otp string) error {
	message := fmt.Sprintf("Your Cashout verification code is: %s. Valid for 10 minutes. Do not share this code with anyone.", otp)
	
	return s.sendSMS(phoneNumber, message)
}

// SendWelcomeMessage sends a welcome message to new users
func (s *SMSService) SendWelcomeMessage(phoneNumber, userName string) error {
	message := fmt.Sprintf("Welcome to Cashout, %s! Your account has been created successfully. You can now send and receive money securely.", userName)
	
	return s.sendSMS(phoneNumber, message)
}

// SendTransactionNotification sends transaction notifications
func (s *SMSService) SendTransactionNotification(phoneNumber, transactionType, amount, transactionID string) error {
	message := fmt.Sprintf("Cashout: Your %s transaction of KES %s (ID: %s) has been completed successfully.", 
		transactionType, amount, transactionID)
	
	return s.sendSMS(phoneNumber, message)
}

// SendSecurityAlert sends security alerts
func (s *SMSService) SendSecurityAlert(phoneNumber, alertType string) error {
	message := fmt.Sprintf("Cashout Security Alert: %s. If this wasn't you, please contact support immediately.", alertType)
	
	return s.sendSMS(phoneNumber, message)
}

// SendAgentNotification sends notifications to agents
func (s *SMSService) SendAgentNotification(phoneNumber, message string) error {
	return s.sendSMS(phoneNumber, message)
}

// sendSMS sends an SMS using Africa's Talking API
func (s *SMSService) sendSMS(phoneNumber, message string) error {
	// Validate phone number format
	formattedPhone := s.formatPhoneNumber(phoneNumber)
	if formattedPhone == "" {
		return fmt.Errorf("invalid phone number format: %s", phoneNumber)
	}

	// Prepare request
	smsRequest := SMSRequest{
		Username: s.username,
		To:       []string{formattedPhone},
		Message:  message,
		From:     s.senderID,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(smsRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal SMS request")
		return fmt.Errorf("failed to prepare SMS request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SMS request")
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", s.apiKey)

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to send SMS")
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status", resp.StatusCode).Str("phone", phoneNumber).Msg("SMS API returned error status")
		return fmt.Errorf("SMS API returned status %d", resp.StatusCode)
	}

	// Parse response
	var smsResponse SMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResponse); err != nil {
		log.Error().Err(err).Msg("Failed to decode SMS response")
		return fmt.Errorf("failed to decode SMS response: %w", err)
	}

	// Check if SMS was sent successfully
	if len(smsResponse.SMS.Recipients) > 0 {
		recipient := smsResponse.SMS.Recipients[0]
		if recipient.Status == "Success" {
			log.Info().
				Str("phone", phoneNumber).
				Str("message_id", recipient.MessageID).
				Str("cost", recipient.Cost).
				Msg("SMS sent successfully")
			return nil
		} else {
			log.Error().
				Str("phone", phoneNumber).
				Str("status", recipient.Status).
				Msg("SMS delivery failed")
			return fmt.Errorf("SMS delivery failed: %s", recipient.Status)
		}
	}

	log.Error().Str("phone", phoneNumber).Msg("No recipients in SMS response")
	return fmt.Errorf("no recipients in SMS response")
}

// formatPhoneNumber formats phone number for Africa's Talking
func (s *SMSService) formatPhoneNumber(phone string) string {
	// Remove any non-digit characters except +
	cleanPhone := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' || char == '+' {
			cleanPhone += string(char)
		}
	}

	// Handle different formats
	if cleanPhone == "" {
		return ""
	}

	// If starts with 0, replace with +254
	if len(cleanPhone) > 0 && cleanPhone[0] == '0' {
		cleanPhone = "+254" + cleanPhone[1:]
	}

	// If starts with 254, add +
	if len(cleanPhone) >= 3 && cleanPhone[:3] == "254" {
		cleanPhone = "+" + cleanPhone
	}

	// If doesn't start with +, add +254
	if len(cleanPhone) > 0 && cleanPhone[0] != '+' {
		cleanPhone = "+254" + cleanPhone
	}

	// Validate final format
	if len(cleanPhone) != 13 || cleanPhone[:4] != "+254" {
		return ""
	}

	return cleanPhone
}

// ValidatePhoneNumber validates if a phone number is in correct format
func (s *SMSService) ValidatePhoneNumber(phone string) bool {
	formatted := s.formatPhoneNumber(phone)
	return formatted != ""
}

// GetSMSCredits gets remaining SMS credits
func (s *SMSService) GetSMSCredits() (string, error) {
	url := "https://api.africastalking.com/version1/user"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get SMS credits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response to get credits
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if user, ok := response["User"].(map[string]interface{}); ok {
		if balance, ok := user["balance"].(string); ok {
			return balance, nil
		}
	}

	return "Unknown", nil
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 