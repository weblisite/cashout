package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cashout/backend/configs"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

// SupabaseService provides methods for interacting with Supabase
type SupabaseService struct {
	client *supabase.Client
}

// NewSupabaseService creates a new Supabase service instance
func NewSupabaseService() *SupabaseService {
	return &SupabaseService{
		client: configs.GetSupabaseClient(),
	}
}

// User represents a user in the database
type User struct {
	ID            uuid.UUID       `json:"id"`
	Email         string          `json:"email"`
	PhoneNumber   string          `json:"phone_number"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	DateOfBirth   *time.Time      `json:"date_of_birth,omitempty"`
	NationalID    *string         `json:"national_id,omitempty"`
	KYCStatus     string          `json:"kyc_status"`
	KYCDocuments  json.RawMessage `json:"kyc_documents,omitempty"`
	WalletBalance float64         `json:"wallet_balance"`
	IsActive      bool            `json:"is_active"`
	IsVerified    bool            `json:"is_verified"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Agent represents an agent in the database
type Agent struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	AgentCode       string     `json:"agent_code"`
	BusinessName    *string    `json:"business_name,omitempty"`
	BusinessAddress *string    `json:"business_address,omitempty"`
	BusinessPhone   *string    `json:"business_phone,omitempty"`
	BusinessEmail   *string    `json:"business_email,omitempty"`
	AgentType       string     `json:"agent_type"`
	CommissionRate  float64    `json:"commission_rate"`
	FloatBalance    float64    `json:"float_balance"`
	MaxFloatLimit   float64    `json:"max_float_limit"`
	Status          string     `json:"status"`
	IsVerified      bool       `json:"is_verified"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Business represents a business in the database
type Business struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	BusinessName    string     `json:"business_name"`
	BusinessType    *string    `json:"business_type,omitempty"`
	BusinessAddress *string    `json:"business_address,omitempty"`
	BusinessPhone   *string    `json:"business_phone,omitempty"`
	BusinessEmail   *string    `json:"business_email,omitempty"`
	TaxID           *string    `json:"tax_id,omitempty"`
	BusinessLicense *string    `json:"business_license,omitempty"`
	QRCode          *string    `json:"qr_code,omitempty"`
	WalletBalance   float64    `json:"wallet_balance"`
	Status          string     `json:"status"`
	IsVerified      bool       `json:"is_verified"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Transaction represents a transaction in the database
type Transaction struct {
	ID              uuid.UUID       `json:"id"`
	TransactionID   string          `json:"transaction_id"`
	UserID          *uuid.UUID      `json:"user_id,omitempty"`
	AgentID         *uuid.UUID      `json:"agent_id,omitempty"`
	BusinessID      *uuid.UUID      `json:"business_id,omitempty"`
	TransactionType string          `json:"transaction_type"`
	Amount          float64         `json:"amount"`
	Fee             float64         `json:"fee"`
	TotalAmount     float64         `json:"total_amount"`
	Currency        string          `json:"currency"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	PaymentProvider *string         `json:"payment_provider,omitempty"`
	PaymentReference *string        `json:"payment_reference,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// QRCode represents a QR code in the database
type QRCode struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	AgentID      *uuid.UUID      `json:"agent_id,omitempty"`
	BusinessID   *uuid.UUID      `json:"business_id,omitempty"`
	QRType       string          `json:"qr_type"`
	Amount       *float64        `json:"amount,omitempty"`
	QRData       string          `json:"qr_data"`
	QRImageURL   *string         `json:"qr_image_url,omitempty"`
	IsActive     bool            `json:"is_active"`
	ExpiresAt    *time.Time      `json:"expires_at,omitempty"`
	UsedAt       *time.Time      `json:"used_at,omitempty"`
	TransactionID *uuid.UUID     `json:"transaction_id,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

// Notification represents a notification in the database
type Notification struct {
	ID             uuid.UUID       `json:"id"`
	UserID         *uuid.UUID      `json:"user_id,omitempty"`
	AgentID        *uuid.UUID      `json:"agent_id,omitempty"`
	BusinessID     *uuid.UUID      `json:"business_id,omitempty"`
	NotificationType string       `json:"notification_type"`
	Title          string          `json:"title"`
	Message        string          `json:"message"`
	Data           json.RawMessage `json:"data,omitempty"`
	IsRead         bool            `json:"is_read"`
	ReadAt         *time.Time      `json:"read_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// OTPCode represents an OTP code in the database
type OTPCode struct {
	ID         uuid.UUID  `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	OTPCode    string     `json:"otp_code"`
	OTPType    string     `json:"otp_type"`
	IsUsed     bool       `json:"is_used"`
	ExpiresAt  time.Time  `json:"expires_at"`
	UsedAt     *time.Time `json:"used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// User Methods
func (s *SupabaseService) CreateUser(ctx context.Context, user *User) error {
	_, err := s.client.DB.From("users").Insert(user).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	_, err := s.client.DB.From("users").Select("*").Eq("id", id).Single().Execute(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *SupabaseService) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	_, err := s.client.DB.From("users").Select("*").Eq("phone_number", phone).Single().Execute(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	return &user, nil
}

func (s *SupabaseService) UpdateUser(ctx context.Context, user *User) error {
	_, err := s.client.DB.From("users").Update(user).Eq("id", user.ID).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Agent Methods
func (s *SupabaseService) CreateAgent(ctx context.Context, agent *Agent) error {
	_, err := s.client.DB.From("agents").Insert(agent).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetAgentByID(ctx context.Context, id uuid.UUID) (*Agent, error) {
	var agent Agent
	_, err := s.client.DB.From("agents").Select("*").Eq("id", id).Single().Execute(ctx, &agent)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return &agent, nil
}

func (s *SupabaseService) GetAgentByUserID(ctx context.Context, userID uuid.UUID) (*Agent, error) {
	var agent Agent
	_, err := s.client.DB.From("agents").Select("*").Eq("user_id", userID).Single().Execute(ctx, &agent)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent by user ID: %w", err)
	}
	return &agent, nil
}

func (s *SupabaseService) UpdateAgent(ctx context.Context, agent *Agent) error {
	_, err := s.client.DB.From("agents").Update(agent).Eq("id", agent.ID).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}
	return nil
}

// Business Methods
func (s *SupabaseService) CreateBusiness(ctx context.Context, business *Business) error {
	_, err := s.client.DB.From("businesses").Insert(business).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetBusinessByID(ctx context.Context, id uuid.UUID) (*Business, error) {
	var business Business
	_, err := s.client.DB.From("businesses").Select("*").Eq("id", id).Single().Execute(ctx, &business)
	if err != nil {
		return nil, fmt.Errorf("failed to get business: %w", err)
	}
	return &business, nil
}

func (s *SupabaseService) GetBusinessByUserID(ctx context.Context, userID uuid.UUID) (*Business, error) {
	var business Business
	_, err := s.client.DB.From("businesses").Select("*").Eq("user_id", userID).Single().Execute(ctx, &business)
	if err != nil {
		return nil, fmt.Errorf("failed to get business by user ID: %w", err)
	}
	return &business, nil
}

func (s *SupabaseService) UpdateBusiness(ctx context.Context, business *Business) error {
	_, err := s.client.DB.From("businesses").Update(business).Eq("id", business.ID).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to update business: %w", err)
	}
	return nil
}

// Transaction Methods
func (s *SupabaseService) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	_, err := s.client.DB.From("transactions").Insert(transaction).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	var transaction Transaction
	_, err := s.client.DB.From("transactions").Select("*").Eq("id", id).Single().Execute(ctx, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

func (s *SupabaseService) GetTransactionByTransactionID(ctx context.Context, transactionID string) (*Transaction, error) {
	var transaction Transaction
	_, err := s.client.DB.From("transactions").Select("*").Eq("transaction_id", transactionID).Single().Execute(ctx, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by transaction ID: %w", err)
	}
	return &transaction, nil
}

func (s *SupabaseService) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	_, err := s.client.DB.From("transactions").Update(transaction).Eq("id", transaction.ID).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Transaction, error) {
	var transactions []Transaction
	_, err := s.client.DB.From("transactions").Select("*").Eq("user_id", userID).Order("created_at", &supabase.OrderOpts{Ascending: false}).Range(offset, offset+limit-1).Execute(ctx, &transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	return transactions, nil
}

// QR Code Methods
func (s *SupabaseService) CreateQRCode(ctx context.Context, qrCode *QRCode) error {
	_, err := s.client.DB.From("qr_codes").Insert(qrCode).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create QR code: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetQRCodeByID(ctx context.Context, id uuid.UUID) (*QRCode, error) {
	var qrCode QRCode
	_, err := s.client.DB.From("qr_codes").Select("*").Eq("id", id).Single().Execute(ctx, &qrCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}
	return &qrCode, nil
}

func (s *SupabaseService) GetActiveQRCodeByData(ctx context.Context, qrData string) (*QRCode, error) {
	var qrCode QRCode
	_, err := s.client.DB.From("qr_codes").Select("*").Eq("qr_data", qrData).Eq("is_active", true).Single().Execute(ctx, &qrCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get active QR code: %w", err)
	}
	return &qrCode, nil
}

func (s *SupabaseService) UpdateQRCode(ctx context.Context, qrCode *QRCode) error {
	_, err := s.client.DB.From("qr_codes").Update(qrCode).Eq("id", qrCode.ID).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to update QR code: %w", err)
	}
	return nil
}

// Notification Methods
func (s *SupabaseService) CreateNotification(ctx context.Context, notification *Notification) error {
	_, err := s.client.DB.From("notifications").Insert(notification).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error) {
	var notifications []Notification
	_, err := s.client.DB.From("notifications").Select("*").Eq("user_id", userID).Order("created_at", &supabase.OrderOpts{Ascending: false}).Range(offset, offset+limit-1).Execute(ctx, &notifications)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	return notifications, nil
}

func (s *SupabaseService) MarkNotificationAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := s.client.DB.From("notifications").Update(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Eq("id", id).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// OTP Methods
func (s *SupabaseService) CreateOTP(ctx context.Context, otp *OTPCode) error {
	_, err := s.client.DB.From("otp_codes").Insert(otp).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}
	return nil
}

func (s *SupabaseService) GetValidOTP(ctx context.Context, phoneNumber, otpCode, otpType string) (*OTPCode, error) {
	var otp OTPCode
	_, err := s.client.DB.From("otp_codes").Select("*").Eq("phone_number", phoneNumber).Eq("otp_code", otpCode).Eq("otp_type", otpType).Eq("is_used", false).Gt("expires_at", time.Now()).Single().Execute(ctx, &otp)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid OTP: %w", err)
	}
	return &otp, nil
}

func (s *SupabaseService) MarkOTPAsUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := s.client.DB.From("otp_codes").Update(map[string]interface{}{
		"is_used": true,
		"used_at": now,
	}).Eq("id", id).Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}
	return nil
}

// Utility Methods
func (s *SupabaseService) HealthCheck(ctx context.Context) error {
	_, err := s.client.DB.From("users").Select("count", false, "", "", "").Execute(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	log.Println("✅ Supabase health check passed")
	return nil
} 