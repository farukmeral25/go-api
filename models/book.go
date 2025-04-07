package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Title     string         `json:"title" binding:"required" gorm:"size:255;not null"`
	Author    string         `json:"author" binding:"required" gorm:"size:255;not null"`
	Summary   string         `json:"summary" binding:"required" gorm:"type:text;not null"`
	ReadDate  time.Time      `json:"read_date" binding:"required"`
	Rating    int            `json:"rating" binding:"required,min=1,max=5" gorm:"not null"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      *User          `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

// BookListResponse kitap listesi için özet response
type BookListResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

// BookResponse detaylı kitap bilgileri için response
type BookResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Summary   string    `json:"summary"`
	ReadDate  time.Time `json:"read_date"`
	Rating    int       `json:"rating"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		ID        uint   `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"user"`
}
