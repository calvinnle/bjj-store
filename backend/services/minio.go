package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/calvinnle/bjj-store/backend/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService struct {
	client     *minio.Client
	bucketName string
}

var MinIOClient *MinIOService

// InitMinIO initializes the MinIO client
func InitMinIO() error {
	cfg := config.AppConfig.MinIO

	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %v", err)
	}

	MinIOClient = &MinIOService{
		client:     client,
		bucketName: cfg.BucketName,
	}

	// Create bucket if it doesn't exist
	if err := MinIOClient.createBucket(); err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	log.Printf("MinIO client initialized successfully")
	log.Printf("Bucket: %s", cfg.BucketName)
	log.Printf("Endpoint: %s", cfg.Endpoint)

	return nil
}

// createBucket creates the bucket if it doesn't exist
func (m *MinIOService) createBucket() error {
	ctx := context.Background()

	// Check if bucket exists
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	if !exists {
		// Create bucket
		err = m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{
			Region: config.AppConfig.MinIO.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("Bucket '%s' created successfully", m.bucketName)

		// Set bucket policy to allow public read access for images
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, m.bucketName)

		err = m.client.SetBucketPolicy(ctx, m.bucketName, policy)
		if err != nil {
			log.Printf("Warning: Failed to set bucket policy: %v", err)
		} else {
			log.Printf("Bucket policy set successfully for public read access")
		}
	} else {
		log.Printf("Bucket '%s' already exists", m.bucketName)
	}

	return nil
}

// UploadImage uploads an image file to MinIO
func (m *MinIOService) UploadImage(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("products/%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)

	// Get content type
	contentType := getContentType(ext)

	// Upload to MinIO
	ctx := context.Background()
	_, err = m.client.PutObject(ctx, m.bucketName, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %v", err)
	}

	// Generate public URL
	imageURL := m.getPublicURL(filename)
	
	log.Printf("Image uploaded successfully: %s", filename)
	return imageURL, nil
}

// DeleteImage deletes an image from MinIO
func (m *MinIOService) DeleteImage(imageURL string) error {
	// Extract object name from URL
	objectName := m.extractObjectName(imageURL)
	if objectName == "" {
		return fmt.Errorf("invalid image URL")
	}

	ctx := context.Background()
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %v", err)
	}

	log.Printf("Image deleted successfully: %s", objectName)
	return nil
}

// getPublicURL generates a public URL for the uploaded image
func (m *MinIOService) getPublicURL(objectName string) string {
	cfg := config.AppConfig.MinIO
	protocol := "http"
	if cfg.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.Endpoint, m.bucketName, objectName)
}

// extractObjectName extracts the object name from a MinIO URL
func (m *MinIOService) extractObjectName(imageURL string) string {
	// Expected format: http://localhost:9000/bjj-store-images/products/123456_abc.jpg
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return ""
	}
	
	// Find bucket name in URL and get everything after it
	bucketIndex := -1
	for i, part := range parts {
		if part == m.bucketName {
			bucketIndex = i
			break
		}
	}
	
	if bucketIndex == -1 || bucketIndex+1 >= len(parts) {
		return ""
	}
	
	// Join remaining parts as object name
	return strings.Join(parts[bucketIndex+1:], "/")
}

// getContentType returns the appropriate content type for the file extension
func getContentType(ext string) string {
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

// generateRandomString generates a random string for unique filenames
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// GetImageStream returns an image stream for serving (if needed for custom serving)
func (m *MinIOService) GetImageStream(objectName string) (io.Reader, error) {
	ctx := context.Background()
	object, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %v", err)
	}
	return object, nil
}