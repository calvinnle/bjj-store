package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Stripe   StripeConfig   `mapstructure:"stripe"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

type AdminConfig struct {
	DefaultEmail    string `mapstructure:"default_email"`
	DefaultPassword string `mapstructure:"default_password"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type StripeConfig struct {
	SecretKey     string `mapstructure:"secret_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

var AppConfig *Config

func LoadConfig() {
	// Set config file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Also try to read .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Set defaults first (fallback values)
	setDefaults()

	// Enable automatic environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read config.yaml first
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config.yaml found: %v", err)

		// Try to read .env file as backup
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("No .env file found either: %v", err)
			log.Println("Using environment variables and defaults only")
		} else {
			log.Println("Configuration loaded from .env file")
		}
	} else {
		log.Println("Configuration loaded from config.yaml")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	AppConfig = &config
	log.Printf("Configuration loaded successfully")
	log.Printf("Database: %s:%s/%s", config.Database.Host, config.Database.Port, config.Database.Name)
	log.Printf("Server: %s (env: %s)", config.Server.Port, config.Server.Environment)
	log.Printf("Admin: %s", config.Admin.DefaultEmail)
}

func setDefaults() {
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "bjj_store")
	viper.SetDefault("database.ssl_mode", "disable")

	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.environment", "development")

	// JWT defaults
	viper.SetDefault("jwt.secret", "change-this-in-production-bjj-store-2024")

	// Admin defaults
	viper.SetDefault("admin.default_email", "admin@bjjstore.com")
	viper.SetDefault("admin.default_password", "admin123")

	viper.SetDefault("stripe.secret_key", "")
	viper.SetDefault("stripe.webhook_secret", "")
}
