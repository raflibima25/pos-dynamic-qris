package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	
	// Create folder structure: products/{uuid}.ext
	objectPath := fmt.Sprintf("products/%s", uniqueFileName)

	// Prepare multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file part
	part, err := writer.CreateFormFile("file", uniqueFileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	writer.Close()

	// Create HTTP request
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucketName, objectPath)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return "", fmt.Errorf("supabase error: %s", errorResp.Message)
		}
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse success response
	var uploadResp UploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		s.logger.Warn("Failed to parse upload response, but upload seemed successful", "response", string(body))
	}

	// Generate public URL
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