package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	NATS     NATSConfig     `mapstructure:"nats"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	MoMo     MoMoConfig     `mapstructure:"momo"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type NATSConfig struct {
	URL string `mapstructure:"url"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

type MoMoConfig struct {
	APIKey      string `mapstructure:"api_key"`
	APISecret   string `mapstructure:"api_secret"`
	Environment string `mapstructure:"environment"`
	CallbackURL string `mapstructure:"callback_url"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	setDefaults()

	// Allow environment variables to override
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	// Override with environment variables
	overrideWithEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	return &config
}

func setDefaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "dongome")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "dongome_db")
	viper.SetDefault("database.ssl_mode", "disable")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("nats.url", "nats://localhost:4222")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration", 24) // 24 hours

	viper.SetDefault("momo.environment", "sandbox")
}

func overrideWithEnv() {
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		viper.Set("server.host", host)
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		viper.Set("database.port", dbPort)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database.user", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.name", dbName)
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		viper.Set("redis.host", redisHost)
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		viper.Set("redis.port", redisPort)
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		viper.Set("redis.password", redisPassword)
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		if db, err := strconv.Atoi(redisDB); err == nil {
			viper.Set("redis.db", db)
		}
	}
	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		viper.Set("nats.url", natsURL)
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("jwt.secret", jwtSecret)
	}
	if momoAPIKey := os.Getenv("MOMO_API_KEY"); momoAPIKey != "" {
		viper.Set("momo.api_key", momoAPIKey)
	}
	if momoAPISecret := os.Getenv("MOMO_API_SECRET"); momoAPISecret != "" {
		viper.Set("momo.api_secret", momoAPISecret)
	}
}
