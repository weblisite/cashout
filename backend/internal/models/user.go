package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PhoneNumber   string    `json:"phone_number" db:"phone_number"`
	HashedID      string    `json:"hashed_id" db:"hashed_id"`
	KYCStatus     string    `json:"kyc_status" db:"kyc_status"`
	WalletBalance float64   `json:"wallet_balance" db:"wallet_balance"`
	UserStatus    string    `json:"user_status" db:"user_status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// UserType represents the type of user
type UserType string

const (
	UserTypeRegular UserType = "regular"
	UserTypeAgent   UserType = "agent"
	UserTypeBusiness UserType = "business"
)

// KYCStatus represents the KYC verification status
type KYCStatus string

const (
	KYCStatusPending   KYCStatus = "pending"
	KYCStatusVerified  KYCStatus = "verified"
	KYCStatusRejected  KYCStatus = "rejected"
	KYCStatusNotStarted KYCStatus = "not_started"
)

// UserStatus represents the user account status
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	HashedID    string `json:"hashed_id" validate:"required"`
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	KYCStatus     *string  `json:"kyc_status,omitempty"`
	WalletBalance *float64 `json:"wallet_balance,omitempty"`
	UserStatus    *string  `json:"user_status,omitempty"`
}

// UserProfile represents the user profile information
type UserProfile struct {
	ID            uuid.UUID `json:"id"`
	PhoneNumber   string    `json:"phone_number"`
	HashedID      string    `json:"hashed_id"`
	KYCStatus     string    `json:"kyc_status"`
	WalletBalance float64   `json:"wallet_balance"`
	UserStatus    string    `json:"user_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Agent represents an agent in the system
type Agent struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AgentCode     string    `json:"agent_code" db:"agent_code"`
	FloatBalance  float64   `json:"float_balance" db:"float_balance"`
	CommissionRate float64  `json:"commission_rate" db:"commission_rate"`
	Status        string    `json:"status" db:"status"`
	Location      string    `json:"location" db:"location"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Business represents a business in the system
type Business struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	BusinessName  string    `json:"business_name" db:"business_name"`
	BusinessType  string    `json:"business_type" db:"business_type"`
	Address       string    `json:"address" db:"address"`
	PhoneNumber   string    `json:"phone_number" db:"phone_number"`
	Email         string    `json:"email" db:"email"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user instance
func NewUser(phoneNumber, hashedID string) *User {
	return &User{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		HashedID:    hashedID,
		KYCStatus:   string(KYCStatusNotStarted),
		UserStatus:  string(UserStatusActive),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsKYCVerified returns true if the user's KYC is verified
func (u *User) IsKYCVerified() bool {
	return u.KYCStatus == string(KYCStatusVerified)
}

// IsActive returns true if the user account is active
func (u *User) IsActive() bool {
	return u.UserStatus == string(UserStatusActive)
}

// HasSufficientBalance checks if user has sufficient balance for a transaction
func (u *User) HasSufficientBalance(amount float64) bool {
	return u.WalletBalance >= amount
}

// UpdateWalletBalance updates the user's wallet balance
func (u *User) UpdateWalletBalance(newBalance float64) {
	u.WalletBalance = newBalance
	u.UpdatedAt = time.Now()
}

// UpdateKYCStatus updates the user's KYC status
func (u *User) UpdateKYCStatus(status KYCStatus) {
	u.KYCStatus = string(status)
	u.UpdatedAt = time.Now()
}

// UpdateUserStatus updates the user's account status
func (u *User) UpdateUserStatus(status UserStatus) {
	u.UserStatus = string(status)
	u.UpdatedAt = time.Now()
} 