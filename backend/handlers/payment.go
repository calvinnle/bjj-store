package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	OrderID    uint    `json:"order_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	CardNumber string  `json:"card_number" binding:"required"`
	ExpiryDate string  `json:"expiry_date" binding:"required"`
	CVV        string  `json:"cvv" binding:"required"`
	NameOnCard string  `json:"name_on_card" binding:"required"`
}

type PaymentResponse struct {
	Success       bool         `json:"success"`
	TransactionID string       `json:"transaction_id,omitempty"`
	Message       string       `json:"message"`
	Order         models.Order `json:"order,omitempty"`
}

// ProcessPayment godoc
// @Summary Process payment for an order
// @Description Process payment for an order using mock payment gateway
// @Tags payment
// @Accept json
// @Produce json
// @Param payment body PaymentRequest true "Payment data"
// @Success 200 {object} PaymentResponse "Payment successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Payment failed"
// @Router /payment/process [post]
func ProcessPayment(c *gin.Context) {
	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid payment data",
			"details": err.Error(),
		})
		return
	}

	// Get the order
	var order models.Order
	if err := models.DB.Preload("Items.Product").First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Order not found",
		})
		return
	}

	// Validate order amount matches payment amount
	if order.TotalAmount != req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Amount mismatch. Expected: %.2f, Received: %.2f", order.TotalAmount, req.Amount),
		})
		return
	}

	// Check if order is already paid
	if order.Status == models.OrderStatusPaid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Order is already paid",
		})
		return
	}

	// Mock payment processing - simulate delay
	time.Sleep(2 * time.Second)

	// Mock different scenarios based on card number
	cardNumber := strings.ReplaceAll(req.CardNumber, " ", "")
	
	// Test card for declined payment
	if cardNumber == "4000000000000002" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Payment declined by bank",
			"message": "Your card was declined. Please try a different payment method.",
		})
		return
	}

	// Test card for processing error
	if cardNumber == "4000000000000119" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Payment processing error",
			"message": "There was an error processing your payment. Please try again.",
		})
		return
	}

	// Successful payment - begin transaction to update order and reduce stock
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update order status
	order.Status = models.OrderStatusPaid
	order.UpdatedAt = time.Now()

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update order status",
		})
		return
	}

	// Now reduce stock for each item since payment was successful
	for _, item := range order.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Product with ID %d not found during stock reduction", item.ProductID),
			})
			return
		}

		// Double-check stock availability (race condition protection)
		if !product.HasSufficientStock(item.Quantity) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Insufficient stock for %s. Available: %d, Requested: %d", 
					product.Name, product.Stock, item.Quantity),
			})
			return
		}

		// Reduce stock
		if err := product.ReduceStock(item.Quantity); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Save updated product stock
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update product stock",
			})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to commit payment transaction",
		})
		return
	}

	// Generate transaction ID
	transactionID := fmt.Sprintf("TXN-%d-%d", time.Now().Unix(), order.ID)

	// Log successful payment
	fmt.Printf("Payment processed successfully - Order: %s, Amount: $%.2f, Transaction: %s\n", 
		order.OrderNumber, order.TotalAmount, transactionID)

	c.JSON(http.StatusOK, PaymentResponse{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Payment processed successfully",
		Order:         order,
	})
}