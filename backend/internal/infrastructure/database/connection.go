package database

import (
	"fmt"
	"time"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg)),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func getLogLevel(cfg config.DatabaseConfig) logger.LogLevel {
	// You can extend this to read from config if needed
	return logger.Info
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Category{},
		&entities.Product{},
		&entities.Transaction{},
		&entities.TransactionItem{},
		&entities.Payment{},
		&entities.QRISCode{},
	)
}

func SeedData(db *gorm.DB) error {
	// Create default categories
	categories := []entities.Category{
		{Name: "Food & Beverages"},
		{Name: "Electronics"},
		{Name: "Clothing"},
		{Name: "Books"},
		{Name: "Health & Beauty"},
		{Name: "Others"},
	}

	for _, category := range categories {
		var existingCategory entities.Category
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return fmt.Errorf("failed to create category %s: %w", category.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check existing category %s: %w", category.Name, err)
			}
		}
	}

	// Create default admin user
	var adminUser entities.User
	if err := db.Where("email = ?", "admin@qrispos.com").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Hash password properly
			hashedPassword := "$2a$12$N9qo8uLOickgx2ZMRZoMye7VFnjZcHZqHVveAJk0/R7/OWlKDENrW" // admin123
			admin := entities.NewUser("admin@qrispos.com", "System Admin", hashedPassword, entities.RoleAdmin)
			if err := db.Create(admin).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing admin user: %w", err)
		}
	}

	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
