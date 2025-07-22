package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cashout/backend/configs"
	"github.com/cashout/backend/internal/handlers"
	"github.com/cashout/backend/internal/middleware"
	"github.com/cashout/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found, using environment variables")
	}

	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Supabase connection
	if err := configs.InitializeSupabase(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Supabase")
	}
	defer configs.CloseSupabase()

	// Initialize Supabase service
	supabaseService := services.NewSupabaseService()

	// Initialize Core Banking service
	coreBankingService := services.NewCoreBankingService(supabaseService)

	// Initialize Phase 2: Settlement, Liquidity, Risk Management
	settlementService := services.NewSettlementService(supabaseService, coreBankingService)
	liquidityService := services.NewLiquidityService(supabaseService, coreBankingService)
	riskManagementService := services.NewRiskManagementService(supabaseService, coreBankingService)

	// Initialize Phase 3: KYC/AML, Enhanced Fraud Detection
	kycAMLService := services.NewKYCAMLService(supabaseService, riskManagementService)
	enhancedFraudDetectionService := services.NewEnhancedFraudDetectionService(supabaseService, riskManagementService, kycAMLService)

	// Initialize Phase 4: High Availability, Analytics
	highAvailabilityService := services.NewHighAvailabilityService(supabaseService, coreBankingService)
	analyticsService := services.NewAnalyticsService(supabaseService, coreBankingService)

	// Initialize existing services
	authService := services.NewAuthService(supabaseService)
	userService := services.NewUserService(supabaseService)
	transactionService := services.NewTransactionService(supabaseService)
	feeService := services.NewFeeService(supabaseService)

	// Initialize handlers with all services
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, feeService)
	
	// Initialize new handlers for Phase 2, 3, 4 services
	settlementHandler := handlers.NewSettlementHandler(settlementService)
	liquidityHandler := handlers.NewLiquidityHandler(liquidityService)
	riskHandler := handlers.NewRiskHandler(riskManagementService)
	kycHandler := handlers.NewKYCHandler(kycAMLService)
	fraudHandler := handlers.NewFraudHandler(enhancedFraudDetectionService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	healthHandler := handlers.NewHealthHandler(highAvailabilityService)

	// Initialize middleware
	authMiddleware := middleware.AuthRequired(authService)
	optionalAuthMiddleware := middleware.OptionalAuth(authService)
	loggerMiddleware := middleware.Logger()
	requestIDMiddleware := middleware.RequestID()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware
	router.Use(loggerMiddleware)
	router.Use(requestIDMiddleware)
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080", "https://yourdomain.com"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authMiddleware, authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// User routes
		users := api.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/wallet", userHandler.GetWallet)
			users.POST("/upload-document", userHandler.UploadDocument)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		transactions.Use(authMiddleware)
		{
			transactions.POST("/transfer", transactionHandler.Transfer)
			transactions.POST("/cash-in", transactionHandler.CashIn)
			transactions.POST("/cash-out", transactionHandler.CashOut)
			transactions.GET("/history", transactionHandler.GetHistory)
			transactions.GET("/:id", transactionHandler.GetTransaction)
		}

		// QR Code routes
		qr := api.Group("/qr")
		qr.Use(authMiddleware)
		{
			qr.POST("/generate", transactionHandler.GenerateQR)
			qr.POST("/scan", transactionHandler.ScanQR)
		}

		// Phase 2: Settlement routes
		settlements := api.Group("/settlements")
		settlements.Use(authMiddleware)
		{
			settlements.GET("/", settlementHandler.GetSettlements)
			settlements.POST("/", settlementHandler.CreateSettlement)
			settlements.GET("/:id", settlementHandler.GetSettlement)
			settlements.POST("/:id/process", settlementHandler.ProcessSettlement)
			settlements.GET("/pending", settlementHandler.GetPendingSettlements)
		}

		// Float management routes
		floats := api.Group("/floats")
		floats.Use(authMiddleware)
		{
			floats.GET("/", settlementHandler.GetFloats)
			floats.POST("/", settlementHandler.CreateFloat)
			floats.PUT("/:id", settlementHandler.UpdateFloat)
			floats.POST("/:id/replenish", settlementHandler.ReplenishFloat)
			floats.GET("/low", settlementHandler.GetLowFloatAgents)
		}

		// Phase 2: Liquidity routes
		liquidity := api.Group("/liquidity")
		liquidity.Use(authMiddleware)
		{
			liquidity.GET("/snapshot", liquidityHandler.GetLiquiditySnapshot)
			liquidity.POST("/snapshot", liquidityHandler.CreateLiquiditySnapshot)
			liquidity.GET("/cash-flow", liquidityHandler.GetCashFlowReport)
			liquidity.GET("/alerts", liquidityHandler.GetActiveAlerts)
			liquidity.POST("/alerts/:id/resolve", liquidityHandler.ResolveAlert)
		}

		// Reserve management routes
		reserves := api.Group("/reserves")
		reserves.Use(authMiddleware)
		{
			reserves.GET("/", liquidityHandler.GetReserves)
			reserves.POST("/", liquidityHandler.CreateReserve)
			reserves.PUT("/:id", liquidityHandler.UpdateReserve)
		}

		// Phase 2: Risk Management routes
		risk := api.Group("/risk")
		risk.Use(authMiddleware)
		{
			risk.POST("/assess", riskHandler.AssessTransactionRisk)
			risk.GET("/alerts", riskHandler.GetActiveFraudAlerts)
			risk.POST("/alerts/:id/resolve", riskHandler.ResolveFraudAlert)
			risk.GET("/limits", riskHandler.GetTransactionLimits)
			risk.POST("/limits", riskHandler.CreateTransactionLimit)
		}

		// Phase 3: KYC/AML routes
		kyc := api.Group("/kyc")
		kyc.Use(authMiddleware)
		{
			kyc.POST("/profile", kycHandler.CreateKYCProfile)
			kyc.GET("/profile", kycHandler.GetKYCProfile)
			kyc.PUT("/profile", kycHandler.UpdateKYCProfile)
			kyc.POST("/documents", kycHandler.AddDocument)
			kyc.PUT("/documents/:id/verify", kycHandler.VerifyDocument)
			kyc.POST("/employment", kycHandler.UpdateEmploymentInfo)
			kyc.POST("/source-of-funds", kycHandler.UpdateSourceOfFunds)
		}

		// KYC Review routes (admin only)
		kycReview := api.Group("/kyc/review")
		kycReview.Use(authMiddleware)
		{
			kycReview.GET("/pending", kycHandler.GetPendingKYCProfiles)
			kycReview.POST("/:id/review", kycHandler.ReviewKYCProfile)
		}

		// AML routes
		aml := api.Group("/aml")
		aml.Use(authMiddleware)
		{
			aml.POST("/check", kycHandler.PerformAMLCheck)
			aml.GET("/alerts", kycHandler.GetActiveAMLAlerts)
			aml.POST("/alerts/:id/resolve", kycHandler.ResolveAMLAlert)
		}

		// Phase 3: Enhanced Fraud Detection routes
		fraud := api.Group("/fraud")
		fraud.Use(authMiddleware)
		{
			fraud.POST("/analyze", fraudHandler.AnalyzeTransaction)
			fraud.GET("/alerts", fraudHandler.GetActiveFraudAlerts)
			fraud.POST("/alerts/:id/resolve", fraudHandler.ResolveFraudAlert)
		}

		// Phase 4: Analytics routes
		analytics := api.Group("/analytics")
		analytics.Use(authMiddleware)
		{
			analytics.POST("/reports", analyticsHandler.GenerateBusinessReport)
			analytics.GET("/reports/:id", analyticsHandler.GetAnalyticsReport)
			analytics.GET("/metrics", analyticsHandler.GetBusinessMetrics)
			analytics.GET("/users", analyticsHandler.GetUserAnalytics)
			analytics.GET("/agents", analyticsHandler.GetAgentAnalytics)
		}

		// Phase 4: High Availability routes
		ha := api.Group("/ha")
		ha.Use(authMiddleware)
		{
			ha.GET("/health", healthHandler.GetServiceHealth)
			ha.GET("/health/all", healthHandler.GetAllServiceHealth)
			ha.GET("/load-balancer", healthHandler.GetLoadBalancerStatus)
			ha.POST("/backup", healthHandler.PerformBackup)
			ha.POST("/recovery", healthHandler.PerformRecovery)
		}
	}

	// Start health monitoring
	ctx := context.Background()
	highAvailabilityService.StartHealthMonitoring(ctx)

	// Start server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("🚀 Cashout Platform starting on port %s", port)
		log.Info().Msg("✅ Core Banking Infrastructure: ACTIVE")
		log.Info().Msg("✅ Settlement Engine: ACTIVE")
		log.Info().Msg("✅ Liquidity Management: ACTIVE")
		log.Info().Msg("✅ Risk Management: ACTIVE")
		log.Info().Msg("✅ KYC/AML System: ACTIVE")
		log.Info().Msg("✅ Enhanced Fraud Detection: ACTIVE")
		log.Info().Msg("✅ High Availability: ACTIVE")
		log.Info().Msg("✅ Analytics & Reporting: ACTIVE")
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("🛑 Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("✅ Server exited gracefully")
} 