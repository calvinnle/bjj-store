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
	log.Printf("=== ENVIRONMENT VARIABLE DEBUG ===")
	
	// Print ALL environment variables to see what Railway is providing
	log.Printf("All environment variables:")
	for _, env := range os.Environ() {
		log.Printf("  %s", env)
	}
	
	log.Printf("=== CHECKING SPECIFIC VARIABLES ===")
	
	// Check Railway-specific variables
	railwayEnv := os.Getenv("RAILWAY_ENVIRONMENT")
	log.Printf("RAILWAY_ENVIRONMENT: %s", railwayEnv)
	
	railwayService := os.Getenv("RAILWAY_SERVICE_NAME")
	log.Printf("RAILWAY_SERVICE_NAME: %s", railwayService)
	
	// Database overrides with detailed logging
	log.Printf("Checking database environment variables...")
	
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		log.Printf("FOUND DATABASE_HOST: %s", host)
		config.Database.Host = host
		log.Printf("✓ Override DATABASE_HOST: %s", host)
	} else {
		log.Printf("✗ DATABASE_HOST not found or empty")
	}
	
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		log.Printf("FOUND DATABASE_PORT: %s", port)
		config.Database.Port = port
		log.Printf("✓ Override DATABASE_PORT: %s", port)
	} else {
		log.Printf("✗ DATABASE_PORT not found or empty")
	}
	
	if user := os.Getenv("DATABASE_USER"); user != "" {
		log.Printf("FOUND DATABASE_USER: %s", user)
		config.Database.User = user
		log.Printf("✓ Override DATABASE_USER: %s", user)
	} else {
		log.Printf("✗ DATABASE_USER not found or empty")
	}
	
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		log.Printf("FOUND DATABASE_PASSWORD: [HIDDEN]")
		config.Database.Password = password
		log.Printf("✓ Override DATABASE_PASSWORD: [HIDDEN]")
	} else {
		log.Printf("✗ DATABASE_PASSWORD not found or empty")
	}
	
	if name := os.Getenv("DATABASE_NAME"); name != "" {
		log.Printf("FOUND DATABASE_NAME: %s", name)
		config.Database.Name = name
		log.Printf("✓ Override DATABASE_NAME: %s", name)
	} else {
		log.Printf("✗ DATABASE_NAME not found or empty")
	}
	
	if sslMode := os.Getenv("DATABASE_SSL_MODE"); sslMode != "" {
		log.Printf("FOUND DATABASE_SSL_MODE: %s", sslMode)
		config.Database.SSLMode = sslMode
		log.Printf("✓ Override DATABASE_SSL_MODE: %s", sslMode)
	} else {
		log.Printf("✗ DATABASE_SSL_MODE not found or empty")
	}
	
	// Try alternative Railway PostgreSQL environment variable names
	log.Printf("=== CHECKING RAILWAY POSTGRES VARIABLES ===")
	if dbUrl := os.Getenv("DATABASE_URL"); dbUrl != "" {
		log.Printf("FOUND DATABASE_URL: %s", dbUrl)
		// Don't parse it here, just log it exists
	} else {
		log.Printf("✗ DATABASE_URL not found")
	}
	
	// Check Railway Postgres reference variables (multiple possible names)
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = os.Getenv("PGHOST")
	}
	
	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		postgresPort = os.Getenv("PGPORT")
	}
	
	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		postgresUser = os.Getenv("PGUSER")
	}
	
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		postgresPassword = os.Getenv("PGPASSWORD")
	}
	
	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresDB == "" {
		postgresDB = os.Getenv("PGDATABASE")
	}
	
	log.Printf("Railway postgres reference variables:")
	log.Printf("  POSTGRES_HOST/PGHOST: %s", postgresHost)
	log.Printf("  POSTGRES_PORT/PGPORT: %s", postgresPort)
	log.Printf("  POSTGRES_USER/PGUSER: %s", postgresUser)
	if postgresPassword != "" {
		log.Printf("  POSTGRES_PASSWORD/PGPASSWORD: [HIDDEN]")
	} else {
		log.Printf("  POSTGRES_PASSWORD/PGPASSWORD: ")
	}
	log.Printf("  POSTGRES_DB/PGDATABASE: %s", postgresDB)
	
	// Use Railway Postgres variables if DATABASE_ variables are not set
	if config.Database.Host == "localhost" && postgresHost != "" {
		config.Database.Host = postgresHost
		log.Printf("✓ Using POSTGRES_HOST for database host: %s", postgresHost)
	}
	
	if config.Database.Port == "5432" && postgresPort != "" {
		config.Database.Port = postgresPort
		log.Printf("✓ Using POSTGRES_PORT for database port: %s", postgresPort)
	}
	
	if config.Database.User == "postgres" && postgresUser != "" {
		config.Database.User = postgresUser
		log.Printf("✓ Using POSTGRES_USER for database user: %s", postgresUser)
	}
	
	if config.Database.Password == "password" && postgresPassword != "" {
		config.Database.Password = postgresPassword
		log.Printf("✓ Using POSTGRES_PASSWORD for database password: [HIDDEN]")
	}
	
	if config.Database.Name == "bjj_store" && postgresDB != "" {
		config.Database.Name = postgresDB
		log.Printf("✓ Using POSTGRES_DB for database name: %s", postgresDB)
	}
	
	// Set SSL mode to require for Railway
	if postgresHost != "" {
		config.Database.SSLMode = "require"
		log.Printf("✓ Set SSL mode to 'require' for Railway Postgres")
	}
	
	// Server overrides
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
		log.Printf("✓ Override SERVER_PORT: %s", port)
	} else if port := os.Getenv("PORT"); port != "" {
		// Railway often uses PORT instead of SERVER_PORT
		config.Server.Port = port
		log.Printf("✓ Override SERVER_PORT from PORT: %s", port)
	}
	
	if env := os.Getenv("SERVER_ENVIRONMENT"); env != "" {
		config.Server.Environment = env
		log.Printf("✓ Override SERVER_ENVIRONMENT: %s", env)
	}
	
	// JWT override
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
		log.Printf("✓ Override JWT_SECRET: [HIDDEN]")
	}
	
	// Admin overrides
	if email := os.Getenv("ADMIN_DEFAULT_EMAIL"); email != "" {
		config.Admin.DefaultEmail = email
		log.Printf("✓ Override ADMIN_DEFAULT_EMAIL: %s", email)
	}
	if password := os.Getenv("ADMIN_DEFAULT_PASSWORD"); password != "" {
		config.Admin.DefaultPassword = password
		log.Printf("✓ Override ADMIN_DEFAULT_PASSWORD: [HIDDEN]")
	}
	
	log.Printf("=== FINAL CONFIG VALUES ===")
	log.Printf("Database Host: %s", config.Database.Host)
	log.Printf("Database Port: %s", config.Database.Port)
	log.Printf("Database User: %s", config.Database.User)
	log.Printf("Database Name: %s", config.Database.Name)
	log.Printf("Database SSL Mode: %s", config.Database.SSLMode)
	log.Printf("Server Port: %s", config.Server.Port)
	log.Printf("Server Environment: %s", config.Server.Environment)
}
