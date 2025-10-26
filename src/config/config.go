package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	JWTSecret   string
	AdminUser   string
	AdminPass   string
	ServerPort  string
	Environment string
	EnableCORS  bool
	CORSOrigins []string
	EnableHTTPS bool
	CertFile    string
	KeyFile     string
}

// LoadConfig loads configuration from environment variables and .env files
func LoadConfig() *Config {
	// Get environment first (can be set via ENV var or default to development)
	env := getEnv("APP_ENV", "development")

	log.Printf("Loading configuration for environment: %s", env)

	// Load environment-specific file first (highest priority after ENV vars)
	envFile := fmt.Sprintf(".env.%s", env)
	// Try multiple possible paths
	loadEnvFileFromPaths(envFile, []string{
		fmt.Sprintf("environments/.env.%s", env),    // From project root
		fmt.Sprintf("../environments/.env.%s", env), // From src/ directory
	})

	// Load general .env file as fallback (lower priority)
	loadEnvFileFromPaths(".env", []string{
		"environments/.env",    // From project root
		"../environments/.env", // From src/ directory
	})

	// Parse CORS origins from environment variable (comma-separated)
	corsOriginsStr := getEnv("CORS_ORIGINS", "")
	var corsOrigins []string
	if corsOriginsStr != "" {
		corsOrigins = strings.Split(corsOriginsStr, ",")
		for i, origin := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(origin)
		}
	}

	config := &Config{
		JWTSecret:   getRequiredEnv("JWT_SECRET"),
		AdminUser:   getEnv("ADMIN_USERNAME", "admin"),
		AdminPass:   getRequiredEnv("ADMIN_PASSWORD"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: env,
		EnableCORS:  getEnv("ENABLE_CORS", "true") == "true",
		CORSOrigins: corsOrigins,
		EnableHTTPS: getEnv("ENABLE_HTTPS", "false") == "true",
		CertFile:    getEnv("CERT_FILE", "ssl/cert.pem"),
		KeyFile:     getEnv("KEY_FILE", "ssl/key.pem"),
	}

	// Log configuration source (without sensitive data)
	log.Printf("Configuration loaded - Environment: %s, Port: %s, Admin User: %s",
		config.Environment, config.ServerPort, config.AdminUser)

	return config
}

// getRequiredEnv gets environment variable and fails if not set
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// loadEnvFileFromPaths tries to load environment file from multiple possible paths
func loadEnvFileFromPaths(filename string, paths []string) {
	for _, path := range paths {
		if loadEnvFile(path) {
			return // Successfully loaded, stop trying other paths
		}
	}
	log.Printf("Environment file not found in any location: %s", filename)
}

// loadEnvFile loads environment variables from a file and returns true if successful
func loadEnvFile(filename string) bool {
	// Clean the file path to prevent directory traversal
	cleanPath := filepath.Clean(filename)

	// Only block paths that try to escape the working directory root
	// Allow relative paths like ../environments/.env but block things like ../../etc/passwd
	if strings.Contains(cleanPath, "../..") {
		log.Printf("Invalid file path detected (too many parent directories): %s", filename)
		return false
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		// File doesn't exist, return false to try next path
		return false
	}
	defer file.Close()

	log.Printf("✓ Loading environment variables from: %s", filename)

	loadedCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set in environment (ENV vars have highest priority)
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				log.Printf("Failed to set environment variable %s: %v", key, err)
				continue
			}
			loadedCount++
		} else {
			log.Printf("  - %s: using environment variable (overriding file)", key)
		}
	}

	if loadedCount > 0 {
		log.Printf("  → Loaded %d variables from %s", loadedCount, filename)
	}

	return true // Successfully loaded
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
