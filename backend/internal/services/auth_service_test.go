package services

import (
	"database/sql"
	"testing"

	"github.com/cashout/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock database for testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*sql.Row)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func TestAuthService_SendOTP(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	req := SendOTPRequest{
		PhoneNumber: "+254700000000",
	}

	err := authService.SendOTP(req)
	assert.NoError(t, err)
}

func TestAuthService_VerifyOTP(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	req := VerifyOTPRequest{
		PhoneNumber: "+254700000000",
		OTP:         "123456",
	}

	// Mock database response for new user
	user := models.NewUser(req.PhoneNumber, "test-hash-id")
	
	mockDB.On("QueryRow", mock.AnythingOfType("string"), req.PhoneNumber).Return(&sql.Row{})
	mockDB.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(sqlmock.NewResult(1, 1), nil)

	response, err := authService.VerifyOTP(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.IsNewUser)
	assert.Equal(t, req.PhoneNumber, response.User.PhoneNumber)
}

func TestAuthService_SetupPIN(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	req := SetupPINRequest{
		PhoneNumber: "+254700000000",
		PIN:         "1234",
	}

	user := models.NewUser(req.PhoneNumber, "test-hash-id")
	
	mockDB.On("QueryRow", mock.AnythingOfType("string"), req.PhoneNumber).Return(&sql.Row{})
	mockDB.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(sqlmock.NewResult(1, 1), nil)

	err := authService.SetupPIN(req)
	assert.NoError(t, err)
}

func TestAuthService_VerifyPIN(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	req := VerifyPINRequest{
		PhoneNumber: "+254700000000",
		PIN:         "1234",
	}

	user := models.NewUser(req.PhoneNumber, "test-hash-id")
	
	mockDB.On("QueryRow", mock.AnythingOfType("string"), req.PhoneNumber).Return(&sql.Row{})
	mockDB.On("QueryRow", mock.AnythingOfType("string"), user.ID).Return(&sql.Row{})

	response, err := authService.VerifyPIN(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.IsNewUser)
}

func TestAuthService_GenerateOTP(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	otp := authService.generateOTP()
	
	assert.Len(t, otp, 6)
	assert.Regexp(t, `^\d{6}$`, otp)
}

func TestAuthService_GenerateHashedID(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	phoneNumber := "+254700000000"
	hashedID := authService.generateHashedID(phoneNumber)
	
	assert.Len(t, hashedID, 18) // 10 chars from hash + 8 chars from salt
	assert.NotEqual(t, phoneNumber, hashedID)
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	userID := uuid.New()
	token, err := authService.generateToken(userID)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	userID := uuid.New()
	token, _ := authService.generateToken(userID)
	
	validatedUserID, err := authService.ValidateToken(token)
	
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	mockDB := new(MockDB)
	authService := NewAuthService(mockDB)

	invalidToken := "invalid-token"
	
	_, err := authService.ValidateToken(invalidToken)
	
	assert.Error(t, err)
} 