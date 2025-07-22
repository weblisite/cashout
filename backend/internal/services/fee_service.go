package services

import (
	"context"
	"fmt"
	"math"
)

// FeeService handles fee calculations for all transaction types
type FeeService struct {
	supabaseService *SupabaseService
}

// FeeStructure represents the fee structure for different transaction types
type FeeStructure struct {
	MinAmount    float64 `json:"min_amount"`
	MaxAmount    float64 `json:"max_amount"`
	P2PFee       float64 `json:"p2p_fee"`
	CashOutFee   float64 `json:"cash_out_fee"`
	BusinessFee  float64 `json:"business_fee"`
}

// FeeResult represents the result of a fee calculation
type FeeResult struct {
	Amount           float64 `json:"amount"`
	Fee              float64 `json:"fee"`
	AgentCommission  float64 `json:"agent_commission,omitempty"`
	PlatformMargin   float64 `json:"platform_margin,omitempty"`
	BusinessFee      float64 `json:"business_fee,omitempty"`
	UserFee          float64 `json:"user_fee,omitempty"`
	Total            float64 `json:"total"`
	TransactionType  string  `json:"transaction_type"`
}

// NewFeeService creates a new fee service instance
func NewFeeService(supabaseService *SupabaseService) *FeeService {
	return &FeeService{
		supabaseService: supabaseService,
	}
}

// CalculateFee calculates the fee for a given amount and transaction type
func (f *FeeService) CalculateFee(ctx context.Context, amount float64, transactionType string) (*FeeResult, error) {
	// Validate amount
	if amount < 50 || amount > 1000000 {
		return nil, fmt.Errorf("amount must be between 50 and 1,000,000 KES")
	}

	// Get fee structure from database
	feeStructure, err := f.getFeeStructure(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee structure: %w", err)
	}

	result := &FeeResult{
		Amount:          amount,
		TransactionType: transactionType,
	}

	switch transactionType {
	case "p2p":
		result.Fee = feeStructure.P2PFee
		result.PlatformMargin = feeStructure.P2PFee
		result.Total = amount + result.Fee

	case "cash_out":
		result.Fee = feeStructure.CashOutFee
		result.AgentCommission = feeStructure.CashOutFee * 0.25 // 25% agent commission
		result.PlatformMargin = feeStructure.CashOutFee * 0.75 // 75% platform margin
		result.Total = amount + result.Fee

	case "business":
		result.Fee = feeStructure.BusinessFee
		result.BusinessFee = feeStructure.BusinessFee * 0.5 // 50% business fee
		result.UserFee = feeStructure.BusinessFee * 0.5     // 50% user fee
		result.Total = amount + result.Fee

	case "cash_in":
		result.Fee = 0.0 // Cash-in is always free
		result.Total = amount

	default:
		return nil, fmt.Errorf("invalid transaction type: %s", transactionType)
	}

	return result, nil
}

// GetFeeStructure retrieves the fee structure for a given amount
func (f *FeeService) getFeeStructure(ctx context.Context, amount float64) (*FeeStructure, error) {
	// This would typically query the database
	// For now, we'll use the hardcoded structure with rounded values
	
	feeStructures := []FeeStructure{
		{MinAmount: 50, MaxAmount: 100, P2PFee: 8, CashOutFee: 8, BusinessFee: 8},
		{MinAmount: 101, MaxAmount: 500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 501, MaxAmount: 1000, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 1001, MaxAmount: 1500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 1501, MaxAmount: 2500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 2501, MaxAmount: 3500, P2PFee: 39, CashOutFee: 39, BusinessFee: 39},
		{MinAmount: 3501, MaxAmount: 5000, P2PFee: 52, CashOutFee: 52, BusinessFee: 52},
		{MinAmount: 5001, MaxAmount: 7500, P2PFee: 65, CashOutFee: 65, BusinessFee: 65},
		{MinAmount: 7501, MaxAmount: 10000, P2PFee: 86, CashOutFee: 86, BusinessFee: 86},
		{MinAmount: 10001, MaxAmount: 15000, P2PFee: 125, CashOutFee: 125, BusinessFee: 125},
		{MinAmount: 15001, MaxAmount: 20000, P2PFee: 139, CashOutFee: 139, BusinessFee: 139},
		{MinAmount: 20001, MaxAmount: 35000, P2PFee: 148, CashOutFee: 148, BusinessFee: 148},
		{MinAmount: 35001, MaxAmount: 50000, P2PFee: 209, CashOutFee: 209, BusinessFee: 209},
		{MinAmount: 50001, MaxAmount: 250000, P2PFee: 232, CashOutFee: 232, BusinessFee: 232},
		{MinAmount: 250001, MaxAmount: 500000, P2PFee: 513, CashOutFee: 513, BusinessFee: 513},
		{MinAmount: 500001, MaxAmount: 1000000, P2PFee: 1076, CashOutFee: 1076, BusinessFee: 1076},
	}

	for _, fs := range feeStructures {
		if amount >= fs.MinAmount && amount <= fs.MaxAmount {
			return &fs, nil
		}
	}

	// Return the highest tier if amount exceeds maximum
	return &feeStructures[len(feeStructures)-1], nil
}

// GetCompleteFeeStructure returns the complete fee structure for display
func (f *FeeService) GetCompleteFeeStructure(ctx context.Context) ([]FeeStructure, error) {
	// Return the complete fee structure
	feeStructures := []FeeStructure{
		{MinAmount: 50, MaxAmount: 100, P2PFee: 8, CashOutFee: 8, BusinessFee: 8},
		{MinAmount: 101, MaxAmount: 500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 501, MaxAmount: 1000, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 1001, MaxAmount: 1500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 1501, MaxAmount: 2500, P2PFee: 22, CashOutFee: 22, BusinessFee: 22},
		{MinAmount: 2501, MaxAmount: 3500, P2PFee: 39, CashOutFee: 39, BusinessFee: 39},
		{MinAmount: 3501, MaxAmount: 5000, P2PFee: 52, CashOutFee: 52, BusinessFee: 52},
		{MinAmount: 5001, MaxAmount: 7500, P2PFee: 65, CashOutFee: 65, BusinessFee: 65},
		{MinAmount: 7501, MaxAmount: 10000, P2PFee: 86, CashOutFee: 86, BusinessFee: 86},
		{MinAmount: 10001, MaxAmount: 15000, P2PFee: 125, CashOutFee: 125, BusinessFee: 125},
		{MinAmount: 15001, MaxAmount: 20000, P2PFee: 139, CashOutFee: 139, BusinessFee: 139},
		{MinAmount: 20001, MaxAmount: 35000, P2PFee: 148, CashOutFee: 148, BusinessFee: 148},
		{MinAmount: 35001, MaxAmount: 50000, P2PFee: 209, CashOutFee: 209, BusinessFee: 209},
		{MinAmount: 50001, MaxAmount: 250000, P2PFee: 232, CashOutFee: 232, BusinessFee: 232},
		{MinAmount: 250001, MaxAmount: 500000, P2PFee: 513, CashOutFee: 513, BusinessFee: 513},
		{MinAmount: 500001, MaxAmount: 1000000, P2PFee: 1076, CashOutFee: 1076, BusinessFee: 1076},
	}

	return feeStructures, nil
}

// ApplyRounding applies the rounding rules (round up if ≥0.50, round down if <0.50)
func (f *FeeService) ApplyRounding(value float64) float64 {
	decimal := value - math.Floor(value)
	if decimal >= 0.50 {
		return math.Ceil(value)
	}
	return math.Floor(value)
}

// ValidateAmount validates if the amount is within acceptable limits
func (f *FeeService) ValidateAmount(amount float64) error {
	if amount < 50 {
		return fmt.Errorf("minimum transaction amount is 50 KES")
	}
	if amount > 1000000 {
		return fmt.Errorf("maximum transaction amount is 1,000,000 KES")
	}
	return nil
} 