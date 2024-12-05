package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	AccessToken string `my_json:"access_token"`
	TokenType   string `my_json:"token_type"`
	ExpiresIn   int64  `my_json:"expires_in"`
}

type CustomClaims struct {
	UserID   int    `my_json:"user_id"`
	UserRole string `my_json:"user_role"`
	jwt.RegisteredClaims
}
