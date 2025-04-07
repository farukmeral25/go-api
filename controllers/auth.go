package controllers

import (
	"net/http"
	"time"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Register godoc
// @Summary      Kullanıcı kaydı
// @Description  Yeni bir kullanıcı kaydı oluşturur
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "Kullanıcı bilgileri"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/auth/register [post]
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz istek",
			"error":   err.Error(),
		})
		return
	}

	// Şifreyi hashle
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Şifre hashleme hatası",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcıyı veritabanına kaydet
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Kullanıcı kaydı oluşturulamadı",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini döndür (şifre hariç)
	response := struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kullanıcı başarıyla oluşturuldu",
		"data":    response,
	})
}

// Login godoc
// @Summary      Kullanıcı girişi
// @Description  Kullanıcı girişi yapar ve JWT token döner
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.LoginRequest  true  "Giriş bilgileri"
// @Success      200         {object}  models.LoginResponse
// @Failure      400         {object}  map[string]interface{}
// @Failure      401         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Router       /api/auth/login [post]
func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz istek",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcıyı veritabanından bul
	var user models.User
	if err := config.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Kullanıcı adı veya şifre hatalı",
			"error":   "Kullanıcı bulunamadı",
		})
		return
	}

	// Şifreyi kontrol et
	if err := user.CheckPassword(loginRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Kullanıcı adı veya şifre hatalı",
			"error":   "Şifre hatalı",
		})
		return
	}

	// JWT token oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Token'ı imzala
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Token oluşturulamadı",
			"error":   err.Error(),
		})
		return
	}

	// Response'u döndür
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Giriş başarılı",
		"data": gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"username":   user.Username,
				"email":      user.Email,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
		},
	})
}
