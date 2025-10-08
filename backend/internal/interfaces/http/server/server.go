package server

import (
	"fmt"
	"net/http"

	"qris-pos-backend/internal/infrastructure/config"
	"qris-pos-backend/internal/infrastructure/database/repositories"
	infraPayment "qris-pos-backend/internal/infrastructure/payment"
	"qris-pos-backend/internal/infrastructure/qrcode"
	"qris-pos-backend/internal/infrastructure/storage"
	"qris-pos-backend/internal/interfaces/http/handlers"
	"qris-pos-backend/internal/interfaces/middleware"
	"qris-pos-backend/internal/usecases/auth"
	usecasePayment "qris-pos-backend/internal/usecases/payment"
	"qris-pos-backend/internal/usecases/product"
	"qris-pos-backend/internal/usecases/transaction"
	pkgAuth "qris-pos-backend/pkg/auth"
	"qris-pos-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger logger.Logger
	router *gin.Engine
}

func NewServer(cfg *config.Config, db *gorm.DB, logger logger.Logger) *Server {
	server := &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	// Set Gin mode based on config
	if s.config.App.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Initialize services
	passwordService := pkgAuth.NewPasswordService()
	jwtService := pkgAuth.NewJWTService(s.config.JWT.Secret, s.config.JWT.ExpiryHour)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Initialize storage client
	storageClient := storage.NewSupabaseClient(s.config.Storage, s.logger)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(s.db)
	productRepo := repositories.NewProductRepository(s.db)
	categoryRepo := repositories.NewCategoryRepository(s.db)
	transactionRepo := repositories.NewTransactionRepository(s.db)
	paymentRepo := repositories.NewPaymentRepository(s.db)

	// Initialize infrastructure services
	midtransClient := infraPayment.NewMidtransClient(s.config.Midtrans)
	qrCodeGenerator := qrcode.NewQRCodeGenerator()

	// Initialize use cases
	authUseCase := auth.NewAuthUseCase(userRepo, passwordService, jwtService, s.logger)
	productUseCase := product.NewProductUseCase(productRepo, categoryRepo, s.logger)
	transactionUseCase := transaction.NewTransactionUseCase(transactionRepo, productRepo, userRepo, s.logger)
	paymentUseCase := usecasePayment.NewPaymentUseCase(paymentRepo, transactionRepo, midtransClient, qrCodeGenerator, s.logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, s.logger)
	productHandler := handlers.NewProductHandler(productUseCase, s.logger)
	transactionHandler := handlers.NewTransactionHandler(transactionUseCase, s.logger)
	paymentHandler := handlers.NewPaymentHandler(paymentUseCase, s.logger)
	imageHandler := handlers.NewImageHandler(storageClient, s.config.Storage, s.logger)

	// Health check endpoint

	// API routes
	api := router.Group("/api/v1")
	api.GET("/health", s.healthCheck)

	{
		// Auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/register", authMiddleware.RequireAdmin(), authHandler.Register)
		}

		// Auth routes (protected)
		authProtected := api.Group("/auth")
		authProtected.Use(authMiddleware.RequireAuth())
		{
			authProtected.GET("/me", authHandler.GetProfile)
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.POST("/change-password", authHandler.ChangePassword)
			authProtected.PUT("/profile", authHandler.UpdateProfile)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.ListProducts)   // Public - can view products
			products.GET("/:id", productHandler.GetProduct) // Public - can view single product
		}

		// Product routes (Admin only)
		productsAdmin := api.Group("/products")
		productsAdmin.Use(authMiddleware.RequireAdmin())
		{
			productsAdmin.POST("", productHandler.CreateProduct)
			productsAdmin.PUT("/:id", productHandler.UpdateProduct)
			productsAdmin.DELETE("/:id", productHandler.DeleteProduct)
			productsAdmin.PATCH("/:id/stock", productHandler.UpdateStock)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", productHandler.ListCategories) // Public
		}

		// Category routes (Admin only)
		categoriesAdmin := api.Group("/categories")
		categoriesAdmin.Use(authMiddleware.RequireAdmin())
		{
			categoriesAdmin.POST("", productHandler.CreateCategory)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		transactions.Use(authMiddleware.RequireAdminOrCashier())
		{
			transactions.GET("", transactionHandler.ListTransactions)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id/cancel", transactionHandler.CancelTransaction)
			transactions.POST("/:id/items", transactionHandler.AddItemToTransaction)
			transactions.DELETE("/:id/items/:item_id", transactionHandler.RemoveItemFromTransaction)
			transactions.PUT("/:id/items/:item_id", transactionHandler.UpdateItemQuantity)
		}

		// QRIS routes (Phase 2 implementation)
		qris := api.Group("/qris")
		qris.Use(authMiddleware.RequireAdminOrCashier())
		{
			qris.POST("/generate", paymentHandler.GenerateQRIS)
			qris.GET("/:transaction_id/status", paymentHandler.GetPaymentStatus)
			qris.POST("/:transaction_id/refresh", paymentHandler.RefreshQRIS)
		}

		// Payment routes (Phase 2 implementation)
		payments := api.Group("/payments")
		{
			payments.POST("/callback", paymentHandler.PaymentCallback) // Public - webhook from Midtrans
			payments.GET("/:transaction_id/status", authMiddleware.RequireAdminOrCashier(), paymentHandler.GetPaymentStatus)
		}

		// Image routes (Admin only)
		images := api.Group("/images")
		images.Use(authMiddleware.RequireAdmin())
		{
			images.POST("/upload", imageHandler.UploadImage)
			images.DELETE("/delete", imageHandler.DeleteImage)
		}
	}

	s.router = router
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": s.config.App.Name,
		"version": s.config.App.Version,
	})
}

func (s *Server) ListenAndServe() error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	return s.router.Run(address)
}

func (s *Server) Shutdown(ctx interface{}) error {
	// Gin doesn't have built-in graceful shutdown, but we can implement it if needed
	return nil
}

// Placeholder handlers - will be implemented later
func (s *Server) login(c *gin.Context)  { c.JSON(200, gin.H{"message": "login endpoint"}) }
func (s *Server) logout(c *gin.Context) { c.JSON(200, gin.H{"message": "logout endpoint"}) }
func (s *Server) getCurrentUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get current user endpoint"})
}
func (s *Server) getProducts(c *gin.Context) { c.JSON(200, gin.H{"message": "get products endpoint"}) }
func (s *Server) createProduct(c *gin.Context) {
	c.JSON(200, gin.H{"message": "create product endpoint"})
}
func (s *Server) getProduct(c *gin.Context) { c.JSON(200, gin.H{"message": "get product endpoint"}) }
func (s *Server) updateProduct(c *gin.Context) {
	c.JSON(200, gin.H{"message": "update product endpoint"})
}
func (s *Server) deleteProduct(c *gin.Context) {
	c.JSON(200, gin.H{"message": "delete product endpoint"})
}
func (s *Server) generateQRIS(c *gin.Context) {
	c.JSON(200, gin.H{"message": "generate qris endpoint"})
}
func (s *Server) getQRISStatus(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get qris status endpoint"})
}
func (s *Server) refreshQRIS(c *gin.Context) { c.JSON(200, gin.H{"message": "refresh qris endpoint"}) }
func (s *Server) paymentCallback(c *gin.Context) {
	c.JSON(200, gin.H{"message": "payment callback endpoint"})
}
func (s *Server) getPaymentStatus(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get payment status endpoint"})
}
