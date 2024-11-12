package config

import (
	"os"
)

type Config struct {
	MongoURI string
	Port     string
	Language string
	I18NPath string
	// Add other configuration fields as needed
}

var AppConfig *Config

func LoadConfig() *Config {
	AppConfig = &Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Port:     getEnv("PORT", "8080"),
		Language: getEnv("LANGUAGE", "en"),
		I18NPath: getEnv("I18N_PATH", "../../locales"),
	}
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
