package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "UNSET_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "returns empty string when both are empty",
			key:          "EMPTY_KEY",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			os.Unsetenv(tt.key)

			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		filename   string
		content    string
		expectLoad bool
		expectVars map[string]string
	}{
		{
			name:       "loads valid env file",
			filename:   "test.env",
			content:    "TEST_VAR1=value1\nTEST_VAR2=value2\n",
			expectLoad: true,
			expectVars: map[string]string{
				"TEST_VAR1": "value1",
				"TEST_VAR2": "value2",
			},
		},
		{
			name:       "skips comments and empty lines",
			filename:   "test_comments.env",
			content:    "# This is a comment\nTEST_VAR3=value3\n\n# Another comment\nTEST_VAR4=value4\n",
			expectLoad: true,
			expectVars: map[string]string{
				"TEST_VAR3": "value3",
				"TEST_VAR4": "value4",
			},
		},
		{
			name:       "returns false for non-existent file",
			filename:   "nonexistent.env",
			content:    "",
			expectLoad: false,
			expectVars: map[string]string{},
		},
		{
			name:       "blocks path traversal attempts",
			filename:   "../../etc/passwd",
			content:    "",
			expectLoad: false,
			expectVars: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variables
			for key := range tt.expectVars {
				os.Unsetenv(key)
			}

			var filePath string
			if tt.content != "" && !strings.Contains(tt.filename, "..") {
				// Create test file
				filePath = filepath.Join(tempDir, tt.filename)
				err := os.WriteFile(filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			} else {
				filePath = tt.filename
			}

			result := loadEnvFile(filePath)
			if result != tt.expectLoad {
				t.Errorf("loadEnvFile() = %v, want %v", result, tt.expectLoad)
			}

			// Check if expected environment variables were set
			for key, expectedValue := range tt.expectVars {
				actualValue := os.Getenv(key)
				if actualValue != expectedValue {
					t.Errorf("Environment variable %s = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"JWT_SECRET", "ADMIN_USERNAME", "ADMIN_PASSWORD", "SERVER_PORT",
		"APP_ENV", "ENABLE_CORS", "CORS_ORIGINS", "ENABLE_HTTPS",
		"CERT_FILE", "KEY_FILE",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("loads config with required environment variables", func(t *testing.T) {
		// Set required environment variables
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("ADMIN_PASSWORD", "test-password")

		config := LoadConfig()

		if config.JWTSecret != "test-secret" {
			t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, "test-secret")
		}

		if config.AdminPass != "test-password" {
			t.Errorf("AdminPass = %v, want %v", config.AdminPass, "test-password")
		}

		// Check defaults
		if config.AdminUser != "admin" {
			t.Errorf("AdminUser = %v, want %v", config.AdminUser, "admin")
		}

		if config.ServerPort != "8080" {
			t.Errorf("ServerPort = %v, want %v", config.ServerPort, "8080")
		}

		if config.Environment != "development" {
			t.Errorf("Environment = %v, want %v", config.Environment, "development")
		}
	})

	t.Run("parses CORS origins correctly", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("ADMIN_PASSWORD", "test-password")
		os.Setenv("CORS_ORIGINS", "http://localhost:3000, https://example.com, https://app.example.com")

		config := LoadConfig()

		expectedOrigins := []string{"http://localhost:3000", "https://example.com", "https://app.example.com"}
		if len(config.CORSOrigins) != len(expectedOrigins) {
			t.Errorf("CORSOrigins length = %v, want %v", len(config.CORSOrigins), len(expectedOrigins))
		}

		for i, expected := range expectedOrigins {
			if i >= len(config.CORSOrigins) || config.CORSOrigins[i] != expected {
				t.Errorf("CORSOrigins[%d] = %v, want %v", i, config.CORSOrigins[i], expected)
			}
		}
	})

	t.Run("handles boolean environment variables", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("ADMIN_PASSWORD", "test-password")
		os.Setenv("ENABLE_CORS", "false")
		os.Setenv("ENABLE_HTTPS", "true")

		config := LoadConfig()

		if config.EnableCORS != false {
			t.Errorf("EnableCORS = %v, want %v", config.EnableCORS, false)
		}

		if config.EnableHTTPS != true {
			t.Errorf("EnableHTTPS = %v, want %v", config.EnableHTTPS, true)
		}
	})
}

func TestGetRequiredEnv(t *testing.T) {
	t.Run("returns value when environment variable is set", func(t *testing.T) {
		key := "TEST_REQUIRED_VAR"
		expectedValue := "required_value"

		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		result := getRequiredEnv(key)
		if result != expectedValue {
			t.Errorf("getRequiredEnv() = %v, want %v", result, expectedValue)
		}
	})

}

func TestLoadEnvFileFromPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	file1 := filepath.Join(tempDir, "test1.env")
	file2 := filepath.Join(tempDir, "test2.env")

	os.WriteFile(file2, []byte("TEST_PATH_VAR=found"), 0644)

	// Clean up environment
	os.Unsetenv("TEST_PATH_VAR")

	paths := []string{file1, file2} // file1 doesn't exist, file2 does
	loadEnvFileFromPaths("test.env", paths)

	// Should have loaded from file2
	if os.Getenv("TEST_PATH_VAR") != "found" {
		t.Errorf("Expected TEST_PATH_VAR to be 'found', got '%s'", os.Getenv("TEST_PATH_VAR"))
	}
}
