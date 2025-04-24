package controllers

import (
	"backend_prodman/db"
	"backend_prodman/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllStockLogs(c *gin.Context) {
	// Query untuk mengambil semua log stok
	rows, err := db.DB.Query("SELECT id, product_id, change_type, amount, note, created_at FROM stock_logs ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error fetching stock logs:", err) // Use the log package correctly
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stock logs"})
		return
	}
	defer rows.Close()

	// Menyimpan hasil query ke dalam slice
	var stockLogs []models.StockLog
	for rows.Next() {
		var stockLog models.StockLog // Rename to stockLog
		if err := rows.Scan(&stockLog.ID, &stockLog.ProductID, &stockLog.ChangeType, &stockLog.Amount, &stockLog.Note, &stockLog.CreatedAt); err != nil {
			log.Println("Error scanning row:", err) // Correct usage of log.Println
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process stock log"})
			return
		}
		stockLogs = append(stockLogs, stockLog)
	}

	// Cek jika tidak ada log ditemukan
	if len(stockLogs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No stock logs found"})
		return
	}

	// Mengembalikan data log stok
	c.JSON(http.StatusOK, stockLogs)
}

func UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var input struct {
		Type   string `json:"type"`   // "in" atau "out"
		Amount int    `json:"amount"` // Harus > 0
		Note   string `json:"note"`   // Optional
	}
	if err := c.ShouldBindJSON(&input); err != nil || (input.Type != "in" && input.Type != "out") || input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ambil stok sekarang
	var currentStock int
	err = db.DB.QueryRow("SELECT stock FROM products WHERE id=$1 AND status='active'", productID).Scan(&currentStock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	newStock := currentStock
	if input.Type == "in" {
		newStock += input.Amount
	} else {
		if input.Amount > currentStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock not enough to reduce"})
			return
		}
		newStock -= input.Amount
	}

	// Update stok produk
	_, err = db.DB.Exec("UPDATE products SET stock=$1, updated_at=$2 WHERE id=$3", newStock, time.Now(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	// Insert log perubahan stok
	_, err = db.DB.Exec(`
		INSERT INTO stock_logs (product_id, change_type, amount, note)
		VALUES ($1, $2, $3, $4)
	`, productID, input.Type, input.Amount, input.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log stock change"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully", "new_stock": newStock})
}
