package middleware

import (
	"net/http"
	"strings"

	"github.com/cashout/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthRequired middleware checks if the request has a valid JWT token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Please provide a valid authentication token",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Get auth service from context (injected in main.go)
		authServiceInterface, exists := c.Get("authService")
		if !exists {
			log.Error().Msg("AuthService not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Authentication service not available",
			})
			c.Abort()
			return
		}

		authService, ok := authServiceInterface.(*services.AuthService)
		if !ok {
			log.Error().Msg("Invalid AuthService type in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Authentication service not available",
			})
			c.Abort()
			return
		}

		// Validate the token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid token provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Please provide a valid authentication token",
			})
			c.Abort()
			return
		}

		// Store user ID in context for handlers to use
		c.Set("userID", userID)
		c.Next()
	}
}

// OptionalAuth middleware checks for JWT token but doesn't require it
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Get auth service from context
		authServiceInterface, exists := c.Get("authService")
		if !exists {
			// Auth service not available, continue without authentication
			c.Next()
			return
		}

		authService, ok := authServiceInterface.(*services.AuthService)
		if !ok {
			// Invalid auth service type, continue without authentication
			c.Next()
			return
		}

		// Validate the token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Store user ID in context for handlers to use
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserIDFromContext extracts the user ID from the gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		return "", false
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return "", false
	}

	return userID, true
} 