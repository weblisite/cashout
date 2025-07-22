package handlers

import (
	"net/http"

	"github.com/cashout/backend/internal/middleware"
	"github.com/cashout/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SendOTP sends an OTP to the provided phone number
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req services.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind SendOTP request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide a valid phone number",
		})
		return
	}

	if err := h.authService.SendOTP(req); err != nil {
		log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Failed to send OTP")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send OTP",
			"message": "Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"phone":   req.PhoneNumber,
	})
}

// VerifyOTP verifies the OTP and returns user information
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req services.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind VerifyOTP request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid phone number and OTP",
		})
		return
	}

	response, err := h.authService.VerifyOTP(req)
	if err != nil {
		log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Failed to verify OTP")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid OTP",
			"message": "Please check your OTP and try again",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SetupPIN sets up PIN for a user
func (h *AuthHandler) SetupPIN(c *gin.Context) {
	var req services.SetupPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind SetupPIN request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid phone number and PIN",
		})
		return
	}

	if err := h.authService.SetupPIN(req); err != nil {
		log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Failed to setup PIN")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to setup PIN",
			"message": "Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "PIN setup successfully",
	})
}

// VerifyPIN verifies the PIN for a user
func (h *AuthHandler) VerifyPIN(c *gin.Context) {
	var req services.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind VerifyPIN request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid phone number and PIN",
		})
		return
	}

	response, err := h.authService.VerifyPIN(req)
	if err != nil {
		log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Failed to verify PIN")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid PIN",
			"message": "Please check your PIN and try again",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SetupBiometric sets up biometric authentication for a user
func (h *AuthHandler) SetupBiometric(c *gin.Context) {
	var req services.SetupBiometricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind SetupBiometric request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid phone number and biometric ID",
		})
		return
	}

	if err := h.authService.SetupBiometric(req); err != nil {
		log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Failed to setup biometric")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to setup biometric",
			"message": "Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Biometric setup successfully",
	})
}

// Logout logs out the user
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not authenticated",
			"message": "Please login to continue",
		})
		return
	}

	log.Info().Str("user_id", userID).Msg("User logged out")

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
} 