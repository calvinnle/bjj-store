package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/calvinnle/bjj-store/backend/services"
	"github.com/gin-gonic/gin"
)

// UploadImage godoc
// @Summary Upload product image
// @Description Upload an image file for a product
// @Tags admin,upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]interface{} "Upload successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 413 {object} map[string]interface{} "File too large"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/upload/image [post]
func UploadImage(c *gin.Context) {
	// Debug: Log request details
	fmt.Printf("Upload request received from: %s\n", c.ClientIP())
	fmt.Printf("Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("Authorization header present: %v\n", c.GetHeader("Authorization") != "")

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Printf("Error getting form file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No image file provided",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("File received: %s, Size: %d bytes\n", file.Filename, file.Size)

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"success": false,
			"error":   "File size too large. Maximum size is 5MB",
		})
		return
	}

	// Validate file type by extension
	ext := strings.ToLower(getFileExtension(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValidType := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isValidType = true
			break
		}
	}

	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file type. Only JPEG, PNG, GIF, and WebP images are allowed",
		})
		return
	}

	// Upload to MinIO
	imageURL, err := services.MinIOClient.UploadImage(file)
	if err != nil {
		fmt.Printf("MinIO upload error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload image",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("Image uploaded successfully to MinIO: %s\n", imageURL)

	// Return the file URL
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"image_url": imageURL,
		"message":   "Image uploaded successfully",
	})
}

// DeleteImage godoc
// @Summary Delete product image
// @Description Delete an uploaded product image
// @Tags admin,upload
// @Accept json
// @Produce json
// @Param request body map[string]string true "Image URL to delete"
// @Success 200 {object} map[string]interface{} "Delete successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/upload/image [delete]
func DeleteImage(c *gin.Context) {
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Delete from MinIO
	if err := services.MinIOClient.DeleteImage(req.ImageURL); err != nil {
		fmt.Printf("MinIO delete error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete image",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("Image deleted successfully from MinIO: %s\n", req.ImageURL)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// getFileExtension extracts file extension from filename
func getFileExtension(filename string) string {
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 {
		return ""
	}
	return filename[lastDot:]
}