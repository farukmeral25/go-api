package controllers

import (
	"net/http"
	"strconv"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// CreateBook godoc
// @Summary      Kitap ekleme
// @Description  Yeni bir kitap ve özet ekler
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      models.Book  true  "Kitap bilgileri"
// @Success      201   {object}  models.BookResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/books [post]
func CreateBook(c *gin.Context) {
	// Kullanıcı ID'sini al
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Yetkilendirme hatası",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz istek",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı ID'sini ata
	book.UserID = userID.(uint)

	// Kitabı veritabanına kaydet
	if err := config.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kitap kaydedilemedi",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini al
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kullanıcı bilgileri alınamadı",
			"error":   err.Error(),
		})
		return
	}

	// Response hazırla
	response := models.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Summary:   book.Summary,
		ReadDate:  book.ReadDate,
		Rating:    book.Rating,
		Notes:     book.Notes,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
	response.User.ID = user.ID
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName
	response.User.Username = user.Username

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kitap başarıyla eklendi",
		"data":    response,
	})
}

// GetBooks godoc
// @Summary      Kitapları listele
// @Description  Kullanıcının kitaplarını listeler
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.BookListResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/books [get]
func GetBooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Yetkilendirme hatası",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	var books []models.Book
	if err := config.DB.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kitaplar alınamadı",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini al
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kullanıcı bilgileri alınamadı",
			"error":   err.Error(),
		})
		return
	}

	// Response hazırla
	var response []models.BookListResponse
	for _, book := range books {
		bookResponse := models.BookListResponse{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Rating: book.Rating,
		}
		response = append(response, bookResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kitaplar başarıyla getirildi",
		"data":    response,
	})
}

// GetBook godoc
// @Summary      Kitap detayı
// @Description  Belirtilen kitabın detaylarını getirir
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Kitap ID"
// @Success      200  {object}  models.BookResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/books/{id} [get]
func GetBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Yetkilendirme hatası",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz kitap ID",
			"error":   err.Error(),
		})
		return
	}

	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kitap bulunamadı",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini al
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kullanıcı bilgileri alınamadı",
			"error":   err.Error(),
		})
		return
	}

	// Response hazırla
	response := models.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Summary:   book.Summary,
		ReadDate:  book.ReadDate,
		Rating:    book.Rating,
		Notes:     book.Notes,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
	response.User.ID = user.ID
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName
	response.User.Username = user.Username

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kitap başarıyla getirildi",
		"data":    response,
	})
}

// UpdateBook godoc
// @Summary      Kitap güncelleme
// @Description  Belirtilen kitabı günceller
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "Kitap ID"
// @Param        book  body      models.Book  true  "Kitap bilgileri"
// @Success      200   {object}  models.BookResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/books/{id} [put]
func UpdateBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Yetkilendirme hatası",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz kitap ID",
			"error":   err.Error(),
		})
		return
	}

	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kitap bulunamadı",
			"error":   err.Error(),
		})
		return
	}

	// Yeni bilgileri bind et
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz istek",
			"error":   err.Error(),
		})
		return
	}

	// ID ve UserID değiştirilemesin
	book.ID = uint(bookID)
	book.UserID = userID.(uint)

	// Güncelle
	if err := config.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kitap güncellenemedi",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini al
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kullanıcı bilgileri alınamadı",
			"error":   err.Error(),
		})
		return
	}

	// Response hazırla
	response := models.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Summary:   book.Summary,
		ReadDate:  book.ReadDate,
		Rating:    book.Rating,
		Notes:     book.Notes,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
	response.User.ID = user.ID
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName
	response.User.Username = user.Username

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kitap başarıyla güncellendi",
		"data":    response,
	})
}

// DeleteBook godoc
// @Summary      Kitap silme
// @Description  Belirtilen kitabı siler
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Kitap ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Yetkilendirme hatası",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz kitap ID",
			"error":   err.Error(),
		})
		return
	}

	result := config.DB.Where("id = ? AND user_id = ?", bookID, userID).Delete(&models.Book{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kitap silinemedi",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kitap bulunamadı",
			"error":   "Kayıt bulunamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kitap başarıyla silindi",
	})
}
