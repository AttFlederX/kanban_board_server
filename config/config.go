package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	DBName    string
	Port      string
	JWTSecret string
}

func Load() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	return &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://admin:admin@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000"),
		DBName:    getEnv("DB_NAME", "kanban_board"),
		Port:      getEnv("PORT", "3000"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
