package transport_models

import (
	"time"
)

type Token struct {
	AccessToken string    `my_json:"access_token"`
	TokenType   string    `my_json:"token_type"`
	ExpiresAt   time.Time `json:"exp"`
}
