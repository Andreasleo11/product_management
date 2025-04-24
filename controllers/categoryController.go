package controllers

import (
	"backend_prodman/db"
	"backend_prodman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories - Ambil semua kategori
func GetCategories(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID - Ambil kategori berdasarkan ID
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	err := db.DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).
		Scan(&cat.ID, &cat.Name)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// CreateCategory - Tambah kategori baru
func CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek duplikat nama
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM categories WHERE name = $1)", cat.Name).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name already exists"})
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO categories (name) VALUES ($1) RETURNING id",
		cat.Name,
	).Scan(&cat.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// UpdateCategory - Edit kategori
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek duplikat tapi exclude current id
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM categories WHERE name = $1 AND id != $2)", cat.Name, id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name already exists"})
		return
	}

	_, err = db.DB.Exec("UPDATE categories SET name=$1 WHERE id=$2", cat.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

// DeleteCategory - Hapus kategori
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
