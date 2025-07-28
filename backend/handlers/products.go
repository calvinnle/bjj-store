package handlers

import (
	"net/http"
	"strconv"

	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products with optional filtering
// @Tags products
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param search query string false "Search in product name and description"
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Get query parameters
	category := c.Query("category")
	search := c.Query("search")

	query := models.DB

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{} "Invalid product ID"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags admin,products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags admin,products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID (Admin only)
// @Tags admin,products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid product ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := models.DB.Delete(&models.Product{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
