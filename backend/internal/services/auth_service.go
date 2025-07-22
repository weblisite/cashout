package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/cashout/backend/internal/models"
	"github.com/cashout/backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AuthService handles authentication operations
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new auth service instance
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// SendOTPRequest represents the request to send OTP
type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// VerifyOTPRequest represents the request to verify OTP
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	OTP         string `json:"otp" validate:"required,len=6"`
}

// SetupPINRequest represents the request to setup PIN
type SetupPINRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	PIN         string `json:"pin" validate:"required,len=4"`
}

// VerifyPINRequest represents the request to verify PIN
type VerifyPINRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	PIN         string `json:"pin" validate:"required,len=4"`
}

// SetupBiometricRequest represents the request to setup biometric
type SetupBiometricRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	BiometricID string `json:"biometric_id" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string         `json:"token"`
	User         models.User    `json:"user"`
	IsNewUser    bool           `json:"is_new_user"`
	RequiresPIN  bool           `json:"requires_pin"`
	RequiresKYC  bool           `json:"requires_kyc"`
}

// SendOTP sends an OTP to the provided phone number
func (s *AuthService) SendOTP(req SendOTPRequest) error {
	// Generate 6-digit OTP
	otp := s.generateOTP()
	
	// Store OTP in database (in production, use Redis for better performance)
	// For now, we'll just log it
	log.Info().Str("phone", req.PhoneNumber).Str("otp", otp).Msg("OTP generated")
	
	// TODO: Integrate with SMS service (Africa's Talking)
	// For demo purposes, we'll just return success
	return nil
}

// VerifyOTP verifies the OTP and returns user information
func (s *AuthService) VerifyOTP(req VerifyOTPRequest) (*AuthResponse, error) {
	// TODO: Verify OTP from database/Redis
	// For demo purposes, accept any 6-digit code
	
	// Check if user exists
	user, err := s.getUserByPhone(req.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}
	
	var isNewUser bool
	if err == sql.ErrNoRows {
		// Create new user
		hashedID := s.generateHashedID(req.PhoneNumber)
		user = models.NewUser(req.PhoneNumber, hashedID)
		
		if err := s.createUser(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		isNewUser = true
	}
	
	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	return &AuthResponse{
		Token:       token,
		User:        *user,
		IsNewUser:   isNewUser,
		RequiresPIN: isNewUser,
		RequiresKYC: !user.IsKYCVerified(),
	}, nil
}

// SetupPIN sets up PIN for a user
func (s *AuthService) SetupPIN(req SetupPINRequest) error {
	user, err := s.getUserByPhone(req.PhoneNumber)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	
	// Hash PIN before storing (in production, use bcrypt)
	hashedPIN := s.hashPIN(req.PIN)
	
	// Store hashed PIN in database
	if err := s.storePIN(user.ID, hashedPIN); err != nil {
		return fmt.Errorf("failed to store PIN: %w", err)
	}
	
	return nil
}

// VerifyPIN verifies the PIN for a user
func (s *AuthService) VerifyPIN(req VerifyPINRequest) (*AuthResponse, error) {
	user, err := s.getUserByPhone(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// Verify PIN
	if err := s.verifyStoredPIN(user.ID, req.PIN); err != nil {
		return nil, fmt.Errorf("invalid PIN: %w", err)
	}
	
	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	return &AuthResponse{
		Token:       token,
		User:        *user,
		IsNewUser:   false,
		RequiresPIN: false,
		RequiresKYC: !user.IsKYCVerified(),
	}, nil
}

// SetupBiometric sets up biometric authentication for a user
func (s *AuthService) SetupBiometric(req SetupBiometricRequest) error {
	user, err := s.getUserByPhone(req.PhoneNumber)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	
	// Store biometric ID in database
	if err := s.storeBiometricID(user.ID, req.BiometricID); err != nil {
		return fmt.Errorf("failed to store biometric ID: %w", err)
	}
	
	return nil
}

// ValidateToken validates a JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.GetJWTSecret()), nil
	})
	
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("invalid user_id in token")
		}
		
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
		}
		
		return userID, nil
	}
	
	return uuid.Nil, fmt.Errorf("invalid token claims")
}

// generateOTP generates a 6-digit OTP
func (s *AuthService) generateOTP() string {
	// Generate random number between 100000 and 999999
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return strconv.FormatInt(n.Int64()+100000, 10)
}

// generateHashedID generates a hashed ID from phone number
func (s *AuthService) generateHashedID(phoneNumber string) string {
	// Generate random salt
	salt := make([]byte, 4)
	rand.Read(salt)
	
	// Hash phone number + salt
	hash := sha256.Sum256([]byte(phoneNumber + hex.EncodeToString(salt)))
	
	// Take first 10 characters of hash + salt
	hashStr := hex.EncodeToString(hash[:])
	return hashStr[:10] + hex.EncodeToString(salt)
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.GetJWTSecret()))
}

// hashPIN hashes a PIN (in production, use bcrypt)
func (s *AuthService) hashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}

// getUserByPhone retrieves a user by phone number
func (s *AuthService) getUserByPhone(phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, phone_number, hashed_id, kyc_status, wallet_balance, user_status, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`
	
	var user models.User
	err := s.db.QueryRow(query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.HashedID,
		&user.KYCStatus,
		&user.WalletBalance,
		&user.UserStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// createUser creates a new user in the database
func (s *AuthService) createUser(user *models.User) error {
	query := `
		INSERT INTO users (id, phone_number, hashed_id, kyc_status, wallet_balance, user_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := s.db.Exec(query,
		user.ID,
		user.PhoneNumber,
		user.HashedID,
		user.KYCStatus,
		user.WalletBalance,
		user.UserStatus,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

// storePIN stores the hashed PIN in the database
func (s *AuthService) storePIN(userID uuid.UUID, hashedPIN string) error {
	query := `
		INSERT INTO user_pins (user_id, hashed_pin, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			hashed_pin = EXCLUDED.hashed_pin,
			updated_at = EXCLUDED.updated_at
	`
	
	now := time.Now()
	_, err := s.db.Exec(query, userID, hashedPIN, now, now)
	return err
}

// verifyStoredPIN verifies the stored PIN
func (s *AuthService) verifyStoredPIN(userID uuid.UUID, pin string) error {
	hashedPIN := s.hashPIN(pin)
	
	query := `
		SELECT hashed_pin FROM user_pins WHERE user_id = $1
	`
	
	var storedHashedPIN string
	err := s.db.QueryRow(query, userID).Scan(&storedHashedPIN)
	if err != nil {
		return err
	}
	
	if hashedPIN != storedHashedPIN {
		return fmt.Errorf("invalid PIN")
	}
	
	return nil
}

// storeBiometricID stores the biometric ID in the database
func (s *AuthService) storeBiometricID(userID uuid.UUID, biometricID string) error {
	query := `
		INSERT INTO user_biometrics (user_id, biometric_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			biometric_id = EXCLUDED.biometric_id,
			updated_at = EXCLUDED.updated_at
	`
	
	now := time.Now()
	_, err := s.db.Exec(query, userID, biometricID, now, now)
	return err
} 