package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppPort          string
	RedisPort        string
	BucketKey        string
	BucketMaxTokens  int64
	BucketRefillRate float64
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		BucketKey:        getEnv("BUCKET_KEY", "KEY1"),
		BucketMaxTokens:  getIntEnv("BUCKET_MAX_TOKENS", 10),
		BucketRefillRate: getFloatEnv("BUCKET_REFILL_RATE", 1.0),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, strconv.FormatInt(defaultValue, 10))
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Warning: Could not parse %s=%s as integer, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func getFloatEnv(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, strconv.FormatFloat(defaultValue, 'f', -1, 64))
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Printf("Warning: Could not parse %s=%s as float, using default: %f", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
