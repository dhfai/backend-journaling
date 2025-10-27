package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	MongoDB  MongoDBConfig
	JWT      JWTConfig
	OTP      OTPConfig
	SMTP     SMTPConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	PrivateKeyPath       string
	PublicKeyPath        string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type OTPConfig struct {
	Pepper      string
	TTL         time.Duration
	MaxAttempts int
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

type ServerConfig struct {
	Port        string
	Host        string
	Environment string
}

func Load() (*Config, error) {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	accessDuration, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_DURATION: %w", err)
	}

	refreshDuration, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_DURATION", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_DURATION: %w", err)
	}

	otpTTL, _ := strconv.Atoi(getEnv("OTP_TTL_MINUTES", "5"))
	otpMaxAttempts, _ := strconv.Atoi(getEnv("OTP_MAX_ATTEMPTS", "5"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "journaling_auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "journaling"),
		},
		JWT: JWTConfig{
			PrivateKeyPath:       getEnv("JWT_PRIVATE_KEY_PATH", "./keys/jwt_private.pem"),
			PublicKeyPath:        getEnv("JWT_PUBLIC_KEY_PATH", "./keys/jwt_public.pem"),
			AccessTokenDuration:  accessDuration,
			RefreshTokenDuration: refreshDuration,
		},
		OTP: OTPConfig{
			Pepper:      getEnv("OTP_PEPPER", "default-pepper-change-me"),
			TTL:         time.Duration(otpTTL) * time.Minute,
			MaxAttempts: otpMaxAttempts,
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:      smtpPort,
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
			FromName:  getEnv("SMTP_FROM_NAME", "Journaling App"),
		},
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
