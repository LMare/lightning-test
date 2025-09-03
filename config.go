package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	BackendPort  string
	BackendUrl   string
	FrontendPort string
	FrontendUrl  string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aucun fichier .env trouvé")
		}

		cfg = &Config{
			BackendPort:  getEnv("BACKEND_PORT", "8080"),
			BackendUrl:   getEnv("BACKEND_URL", "http://localhost"),
			FrontendPort: getEnv("FRONTEND_PORT", "8081"),
			FrontendUrl:  getEnv("FRONTEND_URL", "http://localhost"),
		}
	})

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
