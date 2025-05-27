package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type UploadResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

type Config struct {
	StorageAccoutName string
	StorageAccountKey string
	ContainerName     string
	APIKey            string
	BaseURL           string
}

func loadConfig() *Config {
	return &Config{
		StorageAccoutName: os.Getenv("STORAGE_ACCOUNT_NAME"),
		StorageAccountKey: os.Getenv("STORAGE_ACCOUNT_KEY"),
		ContainerName:     getEnvOrDefault("CONTAINER_NAME", "$web"),
		APIKey:            getEnvOrDefault("API_KEY", ""),
		BaseURL:           os.Getenv("BASE_URL"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	id := uuid.New().String()
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("%s/%s%s", timestamp, id, ext)
}

func upload(w http.ResponseWriter, r *http.Request) {
	config := loadConfig()
	fmt.Fprint(w, config.ContainerName)
}

func main() {
	log.Println("Starting server...")
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc(("/api/upload"), func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			log.Println("Received POST request")
			upload(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
