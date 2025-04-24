package controllers

import (
	"backend_prodman/db"
	"backend_prodman/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	// Ambil parameter query dari URL
	pageStr := c.DefaultQuery("page", "1")           // Halaman default 1
	limitStr := c.DefaultQuery("limit", "10")        // Limit default 10
	sortBy := c.DefaultQuery("sort_by", "name")      // Default sort by name
	sortOrder := c.DefaultQuery("sort_order", "asc") // Default sort order ascending

	// Konversi page dan limit ke integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Hitung offset untuk pagination
	offset := (page - 1) * limit

	// Validasi kolom yang dapat digunakan untuk sorting
	validSortColumns := map[string]bool{
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
	}

	if !validSortColumns[sortBy] {
		sortBy = "name" // Default ke name jika kolom tidak valid
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc" // Default ke ascending jika tidak valid
	}

	// Mulai query untuk produk
	query := fmt.Sprintf(`
		SELECT id, name, description, purchase_price, sell_price, category_id, stock, min_stock_alert, image_url, status, created_at, updated_at
		FROM products
		WHERE status = 'active'
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`, sortBy, sortOrder)

	// Eksekusi query dengan parameter pagination
	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Parsing hasil query ke slice of Product
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PurchasePrice, &p.SellPrice, &p.CategoryID, &p.Stock, &p.MinStockAlert, &p.ImageURL, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	// Jika produk tidak ditemukan
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}

	// Kembalikan hasil pencarian
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	err := db.DB.QueryRow(`
		SELECT id, name, description, purchase_price, sell_price, category_id, stock, min_stock_alert, image_url, status 
		FROM products WHERE id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.PurchasePrice, &p.SellPrice,
		&p.CategoryID, &p.Stock, &p.MinStockAlert, &p.ImageURL, &p.Status,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.QueryRow(`
		INSERT INTO products (name, description, purchase_price, sell_price, category_id, stock, min_stock_alert, image_url, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		p.Name, p.Description, p.PurchasePrice, p.SellPrice,
		p.CategoryID, p.Stock, p.MinStockAlert, p.ImageURL, p.Status,
	).Scan(&p.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(`
		UPDATE products 
		SET name=$1, description=$2, purchase_price=$3, sell_price=$4, category_id=$5,
		    stock=$6, min_stock_alert=$7, image_url=$8, status=$9 
		WHERE id=$10`,
		p.Name, p.Description, p.PurchasePrice, p.SellPrice, p.CategoryID,
		p.Stock, p.MinStockAlert, p.ImageURL, p.Status, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func GetProductsByCategoryID(c *gin.Context) {
	categoryID := c.Param("category_id")

	rows, err := db.DB.Query("SELECT id, name, description, purchase_price, sell_price, category_id, stock, min_stock_alert, image_url, status, created_at, updated_at FROM products WHERE category_id = $1", categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PurchasePrice, &p.SellPrice, &p.CategoryID, &p.Stock, &p.MinStockAlert, &p.ImageURL, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found for this category"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func SearchProducts(c *gin.Context) {
	// Ambil parameter query dari URL (misal ?name=sepatu&category=1&status=active)
	name := c.DefaultQuery("name", "") // Default jika tidak ada parameter
	categoryID := c.DefaultQuery("category_id", "")
	status := c.DefaultQuery("status", "")
	minPrice := c.DefaultQuery("min_price", "") // Tambahan untuk harga minimum
	maxPrice := c.DefaultQuery("max_price", "") // Tambahan untuk harga maksimum
	stock := c.DefaultQuery("stock", "")        // Tambahan untuk stok produk

	// Mulai query
	query := `SELECT id, name, description, purchase_price, sell_price, category_id, stock, min_stock_alert, image_url, status, created_at, updated_at 
			  FROM products WHERE 1=1`
	var args []interface{}
	paramIndex := 1 // Untuk tracking urutan parameter

	// Filter berdasarkan nama produk
	if name != "" {
		query += " AND name ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+name+"%")
		paramIndex++
	}

	// Filter berdasarkan ID kategori
	if categoryID != "" {
		query += " AND category_id = $" + strconv.Itoa(paramIndex)
		args = append(args, categoryID)
		paramIndex++
	}

	// Filter berdasarkan status produk
	if status != "" {
		query += " AND status = $" + strconv.Itoa(paramIndex)
		args = append(args, status)
		paramIndex++
	}

	// Filter berdasarkan harga minimum
	if minPrice != "" {
		query += " AND sell_price >= $" + strconv.Itoa(paramIndex)
		args = append(args, minPrice)
		paramIndex++
	}

	// Filter berdasarkan harga maksimum
	if maxPrice != "" {
		query += " AND sell_price <= $" + strconv.Itoa(paramIndex)
		args = append(args, maxPrice)
		paramIndex++
	}

	// Filter berdasarkan stok produk
	if stock != "" {
		query += " AND stock >= $" + strconv.Itoa(paramIndex)
		args = append(args, stock)
		paramIndex++
	}

	// Eksekusi query dengan parameter yang sudah digabungkan
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Parsing hasil query ke slice of Product
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PurchasePrice, &p.SellPrice, &p.CategoryID, &p.Stock, &p.MinStockAlert, &p.ImageURL, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	// Jika produk tidak ditemukan
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}

	// Kembalikan hasil pencarian
	c.JSON(http.StatusOK, products)
}
