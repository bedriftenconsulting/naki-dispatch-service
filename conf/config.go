package conf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	RedisAddr   string
	RedisPass   string
	KafkaBroker string
	JWTSecret   string
}

var AppConfig *Config

func Load() error {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "8087"))
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}

	AppConfig = &Config{
		Port:        strconv.Itoa(port),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "naki_dispatch"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "redis:6379"),
		RedisPass:   getEnv("REDIS_PASSWORD", ""),
		KafkaBroker: getEnv("KAFKA_BROKER", "kafka:9092"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}

	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
