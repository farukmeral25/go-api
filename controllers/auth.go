package controllers

import (
	"net/http"
	"os"
	"time"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenDuration  = time.Hour
	refreshTokenDuration = time.Hour * 24 * 7 // 7 gün
)

func generateTokens(userID uint) (*models.TokenResponse, error) {
	// Access token oluştur
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	// Refresh token oluştur
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	// Token'ları imzala
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenDuration.Seconds()),
	}, nil
}

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

	// Token'ları oluştur
	tokens, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Token oluşturulamadı",
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
		"data": gin.H{
			"user": response,
			"tokens": gin.H{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
				"token_type":    tokens.TokenType,
				"expires_in":    tokens.ExpiresIn,
			},
		},
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

	// Token'ları oluştur
	tokens, err := generateTokens(user.ID)
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
			"user": gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"username":   user.Username,
				"email":      user.Email,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
			"tokens": gin.H{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
				"token_type":    tokens.TokenType,
				"expires_in":    tokens.ExpiresIn,
			},
		},
	})
}

// RefreshToken godoc
// @Summary      Token yenileme
// @Description  Refresh token ile yeni bir access token alır
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      models.RefreshTokenRequest  true  "Refresh token"
// @Success      200           {object}  models.TokenResponse
// @Failure      400           {object}  map[string]interface{}
// @Failure      401           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /api/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var refreshRequest models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz istek",
			"error":   err.Error(),
		})
		return
	}

	// Refresh token'ı doğrula
	token, err := jwt.ParseWithClaims(refreshRequest.RefreshToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Geçersiz refresh token",
			"error":   err.Error(),
		})
		return
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Geçersiz refresh token",
			"error":   "Token doğrulanamadı",
		})
		return
	}

	// Kullanıcıyı veritabanından bul
	var user models.User
	if err := config.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı",
			"error":   err.Error(),
		})
		return
	}

	// Yeni token'ları oluştur
	tokens, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Token oluşturulamadı",
			"error":   err.Error(),
		})
		return
	}

	// Kullanıcı bilgilerini hazırla
	userResponse := struct {
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

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token başarıyla yenilendi",
		"data": gin.H{
			"user": userResponse,
			"tokens": gin.H{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
				"token_type":    tokens.TokenType,
				"expires_in":    tokens.ExpiresIn,
			},
		},
	})
}
