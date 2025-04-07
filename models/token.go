package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
