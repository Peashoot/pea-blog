package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Frontend    FrontendConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type FrontendConfig struct {
	AutoBuild    bool
	BuildCommand string
	DistPath     string
	SourcePath   string
}

func Load() *Config {
	// Load .env file
	godotenv.Load()

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://user:password@localhost/pea_blog?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	environment := getEnv("ENVIRONMENT", "development")

	// Frontend configuration
	autoBuild, _ := strconv.ParseBool(getEnv("FRONTEND_AUTO_BUILD", "true"))
	buildCommand := getEnv("FRONTEND_BUILD_COMMAND", "npm run build")
	distPath := getEnv("FRONTEND_DIST_PATH", "./frontend/dist")
	sourcePath := getEnv("FRONTEND_SOURCE_PATH", "./frontend")

	return &Config{
		Environment: environment,
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			URL: dbURL,
		},
		JWT: JWTConfig{
			Secret:      jwtSecret,
			ExpireHours: jwtExpireHours,
		},
		Frontend: FrontendConfig{
			AutoBuild:    autoBuild,
			BuildCommand: buildCommand,
			DistPath:     distPath,
			SourcePath:   sourcePath,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
