package config

import (
	"log"
	"os"
	"sync"
	"runtime"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	BackendPort  		string
	BackendUrl			string
	FrontendPort		string
	FrontendUrl			string
	ProjectPath			string
	NodeStorage			string
	NodeNamePattern		string
	LndNetwork			string
	LndDomain			string
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
			BackendPort:			getEnv("BACKEND_PORT", "8080"),
			BackendUrl:				getEnv("BACKEND_URL", "http://localhost"),
			FrontendPort:			getEnv("FRONTEND_PORT", "8081"),
			FrontendUrl:			getEnv("FRONTEND_URL", "http://localhost"),
			ProjectPath:			rootDir(),
			NodeStorage:			getEnv("NODE_STORAGE", "/app/node-storage/"),
			NodeNamePattern:		getEnv("NODE_NAME_PATTERN", "lnd-*"),
			LndNetwork:				getEnv("LND_NETWORK", "simnet"),
			LndDomain:				getEnv("LND_DOMAIN", "lnd-headless.lightning-playground.svc.cluster.local"),
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

// Return the projectDir
func rootDir() string {
    _, file, _, ok := runtime.Caller(0)
    if !ok {
        panic("Impossible de récupérer le chemin")
    }
    return filepath.Dir(file) + "/"
}
