package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Midtrans MidtransConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type AppConfig struct {
	Name     string
	Version  string
	LogLevel string
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

type MidtransConfig struct {
	ServerKey   string
	ClientKey   string
	Environment string
}

type JWTConfig struct {
	Secret     string
	ExpiryHour int
}

type StorageConfig struct {
	SupabaseURL       string
	SupabaseKey       string
	BucketName        string
	MaxFileSizeMB     int
}

func Load() (*Config, error) {
	config := &Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "QRIS POS Backend"),
			Version:  getEnv("APP_VERSION", "1.0.0"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			Database:     getEnv("DB_NAME", "qris_pos"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),
		},
		Midtrans: MidtransConfig{
			ServerKey:   getEnv("MIDTRANS_SERVER_KEY", ""),
			ClientKey:   getEnv("MIDTRANS_CLIENT_KEY", ""),
			Environment: getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			ExpiryHour: getEnvInt("JWT_EXPIRY_HOUR", 24),
		},
		Storage: StorageConfig{
			SupabaseURL:       getEnv("SUPABASE_URL", ""),
			SupabaseKey:       getEnv("SUPABASE_ANON_KEY", ""),
			BucketName:        getEnv("SUPABASE_BUCKET_NAME", "product-images"),
			MaxFileSizeMB:     getEnvInt("MAX_FILE_SIZE_MB", 2),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}