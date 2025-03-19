package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMODE      string
	JWTSecret      string
	RedisURL       string
	RedisPassword  string
	InfluxDBURL    string
	InfluxDBToken  string
	InfluxDBOrg    string
	InfluxDBBucket string
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GetConfig() *Config {
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "myproject"),
		DBSSLMODE:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		InfluxDBURL:    getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBToken:  getEnv("INFLUXDB_TOKEN", "your_token"),
		InfluxDBOrg:    getEnv("INFLUXDB_ORG", "your_org"),
		InfluxDBBucket: getEnv("INFLUXDB_BUCKET", "your_bucket"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
