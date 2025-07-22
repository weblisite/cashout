package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// NewConnection creates a new database connection
func NewConnection() (*sql.DB, error) {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "cashout")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection established successfully")

	// Initialize database schema
	if err := initializeSchema(db); err != nil {
		log.Error().Err(err).Msg("Failed to initialize database schema")
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// initializeSchema creates the database tables if they don't exist
func initializeSchema(db *sql.DB) error {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone_number VARCHAR(20) UNIQUE NOT NULL,
			hashed_id VARCHAR(50) UNIQUE NOT NULL,
			kyc_status VARCHAR(20) DEFAULT 'not_started',
			wallet_balance DECIMAL(15,2) DEFAULT 0.00,
			user_status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// User PINs table
		`CREATE TABLE IF NOT EXISTS user_pins (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			hashed_pin VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// User biometrics table
		`CREATE TABLE IF NOT EXISTS user_biometrics (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			biometric_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Transactions table
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id),
			recipient_id UUID REFERENCES users(id),
			business_id UUID REFERENCES businesses(id),
			agent_id UUID REFERENCES agents(id),
			amount DECIMAL(15,2) NOT NULL,
			fee DECIMAL(15,2) DEFAULT 0.00,
			agent_commission DECIMAL(15,2),
			platform_margin DECIMAL(15,2),
			type VARCHAR(20) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			note TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Agents table
		`CREATE TABLE IF NOT EXISTS agents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID UNIQUE NOT NULL REFERENCES users(id),
			agent_code VARCHAR(20) UNIQUE NOT NULL,
			float_balance DECIMAL(15,2) DEFAULT 0.00,
			commission_rate DECIMAL(5,4) DEFAULT 0.0200,
			status VARCHAR(20) DEFAULT 'active',
			location TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Businesses table
		`CREATE TABLE IF NOT EXISTS businesses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID UNIQUE NOT NULL REFERENCES users(id),
			business_name VARCHAR(100) NOT NULL,
			business_type VARCHAR(50),
			address TEXT,
			phone_number VARCHAR(20),
			email VARCHAR(100),
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// OTP table for storing OTP codes
		`CREATE TABLE IF NOT EXISTS otp_codes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone_number VARCHAR(20) NOT NULL,
			otp_code VARCHAR(6) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create indexes
		`CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number)`,
		`CREATE INDEX IF NOT EXISTS idx_users_hashed_id ON users(hashed_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_user_id ON agents(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_agent_code ON agents(agent_code)`,
		`CREATE INDEX IF NOT EXISTS idx_businesses_user_id ON businesses(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_codes_phone_number ON otp_codes(phone_number)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_codes_expires_at ON otp_codes(expires_at)`,
	}

	// Execute each query
	for i, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %d: %w", i+1, err)
		}
	}

	log.Info().Msg("Database schema initialized successfully")
	return nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseConnection closes the database connection
func CloseConnection(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		} else {
			log.Info().Msg("Database connection closed successfully")
		}
	}
} 