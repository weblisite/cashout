package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a transaction in the system
type Transaction struct {
	ID              uuid.UUID `json:"id" db:"id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	RecipientID     *uuid.UUID `json:"recipient_id,omitempty" db:"recipient_id"`
	BusinessID      *uuid.UUID `json:"business_id,omitempty" db:"business_id"`
	AgentID         *uuid.UUID `json:"agent_id,omitempty" db:"agent_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Fee             float64   `json:"fee" db:"fee"`
	AgentCommission *float64  `json:"agent_commission,omitempty" db:"agent_commission"`
	PlatformMargin  *float64  `json:"platform_margin,omitempty" db:"platform_margin"`
	Type            string    `json:"type" db:"type"`
	Status          string    `json:"status" db:"status"`
	Note            *string   `json:"note,omitempty" db:"note"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeP2P      TransactionType = "p2p"
	TransactionTypeCashIn   TransactionType = "cash_in"
	TransactionTypeCashOut  TransactionType = "cash_out"
	TransactionTypeBusiness TransactionType = "business"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// CreateTransactionRequest represents the request to create a new transaction
type CreateTransactionRequest struct {
	RecipientID *uuid.UUID `json:"recipient_id,omitempty"`
	BusinessID  *uuid.UUID `json:"business_id,omitempty"`
	AgentID     *uuid.UUID `json:"agent_id,omitempty"`
	Amount      float64    `json:"amount" validate:"required,gt=0"`
	Type        string     `json:"type" validate:"required"`
	Note        *string    `json:"note,omitempty"`
}

// P2PTransactionRequest represents a P2P transfer request
type P2PTransactionRequest struct {
	RecipientID string  `json:"recipient_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Note        *string `json:"note,omitempty"`
}

// CashInTransactionRequest represents a cash-in request
type CashInTransactionRequest struct {
	AgentID string  `json:"agent_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
}

// CashOutTransactionRequest represents a cash-out request
type CashOutTransactionRequest struct {
	AgentID string  `json:"agent_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
}

// BusinessTransactionRequest represents a business payment request
type BusinessTransactionRequest struct {
	BusinessID string  `json:"business_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Note       *string `json:"note,omitempty"`
}

// TransactionResponse represents the response for a transaction
type TransactionResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	RecipientID     *uuid.UUID `json:"recipient_id,omitempty"`
	BusinessID      *uuid.UUID `json:"business_id,omitempty"`
	AgentID         *uuid.UUID `json:"agent_id,omitempty"`
	Amount          float64    `json:"amount"`
	Fee             float64    `json:"fee"`
	AgentCommission *float64   `json:"agent_commission,omitempty"`
	PlatformMargin  *float64   `json:"platform_margin,omitempty"`
	Type            string     `json:"type"`
	Status          string     `json:"status"`
	Note            *string    `json:"note,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TransactionHistoryRequest represents the request to get transaction history
type TransactionHistoryRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Type     string `json:"type" form:"type"`
	Status   string `json:"status" form:"status"`
	FromDate string `json:"from_date" form:"from_date"`
	ToDate   string `json:"to_date" form:"to_date"`
}

// TransactionHistoryResponse represents the response for transaction history
type TransactionHistoryResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// NewTransaction creates a new transaction instance
func NewTransaction(userID uuid.UUID, amount float64, transactionType TransactionType) *Transaction {
	return &Transaction{
		ID:        uuid.New(),
		UserID:    userID,
		Amount:    amount,
		Type:      string(transactionType),
		Status:    string(TransactionStatusPending),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsP2P returns true if the transaction is a P2P transfer
func (t *Transaction) IsP2P() bool {
	return t.Type == string(TransactionTypeP2P)
}

// IsCashIn returns true if the transaction is a cash-in
func (t *Transaction) IsCashIn() bool {
	return t.Type == string(TransactionTypeCashIn)
}

// IsCashOut returns true if the transaction is a cash-out
func (t *Transaction) IsCashOut() bool {
	return t.Type == string(TransactionTypeCashOut)
}

// IsBusiness returns true if the transaction is a business payment
func (t *Transaction) IsBusiness() bool {
	return t.Type == string(TransactionTypeBusiness)
}

// IsCompleted returns true if the transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == string(TransactionStatusCompleted)
}

// IsPending returns true if the transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == string(TransactionStatusPending)
}

// IsFailed returns true if the transaction failed
func (t *Transaction) IsFailed() bool {
	return t.Status == string(TransactionStatusFailed)
}

// SetRecipient sets the recipient for the transaction
func (t *Transaction) SetRecipient(recipientID uuid.UUID) {
	t.RecipientID = &recipientID
	t.UpdatedAt = time.Now()
}

// SetBusiness sets the business for the transaction
func (t *Transaction) SetBusiness(businessID uuid.UUID) {
	t.BusinessID = &businessID
	t.UpdatedAt = time.Now()
}

// SetAgent sets the agent for the transaction
func (t *Transaction) SetAgent(agentID uuid.UUID) {
	t.AgentID = &agentID
	t.UpdatedAt = time.Now()
}

// SetFee sets the fee for the transaction
func (t *Transaction) SetFee(fee float64) {
	t.Fee = fee
	t.UpdatedAt = time.Now()
}

// SetCommission sets the commission and margin for the transaction
func (t *Transaction) SetCommission(commission, margin float64) {
	t.AgentCommission = &commission
	t.PlatformMargin = &margin
	t.UpdatedAt = time.Now()
}

// Complete marks the transaction as completed
func (t *Transaction) Complete() {
	t.Status = string(TransactionStatusCompleted)
	t.UpdatedAt = time.Now()
}

// Fail marks the transaction as failed
func (t *Transaction) Fail() {
	t.Status = string(TransactionStatusFailed)
	t.UpdatedAt = time.Now()
}

// Cancel marks the transaction as cancelled
func (t *Transaction) Cancel() {
	t.Status = string(TransactionStatusCancelled)
	t.UpdatedAt = time.Now()
}

// GetTotalAmount returns the total amount including fee
func (t *Transaction) GetTotalAmount() float64 {
	return t.Amount + t.Fee
}

// ToResponse converts the transaction to a response format
func (t *Transaction) ToResponse() TransactionResponse {
	return TransactionResponse{
		ID:              t.ID,
		UserID:          t.UserID,
		RecipientID:     t.RecipientID,
		BusinessID:      t.BusinessID,
		AgentID:         t.AgentID,
		Amount:          t.Amount,
		Fee:             t.Fee,
		AgentCommission: t.AgentCommission,
		PlatformMargin:  t.PlatformMargin,
		Type:            t.Type,
		Status:          t.Status,
		Note:            t.Note,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
} 