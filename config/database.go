package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-api/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("SQLite veritabanına bağlanılamadı:", err)
	}

	// Veritabanı bağlantı havuzu ayarları
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Veritabanı bağlantı havuzu oluşturulamadı:", err)
	}

	// Bağlantı havuzu ayarları
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(0)

	// Otomatik migrasyon
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Veritabanı migrasyonu başarısız:", err)
	}
}

// Kullanıcı kaydetme fonksiyonu
func SaveUser(user *models.User) error {
	return DB.Create(user).Error
}

// Email ile kullanıcı bulma fonksiyonu
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Username ile kullanıcı bulma fonksiyonu
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
