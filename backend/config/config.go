package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Stripe   StripeConfig   `mapstructure:"stripe"`
	Admin    AdminConfig    `mapstructure:"admin"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
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

type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BucketName      string `mapstructure:"bucket_name"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	Region          string `mapstructure:"region"`
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
	
	// Debug: Log environment variables
	log.Printf("DEBUG - DATABASE_HOST env: %s", viper.GetString("database.host"))
	log.Printf("DEBUG - DATABASE_PORT env: %s", viper.GetString("database.port"))
	log.Printf("DEBUG - DATABASE_USER env: %s", viper.GetString("database.user"))

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

	// Override with Railway environment variables if they exist
	overrideWithEnvVars(&config)
	
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

	// MinIO defaults
	viper.SetDefault("minio.endpoint", "localhost:9000")
	viper.SetDefault("minio.access_key_id", "admin")
	viper.SetDefault("minio.secret_access_key", "admin123456")
	viper.SetDefault("minio.bucket_name", "bjj-store-images")
	viper.SetDefault("minio.use_ssl", false)
	viper.SetDefault("minio.region", "us-east-1")
}

// overrideWithEnvVars directly reads Railway environment variables
func overrideWithEnvVars(config *Config) {
	log.Printf("Checking for Railway environment variables...")
	
	// Database overrides
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		config.Database.Host = host
		log.Printf("Override DATABASE_HOST: %s", host)
	}
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		config.Database.Port = port
		log.Printf("Override DATABASE_PORT: %s", port)
	}
	if user := os.Getenv("DATABASE_USER"); user != "" {
		config.Database.User = user
		log.Printf("Override DATABASE_USER: %s", user)
	}
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		config.Database.Password = password
		log.Printf("Override DATABASE_PASSWORD: [HIDDEN]")
	}
	if name := os.Getenv("DATABASE_NAME"); name != "" {
		config.Database.Name = name
		log.Printf("Override DATABASE_NAME: %s", name)
	}
	if sslMode := os.Getenv("DATABASE_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
		log.Printf("Override DATABASE_SSL_MODE: %s", sslMode)
	}
	
	// Server overrides
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
		log.Printf("Override SERVER_PORT: %s", port)
	}
	if env := os.Getenv("SERVER_ENVIRONMENT"); env != "" {
		config.Server.Environment = env
		log.Printf("Override SERVER_ENVIRONMENT: %s", env)
	}
	
	// JWT override
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
		log.Printf("Override JWT_SECRET: [HIDDEN]")
	}
	
	// Admin overrides
	if email := os.Getenv("ADMIN_DEFAULT_EMAIL"); email != "" {
		config.Admin.DefaultEmail = email
		log.Printf("Override ADMIN_DEFAULT_EMAIL: %s", email)
	}
	if password := os.Getenv("ADMIN_DEFAULT_PASSWORD"); password != "" {
		config.Admin.DefaultPassword = password
		log.Printf("Override ADMIN_DEFAULT_PASSWORD: [HIDDEN]")
	}
}
