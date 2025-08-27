package server

import (
	"fmt"
	"net/http"

	"qris-pos-backend/internal/infrastructure/config"
	"qris-pos-backend/internal/infrastructure/database/repositories"
	"qris-pos-backend/internal/interfaces/http/handlers"
	"qris-pos-backend/internal/interfaces/middleware"
	"qris-pos-backend/internal/usecases/auth"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(s.db)
	productRepo := repositories.NewProductRepository(s.db)
	categoryRepo := repositories.NewCategoryRepository(s.db)
	transactionRepo := repositories.NewTransactionRepository(s.db)

	// Initialize use cases
	authUseCase := auth.NewAuthUseCase(userRepo, passwordService, jwtService, s.logger)
	productUseCase := product.NewProductUseCase(productRepo, categoryRepo, s.logger)
	transactionUseCase := transaction.NewTransactionUseCase(transactionRepo, productRepo, userRepo, s.logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, s.logger)
	productHandler := handlers.NewProductHandler(productUseCase, s.logger)
	transactionHandler := handlers.NewTransactionHandler(transactionUseCase, s.logger)

	// Health check endpoint
	router.GET("/health", s.healthCheck)

	// API routes
	api := router.Group("/api/v1")
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

		// QRIS routes (placeholder for Phase 2)
		qris := api.Group("/qris")
		qris.Use(authMiddleware.RequireAdminOrCashier())
		{
			qris.POST("/generate", s.generateQRIS)
			qris.GET("/:transaction_id/status", s.getQRISStatus)
			qris.POST("/refresh", s.refreshQRIS)
		}

		// Payment routes (placeholder for Phase 2)
		payments := api.Group("/payments")
		{
			payments.POST("/callback", s.paymentCallback) // Public - webhook from Midtrans
			payments.GET("/:id/status", authMiddleware.RequireAdminOrCashier(), s.getPaymentStatus)
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
