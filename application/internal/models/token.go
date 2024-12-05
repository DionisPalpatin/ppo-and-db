package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}
