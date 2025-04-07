package models

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID        uint      `json:"id"`
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"user"`
	} `json:"data"`
}
