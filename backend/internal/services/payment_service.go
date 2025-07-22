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

// PaymentService handles payment operations using Intasend
type PaymentService struct {
	apiKey           string
	publishableKey   string
	webhookSecret    string
	baseURL          string
	httpClient       *http.Client
}

// IntasendPaymentRequest represents a payment request to Intasend
type IntasendPaymentRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	PaymentType string  `json:"payment_type"`
	PhoneNumber string  `json:"phone_number"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Reference   string  `json:"reference"`
	Description string  `json:"description"`
}

// IntasendPaymentResponse represents the response from Intasend
type IntasendPaymentResponse struct {
	State       string `json:"state"`
	Invoice     string `json:"invoice"`
	URL         string `json:"url"`
	Reference   string `json:"reference"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PaymentType string `json:"payment_type"`
}

// IntasendWebhookPayload represents webhook data from Intasend
type IntasendWebhookPayload struct {
	State       string `json:"state"`
	Invoice     string `json:"invoice"`
	Reference   string `json:"reference"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PaymentType string `json:"payment_type"`
	Signature   string `json:"signature"`
}

// PaymentResult represents the result of a payment operation
type PaymentResult struct {
	Success     bool    `json:"success"`
	TransactionID string `json:"transaction_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	PaymentURL  string  `json:"payment_url,omitempty"`
}

// NewPaymentService creates a new payment service instance
func NewPaymentService() *PaymentService {
	return &PaymentService{
		apiKey:         getEnv("INTASEND_API_KEY", ""),
		publishableKey: getEnv("INTASEND_PUBLISHABLE_KEY", ""),
		webhookSecret:  getEnv("INTASEND_WEBHOOK_SECRET", ""),
		baseURL:        "https://api.intasend.com/v1",
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}
}

// InitiatePayment initiates a payment request
func (p *PaymentService) InitiatePayment(req IntasendPaymentRequest) (*PaymentResult, error) {
	// Validate request
	if err := p.validatePaymentRequest(req); err != nil {
		return nil, fmt.Errorf("invalid payment request: %w", err)
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"amount":       req.Amount,
		"currency":     req.Currency,
		"payment_type": req.PaymentType,
		"phone_number": req.PhoneNumber,
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"email":        req.Email,
		"reference":    req.Reference,
		"description":  req.Description,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal payment request")
		return nil, fmt.Errorf("failed to prepare payment request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", p.baseURL+"/payment/requests/", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment request")
		return nil, fmt.Errorf("failed to create payment request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Token "+p.apiKey)

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send payment request")
		return nil, fmt.Errorf("failed to send payment request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status", resp.StatusCode).Msg("Payment API returned error status")
		return nil, fmt.Errorf("payment API returned status %d", resp.StatusCode)
	}

	// Parse response
	var paymentResponse IntasendPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		log.Error().Err(err).Msg("Failed to decode payment response")
		return nil, fmt.Errorf("failed to decode payment response: %w", err)
	}

	// Create result
	result := &PaymentResult{
		Success:       paymentResponse.State == "PENDING",
		TransactionID: paymentResponse.Invoice,
		Amount:        req.Amount,
		Status:        paymentResponse.State,
		Message:       "Payment initiated successfully",
		PaymentURL:    paymentResponse.URL,
	}

	log.Info().
		Str("invoice", paymentResponse.Invoice).
		Str("reference", paymentResponse.Reference).
		Str("state", paymentResponse.State).
		Msg("Payment initiated successfully")

	return result, nil
}

// CheckPaymentStatus checks the status of a payment
func (p *PaymentService) CheckPaymentStatus(invoiceID string) (*PaymentResult, error) {
	// Create HTTP request
	httpReq, err := http.NewRequest("GET", p.baseURL+"/payment/requests/"+invoiceID+"/", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create status check request")
		return nil, fmt.Errorf("failed to create status check request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Token "+p.apiKey)

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Str("invoice", invoiceID).Msg("Failed to check payment status")
		return nil, fmt.Errorf("failed to check payment status: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Str("invoice", invoiceID).Msg("Payment status API returned error")
		return nil, fmt.Errorf("payment status API returned status %d", resp.StatusCode)
	}

	// Parse response
	var paymentResponse IntasendPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		log.Error().Err(err).Msg("Failed to decode payment status response")
		return nil, fmt.Errorf("failed to decode payment status response: %w", err)
	}

	// Create result
	result := &PaymentResult{
		Success:       paymentResponse.State == "COMPLETED",
		TransactionID: paymentResponse.Invoice,
		Status:        paymentResponse.State,
		Message:       "Payment status retrieved",
	}

	log.Info().
		Str("invoice", paymentResponse.Invoice).
		Str("state", paymentResponse.State).
		Msg("Payment status checked")

	return result, nil
}

// ProcessWebhook processes webhook notifications from Intasend
func (p *PaymentService) ProcessWebhook(payload []byte, signature string) (*PaymentResult, error) {
	// Verify webhook signature
	if !p.verifyWebhookSignature(payload, signature) {
		log.Error().Msg("Invalid webhook signature")
		return nil, fmt.Errorf("invalid webhook signature")
	}

	// Parse webhook payload
	var webhookPayload IntasendWebhookPayload
	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		log.Error().Err(err).Msg("Failed to parse webhook payload")
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Create result
	result := &PaymentResult{
		Success:       webhookPayload.State == "COMPLETED",
		TransactionID: webhookPayload.Invoice,
		Status:        webhookPayload.State,
		Message:       "Webhook processed successfully",
	}

	log.Info().
		Str("invoice", webhookPayload.Invoice).
		Str("state", webhookPayload.State).
		Str("reference", webhookPayload.Reference).
		Msg("Webhook processed successfully")

	return result, nil
}

// RefundPayment refunds a payment
func (p *PaymentService) RefundPayment(invoiceID string, amount float64, reason string) (*PaymentResult, error) {
	// Prepare request body
	requestBody := map[string]interface{}{
		"amount": amount,
		"reason": reason,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal refund request")
		return nil, fmt.Errorf("failed to prepare refund request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", p.baseURL+"/payment/requests/"+invoiceID+"/refund/", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create refund request")
		return nil, fmt.Errorf("failed to create refund request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Token "+p.apiKey)

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Str("invoice", invoiceID).Msg("Failed to process refund")
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status", resp.StatusCode).Str("invoice", invoiceID).Msg("Refund API returned error status")
		return nil, fmt.Errorf("refund API returned status %d", resp.StatusCode)
	}

	// Parse response
	var refundResponse IntasendPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&refundResponse); err != nil {
		log.Error().Err(err).Msg("Failed to decode refund response")
		return nil, fmt.Errorf("failed to decode refund response: %w", err)
	}

	// Create result
	result := &PaymentResult{
		Success:       refundResponse.State == "COMPLETED",
		TransactionID: refundResponse.Invoice,
		Amount:        amount,
		Status:        refundResponse.State,
		Message:       "Refund processed successfully",
	}

	log.Info().
		Str("invoice", refundResponse.Invoice).
		Str("state", refundResponse.State).
		Float64("amount", amount).
		Msg("Refund processed successfully")

	return result, nil
}

// validatePaymentRequest validates a payment request
func (p *PaymentService) validatePaymentRequest(req IntasendPaymentRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if req.PaymentType == "" {
		return fmt.Errorf("payment type is required")
	}

	if req.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}

	if req.Reference == "" {
		return fmt.Errorf("reference is required")
	}

	return nil
}

// verifyWebhookSignature verifies the webhook signature
func (p *PaymentService) verifyWebhookSignature(payload []byte, signature string) bool {
	// TODO: Implement proper signature verification
	// For now, return true for demo purposes
	return true
}

// GetPaymentHistory gets payment history for a user
func (p *PaymentService) GetPaymentHistory(phoneNumber string, limit int) ([]PaymentResult, error) {
	// Create HTTP request
	url := fmt.Sprintf("%s/payment/requests/?phone_number=%s&limit=%d", p.baseURL, phoneNumber, limit)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment history request")
		return nil, fmt.Errorf("failed to create payment history request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Token "+p.apiKey)

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to get payment history")
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Str("phone", phoneNumber).Msg("Payment history API returned error")
		return nil, fmt.Errorf("payment history API returned status %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Results []IntasendPaymentResponse `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error().Err(err).Msg("Failed to decode payment history response")
		return nil, fmt.Errorf("failed to decode payment history response: %w", err)
	}

	// Convert to PaymentResult slice
	var results []PaymentResult
	for _, payment := range response.Results {
		result := PaymentResult{
			Success:       payment.State == "COMPLETED",
			TransactionID: payment.Invoice,
			Status:        payment.State,
			Message:       "Payment history retrieved",
		}
		results = append(results, result)
	}

	log.Info().
		Str("phone", phoneNumber).
		Int("count", len(results)).
		Msg("Payment history retrieved")

	return results, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 