package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	GuestEmail      string            `json:"guest_email" binding:"required,email"`
	ShippingAddress models.Address    `json:"shipping_address" binding:"required"`
	Items           []models.CartItem `json:"items" binding:"required,min=1"`
}

type OrderResponse struct {
	Order   models.Order `json:"order"`
	Message string       `json:"message"`
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with guest checkout
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order data"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid order request",
			"details": err.Error(),
		})
		return
	}

	// Validate shipping address
	if req.ShippingAddress.FirstName == "" || req.ShippingAddress.LastName == "" ||
		req.ShippingAddress.Address1 == "" || req.ShippingAddress.City == "" ||
		req.ShippingAddress.ZipCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incomplete shipping address",
		})
		return
	}

	// Begin database transaction
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order := models.Order{
		OrderNumber:     generateOrderNumber(),
		GuestEmail:      req.GuestEmail,
		ShippingAddress: req.ShippingAddress,
		Status:          models.OrderStatusPending,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	// Process each cart item
	for _, item := range req.Items {
		// Get product details
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Product with ID %d not found", item.ProductID),
			})
			return
		}

		// Check stock availability (but don't reduce yet - only on successful payment)
		if !product.HasSufficientStock(item.Quantity) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insufficient stock for %s. Available: %d, Requested: %d",
					product.Name, product.Stock, item.Quantity),
			})
			return
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Product:   product,
			Quantity:  item.Quantity,
			Price:     product.Price, // Use current product price
			Size:      item.Size,
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)
	}

	// Save all order items
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order items",
		})
		return
	}

	// Update order total
	order.TotalAmount = totalAmount
	order.Items = orderItems

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order total",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit order",
		})
		return
	}

	c.JSON(http.StatusCreated, OrderResponse{
		Order:   order,
		Message: "Order created successfully",
	})
}

// TrackOrder godoc
// @Summary Track an order
// @Description Track an order by order number
// @Tags orders
// @Accept json
// @Produce json
// @Param orderNumber path string true "Order Number"
// @Success 200 {object} map[string]interface{} "Order details"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /orders/track/{orderNumber} [get]
func TrackOrder(c *gin.Context) {
	orderNumber := c.Param("orderNumber")

	var order models.Order
	if err := models.DB.Preload("Items.Product").Where("order_number = ?", orderNumber).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":        "Order not found",
			"order_number": orderNumber,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"order":   order,
	})
}

// GetOrdersByEmail godoc
// @Summary Get orders by email
// @Description Get all orders for a specific email address
// @Tags orders
// @Accept json
// @Produce json
// @Param email path string true "Customer Email"
// @Success 200 {object} map[string]interface{} "Orders list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/email/{email} [get]
func GetOrdersByEmail(c *gin.Context) {
	email := c.Param("email")

	var orders []models.Order
	if err := models.DB.Preload("Items.Product").Where("guest_email = ?", email).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"email":   email,
		"orders":  orders,
		"total":   len(orders),
	})
}

// GetAllOrders godoc
// @Summary Get all orders (Admin only)
// @Description Get all orders with pagination
// @Tags admin,orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "Orders list with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/orders [get]
func GetAllOrders(c *gin.Context) {
	var orders []models.Order

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	if err := models.DB.Preload("Items.Product").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch orders",
		})
		return
	}

	// Get total count
	var total int64
	models.DB.Model(&models.Order{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"orders":  orders,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status (Admin only)
// @Description Update the status of an order
// @Tags admin,orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body map[string]string true "New status"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid status",
			"details": err.Error(),
		})
		return
	}

	// Validate status
	validStatuses := []models.OrderStatus{
		models.OrderStatusPending,
		models.OrderStatusPaid,
		models.OrderStatusShipped,
		models.OrderStatusDelivered,
		models.OrderStatusCancelled,
	}

	valid := false
	for _, status := range validStatuses {
		if req.Status == status {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Invalid order status",
			"valid_statuses": validStatuses,
		})
		return
	}

	var order models.Order
	if err := models.DB.First(&order, uint(orderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	order.Status = req.Status
	order.UpdatedAt = time.Now()

	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order status updated successfully",
		"order":   order,
	})
}

func generateOrderNumber() string {
	return fmt.Sprintf("BJJ-%d", time.Now().Unix())
}
