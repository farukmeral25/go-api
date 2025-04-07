package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-api/docs" // docs klasörü için import

	"go-api/config"
	"go-api/routes"
)

// @title           Go API Auth Example
// @version         1.0
// @description     Bu API, JWT tabanlı kimlik doğrulama sistemi örneğidir.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api
// @schemes   http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	router := gin.Default()

	// CORS ayarları
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger endpoint'i
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // Swagger JSON dosyasının URL'i
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Route'ları ayarla
	routes.SetupAuthRoutes(router)

	// Port ayarı
	port := ":8000"
	log.Printf("Uygulama %s portunda başlatılıyor...", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Uygulama başlatılamadı: %v", err)
	}
}
