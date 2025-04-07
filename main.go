package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-api/docs" // docs klasörü için import

	"go-api/config"
	"go-api/models"
	"go-api/routes"
)

// @title           Go API
// @version         1.0
// @description     Go API örneği
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

	// Komut satırı argümanlarını kontrol et
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		config.ConnectDatabase()
		config.DB.AutoMigrate(&models.User{}, &models.Book{})
		fmt.Println("Veritabanı tabloları oluşturuldu!")
		return
	}

	router := gin.Default()

	// CORS ayarları
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Veritabanı bağlantısı
	config.ConnectDatabase()

	// Swagger endpoint'i
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // Swagger JSON dosyasının URL'i
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Route'ları ayarla
	routes.SetupAuthRoutes(router)
	routes.SetupBookRoutes(router)

	// Port ayarı
	port := ":8000"
	fmt.Printf("Uygulama %s portunda başlatılıyor...\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Uygulama başlatılamadı: %v", err)
	}
}
