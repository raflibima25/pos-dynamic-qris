package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"qris-pos-backend/internal/infrastructure/config"
	"qris-pos-backend/pkg/logger"

	"github.com/google/uuid"
)

type SupabaseClient struct {
	baseURL    string
	apiKey     string
	bucketName string
	httpClient *http.Client
	logger     logger.Logger
}

type UploadResponse struct {
	Key      string `json:"Key"`
	Id       string `json:"Id"`
	FullPath string `json:"fullPath"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewSupabaseClient(cfg config.StorageConfig, logger logger.Logger) *SupabaseClient {
	return &SupabaseClient{
		baseURL:    cfg.SupabaseURL,
		apiKey:     cfg.SupabaseKey,
		bucketName: cfg.BucketName,
		logger:     logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SupabaseClient) UploadImage(file io.Reader, fileName string, contentType string) (string, error) {
	// Generate UUID filename
	fileExtension := getFileExtension(fileName)
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	objectPath := fmt.Sprintf("products/%s", uniqueFileName)

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucketName, objectPath)

	// Baca file ke buffer
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 🔑 set content type sesuai mimetype (image/jpeg, image/png, dll)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "false")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	publicURL := s.GetPublicURL(objectPath)
	s.logger.Info("Image uploaded successfully", "path", objectPath, "url", publicURL)
	return publicURL, nil
}

func (s *SupabaseClient) GetPublicURL(objectPath string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.baseURL, s.bucketName, objectPath)
}

func (s *SupabaseClient) DeleteImage(objectPath string) error {
	// Extract path from full URL if needed
	if len(objectPath) > len(s.baseURL) && objectPath[:len(s.baseURL)] == s.baseURL {
		// This is a full URL, extract the object path
		prefix := fmt.Sprintf("%s/storage/v1/object/public/%s/", s.baseURL, s.bucketName)
		if len(objectPath) > len(prefix) && objectPath[:len(prefix)] == prefix {
			objectPath = objectPath[len(prefix):]
		}
	}

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucketName, objectPath)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Info("Image deleted successfully", "path", objectPath)
	return nil
}

func getFileExtension(fileName string) string {
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ""
}

// ValidateImageFile validates if the uploaded file is a valid image
func ValidateImageFile(contentType string, size int64, maxSizeMB int) error {
	// Check content type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}

	if !allowedTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s. Allowed types: JPEG, PNG, WebP, GIF", contentType)
	}

	// Check file size
	maxSize := int64(maxSizeMB) * 1024 * 1024 // Convert MB to bytes
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d MB", size, maxSizeMB)
	}

	return nil
}
