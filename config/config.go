package config

import "os"

type Config struct {
	MongoURI string
	DBName   string
	Port     string
}

func Load() *Config {
	return &Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://admin:admin@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000"),
		DBName:   getEnv("DB_NAME", "kanban_board"),
		Port:     getEnv("PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
