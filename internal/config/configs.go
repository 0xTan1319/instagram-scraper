package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDB       string
	UserAgent     string
	RequestTimeout time.Duration
	ProxyFile     string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	timeout, _ := strconv.Atoi(getEnv("REQUEST_TIMEOUT", "10"))

	AppConfig = &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGO_DB", "instagram_scraper"),
		UserAgent:     getEnv("USER_AGENT", "Mozilla/5.0 ..."),
		RequestTimeout: time.Duration(timeout) * time.Second,
		ProxyFile:     getEnv("PROXY_FILE", "proxies.txt"),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
