package handlers

import (
	"net/http"

	"qris-pos-backend/internal/infrastructure/config"
	"qris-pos-backend/internal/infrastructure/storage"
	"qris-pos-backend/pkg/logger"
	"qris-pos-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	storageClient *storage.SupabaseClient
	config        config.StorageConfig
	logger        logger.Logger
}

func NewImageHandler(storageClient *storage.SupabaseClient, config config.StorageConfig, logger logger.Logger) *ImageHandler {
	return &ImageHandler{
		storageClient: storageClient,
		config:        config,
		logger:        logger,
	}
}

type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

// UploadImage godoc
// @Summary Upload product image
// @Description Upload an image for a product (Admin only)
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "Image file (JPEG, PNG, WebP, GIF, max 2MB)"
// @Success 200 {object} response.Response{data=UploadImageResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 413 {object} response.Response
// @Router /images/upload [post]
func (h *ImageHandler) UploadImage(c *gin.Context) {
	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get uploaded file", "error", err)
		response.BadRequest(c, "No file provided or invalid file", err.Error())
		return
	}
	defer file.Close()

	// Validate file
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect content type from filename
		if ext := getFileExtension(header.Filename); ext != "" {
			contentType = getContentTypeFromExtension(ext)
		}
	}

	if err := storage.ValidateImageFile(contentType, header.Size, h.config.MaxFileSizeMB); err != nil {
		h.logger.Warn("Invalid file upload attempt", "error", err, "filename", header.Filename, "size", header.Size, "content_type", contentType)
		if header.Size > int64(h.config.MaxFileSizeMB)*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"message": "File too large",
				"error":   err.Error(),
			})
		} else {
			response.BadRequest(c, "Invalid file", err.Error())
		}
		return
	}

	// Upload to Supabase
	imageURL, err := h.storageClient.UploadImage(file, header.Filename, contentType)
	if err != nil {
		h.logger.Error("Failed to upload image to Supabase", "error", err, "filename", header.Filename)
		response.InternalError(c, "Failed to upload image", err.Error())
		return
	}

	response.Success(c, "Image uploaded successfully", UploadImageResponse{
		ImageURL: imageURL,
		Message:  "Image uploaded successfully",
	})
}

// DeleteImage godoc
// @Summary Delete product image
// @Description Delete an image from storage (Admin only)
// @Tags images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body map[string]string true "Image URL to delete"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /images/delete [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	// Delete from Supabase
	if err := h.storageClient.DeleteImage(req.ImageURL); err != nil {
		h.logger.Error("Failed to delete image from Supabase", "error", err, "image_url", req.ImageURL)
		response.InternalError(c, "Failed to delete image", err.Error())
		return
	}

	response.Success(c, "Image deleted successfully", nil)
}

func getFileExtension(fileName string) string {
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ""
}

func getContentTypeFromExtension(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
