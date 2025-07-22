package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// KYCAMLService handles all KYC/AML operations
type KYCAMLService struct {
	supabaseService *SupabaseService
	riskManagement  *RiskManagementService
}

// NewKYCAMLService creates a new KYC/AML service
func NewKYCAMLService(supabaseService *SupabaseService, riskManagement *RiskManagementService) *KYCAMLService {
	return &KYCAMLService{
		supabaseService: supabaseService,
		riskManagement:  riskManagement,
	}
}

// KYCProfile represents a KYC profile
type KYCProfile struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	Status            string    `json:"status"` // pending, approved, rejected, under_review
	VerificationLevel string    `json:"verification_level"` // basic, enhanced, corporate
	RiskLevel         string    `json:"risk_level"` // low, medium, high
	Documents         []KYCDocument `json:"documents"`
	PersonalInfo      PersonalInfo `json:"personal_info"`
	AddressInfo       AddressInfo `json:"address_info"`
	EmploymentInfo    *EmploymentInfo `json:"employment_info,omitempty"`
	SourceOfFunds     *SourceOfFunds `json:"source_of_funds,omitempty"`
	ReviewNotes       string    `json:"review_notes"`
	ReviewedBy        *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// KYCDocument represents a KYC document
type KYCDocument struct {
	ID           uuid.UUID `json:"id"`
	KYCProfileID uuid.UUID `json:"kyc_profile_id"`
	DocumentType string    `json:"document_type"` // national_id, passport, driving_license, utility_bill, bank_statement
	DocumentNumber string  `json:"document_number"`
	IssuingCountry string  `json:"issuing_country"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	FileURL       string    `json:"file_url"`
	Status        string    `json:"status"` // pending, verified, rejected
	VerificationNotes string `json:"verification_notes"`
	VerifiedBy    *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// PersonalInfo represents personal information
type PersonalInfo struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	MiddleName    *string   `json:"middle_name,omitempty"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Gender        string    `json:"gender"`
	Nationality   string    `json:"nationality"`
	PhoneNumber   string    `json:"phone_number"`
	EmailAddress  string    `json:"email_address"`
	TaxID         *string   `json:"tax_id,omitempty"`
}

// AddressInfo represents address information
type AddressInfo struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	AddressType   string `json:"address_type"` // residential, business, mailing
	IsPermanent   bool   `json:"is_permanent"`
}

// EmploymentInfo represents employment information
type EmploymentInfo struct {
	EmployerName    string `json:"employer_name"`
	JobTitle        string `json:"job_title"`
	Industry        string `json:"industry"`
	EmploymentType  string `json:"employment_type"` // full_time, part_time, self_employed, retired
	AnnualIncome    *float64 `json:"annual_income,omitempty"`
	WorkPhone       *string `json:"work_phone,omitempty"`
	WorkEmail       *string `json:"work_email,omitempty"`
}

// SourceOfFunds represents source of funds information
type SourceOfFunds struct {
	PrimarySource   string `json:"primary_source"` // employment, business, investment, inheritance, other
	SecondarySource *string `json:"secondary_source,omitempty"`
	EmployerName    *string `json:"employer_name,omitempty"`
	BusinessName    *string `json:"business_name,omitempty"`
	AnnualIncome    *float64 `json:"annual_income,omitempty"`
	ExpectedMonthlyVolume *float64 `json:"expected_monthly_volume,omitempty"`
}

// AMLAlert represents an AML alert
type AMLAlert struct {
	ID            uuid.UUID `json:"id"`
	KYCProfileID  uuid.UUID `json:"kyc_profile_id"`
	AlertType     string    `json:"alert_type"` // suspicious_activity, high_risk_country, pep, sanctions
	Severity      string    `json:"severity"` // low, medium, high, critical
	Description   string    `json:"description"`
	RiskScore     float64   `json:"risk_score"`
	Status        string    `json:"status"` // active, investigated, resolved, false_positive
	InvestigationNotes string `json:"investigation_notes"`
	CreatedAt     time.Time `json:"created_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
}

// KYC Methods
func (k *KYCAMLService) CreateKYCProfile(ctx context.Context, userID uuid.UUID, personalInfo PersonalInfo, addressInfo AddressInfo) (*KYCProfile, error) {
	// Validate personal information
	if err := k.validatePersonalInfo(personalInfo); err != nil {
		return nil, fmt.Errorf("invalid personal information: %w", err)
	}

	// Validate address information
	if err := k.validateAddressInfo(addressInfo); err != nil {
		return nil, fmt.Errorf("invalid address information: %w", err)
	}

	// Calculate initial risk level
	riskLevel := k.calculateInitialRiskLevel(personalInfo, addressInfo)

	profile := &KYCProfile{
		ID:                uuid.New(),
		UserID:            userID,
		Status:            "pending",
		VerificationLevel: "basic",
		RiskLevel:         riskLevel,
		PersonalInfo:      personalInfo,
		AddressInfo:       addressInfo,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := k.supabaseService.CreateKYCProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create KYC profile: %w", err)
	}

	log.Printf("Created KYC profile for user %s with risk level: %s", userID, riskLevel)
	return profile, nil
}

func (k *KYCAMLService) validatePersonalInfo(info PersonalInfo) error {
	// Validate required fields
	if strings.TrimSpace(info.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	if strings.TrimSpace(info.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	if info.DateOfBirth.IsZero() {
		return fmt.Errorf("date of birth is required")
	}

	// Validate age (must be 18+)
	age := time.Now().Year() - info.DateOfBirth.Year()
	if age < 18 {
		return fmt.Errorf("user must be at least 18 years old")
	}

	// Validate phone number
	if err := k.validatePhoneNumber(info.PhoneNumber); err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	// Validate email
	if err := k.validateEmail(info.EmailAddress); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

func (k *KYCAMLService) validateAddressInfo(info AddressInfo) error {
	// Validate required fields
	if strings.TrimSpace(info.StreetAddress) == "" {
		return fmt.Errorf("street address is required")
	}
	if strings.TrimSpace(info.City) == "" {
		return fmt.Errorf("city is required")
	}
	if strings.TrimSpace(info.Country) == "" {
		return fmt.Errorf("country is required")
	}

	return nil
}

func (k *KYCAMLService) validatePhoneNumber(phone string) error {
	// Basic phone number validation
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

func (k *KYCAMLService) validateEmail(email string) error {
	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (k *KYCAMLService) calculateInitialRiskLevel(personalInfo PersonalInfo, addressInfo AddressInfo) string {
	riskScore := 0

	// Age risk
	age := time.Now().Year() - personalInfo.DateOfBirth.Year()
	if age < 25 || age > 65 {
		riskScore += 1
	}

	// Country risk
	if k.isHighRiskCountry(addressInfo.Country) {
		riskScore += 2
	}

	// Determine risk level
	if riskScore >= 2 {
		return "high"
	} else if riskScore >= 1 {
		return "medium"
	}
	return "low"
}

func (k *KYCAMLService) isHighRiskCountry(country string) bool {
	// High-risk countries (simplified list)
	highRiskCountries := map[string]bool{
		"AF": true, // Afghanistan
		"IR": true, // Iran
		"KP": true, // North Korea
		"SD": true, // Sudan
		"SY": true, // Syria
	}

	return highRiskCountries[strings.ToUpper(country)]
}

func (k *KYCAMLService) AddDocument(ctx context.Context, kycProfileID uuid.UUID, documentType, documentNumber, issuingCountry, fileURL string, expiryDate *time.Time) (*KYCDocument, error) {
	document := &KYCDocument{
		ID:             uuid.New(),
		KYCProfileID:   kycProfileID,
		DocumentType:   documentType,
		DocumentNumber: documentNumber,
		IssuingCountry: issuingCountry,
		ExpiryDate:     expiryDate,
		FileURL:        fileURL,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	if err := k.supabaseService.CreateKYCDocument(ctx, document); err != nil {
		return nil, fmt.Errorf("failed to create KYC document: %w", err)
	}

	log.Printf("Added %s document for KYC profile %s", documentType, kycProfileID)
	return document, nil
}

func (k *KYCAMLService) VerifyDocument(ctx context.Context, documentID uuid.UUID, status string, notes string, verifiedBy uuid.UUID) error {
	document, err := k.supabaseService.GetKYCDocumentByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to get KYC document: %w", err)
	}

	document.Status = status
	document.VerificationNotes = notes
	document.VerifiedBy = &verifiedBy
	now := time.Now()
	document.VerifiedAt = &now

	if err := k.supabaseService.UpdateKYCDocument(ctx, document); err != nil {
		return fmt.Errorf("failed to update KYC document: %w", err)
	}

	log.Printf("Verified document %s: %s", documentID, status)
	return nil
}

func (k *KYCAMLService) UpdateEmploymentInfo(ctx context.Context, kycProfileID uuid.UUID, employmentInfo EmploymentInfo) error {
	profile, err := k.supabaseService.GetKYCProfileByID(ctx, kycProfileID)
	if err != nil {
		return fmt.Errorf("failed to get KYC profile: %w", err)
	}

	profile.EmploymentInfo = &employmentInfo
	profile.UpdatedAt = time.Now()

	// Recalculate risk level
	profile.RiskLevel = k.recalculateRiskLevel(profile)

	if err := k.supabaseService.UpdateKYCProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update KYC profile: %w", err)
	}

	log.Printf("Updated employment info for KYC profile %s", kycProfileID)
	return nil
}

func (k *KYCAMLService) UpdateSourceOfFunds(ctx context.Context, kycProfileID uuid.UUID, sourceOfFunds SourceOfFunds) error {
	profile, err := k.supabaseService.GetKYCProfileByID(ctx, kycProfileID)
	if err != nil {
		return fmt.Errorf("failed to get KYC profile: %w", err)
	}

	profile.SourceOfFunds = &sourceOfFunds
	profile.UpdatedAt = time.Now()

	// Recalculate risk level
	profile.RiskLevel = k.recalculateRiskLevel(profile)

	if err := k.supabaseService.UpdateKYCProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update KYC profile: %w", err)
	}

	log.Printf("Updated source of funds for KYC profile %s", kycProfileID)
	return nil
}

func (k *KYCAMLService) recalculateRiskLevel(profile *KYCProfile) string {
	riskScore := 0

	// Base risk from personal info
	if k.isHighRiskCountry(profile.AddressInfo.Country) {
		riskScore += 2
	}

	// Employment risk
	if profile.EmploymentInfo != nil {
		if profile.EmploymentInfo.AnnualIncome != nil && *profile.EmploymentInfo.AnnualIncome < 50000 {
			riskScore += 1
		}
		if profile.EmploymentInfo.EmploymentType == "self_employed" {
			riskScore += 1
		}
	}

	// Source of funds risk
	if profile.SourceOfFunds != nil {
		if profile.SourceOfFunds.PrimarySource == "inheritance" {
			riskScore += 1
		}
		if profile.SourceOfFunds.ExpectedMonthlyVolume != nil && *profile.SourceOfFunds.ExpectedMonthlyVolume > 1000000 {
			riskScore += 1
		}
	}

	// Document verification risk
	verifiedDocuments := 0
	for _, doc := range profile.Documents {
		if doc.Status == "verified" {
			verifiedDocuments++
		}
	}
	if verifiedDocuments < 2 {
		riskScore += 1
	}

	// Determine risk level
	if riskScore >= 3 {
		return "high"
	} else if riskScore >= 1 {
		return "medium"
	}
	return "low"
}

func (k *KYCAMLService) ReviewKYCProfile(ctx context.Context, kycProfileID uuid.UUID, status string, notes string, reviewedBy uuid.UUID) error {
	profile, err := k.supabaseService.GetKYCProfileByID(ctx, kycProfileID)
	if err != nil {
		return fmt.Errorf("failed to get KYC profile: %w", err)
	}

	profile.Status = status
	profile.ReviewNotes = notes
	profile.ReviewedBy = &reviewedBy
	now := time.Now()
	profile.ReviewedAt = &now
	profile.UpdatedAt = now

	if err := k.supabaseService.UpdateKYCProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update KYC profile: %w", err)
	}

	log.Printf("Reviewed KYC profile %s: %s", kycProfileID, status)
	return nil
}

func (k *KYCAMLService) GetKYCProfileByUserID(ctx context.Context, userID uuid.UUID) (*KYCProfile, error) {
	return k.supabaseService.GetKYCProfileByUserID(ctx, userID)
}

func (k *KYCAMLService) GetPendingKYCProfiles(ctx context.Context) ([]KYCProfile, error) {
	return k.supabaseService.GetKYCProfilesByStatus(ctx, "pending")
}

// AML Methods
func (k *KYCAMLService) PerformAMLCheck(ctx context.Context, kycProfileID uuid.UUID) error {
	profile, err := k.supabaseService.GetKYCProfileByID(ctx, kycProfileID)
	if err != nil {
		return fmt.Errorf("failed to get KYC profile: %w", err)
	}

	var alerts []AMLAlert

	// Check for PEP (Politically Exposed Person)
	if k.isPEP(profile.PersonalInfo) {
		alert := AMLAlert{
			ID:           uuid.New(),
			KYCProfileID: kycProfileID,
			AlertType:    "pep",
			Severity:     "high",
			Description:  "Potential Politically Exposed Person detected",
			RiskScore:    0.8,
			Status:       "active",
			CreatedAt:    time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check for sanctions
	if k.isSanctioned(profile.PersonalInfo, profile.AddressInfo) {
		alert := AMLAlert{
			ID:           uuid.New(),
			KYCProfileID: kycProfileID,
			AlertType:    "sanctions",
			Severity:     "critical",
			Description:  "Individual or country under sanctions",
			RiskScore:    1.0,
			Status:       "active",
			CreatedAt:    time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check for high-risk country
	if k.isHighRiskCountry(profile.AddressInfo.Country) {
		alert := AMLAlert{
			ID:           uuid.New(),
			KYCProfileID: kycProfileID,
			AlertType:    "high_risk_country",
			Severity:     "medium",
			Description:  fmt.Sprintf("High-risk country: %s", profile.AddressInfo.Country),
			RiskScore:    0.6,
			Status:       "active",
			CreatedAt:    time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Save alerts
	for _, alert := range alerts {
		if err := k.supabaseService.CreateAMLAlert(ctx, &alert); err != nil {
			return fmt.Errorf("failed to create AML alert: %w", err)
		}
		log.Printf("🚨 AML alert created: %s for profile %s", alert.AlertType, kycProfileID)
	}

	return nil
}

func (k *KYCAMLService) isPEP(personalInfo PersonalInfo) bool {
	// Simplified PEP check (in production, would use external PEP databases)
	// This is a mock implementation
	return false
}

func (k *KYCAMLService) isSanctioned(personalInfo PersonalInfo, addressInfo AddressInfo) bool {
	// Simplified sanctions check (in production, would use external sanctions databases)
	// This is a mock implementation
	return false
}

func (k *KYCAMLService) GetActiveAMLAlerts(ctx context.Context) ([]AMLAlert, error) {
	return k.supabaseService.GetAMLAlertsByStatus(ctx, "active")
}

func (k *KYCAMLService) ResolveAMLAlert(ctx context.Context, alertID uuid.UUID, status string, notes string) error {
	alert, err := k.supabaseService.GetAMLAlertByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get AML alert: %w", err)
	}

	alert.Status = status
	alert.InvestigationNotes = notes
	if status != "active" {
		now := time.Now()
		alert.ResolvedAt = &now
	}

	if err := k.supabaseService.UpdateAMLAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to update AML alert: %w", err)
	}

	log.Printf("Resolved AML alert: %s", alert.Description)
	return nil
}

// Health Check
func (k *KYCAMLService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := k.supabaseService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("KYC/AML service health check failed: %w", err)
	}

	log.Println("✅ KYC/AML service health check passed")
	return nil
} 