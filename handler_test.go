package main

import (
	"os"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			expectError: false,
		},
		{
			name: "missing storage account name",
			config: &Config{
				StorageAccoutName: "",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			expectError: true,
			errorMsg:    "STORAGE_ACCOUNT_NAME is required",
		},
		{
			name: "missing storage account key",
			config: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "",
				ContainerName:     "testcontainer",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			expectError: true,
			errorMsg:    "STORAGE_ACCOUNT_KEY is required",
		},
		{
			name: "missing container name",
			config: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			expectError: true,
			errorMsg:    "CONTAINER_NAME is required",
		},
		{
			name: "missing API key",
			config: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "",
				BaseURL:           "https://test.com",
			},
			expectError: true,
			errorMsg:    "API_KEY is required",
		},
		{
			name: "missing base URL",
			config: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "testapikey",
				BaseURL:           "",
			},
			expectError: true,
			errorMsg:    "BASE_URL is required",
		},
		{
			name: "all fields missing",
			config: &Config{
				StorageAccoutName: "",
				StorageAccountKey: "",
				ContainerName:     "",
				APIKey:            "",
				BaseURL:           "",
			},
			expectError: true,
			errorMsg:    "STORAGE_ACCOUNT_NAME is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	originalVars := map[string]string{
		"STORAGE_ACCOUNT_NAME": os.Getenv("STORAGE_ACCOUNT_NAME"),
		"STORAGE_ACCOUNT_KEY":  os.Getenv("STORAGE_ACCOUNT_KEY"),
		"CONTAINER_NAME":       os.Getenv("CONTAINER_NAME"),
		"API_KEY":              os.Getenv("API_KEY"),
		"BASE_URL":             os.Getenv("BASE_URL"),
	}

	cleanup := func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	tests := []struct {
		name        string
		envVars     map[string]string
		expected    *Config
		description string
	}{
		{
			name: "all environment variables set",
			envVars: map[string]string{
				"STORAGE_ACCOUNT_NAME": "testaccount",
				"STORAGE_ACCOUNT_KEY":  "testkey",
				"CONTAINER_NAME":       "testcontainer",
				"API_KEY":              "testapikey",
				"BASE_URL":             "https://test.com",
			},
			expected: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			description: "should load all values from environment variables",
		},
		{
			name: "container name uses default when not set",
			envVars: map[string]string{
				"STORAGE_ACCOUNT_NAME": "testaccount",
				"STORAGE_ACCOUNT_KEY":  "testkey",
				"CONTAINER_NAME":       "",
				"API_KEY":              "testapikey",
				"BASE_URL":             "https://test.com",
			},
			expected: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "$web",
				APIKey:            "testapikey",
				BaseURL:           "https://test.com",
			},
			description: "should use default container name '$web' when CONTAINER_NAME is empty",
		},
		{
			name: "API key uses default when not set",
			envVars: map[string]string{
				"STORAGE_ACCOUNT_NAME": "testaccount",
				"STORAGE_ACCOUNT_KEY":  "testkey",
				"CONTAINER_NAME":       "testcontainer",
				"API_KEY":              "",
				"BASE_URL":             "https://test.com",
			},
			expected: &Config{
				StorageAccoutName: "testaccount",
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
				APIKey:            "",
				BaseURL:           "https://test.com",
			},
			description: "should use empty string as default for API_KEY when not set",
		},
		{
			name: "no environment variables set",
			envVars: map[string]string{
				"STORAGE_ACCOUNT_NAME": "",
				"STORAGE_ACCOUNT_KEY":  "",
				"CONTAINER_NAME":       "",
				"API_KEY":              "",
				"BASE_URL":             "",
			},
			expected: &Config{
				StorageAccoutName: "",
				StorageAccountKey: "",
				ContainerName:     "$web",
				APIKey:            "",
				BaseURL:           "",
			},
			description: "should use defaults when no environment variables are set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key := range tt.envVars {
				os.Unsetenv(key)
			}

			for key, value := range tt.envVars {
				if value != "" {
					os.Setenv(key, value)
				}
			}

			config := loadConfig()

			if config.StorageAccoutName != tt.expected.StorageAccoutName {
				t.Errorf("StorageAccoutName: expected '%s', got '%s'", tt.expected.StorageAccoutName, config.StorageAccoutName)
			}
			if config.StorageAccountKey != tt.expected.StorageAccountKey {
				t.Errorf("StorageAccountKey: expected '%s', got '%s'", tt.expected.StorageAccountKey, config.StorageAccountKey)
			}
			if config.ContainerName != tt.expected.ContainerName {
				t.Errorf("ContainerName: expected '%s', got '%s'", tt.expected.ContainerName, config.ContainerName)
			}
			if config.APIKey != tt.expected.APIKey {
				t.Errorf("APIKey: expected '%s', got '%s'", tt.expected.APIKey, config.APIKey)
			}
			if config.BaseURL != tt.expected.BaseURL {
				t.Errorf("BaseURL: expected '%s', got '%s'", tt.expected.BaseURL, config.BaseURL)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	originalValue := os.Getenv("TEST_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_VAR")
		} else {
			os.Setenv("TEST_VAR", originalValue)
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			setEnv:       true,
			expected:     "custom",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			setEnv:       false,
			expected:     "default",
		},
		{
			name:         "environment variable set to empty string",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			expected:     "default",
		},
		{
			name:         "empty default value",
			key:          "TEST_VAR",
			defaultValue: "",
			envValue:     "",
			setEnv:       false,
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("TEST_VAR")

			if tt.setEnv {
				os.Setenv("TEST_VAR", tt.envValue)
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
