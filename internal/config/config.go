package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Email    EmailConfig
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret string
}

// EmailConfig represents email API configuration
type EmailConfig struct {
	APIURL string
	APIKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "bioskop_db")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_ENV", "development")
	viper.SetDefault("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")
	viper.SetDefault("EMAIL_API_URL", "https://lumoshive-academy-email-api.vercel.app/send-email")
	viper.SetDefault("EMAIL_API_KEY", "")

	// Read .env file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using default values and environment variables")
		} else {
			log.Fatalf("Error reading config: %v", err)
		}
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  "disable",
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Env:  viper.GetString("SERVER_ENV"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Email: EmailConfig{
			APIURL: viper.GetString("EMAIL_API_URL"),
			APIKey: viper.GetString("EMAIL_API_KEY"),
		},
	}
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}
