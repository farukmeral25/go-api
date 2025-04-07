package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement"`
	FirstName string         `json:"first_name" binding:"required" form:"first_name"`
	LastName  string         `json:"last_name" binding:"required" form:"last_name"`
	Username  string         `json:"username" binding:"required" gorm:"unique" form:"username"`
	Password  string         `json:"password" binding:"required" form:"password"`
	Email     string         `json:"email" binding:"required,email" gorm:"unique" form:"email"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
